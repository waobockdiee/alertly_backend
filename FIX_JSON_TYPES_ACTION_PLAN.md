# Fix JSON Type Serialization - Action Plan

**Date:** 2026-01-17
**Issue:** PostgreSQL migration broke JSON serialization
**Estimated Time:** 2-3 hours
**Priority:** 🔴 CRITICAL

---

## Step 1: Create Custom Wrapper Types (30 minutes)

### File: `/internal/common/nulltypes.go` (NEW FILE)

Create this file with the following content:

```go
package common

import (
	"database/sql"
	"encoding/json"
)

// NullString wraps sql.NullString with proper JSON serialization
// Serializes to: "value" or null (not {"String":"value","Valid":true})
type NullString struct {
	sql.NullString
}

// MarshalJSON implements json.Marshaler interface
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (ns *NullString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ns.Valid = false
		ns.String = ""
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
// Serializes to: 123 or null (not {"Int64":123,"Valid":true})
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON implements json.Marshaler interface
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ni.Valid = false
		ni.Int64 = 0
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
// Serializes to: 123.45 or null (not {"Float64":123.45,"Valid":true})
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON implements json.Marshaler interface
func (nf NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nf.Valid = false
		nf.Float64 = 0
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

### Verification

Test the new types:

```bash
cd /Users/garyeikoow/Desktop/alertly/backend
go test -v ./internal/common/... -run TestNull
```

---

## Step 2: Update Model Files (60 minutes)

### Priority 1: Critical Map Endpoint

#### File: `/internal/getclustersbylocation/model.go`

**BEFORE:**
```go
type Subcategory struct {
	InsuId             int64          `json:"insu_id"`
	IncaId             int64          `json:"inca_id"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	Icon               sql.NullString `json:"icon"`
	IconURI            sql.NullString `json:"icon_uri"`
	Code               string         `json:"code"`
	MinCircleRange     sql.NullInt64  `json:"min_circle_range"`
	MaxCircleRange     sql.NullInt64  `json:"max_circle_range"`
	DefaultCircleRange sql.NullInt64  `json:"default_circle_range"`
	CategoryCode       string         `json:"category_code"`
	SubcategoryCode    string         `json:"subcategory_code"`
}
```

**AFTER:**
```go
import (
	"alertly/internal/common"
)

type Subcategory struct {
	InsuId             int64              `json:"insu_id"`
	IncaId             int64              `json:"inca_id"`
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	Icon               common.NullString  `json:"icon"`        // ✅ FIXED
	IconURI            common.NullString  `json:"icon_uri"`    // ✅ FIXED
	Code               string             `json:"code"`
	MinCircleRange     common.NullInt64   `json:"min_circle_range"`      // ✅ FIXED
	MaxCircleRange     common.NullInt64   `json:"max_circle_range"`      // ✅ FIXED
	DefaultCircleRange common.NullInt64   `json:"default_circle_range"`  // ✅ FIXED
	CategoryCode       string             `json:"category_code"`
	SubcategoryCode    string             `json:"subcategory_code"`
}
```

**Changes:**
- Change import: Add `"alertly/internal/common"`
- Replace `sql.NullString` → `common.NullString`
- Replace `sql.NullInt64` → `common.NullInt64`

---

### Priority 2: Cluster Details Endpoint

#### File: `/internal/getclusterby/model.go`

**BEFORE:**
```go
type Cluster struct {
	InclId                 int64              `json:"incl_id"`
	CreatedAt              sql.NullTime       `json:"created_at"`
	StartTime              sql.NullTime       `json:"start_time"`
	EndTime                sql.NullTime       `json:"end_time"`
	// ... other fields
	AccountId              sql.NullInt64      `json:"account_id"`
}

type Incident struct {
	// ... other fields
	Status           sql.NullString    `json:"status"`
}
```

**AFTER:**
```go
import (
	"alertly/internal/common"
	"alertly/internal/comments"
)

type Cluster struct {
	InclId                 int64              `json:"incl_id"`
	CreatedAt              common.NullTime    `json:"created_at"`      // ✅ FIXED
	StartTime              common.NullTime    `json:"start_time"`      // ✅ FIXED
	EndTime                common.NullTime    `json:"end_time"`        // ✅ FIXED
	// ... other fields
	AccountId              common.NullInt64   `json:"account_id"`      // ✅ FIXED
}

type Incident struct {
	// ... other fields
	Status           common.NullString    `json:"status"`      // ✅ FIXED
}
```

**Changes:**
- Replace `sql.NullTime` → `common.NullTime` (already exists, verify import)
- Replace `sql.NullInt64` → `common.NullInt64`
- Replace `sql.NullString` → `common.NullString`

---

### Priority 3: Authentication Endpoint

#### File: `/internal/auth/model.go`

**BEFORE:**
```go
import (
	"database/sql"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	AccountID           int64          `json:"account_id"`
	Email               string         `json:"email"`
	FirstName           string         `json:"first_name"`
	LastName            string         `json:"last_name"`
	Password            string         `json:"password"`
	PhoneNumber         sql.NullString `json:"phone_number"`
	BirthYear           int            `json:"birth_year"`
	BirthMonth          int            `json:"birth_month"`
	BirthDay            int            `json:"birth_day"`
	Status              string         `json:"status"`
	IsPremium           bool           `json:"is_premium"`
	HasFinishedTutorial bool           `json:"has_finished_tutorial"`
}
```

**AFTER:**
```go
import (
	"alertly/internal/common"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	AccountID           int64             `json:"account_id"`
	Email               string            `json:"email"`
	FirstName           string            `json:"first_name"`
	LastName            string            `json:"last_name"`
	Password            string            `json:"password"`
	PhoneNumber         common.NullString `json:"phone_number"`  // ✅ FIXED
	BirthYear           int               `json:"birth_year"`
	BirthMonth          int               `json:"birth_month"`
	BirthDay            int               `json:"birth_day"`
	Status              string            `json:"status"`
	IsPremium           bool              `json:"is_premium"`
	HasFinishedTutorial bool              `json:"has_finished_tutorial"`
}
```

**Changes:**
- Remove `"database/sql"` import
- Add `"alertly/internal/common"` import
- Replace `sql.NullString` → `common.NullString`

---

### Priority 4: Categories Endpoint

#### File: `/internal/getcategories/model.go`

**BEFORE:**
```go
import "database/sql"

type Category struct {
	IncaId      int64           `json:"inca_id"`
	Name        string          `json:"name"`
	Code        string          `json:"code"`
	Icon        *sql.NullString `json:"icon"`
	Description string          `json:"description"`
}
```

**AFTER:**
```go
import "alertly/internal/common"

type Category struct {
	IncaId      int64              `json:"inca_id"`
	Name        string             `json:"name"`
	Code        string             `json:"code"`
	Icon        *common.NullString `json:"icon"`  // ✅ FIXED (pointer preserved)
	Description string             `json:"description"`
}
```

**Changes:**
- Replace import `"database/sql"` → `"alertly/internal/common"`
- Replace `*sql.NullString` → `*common.NullString` (keep the pointer)

---

### Priority 5: Subcategories Endpoint

#### File: `/internal/getsubcategoriesbycategoryid/model.go`

**BEFORE:**
```go
import "database/sql"

type Subcategory struct {
	InsuId             int64           `json:"insu_id"`
	IncaId             int64           `json:"inca_id"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Icon               *sql.NullString `json:"icon"`
	Code               string          `json:"code"`
	MinCircleRange     *sql.NullInt64  `json:"min_circle_range"`
	MaxCircleRange     *sql.NullInt64  `json:"max_circle_range"`
	DefaultCircleRange *sql.NullInt64  `json:"default_circle_range"`
	CategoryCode       string          `json:"category_code"`
	SubcategoryCode    string          `json:"subcategory_code"`
}
```

**AFTER:**
```go
import "alertly/internal/common"

