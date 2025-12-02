# Bot Creator (Data Seeder) - Implementation Guide

## Overview

The `cjbot_creator` cronjob system automatically ingests real-time incident data from 5 Toronto data sources to populate Alertly's database and prevent "Cold Start" issues.

## Architecture

### Directory Structure

```
backend/internal/cronjobs/cjbot_creator/
├── model.go               # Data models and OfficialAssetsMap
├── repository.go          # Database operations
├── normalizer.go          # Category mapping logic
├── service.go             # Main orchestrator + geocoding
└── scrapers/
    ├── types.go           # ScrapedIncident and Point types
    ├── tfs.go             # Toronto Fire Services scraper
    ├── hydro.go           # Toronto Hydro outages scraper
    ├── tps.go             # (TODO) Toronto Police scraper
    ├── ttc.go             # (TODO) TTC transit alerts scraper
    └── weather.go         # (TODO) Environment Canada scraper
```

### Database Tables

Two new tables have been created:

1. **`bot_incident_hashes`** - Deduplication using SHA256 hashes
2. **`geocoding_cache`** - Nominatim geocoding cache (avoids API rate limits)

**Migration file:** `backend/assets/db/bot_seeder_tables.sql`

**To apply migration:**
```bash
mysql -u root -p alertly < backend/assets/db/bot_seeder_tables.sql
```

## Features Implemented

### ✅ Completed (2 Scrapers)

1. **TFS (Toronto Fire Services)** - `scrapers/tfs.go`
   - HTML parsing with goquery
   - Mock data for testing
   - Geocoding support

2. **Toronto Hydro** - `scrapers/hydro.go`
   - JSON API consumption
   - Polygon to centroid calculation
   - Estimated Time of Restoration (ETR) tracking

### ⏳ Pending Implementation (3 Scrapers)

3. **TPS (Toronto Police)** - Placeholder created
4. **TTC (Transit)** - Placeholder created
5. **Environment Canada Weather** - Placeholder created

## How It Works

### 1. Scraping Pipeline

```
Source Website → Scraper → ScrapedIncident → Normalizer → NormalizedIncident → Database
                              ↓
                         Deduplication (SHA256 hash)
                              ↓
                         Geocoding (Nominatim + Cache)
                              ↓
                         Save to incident_reports
```

### 2. Geocoding with Cache

- **Provider:** Nominatim (OpenStreetMap) - FREE
- **Rate Limit:** 1 request/second (respected via `golang.org/x/time/rate`)
- **Cache:** MySQL `geocoding_cache` table
- **Cache TTL:** 30 days (auto-cleanup)

### 3. Deduplication

Each incident generates a unique hash:
```go
SHA256(source + external_id + timestamp)
```

Before inserting, checks `bot_incident_hashes` table.

### 4. Category Normalization

Intelligent pattern matching maps source categories to Alertly categories:

**Example mappings:**
- TPS "ASSAULT" → `crime` / `assault`
- TFS "MEDICAL CALL" → `medical_emergency` / `trauma`
- Hydro "UNPLANNED OUTAGE" → `infrastructure_issues` / `utility_issues`
- Weather "SEVERE THUNDERSTORM" → `extreme_weather` / `thunderstorm`

**See:** `normalizer.go` for full mapping rules

### 5. Official Report Images

Static images hosted on S3:

```
https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/{category}.webp
```

**Required images (900x1200px WebP):**
- `crime.webp`
- `traffic_accident.webp`
- `medical_emergency.webp`
- `fire_incident.webp`
- `vandalism.webp`
- `suspicious_activity.webp`
- `infrastructure_issues.webp`
- `extreme_weather.webp`
- `community_events.webp`
- `dangerous_wildlife_sighting.webp`
- `positive_actions.webp`
- `lost_pet.webp`

## AWS Lambda Invocation

### Lambda Event Format

```json
{
  "task": "bot_creator_tfs"
}
```

### Available Tasks

| Task Name              | Source                | Frequency     | Status         |
|------------------------|-----------------------|---------------|----------------|
| `bot_creator_tps`      | Toronto Police        | Every 5 min   | ⏳ Pending     |
| `bot_creator_tfs`      | Toronto Fire          | Every 10 min  | ✅ Implemented |
| `bot_creator_ttc`      | TTC Transit           | Every 15 min  | ⏳ Pending     |
| `bot_creator_hydro`    | Toronto Hydro         | Every 30 min  | ✅ Implemented |
| `bot_creator_weather`  | Environment Canada    | Every 1 hour  | ⏳ Pending     |

### EventBridge Schedule (Example)

