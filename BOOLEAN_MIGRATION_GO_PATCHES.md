# Parches de Código Go para Migración de Columnas BOOLEAN

**Fecha:** 2026-01-18
**Prerequisito:** Ejecutar `001_standardize_boolean_columns.sql` ANTES de aplicar estos cambios

---

## ARCHIVO 1: internal/cronjob/premium_expiration.go

**Cambios:** 6 líneas (líneas 65, 108, 143, 153, 167, 188)

### Línea 65
```diff
- WHERE is_premium = 1
+ WHERE is_premium = true
```

### Línea 108
```diff
- SET is_premium = 0, premium_expired_date = NULL
+ SET is_premium = false, premium_expired_date = NULL
```

### Línea 143
```diff
- err := s.db.QueryRow("SELECT COUNT(*) FROM account WHERE is_premium = 1").Scan(&activePremium)
+ err := s.db.QueryRow("SELECT COUNT(*) FROM account WHERE is_premium = true").Scan(&activePremium)
```

### Línea 153
```diff
- WHERE is_premium = 1
+ WHERE is_premium = true
```

### Línea 167
```diff
- WHERE is_premium = 1
+ WHERE is_premium = true
```

### Línea 188
```diff
- WHERE is_premium = 1
+ WHERE is_premium = true
```

---

## ARCHIVO 2: internal/cronjobs/cjuserank/repository.go

**Cambios:** 1 línea (línea 28)

### Línea 28
```diff
- status = 'active' AND receive_notifications = 1
+ status = 'active' AND receive_notifications = true
```

---

## ARCHIVO 3: internal/cronjobs/cjbadgeearn/repository.go

**Cambios:** 1 línea (línea 41)

### Línea 41
```diff
- status = 'active' AND receive_notifications = 1
+ status = 'active' AND receive_notifications = true
```

---

## ARCHIVO 4: internal/cronjobs/cjcomments/repository.go

**Cambios:** 1 línea (línea 140)

### Línea 140
```diff
- WHERE dt.account_id IN (%s) AND a.status = 'active' AND a.is_premium = 1 AND a.receive_notifications = 1
+ WHERE dt.account_id IN (%s) AND a.status = 'active' AND a.is_premium = true AND a.receive_notifications = true
```

---

## ARCHIVO 5: internal/cronjobs/cjincidentupdate/repository.go

**Cambios:** 2 líneas (líneas 140-141)

### Líneas 140-141
```diff
- AND a.is_premium = 1
- AND a.receive_notifications = 1
+ AND a.is_premium = true
+ AND a.receive_notifications = true
```

---

## ARCHIVO 6: internal/editprofile/repository.go

**Cambios:** Eliminar `dbtypes.BoolToInt()` en múltiples líneas

### Línea 133 (función UpdateEmail)
```diff
- _, err := r.db.Exec(query, email, dbtypes.BoolToInt(false), accountID)
+ _, err := r.db.Exec(query, email, false, accountID)
```

### Línea 181 (función UpdateNickname)
```diff
- _, err := r.db.Exec(query, nickname, dbtypes.BoolToInt(false), accountID, dbtypes.BoolToInt(true))
+ _, err := r.db.Exec(query, nickname, false, accountID, true)
```

### Línea 206 (función UpdateFullName)
```diff
- _, err := r.db.Exec(query, firstName, lastName, dbtypes.BoolToInt(false), accountID, dbtypes.BoolToInt(true))
+ _, err := r.db.Exec(query, firstName, lastName, false, accountID, true)
```

### Línea 217 (función UpdateIsPrivateProfile)
```diff
- _, err := r.db.Exec(query, dbtypes.BoolToInt(isPrivateProfile), accountID)
+ _, err := r.db.Exec(query, isPrivateProfile, accountID)
```

### Línea 230 (función UpdateBirthDate)
```diff
- _, err := r.db.Exec(query, year, month, day, dbtypes.BoolToInt(false), accountID, dbtypes.BoolToInt(true))
+ _, err := r.db.Exec(query, year, month, day, false, accountID, true)
```

---

## ARCHIVO 7: internal/tutorial/repository.go

**Cambios:** 1 línea (línea 24)

### Línea 24
```diff
- _, err := r.db.Exec(query, dbtypes.BoolToInt(true), accountID)
+ _, err := r.db.Exec(query, true, accountID)
```

---

## ARCHIVO 8: internal/account/repository.go

**Cambios:** 3 líneas en la función `UpdatePremiumStatus` (~líneas 306, 310, 314)

### Función UpdatePremiumStatus (líneas 306-314)
```diff
if isPremium && expirationDate != nil {
    updateAccountQuery = "UPDATE account SET is_premium = $1, premium_expired_date = $2 WHERE account_id = $3"
-   args = []interface{}{dbtypes.BoolToInt(isPremium), expirationDate, accountID}
+   args = []interface{}{isPremium, expirationDate, accountID}
} else if !isPremium {
    // When cancelling or expiring, set is_premium to false and clear expiration date
    updateAccountQuery = "UPDATE account SET is_premium = $1, premium_expired_date = NULL WHERE account_id = $2"
-   args = []interface{}{dbtypes.BoolToInt(isPremium), accountID}
+   args = []interface{}{isPremium, accountID}
} else {
    // Fallback for safety, though should not be reached in normal flow
    updateAccountQuery = "UPDATE account SET is_premium = $1 WHERE account_id = $2"
-   args = []interface{}{dbtypes.BoolToInt(isPremium), accountID}
+   args = []interface{}{isPremium, accountID}
}
```