type Subcategory struct {
	InsuId             int64              `json:"insu_id"`
	IncaId             int64              `json:"inca_id"`
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	Icon               *common.NullString `json:"icon"`               // ✅ FIXED
	Code               string             `json:"code"`
	MinCircleRange     *common.NullInt64  `json:"min_circle_range"`   // ✅ FIXED
	MaxCircleRange     *common.NullInt64  `json:"max_circle_range"`   // ✅ FIXED
	DefaultCircleRange *common.NullInt64  `json:"default_circle_range"` // ✅ FIXED
	CategoryCode       string             `json:"category_code"`
	SubcategoryCode    string             `json:"subcategory_code"`
}
```

**Changes:**
- Replace import `"database/sql"` → `"alertly/internal/common"`
- Replace `*sql.NullString` → `*common.NullString`
- Replace `*sql.NullInt64` → `*common.NullInt64`

---

## Step 3: Update Repository Files (30 minutes)

Repository files that **scan** into these types don't need changes because `common.NullString` embeds `sql.NullString`, so the `Scan()` method is inherited.

**Example - No changes needed:**
```go
// This still works because common.NullString embeds sql.NullString
var icon common.NullString
err := row.Scan(&icon)
// ✅ Works - Scan method is inherited from sql.NullString
```

**Verify all repository files still compile:**
```bash
cd /Users/garyeikoow/Desktop/alertly/backend
go build ./internal/getclustersbylocation/...
go build ./internal/getclusterby/...
go build ./internal/auth/...
go build ./internal/getcategories/...
go build ./internal/getsubcategoriesbycategoryid/...
```

---

## Step 4: Test JSON Responses (30 minutes)

### Start the Server

```bash
cd /Users/garyeikoow/Desktop/alertly/backend
go run cmd/app/main.go
```

### Test Each Endpoint

Create test script:

```bash
cat > /tmp/test_json_responses.sh << 'EOF'
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "Testing JSON serialization after fixes..."
echo ""

