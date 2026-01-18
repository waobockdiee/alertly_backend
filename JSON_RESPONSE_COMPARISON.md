# JSON Response Comparison - MySQL vs PostgreSQL

**Date:** 2026-01-17
**Issue:** Frontend expects MySQL-style JSON but receives PostgreSQL-style JSON

---

## Example 1: GetClustersByLocation Endpoint

### MySQL Response (WORKING - Before Migration)
```json
{
  "clusters": [
    {
      "incl_id": 123,
      "latitude": 45.5017,
      "longitude": -73.5673,
      "insu_id": 1,
      "category_code": "crime",
      "subcategory_code": "theft",
      "subcategory": {
        "insu_id": 1,
        "inca_id": 1,
        "name": "Theft",
        "description": "Theft or burglary",
        "icon": "theft-icon.png",
        "icon_uri": "uri://icons/theft",
        "code": "theft",
        "min_circle_range": 50,
        "max_circle_range": 500,
        "default_circle_range": 100,
        "category_code": "crime",
        "subcategory_code": "theft"
      }
    }
  ]
}
```

### PostgreSQL Response (BROKEN - After Migration)
```json
{
  "clusters": [
    {
      "incl_id": 123,
      "latitude": 45.5017,
      "longitude": -73.5673,
      "insu_id": 1,
      "category_code": "crime",
      "subcategory_code": "theft",
      "subcategory": {
        "insu_id": 1,
        "inca_id": 1,
        "name": "Theft",
        "description": "Theft or burglary",
        "icon": {
          "String": "theft-icon.png",
          "Valid": true
        },
        "icon_uri": {
          "String": "uri://icons/theft",
          "Valid": true
        },
        "code": "theft",
        "min_circle_range": {
          "Int64": 50,
          "Valid": true
        },
        "max_circle_range": {
          "Int64": 500,
          "Valid": true
        },
        "default_circle_range": {
          "Int64": 100,
          "Valid": true
        },
        "category_code": "crime",
        "subcategory_code": "theft"
      }
    }
  ]
}
```

### Frontend Impact
```javascript
// ❌ FRONTEND CODE (React Native)
const iconUri = cluster.subcategory.icon_uri;
// Expected: "uri://icons/theft"
// Actual: {String: "uri://icons/theft", Valid: true}

// This will fail:
<Image source={{ uri: iconUri }} />
// Error: uri must be a string, received object

// This will also fail:
const radius = cluster.subcategory.default_circle_range;
// Expected: 100 (number)
// Actual: {Int64: 100, Valid: true} (object)

// Math operations fail:
const adjustedRadius = radius * 2;
// NaN (can't multiply object by number)
```

---

## Example 2: GetClusterBy Endpoint (Cluster Details)

### MySQL Response (WORKING)
```json
{
  "incl_id": 456,
  "created_at": "2026-01-17T10:30:00Z",
  "start_time": "2026-01-17T10:00:00Z",
  "end_time": "2026-01-18T10:00:00Z",
  "media_url": "uploads/image123.jpg",
  "center_latitude": 45.5017,
  "center_longitude": -73.5673,
  "is_active": true,
  "credibility": 8.5,
  "account_id": 789,
  "incidents": [
    {
      "inre_id": 111,
      "media_url": "uploads/incident1.jpg",
      "description": "Suspicious activity",
      "is_anonymous": "0",
      "status": "verified",
      "created_at": "2026-01-17T10:30:00"
    }
  ]
}
```

### PostgreSQL Response (BROKEN)
```json
{
  "incl_id": 456,
  "created_at": {
    "Time": "2026-01-17T10:30:00Z",
    "Valid": true
  },
  "start_time": {
    "Time": "2026-01-17T10:00:00Z",
    "Valid": true
  },
  "end_time": {
    "Time": "2026-01-18T10:00:00Z",
    "Valid": true
  },
  "media_url": "uploads/image123.jpg",
  "center_latitude": 45.5017,
  "center_longitude": -73.5673,
  "is_active": true,
  "credibility": 8.5,
  "account_id": {
    "Int64": 789,
    "Valid": true
  },
  "incidents": [
    {
      "inre_id": 111,
      "media_url": "uploads/incident1.jpg",
      "description": "Suspicious activity",
      "is_anonymous": "0",
      "status": {
        "String": "verified",
        "Valid": true
      },
      "created_at": "2026-01-17T10:30:00"
    }
  ]
}
```

