# PostgreSQL JSON Type Fix - Documentation Index

**Date:** 2026-01-17
**Issue:** Backend sending incorrect JSON format after PostgreSQL migration
**Status:** 🔴 CRITICAL - Requires immediate fix

---

## Start Here

If you need to fix this issue RIGHT NOW, read these files in order:

1. **EXECUTIVE_SUMMARY_JSON_ISSUE.md** (3.5 KB)
   - Quick overview of the problem
   - What broke and why
   - Estimated time to fix: 2-3 hours

2. **FIX_JSON_TYPES_ACTION_PLAN.md** (19 KB)
   - Complete step-by-step implementation guide
   - Copy/paste ready code
   - Testing instructions

3. **QUICK_FIX_COMMANDS.sh** (5.1 KB)
   - Automated script to apply fixes
   - Run with: `./QUICK_FIX_COMMANDS.sh`

---

## Detailed Documentation

For deeper understanding:

4. **POSTGRESQL_JSON_TYPE_ISSUES.md** (12 KB)
   - Technical deep dive
   - Complete list of affected files
   - Recommended solutions with pros/cons

5. **JSON_RESPONSE_COMPARISON.md** (11 KB)
   - Before/after JSON examples
   - Frontend impact analysis
   - Testing commands

---

## The Problem (TL;DR)

### What Broke
After migrating from MySQL to PostgreSQL, the backend started sending:

```json
{
  "icon": {"String": "icon.png", "Valid": true}
}
```

Instead of:

```json
{
  "icon": "icon.png"
}
```

### Why It Broke
- Go's `sql.NullString`, `sql.NullInt64`, `sql.NullTime` don't have custom JSON serialization
- They serialize as complete structs: `{Value, Valid}`
- Frontend expects primitive types (string, number, boolean)
- Result: Frontend crashes with "Expected string, got object"

### The Fix
Create custom wrapper types with `MarshalJSON` methods (you already have this pattern for `common.NullTime` and `dbtypes.NullBool`).

---

## Quick Reference

### Files to Create
- `/internal/common/nulltypes.go` (NEW - 200 lines)

### Files to Modify
- `/internal/getclustersbylocation/model.go` (Replace sql.Null* → common.Null*)
- `/internal/getclusterby/model.go` (Replace sql.Null* → common.Null*)
- `/internal/auth/model.go` (Replace sql.Null* → common.Null*)
- `/internal/getcategories/model.go` (Replace sql.Null* → common.Null*)
- `/internal/getsubcategoriesbycategoryid/model.go` (Replace sql.Null* → common.Null*)

### Affected Endpoints
- `/api/getclustersbylocation` - Map view (CRITICAL)
- `/api/getclusterby/:id` - Cluster details (CRITICAL)
- `/api/login` - Authentication (HIGH)
- `/api/categories` - Category list (HIGH)
- `/api/subcategories/:id` - Subcategories (HIGH)

### Testing Commands
```bash
# Start server
go run cmd/app/main.go

# Test categories endpoint
curl "http://localhost:8080/api/categories" | jq '.data[0].icon'
# Expected: "icon.png" (string)
# NOT: {"String":"icon.png","Valid":true} (object)

# Check for nested objects
curl "http://localhost:8080/api/categories" | grep '"String"'
# Should return nothing if fixed correctly
```

---

## Implementation Checklist

Use this checklist as you implement the fix:

- [ ] Read EXECUTIVE_SUMMARY_JSON_ISSUE.md
- [ ] Read FIX_JSON_TYPES_ACTION_PLAN.md
- [ ] Create `/internal/common/nulltypes.go`
- [ ] Update model files (5 files)
- [ ] Add `"alertly/internal/common"` imports
- [ ] Test compilation: `go build ./internal/...`
- [ ] Run server: `go run cmd/app/main.go`
- [ ] Test JSON responses with curl
- [ ] Test frontend on device/simulator
- [ ] Write unit tests (optional)
- [ ] Commit changes
- [ ] Deploy to staging
- [ ] Monitor production logs

---