# Test 1: Categories
echo "1. Testing /api/categories"
RESPONSE=$(curl -s "$BASE_URL/api/categories" | jq '.data[0].icon')
if [[ $RESPONSE == *"String"* ]]; then
  echo "   ❌ FAILED - Still has nested object"
  echo "   Response: $RESPONSE"
else
  echo "   ✅ PASSED - Icon is a string: $RESPONSE"
fi
echo ""

# Test 2: Subcategories
echo "2. Testing /api/subcategories/1"
RESPONSE=$(curl -s "$BASE_URL/api/subcategories/1" | jq '.data[0].min_circle_range')
if [[ $RESPONSE == *"Int64"* ]]; then
  echo "   ❌ FAILED - Still has nested object"
  echo "   Response: $RESPONSE"
else
  echo "   ✅ PASSED - Range is a number: $RESPONSE"
fi
echo ""

# Test 3: Cluster details (need valid incl_id)
echo "3. Testing /api/cluster/1"
RESPONSE=$(curl -s "$BASE_URL/api/cluster/1" | jq '.account_id')
if [[ $RESPONSE == *"Int64"* ]]; then
  echo "   ❌ FAILED - Still has nested object"
  echo "   Response: $RESPONSE"
else
  echo "   ✅ PASSED - Account ID is correct: $RESPONSE"
fi
echo ""

# Test 4: Check for any remaining "String" or "Int64" keys
echo "4. Scanning all responses for nested objects..."
CLUSTERS=$(curl -s "$BASE_URL/api/clusters/..." 2>/dev/null || echo "{}")
if echo "$CLUSTERS" | grep -q '"String"'; then
  echo "   ⚠️  WARNING - Found 'String' key in response"
elif echo "$CLUSTERS" | grep -q '"Int64"'; then
  echo "   ⚠️  WARNING - Found 'Int64' key in response"
else
  echo "   ✅ PASSED - No nested objects found"
fi

echo ""
echo "Testing complete!"
EOF

chmod +x /tmp/test_json_responses.sh
/tmp/test_json_responses.sh
```

### Manual Testing

```bash
# Test categories - icon should be a string
curl "http://localhost:8080/api/categories" | jq '.data[0].icon'
# Expected: "icon.png" or null
# NOT: {"String":"icon.png","Valid":true}

# Test subcategories - ranges should be numbers
curl "http://localhost:8080/api/subcategories/1" | jq '.data[0].min_circle_range'
# Expected: 100 or null
# NOT: {"Int64":100,"Valid":true}

# Test cluster details - account_id should be number
curl "http://localhost:8080/api/cluster/123" | jq '.account_id'
# Expected: 789 or null
# NOT: {"Int64":789,"Valid":true}
```

---

## Step 5: Frontend Verification (30 minutes)

### Test on Real Device/Simulator

1. **Build and run frontend:**
   ```bash
   cd /Users/garyeikoow/Desktop/alertly/frontend
   npm start
   ```

2. **Test map view:**
   - Open app
   - Navigate to map
   - Verify clusters render correctly
   - Check that icons display (not "[object Object]")

3. **Test cluster details:**
   - Tap on a cluster
   - Verify all data displays correctly
   - Check timestamps are formatted properly
   - Verify account information shows correctly

4. **Test login:**
   - Login with an account
   - Verify phone number displays correctly
   - Check profile data

### Common Issues to Watch For

```javascript
// ❌ If you see this in React Native debugger:
console.log(cluster.subcategory.icon);
// Output: [object Object]
// FIX FAILED - icon is still nested object

