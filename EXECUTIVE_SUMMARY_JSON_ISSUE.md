# Executive Summary - PostgreSQL Migration JSON Issue

**Date:** 2026-01-17
**Severity:** 🔴 CRITICAL
**Impact:** Frontend crashes and data display failures
**Root Cause:** Database type mismatch in JSON serialization
**ETA to Fix:** 2-3 hours

---

## The Problem in One Sentence

After migrating from MySQL to PostgreSQL, the backend is sending JSON objects like `{"String":"value","Valid":true}` instead of simple strings like `"value"`, causing the frontend to crash.

---

## What Changed During Migration

### Before (MySQL) ✅
```json
{
  "icon": "icon.png",
  "phone_number": "+1234567890",
  "min_circle_range": 100
}
```

### After (PostgreSQL) ❌
```json
{
  "icon": {"String": "icon.png", "Valid": true},
  "phone_number": {"String": "+1234567890", "Valid": true},
  "min_circle_range": {"Int64": 100, "Valid": true}
}
```

---

## Why Frontend Crashes

```javascript
// Frontend expects:
<Image source={{ uri: cluster.icon }} />

// But receives:
cluster.icon = {String: "icon.png", Valid: true}

// Result: ❌ Error: uri must be a string, not object
```

---

## Root Cause

Go's `sql.NullString`, `sql.NullInt64`, and `sql.NullTime` types **do not have custom JSON serialization**, so they serialize as complete structs with `{Value, Valid}` fields.

MySQL driver converted NULLs to zero values automatically, hiding this issue.
PostgreSQL driver is stricter and preserves the wrapper types.

---

## The Fix

Create custom wrapper types with `MarshalJSON` methods (like you already have for `common.NullTime` and `dbtypes.NullBool`).

### Step 1: Create `/internal/common/nulltypes.go`
Define `NullString`, `NullInt64`, `NullFloat64` with custom JSON serialization.

### Step 2: Replace types in 5 critical model files
- `/internal/getclustersbylocation/model.go` (map view)
- `/internal/getclusterby/model.go` (cluster details)
- `/internal/auth/model.go` (login/signup)
- `/internal/getcategories/model.go` (categories)
- `/internal/getsubcategoriesbycategoryid/model.go` (subcategories)

### Step 3: Test
Verify JSON responses return primitive values, not objects.

---

## Affected Endpoints

| Endpoint | Impact | Priority |
|----------|--------|----------|
| `/api/getclustersbylocation` | Map won't load | 🔴 Critical |
| `/api/getclusterby/:id` | Cluster details broken | 🔴 Critical |
| `/api/login` | Phone number display fails | 🟡 High |
| `/api/categories` | Icons don't render | 🟡 High |
| `/api/subcategories/:id` | Radius calculations fail | 🟡 High |

---

## Timeline

| Task | Duration | Blocker? |
|------|----------|----------|
| Create custom wrapper types | 30 min | No |
| Update model files | 60 min | No |
| Test JSON responses | 30 min | Yes |
| Frontend verification | 30 min | Yes |
| **Total** | **2-3 hours** | |

---

## Risks if Not Fixed

1. **Frontend crashes** on map view (most used feature)
2. **Data corruption** from incorrect parsing
3. **User complaints** about broken UI
4. **App store reviews** mentioning bugs
5. **Support tickets** from confused users

---

## Documentation Created

Three detailed documents have been created:

1. **POSTGRESQL_JSON_TYPE_ISSUES.md** - Technical deep dive
2. **JSON_RESPONSE_COMPARISON.md** - Before/after examples
3. **FIX_JSON_TYPES_ACTION_PLAN.md** - Step-by-step implementation guide

---

## Immediate Action Required

**START HERE:** Read `FIX_JSON_TYPES_ACTION_PLAN.md` and follow the step-by-step guide.

Estimated time to production-ready: 2-3 hours.

---

**This is a migration-breaking bug that must be fixed before production deployment.**