### Frontend Impact
```javascript
// ❌ Date parsing fails
const createdAt = new Date(cluster.created_at);
// Expected: Valid Date object
// Actual: Invalid Date (trying to parse object)

// ❌ Account ID comparison fails
if (cluster.account_id === currentUserId) {
  // This never executes because:
  // cluster.account_id = {Int64: 789, Valid: true}
  // currentUserId = 789
}

// ❌ Status display fails
<Text>{incident.status}</Text>
// Expected: "verified"
// Actual: [object Object]
```

---

## Example 3: Auth Endpoint (Login/Signup)

### MySQL Response (WORKING)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "account_id": 123,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+1234567890",
    "is_premium": true,
    "has_finished_tutorial": false
  }
}
```

### PostgreSQL Response (BROKEN)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "account_id": 123,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": {
      "String": "+1234567890",
      "Valid": true
    },
    "is_premium": true,
    "has_finished_tutorial": false
  }
}
```

### Frontend Impact
```javascript
// ❌ Phone number display fails
<Text>Call: {user.phone_number}</Text>
// Expected: "Call: +1234567890"
// Actual: "Call: [object Object]"

// ❌ Phone number validation fails
const isPhoneValid = user.phone_number.startsWith('+');
// Error: phone_number.startsWith is not a function
```

---

## Example 4: Profile Endpoint

### Current Code (POTENTIAL ISSUE)
```go
// File: /internal/profile/model.go
type Incident struct {
    InreId          int64   `json:"inre_id"`
    MediaUrl        string  `json:"media_url"`
    Description     string  `json:"description"`
    EventType       string  `json:"event_type"`
    SubcategoryName string  `json:"subcategory_name"`
    Credibility     float32 `json:"credibility"`
    InclID          int64   `json:"incl_id"`
    IsAnonymous     string  `json:"is_anonymous"`  // ✅ Already string (CORRECT)
    CreatedAt       string  `json:"created_at"`    // ✅ Already string (CORRECT)
}
```

**Note:** This model uses `string` for dates, which is correct. However, verify the SQL query uses `CAST` correctly:

```sql
-- File: /internal/profile/repository.go line 75
'is_anonymous', COALESCE(CAST(i.is_anonymous AS TEXT), '0'),
```

This is **CORRECT** - it converts to TEXT in SQL, so the string is already formatted before Go receives it.

---

## Example 5: Categories Endpoint

### PostgreSQL Response (BROKEN)
```json
{
  "categories": [
    {
      "inca_id": 1,
      "name": "Crime",
      "code": "crime",
      "icon": {
        "String": "crime-icon.png",
        "Valid": true
      }
    }
  ]
}
```

### Frontend Impact
```javascript
// ❌ Icon rendering fails
categories.map(cat => (
  <Icon source={cat.icon} />
))
// Expected: cat.icon = "crime-icon.png"
// Actual: cat.icon = {String: "crime-icon.png", Valid: true}
```

---

## Example 6: NULL Values

### MySQL Behavior (Converts NULL to Zero Values)
```json
{
  "phone_number": "",           // NULL string → empty string
  "min_circle_range": 0,        // NULL int → 0
  "created_at": null            // NULL time → null (depends on driver)
}
```

### PostgreSQL Behavior (Uses sql.Null* Wrapper)
```json
{
  "phone_number": {
    "String": "",
    "Valid": false
  },
  "min_circle_range": {
    "Int64": 0,
    "Valid": false
  },
  "created_at": {
    "Time": "0001-01-01T00:00:00Z",
    "Valid": false
  }
}
```

