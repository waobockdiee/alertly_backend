#!/bin/bash

# Quick Fix Commands for PostgreSQL JSON Type Issue
# Date: 2026-01-17
# Run this script to automatically apply all fixes

set -e  # Exit on error

BACKEND_DIR="/Users/garyeikoow/Desktop/alertly/backend"
cd "$BACKEND_DIR"

echo "=========================================="
echo "PostgreSQL JSON Type Fix - Automated"
echo "=========================================="
echo ""

# Step 1: Backup original files
echo "Step 1: Creating backups..."
mkdir -p /tmp/alertly_backup_$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/tmp/alertly_backup_$(date +%Y%m%d_%H%M%S)"

cp internal/getclustersbylocation/model.go "$BACKUP_DIR/"
cp internal/getclusterby/model.go "$BACKUP_DIR/"
cp internal/auth/model.go "$BACKUP_DIR/"
cp internal/getcategories/model.go "$BACKUP_DIR/"
cp internal/getsubcategoriesbycategoryid/model.go "$BACKUP_DIR/"

echo "✅ Backups created in: $BACKUP_DIR"
echo ""

# Step 2: Show files that need to be updated
echo "Step 2: Files using sql.Null* types:"
echo ""
grep -l "sql\.NullString\|sql\.NullInt64\|sql\.NullTime\|sql\.NullFloat" internal/*/model.go | head -10
echo ""

# Step 3: Create nulltypes.go (manual - code provided)
echo "Step 3: You need to manually create internal/common/nulltypes.go"
echo "Copy the code from FIX_JSON_TYPES_ACTION_PLAN.md Step 1"
echo ""
read -p "Press ENTER when you've created internal/common/nulltypes.go..."

# Step 4: Verify nulltypes.go exists
if [ ! -f "internal/common/nulltypes.go" ]; then
    echo "❌ ERROR: internal/common/nulltypes.go not found!"
    echo "Please create it first using the code from FIX_JSON_TYPES_ACTION_PLAN.md"
    exit 1
fi

echo "✅ Found internal/common/nulltypes.go"
echo ""

# Step 5: Show replacement commands
echo "Step 5: Replacement commands to run:"
echo ""
echo "# Replace sql.NullString with common.NullString"
echo "find internal/*/model.go -type f -exec sed -i '' 's/sql\.NullString/common.NullString/g' {} +"
echo ""
echo "# Replace sql.NullInt64 with common.NullInt64"
echo "find internal/*/model.go -type f -exec sed -i '' 's/sql\.NullInt64/common.NullInt64/g' {} +"
echo ""
echo "# Replace sql.NullFloat64 with common.NullFloat64"
echo "find internal/*/model.go -type f -exec sed -i '' 's/sql\.NullFloat64/common.NullFloat64/g' {} +"
echo ""
echo "# Replace sql.NullTime with common.NullTime (if not already using common.NullTime)"
echo "find internal/*/model.go -type f -exec sed -i '' 's/sql\.NullTime/common.NullTime/g' {} +"
echo ""

read -p "Run these replacements? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Running replacements..."
    find internal/*/model.go -type f -exec sed -i '' 's/sql\.NullString/common.NullString/g' {} +
    find internal/*/model.go -type f -exec sed -i '' 's/sql\.NullInt64/common.NullInt64/g' {} +
    find internal/*/model.go -type f -exec sed -i '' 's/sql\.NullFloat64/common.NullFloat64/g' {} +
    # Only replace sql.NullTime if it's NOT already common.NullTime
    find internal/*/model.go -type f -exec sed -i '' 's/\bsql\.NullTime\b/common.NullTime/g' {} +
    echo "✅ Replacements complete"
else
    echo "⚠️  Skipped automatic replacements. Run them manually."
fi

echo ""

# Step 6: Add common import if missing
echo "Step 6: Checking imports..."
for file in internal/getclustersbylocation/model.go \
            internal/getclusterby/model.go \
            internal/auth/model.go \
            internal/getcategories/model.go \
            internal/getsubcategoriesbycategoryid/model.go; do
    if ! grep -q '"alertly/internal/common"' "$file"; then
        echo "⚠️  Missing common import in $file - needs manual fix"
    else
        echo "✅ $file has common import"
    fi
done

echo ""

# Step 7: Verify compilation
echo "Step 7: Testing compilation..."
if go build ./internal/getclustersbylocation/...; then
    echo "✅ getclustersbylocation compiles"
else
    echo "❌ getclustersbylocation compilation failed"
fi

if go build ./internal/getclusterby/...; then
    echo "✅ getclusterby compiles"
else
    echo "❌ getclusterby compilation failed"
fi

if go build ./internal/auth/...; then
    echo "✅ auth compiles"
else
    echo "❌ auth compilation failed"
fi

echo ""

# Step 8: Run tests
echo "Step 8: Running tests..."
if go test ./internal/common/... -run TestNull; then
    echo "✅ Common tests pass"
else
    echo "⚠️  Common tests failed or don't exist yet"
fi

echo ""

# Step 9: Show next steps
echo "=========================================="
echo "Next Steps:"
echo "=========================================="
echo ""
echo "1. Start the server:"
echo "   go run cmd/app/main.go"
echo ""
echo "2. Test JSON responses:"
echo "   curl 'http://localhost:8080/api/categories' | jq '.data[0].icon'"
echo "   (Should be a string, not an object)"
echo ""
echo "3. Test frontend:"
echo "   cd ../frontend && npm start"
echo ""
echo "4. If everything works, commit changes:"
echo "   git add ."
echo "   git commit -m 'Fix PostgreSQL JSON serialization for sql.Null* types'"
echo ""
echo "5. If something breaks, restore backups:"
echo "   cp $BACKUP_DIR/*.go internal/*/"
echo ""
echo "=========================================="
echo "Fix script complete!"
echo "=========================================="