```yaml
# AWS EventBridge Rules
- Name: BotCreator-TFS
  Schedule: rate(10 minutes)
  Target: cronjob-lambda
  Input: { "task": "bot_creator_tfs" }

- Name: BotCreator-Hydro
  Schedule: rate(30 minutes)
  Target: cronjob-lambda
  Input: { "task": "bot_creator_hydro" }
```

## Local Testing

### Run a Specific Scraper

```bash
cd backend

# Test TFS scraper
go run -ldflags "-X main.task=bot_creator_tfs" cmd/cronjob/main.go

# Test Hydro scraper
go run -ldflags "-X main.task=bot_creator_hydro" cmd/cronjob/main.go
```

### Using Mock Data

Both implemented scrapers fall back to mock data if real endpoints fail:

```go
// In TFS scraper
incidents, err := scraper.Scrape()
if err != nil {
    log.Printf("Real scraping failed, using mock data")
    incidents = scraper.ScrapeMockData()
}
```

## Next Steps

### Phase 1: Complete Remaining Scrapers

1. **TPS (Police)** - Investigate actual endpoint
   - Likely URL: `https://data.torontopolice.on.ca/pages/calls-for-service`
   - May require API key or scraping setup

2. **TTC (Transit)** - RSS Feed
   - URL: `https://www.ttc.ca/Service_Advisories/all_service_alerts.rss`
   - Parse XML with `encoding/xml`

3. **Weather (Environment Canada)** - CAP Protocol
   - URL: `https://dd.weather.gc.ca/alerts/cap/`
   - Filter by Toronto area code

### Phase 2: Upload Static Images to S3

```bash
# Example using AWS CLI
aws s3 cp crime.webp s3://alertly-images-production/incidents/crime.webp --acl public-read
aws s3 cp fire_incident.webp s3://alertly-images-production/incidents/fire_incident.webp --acl public-read
# ... repeat for all 12 categories
```

### Phase 3: API Endpoint Research

Before implementing remaining scrapers, you'll need to provide:

1. **Real API/RSS URLs** for each source
2. **Example API responses** (JSON/XML structure)
3. **Authentication requirements** (API keys, tokens)

**Helpful commands for research:**
```bash
# Inspect TPS API calls
curl 'https://data.torontopolice.on.ca/api/...' -v | jq

# Test TTC RSS feed
curl 'https://www.ttc.ca/Service_Advisories/all_service_alerts.rss' | xmllint --format -

# Check Environment Canada CAP alerts
curl 'https://dd.weather.gc.ca/alerts/cap/' | grep Toronto
```

### Phase 4: Deployment

1. Apply database migration
2. Build and deploy Lambda function
3. Configure EventBridge schedules
4. Upload static images to S3
5. Monitor CloudWatch logs

## Monitoring & Maintenance

### Key Metrics to Monitor

- **Incident creation rate** (incidents/hour by source)
- **Geocoding cache hit ratio** (should be >70%)
- **Deduplication effectiveness** (blocked duplicates)
- **Scraping failures** (check CloudWatch logs)

### Cleanup Cronjobs (Optional)

Consider adding these periodic tasks:

```go
// In cmd/cronjob/main.go
case "bot_cleanup_hashes":
    repo := cjbot_creator.NewRepository(database.DB)
    deleted, _ := repo.CleanupExpiredHashes()
    log.Printf("Deleted %d expired hashes", deleted)

case "bot_cleanup_geocache":
    repo := cjbot_creator.NewRepository(database.DB)
    deleted, _ := repo.CleanupOldGeocodeCache()
    log.Printf("Deleted %d old geocode cache entries", deleted)
```

## Dependencies

New dependencies added to `go.mod`:

```go
github.com/PuerkitoBio/goquery v1.11.0  // HTML parsing
golang.org/x/sync v0.18.0               // errgroup for concurrency
golang.org/x/time v0.8.0                // rate limiter
```

## Troubleshooting

### Issue: "Geocoding failed"

**Cause:** Nominatim rate limit exceeded
**Solution:** Increase rate limiter interval or check cache

### Issue: "Subcategory not found"

**Cause:** Category code doesn't exist in `incident_subcategories` table
**Solution:** Add missing subcategory to database or update normalizer mapping

### Issue: "Invalid UTF-8 encoding"

**Cause:** Emoji characters in log statements
**Solution:** Already fixed - avoid emojis in production code

## Contact

For questions about implementation:
- Review `normalizer.go` for category mapping logic
- Review `service.go` for geocoding and processing pipeline
- Check CloudWatch logs for runtime errors

---

**Generated:** 2025-01-22
**Version:** 1.0
**Status:** 2/5 Scrapers Implemented
