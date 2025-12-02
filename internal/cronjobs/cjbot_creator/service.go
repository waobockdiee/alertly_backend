package cjbot_creator

import (
	"alertly/internal/cronjobs/cjbot_creator/scrapers"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

const (
	NOMINATIM_URL        = "https://nominatim.openstreetmap.org/search"
	NOMINATIM_USER_AGENT = "Alertly/1.0 (https://alertly.app)" // Required by Nominatim usage policy
	MAX_CONCURRENT_JOBS  = 5                                   // Limit concurrent geocoding requests
)

// Service orchestrates bot incident creation from all sources
type Service struct {
	repo         *Repository
	normalizer   *Normalizer
	geocoder     *Geocoder
	httpClient   *http.Client
	rateLimiter  *rate.Limiter
}

// NewService creates a new bot creator service
func NewService(repo *Repository) *Service {
	return &Service{
		repo:        repo,
		normalizer:  NewNormalizer(),
		geocoder:    NewGeocoder(repo),
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		rateLimiter: rate.NewLimiter(rate.Every(1*time.Second), 1), // 1 req/sec for Nominatim
	}
}

// ============================================
// SCRAPER RUN METHODS (One per source)
// ============================================

// RunTFS executes Toronto Fire Services scraper
func (s *Service) RunTFS() {
	log.Printf("🚒 [TFS] Starting bot creator job...")
	startTime := time.Now()

	scraper := scrapers.NewTFSScraper()

	// Try real scraping first, fallback to mock data
	incidents, err := scraper.Scrape()
	if err != nil {
		log.Printf("⚠️ [TFS] Real scraping failed: %v. Using mock data.", err)
		incidents = scraper.ScrapeMockData()
	}

	if len(incidents) == 0 {
		log.Printf("ℹ️ [TFS] No active incidents found")
		return
	}

	// Process incidents concurrently with rate limiting
	processed := s.processIncidents(incidents)

	duration := time.Since(startTime)
	log.Printf("✅ [TFS] Job completed in %v. Processed %d/%d incidents", duration, processed, len(incidents))
}

// RunHydro executes Toronto Hydro outages scraper
func (s *Service) RunHydro() {
	log.Printf("⚡ [Hydro] Starting bot creator job...")
	startTime := time.Now()

	scraper := scrapers.NewHydroScraper()

	// Try real scraping first, fallback to mock data
	incidents, err := scraper.Scrape()
	if err != nil {
		log.Printf("⚠️ [Hydro] Real scraping failed: %v. Using mock data.", err)
		incidents = scraper.ScrapeMockData()
	}

	if len(incidents) == 0 {
		log.Printf("ℹ️ [Hydro] No active outages found")
		return
	}

	// Process incidents concurrently with rate limiting
	processed := s.processIncidents(incidents)

	duration := time.Since(startTime)
	log.Printf("✅ [Hydro] Job completed in %v. Processed %d/%d incidents", duration, processed, len(incidents))
}

// RunTPS executes Toronto Police Service scraper
func (s *Service) RunTPS() {
	log.Printf("🚔 [TPS] Starting bot creator job...")
	startTime := time.Now()

	scraper := scrapers.NewTPSScraper()

	// Fetch real-time calls for service
	incidents, err := scraper.Scrape()
	if err != nil {
		log.Printf("❌ [TPS] Scraping failed: %v", err)
		return
	}

	if len(incidents) == 0 {
		log.Printf("ℹ️ [TPS] No active calls for service found")
		return
	}

	// Process incidents concurrently with rate limiting
	processed := s.processIncidents(incidents)

	duration := time.Since(startTime)
	log.Printf("✅ [TPS] Job completed in %v. Processed %d/%d incidents", duration, processed, len(incidents))
}

// RunTTC executes TTC transit alerts scraper (placeholder)
func (s *Service) RunTTC() {
	log.Printf("🚇 [TTC] Starting bot creator job...")
	log.Printf("⚠️ [TTC] Scraper not yet implemented. RSS/API endpoint needed.")
}

// RunWeather executes Environment Canada weather alerts scraper (placeholder)
func (s *Service) RunWeather() {
	log.Printf("🌤️ [Weather] Starting bot creator job...")
	log.Printf("⚠️ [Weather] Scraper not yet implemented. CAP RSS endpoint needed.")
}

// ============================================
// INCIDENT PROCESSING PIPELINE
// ============================================

// processIncidents handles concurrent processing of scraped incidents
func (s *Service) processIncidents(incidents []scrapers.ScrapedIncident) int {
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(MAX_CONCURRENT_JOBS)

	successCount := 0

	for i := range incidents {
		incident := incidents[i] // Capture variable
		g.Go(func() error {
			if s.processAndSaveIncident(ctx, incident) {
				successCount++
			}
			return nil // Never fail the entire batch
		})
	}

	g.Wait()
	return successCount
}

// processAndSaveIncident handles the full pipeline for a single incident
func (s *Service) processAndSaveIncident(ctx context.Context, scraped scrapers.ScrapedIncident) bool {
	// 1. Check for duplicates
	hash := GenerateIncidentHash(scraped.Source, scraped.ExternalID, scraped.Timestamp)
	exists, err := s.repo.CheckHashExists(hash)
	if err != nil {
		log.Printf("⚠️ [%s] Error checking hash for %s: %v", scraped.Source, scraped.ExternalID, err)
		return false
	}
	if exists {
		log.Printf("⏭️ [%s] Skipping duplicate incident: %s", scraped.Source, scraped.ExternalID)
		return false
	}

	// 2. Normalize incident to Alertly schema
	normalized, err := s.normalizer.Normalize(scraped)
	if err != nil {
		log.Printf("⚠️ [%s] Failed to normalize %s: %v", scraped.Source, scraped.ExternalID, err)
		return false
	}

	// 3. Geocode address if coordinates not available
	if normalized.Latitude == 0 && normalized.Longitude == 0 && scraped.Address != "" {
		lat, lng, addr, err := s.geocoder.Geocode(ctx, scraped.Address)
		if err != nil {
			log.Printf("⚠️ [%s] Geocoding failed for '%s': %v", scraped.Source, scraped.Address, err)
			return false
		}
		normalized.Latitude = lat
		normalized.Longitude = lng
		normalized.Address = addr
		normalized.City = "Toronto"
		normalized.Province = "ON"
	}

	// 4. Save to database
	reportID, err := s.repo.SaveBotIncident(*normalized)
	if err != nil {
		log.Printf("⚠️ [%s] Failed to save incident %s: %v", scraped.Source, scraped.ExternalID, err)
		return false
	}

	// 5. Save hash to prevent duplicates
	hashRecord := BotIncidentHash{
		Hash:         hash,
		Source:       scraped.Source,
		ExternalID:   scraped.ExternalID,
		CategoryCode: normalized.CategoryCode,
		CreatedAt:    time.Now(),
		ExpiresAt:    getExpirationTime(normalized.CategoryCode),
	}
	if err := s.repo.SaveIncidentHash(hashRecord); err != nil {
		log.Printf("⚠️ [%s] Failed to save hash (incident saved): %v", scraped.Source, err)
	}

	log.Printf("✅ [%s] Saved incident #%d: %s at (%.4f, %.4f)",
		scraped.Source, reportID, normalized.Title, normalized.Latitude, normalized.Longitude)

	return true
}

// getExpirationTime calculates TTL for bot incidents
func getExpirationTime(categoryCode string) *time.Time {
	var hours int
	switch categoryCode {
	case "infrastructure_issues": // Power outages may last longer
		hours = 24
	case "extreme_weather": // Weather alerts
		hours = 12
	case "fire_incident":
		hours = 6
	default:
		hours = 6
	}
	expiration := time.Now().Add(time.Duration(hours) * time.Hour)
	return &expiration
}

// ============================================
// GEOCODER (Nominatim with Cache)
// ============================================

type Geocoder struct {
	repo        *Repository
	httpClient  *http.Client
	rateLimiter *rate.Limiter
}

func NewGeocoder(repo *Repository) *Geocoder {
	return &Geocoder{
		repo:        repo,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		rateLimiter: rate.NewLimiter(rate.Every(1*time.Second), 1),
	}
}

type NominatimResponse struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

// Geocode converts address to coordinates using Nominatim + cache
func (g *Geocoder) Geocode(ctx context.Context, address string) (lat, lng float64, fullAddress string, err error) {
	// 1. Normalize address
	normalized := NormalizeAddress(address)
	addressHash := GenerateAddressHash(normalized)

	// 2. Check cache
	cached, err := g.repo.GetCachedGeocode(addressHash)
	if err == nil && cached != nil {
		log.Printf("🎯 Cache HIT for: %s", address)
		return cached.Latitude, cached.Longitude, cached.NormalizedAddress, nil
	}

	// 3. Rate limit (respect Nominatim's 1 req/sec policy)
	if err := g.rateLimiter.Wait(ctx); err != nil {
		return 0, 0, "", fmt.Errorf("rate limiter error: %w", err)
	}

	// 4. Call Nominatim API
	log.Printf("🌐 Geocoding (cache MISS): %s", address)

	apiURL := fmt.Sprintf("%s?q=%s&format=json&limit=1", NOMINATIM_URL, url.QueryEscape(normalized))
	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	req.Header.Set("User-Agent", NOMINATIM_USER_AGENT)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return 0, 0, "", fmt.Errorf("nominatim request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var results []NominatimResponse
	if err := json.Unmarshal(body, &results); err != nil {
		return 0, 0, "", fmt.Errorf("parsing nominatim response: %w", err)
	}

	if len(results) == 0 {
		return 0, 0, "", fmt.Errorf("no results found for address: %s", address)
	}

	// 5. Parse coordinates
	result := results[0]
	var latitude, longitude float64
	fmt.Sscanf(result.Lat, "%f", &latitude)
	fmt.Sscanf(result.Lon, "%f", &longitude)

	// 6. Save to cache
	cache := GeocodingCache{
		AddressHash:       addressHash,
		OriginalAddress:   address,
		NormalizedAddress: normalized,
		Latitude:          latitude,
		Longitude:         longitude,
		Source:            "nominatim",
		CreatedAt:         time.Now(),
	}
	if err := g.repo.SaveGeocodeCache(cache); err != nil {
		log.Printf("⚠️ Failed to cache geocode result: %v", err)
	}

	return latitude, longitude, result.DisplayName, nil
}