// ✅ If you see this:
console.log(cluster.subcategory.icon);
// Output: "icon.png" or null
// FIX SUCCESSFUL
```

---

## Step 6: Write Tests (Optional but Recommended)

### Create test file: `/internal/common/nulltypes_test.go`

```go
package common

import (
	"encoding/json"
	"testing"
)

func TestNullString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		ns       NullString
		expected string
	}{
		{
			name:     "valid string",
			ns:       NullString{NullString: sql.NullString{String: "test", Valid: true}},
			expected: `"test"`,
		},
		{
			name:     "null string",
			ns:       NullString{NullString: sql.NullString{Valid: false}},
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.ns)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("MarshalJSON() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}

func TestNullInt64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		ni       NullInt64
		expected string
	}{
		{
			name:     "valid int64",
			ni:       NullInt64{NullInt64: sql.NullInt64{Int64: 123, Valid: true}},
			expected: `123`,
		},
		{
			name:     "null int64",
			ni:       NullInt64{NullInt64: sql.NullInt64{Valid: false}},
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.ni)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("MarshalJSON() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}
```

Run tests:
```bash
cd /Users/garyeikoow/Desktop/alertly/backend
go test -v ./internal/common/ -run TestNull
```

---

## Rollback Plan (If Something Breaks)

If the changes cause issues, you can quickly rollback:

```bash
cd /Users/garyeikoow/Desktop/alertly/backend
git checkout internal/getclustersbylocation/model.go
git checkout internal/getclusterby/model.go
git checkout internal/auth/model.go
git checkout internal/getcategories/model.go
git checkout internal/getsubcategoriesbycategoryid/model.go
rm internal/common/nulltypes.go
```

---

## Summary Checklist

- [ ] Step 1: Create `/internal/common/nulltypes.go` with custom wrappers
- [ ] Step 2.1: Fix `/internal/getclustersbylocation/model.go`
- [ ] Step 2.2: Fix `/internal/getclusterby/model.go`
- [ ] Step 2.3: Fix `/internal/auth/model.go`
- [ ] Step 2.4: Fix `/internal/getcategories/model.go`
- [ ] Step 2.5: Fix `/internal/getsubcategoriesbycategoryid/model.go`
- [ ] Step 3: Verify all repository files compile
- [ ] Step 4: Test JSON responses with curl/jq
- [ ] Step 5: Test frontend on device/simulator
- [ ] Step 6: Write unit tests for new types
- [ ] Step 7: Commit changes with descriptive message

---

## Estimated Time Breakdown

| Task | Time | Status |
|------|------|--------|
| Create nulltypes.go | 30 min | ⏳ Pending |
| Update model files | 60 min | ⏳ Pending |
| Verify compilation | 15 min | ⏳ Pending |
| Test JSON responses | 30 min | ⏳ Pending |
| Frontend testing | 30 min | ⏳ Pending |
| Write tests | 30 min | ⏳ Optional |
| **Total** | **2-3 hours** | |

---

## Questions to Ask Before Starting

1. **Do you have a staging environment?**
   - If yes, test there first before production

2. **Is the frontend expecting any NULL values to be empty strings instead of null?**
   - Check frontend code for null handling
   - May need to adjust MarshalJSON to return `""` instead of `null`

3. **Are there other model files not listed here that use sql.Null* types?**
   - Run: `grep -r "sql.Null" internal/*/model.go`
   - Apply same fixes to those files

4. **Do you want to handle NULL values as empty strings or as JSON null?**
   - Current implementation: NULL → `null` in JSON
   - Alternative: NULL → `""` for strings, `0` for numbers
   - Frontend may expect one or the other

---

## Next Steps After Completion

1. **Monitor production logs** for JSON parsing errors
2. **Check frontend crash reports** (Sentry, Firebase Crashlytics)
3. **Update API documentation** if NULL handling changed
4. **Add integration tests** for critical endpoints
5. **Document the fix** in CHANGELOG or migration notes

---

**Good luck! This fix should restore frontend functionality after the PostgreSQL migration.**
