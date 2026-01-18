# PostgreSQL Migration - JSON Type Issues Report

**Date:** 2026-01-17
**Migration:** MySQL → PostgreSQL
**Issue:** Incorrect JSON serialization of database types causing frontend crashes
**Status:** 🔴 CRITICAL - REQUIRES IMMEDIATE FIX

---

## Executive Summary

After migrating from MySQL to PostgreSQL, the backend is sending **incorrectly formatted JSON** to the frontend, which expects primitive types (strings, numbers, booleans) but is receiving complex objects. This is causing the frontend to crash or display incorrect data.

**Root Cause:** Go's `sql.NullString`, `sql.NullInt64`, and `sql.NullTime` types serialize to JSON as objects `{"String":"value","Valid":true}` instead of primitive values.

---

## Critical Problem: JSON Serialization Mismatch

### What Frontend Expects (MySQL Format)
```json
{
  "icon": "icon.png",
  "icon_uri": "uri://icon",
  "min_circle_range": 100,
  "phone_number": "+1234567890"
}
```

### What Backend is Sending (PostgreSQL Format)
```json
{
  "icon": {"String": "icon.png", "Valid": true},
  "icon_uri": {"String": "uri://icon", "Valid": true},
  "min_circle_range": {"Int64": 100, "Valid": true},
  "phone_number": {"String": "+1234567890", "Valid": true}
}
```

**Impact:** Frontend cannot parse these objects. It expects `icon` to be a string but receives an object with `String` and `Valid` fields.

---

## Affected Files and Fields

### 1. 🔴 `/internal/getclustersbylocation/model.go`

**CRITICAL:** This endpoint returns cluster data to the map view.

**Problematic Fields:**
```go
type Subcategory struct {
    InsuId             int64          `json:"insu_id"`
    IncaId             int64          `json:"inca_id"`
    Name               string         `json:"name"`
    Description        string         `json:"description"`
    Icon               sql.NullString `json:"icon"`        // ❌ WRONG TYPE
    IconURI            sql.NullString `json:"icon_uri"`    // ❌ WRONG TYPE
    Code               string         `json:"code"`
    MinCircleRange     sql.NullInt64  `json:"min_circle_range"`      // ❌ WRONG TYPE
    MaxCircleRange     sql.NullInt64  `json:"max_circle_range"`      // ❌ WRONG TYPE
    DefaultCircleRange sql.NullInt64  `json:"default_circle_range"`  // ❌ WRONG TYPE
    CategoryCode       string         `json:"category_code"`
    SubcategoryCode    string         `json:"subcategory_code"`
}
```

**Expected JSON:**
```json
{
  "icon": "icon.png",
  "icon_uri": "uri://icon",
  "min_circle_range": 100
}
```

**Actual JSON:**
```json
{
  "icon": {"String": "icon.png", "Valid": true},
  "icon_uri": {"String": "uri://icon", "Valid": true},
  "min_circle_range": {"Int64": 100, "Valid": true}
}
```

---

### 2. 🔴 `/internal/getclusterby/model.go`

**CRITICAL:** This endpoint returns detailed cluster information.

**Problematic Fields:**
```go
type Cluster struct {
    InclId                 int64              `json:"incl_id"`
    CreatedAt              sql.NullTime       `json:"created_at"`      // ❌ WRONG TYPE
    StartTime              sql.NullTime       `json:"start_time"`      // ❌ WRONG TYPE
    EndTime                sql.NullTime       `json:"end_time"`        // ❌ WRONG TYPE
    AccountId              sql.NullInt64      `json:"account_id"`      // ❌ WRONG TYPE
    // ... other fields
}

type Incident struct {
    Status           sql.NullString    `json:"status"`      // ❌ WRONG TYPE
    // ... other fields
}

type Comment struct {
    CreatedAt        common.NullTime   `json:"created_at"`  // ✅ CORRECT (custom type)
    // ... other fields
}
```

**Note:** `common.NullTime` has custom MarshalJSON and works correctly. `sql.NullTime` does NOT.

---

### 3. 🟡 `/internal/auth/model.go`

**Problematic Fields:**
```go
type User struct {
    AccountID           int64          `json:"account_id"`
    Email               string         `json:"email"`
    PhoneNumber         sql.NullString `json:"phone_number"`  // ❌ WRONG TYPE
    // ... other fields
}
```