## Time Estimate

| Phase | Duration | Difficulty |
|-------|----------|-----------|
| Understanding the problem | 15 min | Easy |
| Creating nulltypes.go | 30 min | Medium |
| Updating model files | 60 min | Easy |
| Testing backend | 30 min | Medium |
| Testing frontend | 30 min | Medium |
| **Total** | **2-3 hours** | **Medium** |

---

## Risk Assessment

### High Risk if Not Fixed
1. Frontend crashes on map view (primary feature)
2. Data display failures across the app
3. User frustration and negative reviews
4. Potential data corruption from incorrect parsing

### Low Risk if Fixed Correctly
- Changes are isolated to model definitions
- No database schema changes required
- Easy to rollback if needed
- Existing code patterns (common.NullTime) prove the approach works

---

## Rollback Plan

If something goes wrong:

```bash
# Restore backups
cp /tmp/alertly_backup_*/model.go internal/*/

# Remove new file
rm internal/common/nulltypes.go

# Restart server
go run cmd/app/main.go
```

---

## Support

If you encounter issues during implementation:

1. **Check compilation errors:**
   ```bash
   go build ./internal/...
   ```

2. **Verify JSON output:**
   ```bash
   curl "http://localhost:8080/api/categories" | jq '.'
   ```

3. **Check frontend logs:**
   - React Native debugger
   - Console errors mentioning "Expected string, got object"

4. **Review existing working examples:**
   - `/internal/common/types.go` (NullTime implementation)
   - `/internal/dbtypes/nullbool.go` (NullBool implementation)

---

## Additional Context

### Why This Worked in MySQL

MySQL's Go driver (`go-sql-driver/mysql`) automatically converts NULL values to zero values:
- NULL string → `""`
- NULL int64 → `0`
- NULL time → zero time

PostgreSQL's driver (`lib/pq`) preserves NULL state in wrapper types, which exposes the JSON serialization issue.

### Why We Use sql.Null* Types

To distinguish between:
- NULL (value not set)
- Zero value (value explicitly set to "", 0, etc.)

This is important for optional fields in the API.

### Alternative Approaches (Not Recommended)

1. **Use primitive types + COALESCE in SQL:**
   - Pros: Simple, no custom code
   - Cons: Can't distinguish NULL from zero value

2. **Use pointers (*string, *int64):**
   - Pros: Built-in JSON support
   - Cons: Cumbersome to work with, easy to cause nil pointer panics

3. **Custom JSON marshaling in handlers:**
   - Pros: Keeps models simple
   - Cons: Scattered logic, hard to maintain

**Best approach:** Custom wrapper types with embedded sql.Null* and MarshalJSON (recommended in this fix).

---

## Documentation Files Summary

| File | Size | Purpose |
|------|------|---------|
| EXECUTIVE_SUMMARY_JSON_ISSUE.md | 3.5 KB | High-level overview |
| FIX_JSON_TYPES_ACTION_PLAN.md | 19 KB | Step-by-step implementation |
| POSTGRESQL_JSON_TYPE_ISSUES.md | 12 KB | Technical deep dive |
| JSON_RESPONSE_COMPARISON.md | 11 KB | Before/after examples |
| QUICK_FIX_COMMANDS.sh | 5.1 KB | Automated fix script |
| README_JSON_FIX.md (this file) | 6 KB | Documentation index |

---

## Next Steps After Fix

1. Update API documentation
2. Add integration tests for JSON serialization
3. Monitor production logs for JSON parsing errors
4. Consider adding linting rules to prevent sql.Null* in models
5. Update developer onboarding docs with this pattern

---

## Questions?

If you're unsure about any step:

1. Read the relevant documentation file above
2. Check existing code for similar patterns (common.NullTime, dbtypes.NullBool)
3. Test in a staging environment first
4. Start with one model file (getclustersbylocation) and verify it works before updating others

---

**Remember:** This fix follows patterns already proven to work in your codebase (common.NullTime, dbtypes.NullBool). You're just extending that pattern to other types.

**Good luck!**