### Frontend Impact
```javascript
// ❌ NULL checks fail
if (!user.phone_number) {
  // This doesn't work because phone_number is an object (truthy)
}

// ✅ Correct way (but frontend doesn't know to do this)
if (!user.phone_number || !user.phone_number.Valid) {
  // Frontend code doesn't check .Valid property
}
```

---

## Root Cause Analysis

### Why This Happens

1. **MySQL Driver Behavior:**
   - Returns primitive Go types (string, int64, time.Time)
   - NULL values converted to zero values (empty string, 0, zero time)
   - `sql.NullString` was used for NULL detection, but JSON encoder still saw the inner `.String` value

2. **PostgreSQL Driver Behavior:**
   - Returns `sql.Null*` wrapper types
   - JSON encoder serializes the **entire struct** (both `Value` and `Valid` fields)
   - Results in nested objects instead of primitive values

### Go's JSON Encoding Default Behavior

```go
type sql.NullString struct {
    String string
    Valid  bool
}

// Default JSON encoding (NO custom MarshalJSON):
json.Marshal(sql.NullString{String: "test", Valid: true})
// Output: {"String":"test","Valid":true}

// Custom JSON encoding (WITH MarshalJSON):
func (ns NullString) MarshalJSON() ([]byte, error) {
    if !ns.Valid {
        return []byte("null"), nil
    }
    return json.Marshal(ns.String)
}
// Output: "test"
```

---

## Testing Commands

### Test JSON Response from Actual Endpoint

```bash
# Start backend server
cd /Users/garyeikoow/Desktop/alertly/backend
go run cmd/app/main.go

# Test categories endpoint
curl -X GET "http://localhost:8080/api/categories" | jq '.'

# Look for nested objects like:
# "icon": {"String": "...", "Valid": true}
# Instead of:
# "icon": "..."
```

### Quick Fix Verification Script

```bash
# Create test script
cat > /tmp/test_json.sh << 'EOF'
#!/bin/bash
echo "Testing JSON responses..."

# Test categories
echo "Categories endpoint:"
curl -s "http://localhost:8080/api/categories" | jq '.data[0].icon' || echo "Failed"

# Test subcategories
echo "Subcategories endpoint:"
curl -s "http://localhost:8080/api/subcategories/1" | jq '.data[0].icon' || echo "Failed"

# Look for objects with "String" and "Valid" keys
echo "Searching for incorrect JSON format..."
curl -s "http://localhost:8080/api/categories" | grep -o '"String"' && echo "❌ FOUND INCORRECT FORMAT" || echo "✅ Format OK"
EOF

chmod +x /tmp/test_json.sh
/tmp/test_json.sh
```

---

## Summary Table

| Field Type | MySQL Response | PostgreSQL Response | Frontend Expects | Status |
|------------|---------------|---------------------|------------------|--------|
| `sql.NullString` | `"value"` | `{"String":"value","Valid":true}` | `string` | 🔴 BROKEN |
| `sql.NullInt64` | `123` | `{"Int64":123,"Valid":true}` | `number` | 🔴 BROKEN |
| `sql.NullTime` | `"2026-01-17T10:00:00Z"` | `{"Time":"2026-01-17T10:00:00Z","Valid":true}` | `string` | 🔴 BROKEN |
| `common.NullTime` | `"2026-01-17T10:00:00Z"` | `"2026-01-17T10:00:00Z"` | `string` | ✅ WORKING |
| `dbtypes.NullBool` | `true` | `true` | `boolean` | ✅ WORKING |
| `string` | `"value"` | `"value"` | `string` | ✅ WORKING |
| `int64` | `123` | `123` | `number` | ✅ WORKING |

---

## Conclusion

**ALL instances of `sql.NullString`, `sql.NullInt64`, `sql.NullTime`, and `sql.NullFloat64` in model structs MUST be replaced with custom wrapper types that implement `MarshalJSON` and `UnmarshalJSON` methods.**

Follow the existing patterns in:
- `/internal/common/types.go` (NullTime, NullBool)
- `/internal/dbtypes/nullbool.go` (NullBool)

**Next Step:** Implement custom wrappers for NullString, NullInt64, NullFloat64 and replace all usages in model files.