**Impact:** Login/signup responses include phone_number as object instead of string.

---

### 4. 🟡 `/internal/getcategories/model.go`

**Problematic Fields:**
```go
type Category struct {
    Icon        *sql.NullString `json:"icon"`  // ❌ WRONG TYPE (pointer to NullString)
    // ... other fields
}
```

**Double Problem:** Not only is it `sql.NullString`, but it's also a **pointer**, making JSON even more complex.

---

### 5. 🟡 `/internal/getsubcategoriesbycategoryid/model.go`

**Problematic Fields:**
```go
type Subcategory struct {
    Icon               *sql.NullString `json:"icon"`        // ❌ WRONG TYPE
    MinCircleRange     *sql.NullInt64  `json:"min_circle_range"`      // ❌ WRONG TYPE
    MaxCircleRange     *sql.NullInt64  `json:"max_circle_range"`      // ❌ WRONG TYPE
    DefaultCircleRange *sql.NullInt64  `json:"default_circle_range"`  // ❌ WRONG TYPE
}
```

---

## Why This Worked in MySQL

In MySQL, the Go driver **always returns concrete values** for these fields:
- NULL strings → empty string `""`
- NULL integers → `0`
- NULL times → zero time

The structs used `sql.Null*` types to handle NULL detection, but JSON serialization still worked because Go's JSON encoder would serialize the **underlying value**, not the wrapper object.

**PostgreSQL's driver is stricter:** It returns `sql.Null*` types that serialize as objects with `{Value, Valid}` structure.

---

## Solutions

### ✅ Option 1: Use Custom Wrapper Types (RECOMMENDED)

You already have `common.NullTime` which implements `MarshalJSON` correctly. Create similar wrappers for all types.

**Example: Create `/internal/common/nulltypes.go`**

```go
package common

import (
    "database/sql"
    "encoding/json"
)

// NullString wraps sql.NullString with proper JSON serialization
type NullString struct {
    sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
    if !ns.Valid {
        return []byte("null"), nil
    }
    return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
    if string(b) == "null" {
        ns.Valid = false
        return nil
    }
    var s string
    if err := json.Unmarshal(b, &s); err != nil {
        return err
    }
    ns.String = s
    ns.Valid = true
    return nil
}

// NullInt64 wraps sql.NullInt64 with proper JSON serialization
type NullInt64 struct {
    sql.NullInt64
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
    if !ni.Valid {
        return []byte("null"), nil
    }
    return json.Marshal(ni.Int64)
}

func (ni *NullInt64) UnmarshalJSON(b []byte) error {
    if string(b) == "null" {
        ni.Valid = false
        return nil
    }
    var i int64
    if err := json.Unmarshal(b, &i); err != nil {
        return err
    }
    ni.Int64 = i
    ni.Valid = true
    return nil
}

// NullFloat64 wraps sql.NullFloat64 with proper JSON serialization
type NullFloat64 struct {
    sql.NullFloat64
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) {
    if !nf.Valid {
        return []byte("null"), nil
    }
    return json.Marshal(nf.Float64)
}

func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
    if string(b) == "null" {
        nf.Valid = false
        return nil
    }
    var f float64
    if err := json.Unmarshal(b, &f); err != nil {
        return err
    }
    nf.Float64 = f
    nf.Valid = true
    return nil
}
```

**Then replace all `sql.Null*` types in models:**

```go
// ❌ BEFORE
type Subcategory struct {
    Icon               sql.NullString `json:"icon"`
    MinCircleRange     sql.NullInt64  `json:"min_circle_range"`
}

// ✅ AFTER
type Subcategory struct {
    Icon               common.NullString `json:"icon"`
    MinCircleRange     common.NullInt64  `json:"min_circle_range"`
}
```

---

### ⚠️ Option 2: Use Primitive Types with COALESCE (SIMPLE BUT LOSES NULL DETECTION)

Replace `sql.Null*` types with primitive types and use `COALESCE` in SQL queries.