### Función SetHasFinishedTutorial (línea 284)
```diff
- _, err := r.db.Exec(query, dbtypes.BoolToInt(true), accountID)
+ _, err := r.db.Exec(query, true, accountID)
```

---

## RESUMEN DE CAMBIOS

| Archivo | Líneas Modificadas | Tipo de Cambio |
|---------|-------------------|----------------|
| `internal/cronjob/premium_expiration.go` | 6 | `= 1` → `= true` |
| `internal/cronjobs/cjuserank/repository.go` | 1 | `= 1` → `= true` |
| `internal/cronjobs/cjbadgeearn/repository.go` | 1 | `= 1` → `= true` |
| `internal/cronjobs/cjcomments/repository.go` | 1 | `= 1` → `= true` |
| `internal/cronjobs/cjincidentupdate/repository.go` | 2 | `= 1` → `= true` |
| `internal/editprofile/repository.go` | 5 | Eliminar `BoolToInt()` |
| `internal/tutorial/repository.go` | 1 | Eliminar `BoolToInt()` |
| `internal/account/repository.go` | 4 | Eliminar `BoolToInt()` |
| **TOTAL** | **21 líneas** | **8 archivos** |

---

## SCRIPT DE APLICACIÓN AUTOMÁTICA

Puedes usar este script bash para aplicar los cambios automáticamente:

```bash
#!/bin/bash
# apply_boolean_patches.sh

cd /Users/garyeikoow/Desktop/alertly/backend

echo "Aplicando parches de migración BOOLEAN..."

# Archivo 1: premium_expiration.go
sed -i '' 's/WHERE is_premium = 1/WHERE is_premium = true/g' internal/cronjob/premium_expiration.go
sed -i '' 's/SET is_premium = 0,/SET is_premium = false,/g' internal/cronjob/premium_expiration.go
sed -i '' 's/WHERE is_premium = 1/WHERE is_premium = true/g' internal/cronjob/premium_expiration.go

# Archivo 2-5: Queries con receive_notifications y is_premium
sed -i '' 's/receive_notifications = 1/receive_notifications = true/g' internal/cronjobs/cjuserank/repository.go
sed -i '' 's/receive_notifications = 1/receive_notifications = true/g' internal/cronjobs/cjbadgeearn/repository.go
sed -i '' 's/a.is_premium = 1 AND a.receive_notifications = 1/a.is_premium = true AND a.receive_notifications = true/g' internal/cronjobs/cjcomments/repository.go
sed -i '' 's/a.is_premium = 1/a.is_premium = true/g' internal/cronjobs/cjincidentupdate/repository.go
sed -i '' 's/a.receive_notifications = 1/a.receive_notifications = true/g' internal/cronjobs/cjincidentupdate/repository.go

# Archivo 6-8: Eliminar dbtypes.BoolToInt()
sed -i '' 's/dbtypes\.BoolToInt(false)/false/g' internal/editprofile/repository.go
sed -i '' 's/dbtypes\.BoolToInt(true)/true/g' internal/editprofile/repository.go
sed -i '' 's/dbtypes\.BoolToInt(isPrivateProfile)/isPrivateProfile/g' internal/editprofile/repository.go
sed -i '' 's/dbtypes\.BoolToInt(true)/true/g' internal/tutorial/repository.go
sed -i '' 's/dbtypes\.BoolToInt(isPremium)/isPremium/g' internal/account/repository.go
sed -i '' 's/dbtypes\.BoolToInt(true)/true/g' internal/account/repository.go

echo "✅ Parches aplicados exitosamente"
echo "⚠️  Recuerde ejecutar 'go build' para verificar que no hay errores de compilación"
```

**IMPORTANTE:** Revisa manualmente los cambios antes de commitear.

---

## VERIFICACIÓN POST-MIGRACIÓN

Después de aplicar los parches, ejecuta:

```bash
# 1. Compilación
cd /Users/garyeikoow/Desktop/alertly/backend
go build ./cmd/app
go build ./cmd/cronjob

# 2. Tests unitarios
go test ./internal/cronjobs/cjincidentexpiration/...
go test ./internal/analytics/...
go test ./internal/common/...

# 3. Verificación manual de queries
grep -rn "= 1" internal/ | grep -E "(is_premium|receive_notifications)"
grep -rn "= 0" internal/ | grep -E "(is_premium|receive_notifications)"
# No debe devolver resultados en archivos modificados

# 4. Test de integración (requiere base de datos activa)
go run cmd/app/main.go &
# Probar login, editar perfil, cronjobs, etc.
```

---

## NOTAS IMPORTANTES

1. **NO ejecutar en producción sin testing previo**
2. **Aplicar migración SQL ANTES que cambios en Go**
3. **Crear backup de base de datos antes de migrar**
4. **Estos cambios son OBLIGATORIOS después de ejecutar `001_standardize_boolean_columns.sql`**
5. **La función `dbtypes.BoolToInt()` puede ser deprecada** (pero mantenerla por compatibilidad con código legacy)

---

## ROLLBACK

Si necesitas revertir los cambios:

1. Ejecutar `001_rollback_boolean_columns.sql` en la base de datos
2. Usar `git checkout` para revertir los archivos Go modificados
3. Verificar que la aplicación funciona correctamente

```bash
git checkout internal/cronjob/premium_expiration.go
git checkout internal/cronjobs/cjuserank/repository.go
git checkout internal/cronjobs/cjbadgeearn/repository.go
git checkout internal/cronjobs/cjcomments/repository.go
git checkout internal/cronjobs/cjincidentupdate/repository.go
git checkout internal/editprofile/repository.go
git checkout internal/tutorial/repository.go
git checkout internal/account/repository.go
```

---

**Fin del documento**
