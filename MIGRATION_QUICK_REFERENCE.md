# PostgreSQL Migration Quick Reference

## Files Modified (9 total)

### 1. newincident/service.go
**Line 76-80:** Media type empty string fix
```go
mediaType := incident.MediaType
if mediaType == "" {
    mediaType = "image" // Default to image if not specified
}
```

### 2. getclustersbylocation/repository.go
**Line 35:** `AND TRIM(t1.is_active) = '1'`

### 3. getclusterby/repository.go
**Line 44:** `activeFilter = "AND TRIM(c.is_active) = '1'"`
**Line 124:** `AND TRIM(COALESCE(r.is_active, '0')) = '1'`

### 4. cronjobs/cjincidentexpiration/repository.go
**Line 51:** `TRIM(ic.is_active) = '1'`

### 5. getclusterbyradius/repository.go
**Line 53:** `AND TRIM(t1.is_active) = '1'`

### 6. getincidentsasreels/repository.go
**Line 41:** `AND TRIM(c.is_active) = '1'`

### 7. referrals/repository.go
**Lines 117, 402, 464, 483, 604:** `TRIM(is_active) = '1'` or `TRIM(i.is_active) = '1'`
**Lines 314, 320, 350, 355:** `DATE()` → `::date`

### 8. cronjobs/cronjobs/repository.go
**Lines 54, 64:** `TRIM(ic.is_active) = '1'`

### 9. cronjobs/cjbot_creator/repository.go
**Line 290:** `AND TRIM(is_active) = '1'`

---

## Pattern Changes

### Boolean Comparisons (CHAR(1))
```sql
-- BEFORE
WHERE is_active = '1'

-- AFTER
WHERE TRIM(is_active) = '1'
```

### Date Casting
```sql
-- BEFORE (MySQL)
DATE(registered_at)
GROUP BY DATE(registered_at)

-- AFTER (PostgreSQL)
registered_at::date
GROUP BY registered_at::date
```

### Media Type Handling
```go
// BEFORE
MediaType: incident.MediaType  // Could be ""

// AFTER
MediaType: mediaType  // Defaults to "image" if empty
```

---

## Quick Test Commands

```bash
# Test incident creation
curl -X POST http://localhost:8080/newincident -F "description=Test"

# Test cluster filtering
curl "http://localhost:8080/getclustersbylocation?lat=43.65&lng=-79.38"

# Test referral metrics
curl "http://localhost:8080/referrals/metrics/daily?code=TEST&days=7"

# Run cronjob
go run cmd/cronjob/main.go
```

---

## Deployment Checklist

- [ ] Pull latest code with fixes
- [ ] Run `go build cmd/app/main.go`
- [ ] Stop current service
- [ ] Deploy new binary
- [ ] Start service
- [ ] Monitor logs for errors
- [ ] Test critical endpoints
- [ ] Verify cronjob execution

---

## Rollback

```bash
git revert HEAD
go build -o alertly-backend cmd/app/main.go
systemctl restart alertly-backend
```