**Model Changes:**
```go
// ❌ BEFORE
type Subcategory struct {
    Icon               sql.NullString `json:"icon"`
    MinCircleRange     sql.NullInt64  `json:"min_circle_range"`
}

// ✅ AFTER
type Subcategory struct {
    Icon               string `json:"icon"`
    MinCircleRange     int64  `json:"min_circle_range"`
}
```

**Query Changes:**
```sql
-- ❌ BEFORE
SELECT icon, min_circle_range FROM subcategories

-- ✅ AFTER
SELECT COALESCE(icon, '') as icon, COALESCE(min_circle_range, 0) as min_circle_range FROM subcategories
```

**Downside:** You lose the ability to distinguish between NULL and empty string/0.

---

## Recommended Action Plan

### Phase 1: Fix Critical Endpoints (IMMEDIATE)

1. Create `/internal/common/nulltypes.go` with custom wrappers for:
   - `NullString`
   - `NullInt64`
   - `NullFloat64`
   - `NullTime` (already exists, verify it's correct)

2. Replace `sql.Null*` types in these critical files:
   - `/internal/getclustersbylocation/model.go` (map view data)
   - `/internal/getclusterby/model.go` (cluster detail view)
   - `/internal/auth/model.go` (login/signup)

### Phase 2: Fix Remaining Endpoints

3. Replace `sql.Null*` in:
   - `/internal/getcategories/model.go`
   - `/internal/getsubcategoriesbycategoryid/model.go`
   - `/internal/profile/model.go` (if affected)

### Phase 3: Testing

4. Test JSON responses from each endpoint:
   ```bash
   # Test map endpoint
   curl "http://localhost:8080/getclustersbylocation/..." | jq .

   # Verify icon is string, not object
   # Expected: "icon": "icon.png"
   # NOT: "icon": {"String": "icon.png", "Valid": true}
   ```

---

## Verification Commands

After fixes, verify JSON format:

```bash
# Test that NullString serializes correctly
cd /Users/garyeikoow/Desktop/alertly/backend
go run -exec echo "Testing JSON serialization" cmd/app/main.go

# Start server
go run cmd/app/main.go

# Test endpoints
curl -X GET "http://localhost:8080/api/categories" | jq '.data[0].icon'
# Expected output: "icon.png" (string)
# NOT: {"String":"icon.png","Valid":true} (object)
```

---

## Files That Need Modification

### High Priority (CRITICAL - Breaks Frontend)
1. `/internal/getclustersbylocation/model.go` - Map view
2. `/internal/getclusterby/model.go` - Cluster details
3. `/internal/auth/model.go` - Authentication

### Medium Priority (Affects UI/UX)
4. `/internal/getcategories/model.go` - Category list
5. `/internal/getsubcategoriesbycategoryid/model.go` - Subcategory data
6. `/internal/profile/model.go` - User profiles

### Low Priority (Edge Cases)
7. Any other files using `sql.Null*` types in JSON responses

---

## Additional Notes

### Why `common.NullTime` Works Correctly

The existing `common.NullTime` already has custom JSON serialization:

```go
// ✅ CORRECT IMPLEMENTATION
type NullTime struct {
    sql.NullTime
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
    if !nt.Valid {
        return []byte("null"), nil
    }
    return json.Marshal(nt.Time.Format("2006-01-02T15:04:05Z07:00"))
}
```

This ensures it serializes as a string (or null), not as an object.

### Why `dbtypes.NullBool` Works Correctly

The `dbtypes.NullBool` also has proper JSON serialization:

```go
// ✅ CORRECT IMPLEMENTATION
func (nb NullBool) MarshalJSON() ([]byte, error) {
    if !nb.Valid {
        return []byte("null"), nil
    }
    return json.Marshal(nb.Bool)
}
```

**Apply the same pattern to all `sql.Null*` types.**

---

## Summary

**Problem:** `sql.NullString`, `sql.NullInt64`, `sql.NullTime` serialize to JSON as objects instead of primitive values.

**Solution:** Create custom wrapper types with `MarshalJSON/UnmarshalJSON` methods (following the pattern of `common.NullTime` and `dbtypes.NullBool`).

**Impact:** This is a **CRITICAL** issue causing frontend crashes. Must be fixed before production deployment.

**Estimated Effort:** 2-4 hours to implement custom wrappers and update all affected models.
