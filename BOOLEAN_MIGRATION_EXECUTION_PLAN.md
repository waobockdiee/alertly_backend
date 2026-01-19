# Plan de Ejecución: Migración de Columnas BOOLEAN

**Fecha de creación:** 2026-01-18
**Base de datos:** Railway PostgreSQL (`metro.proxy.rlwy.net:48204/railway`)
**Duración estimada:** 1-2 horas (incluye testing)
**Riesgo:** Bajo (cambios quirúrgicos y bien localizados)
**Rollback disponible:** Sí (script incluido)

---

## RESUMEN EJECUTIVO

La migración MySQL → PostgreSQL ha creado una **inconsistencia en tipos booleanos**:

- **Problema:** Columnas mixtas (BOOLEAN, SMALLINT, CHAR) + queries SQL hardcodeadas con `= 1` y `= true`
- **Solución:** Estandarizar TODO a tipo BOOLEAN nativo de PostgreSQL
- **Impacto:** 8 columnas en 2 tablas, 21 líneas de código en 8 archivos Go
- **Beneficio:** Elimina errores de tipo, mejora performance, código más limpio

**Archivos de referencia:**
- `BOOLEAN_MIGRATION_ANALYSIS.md` - Análisis exhaustivo completo
- `BOOLEAN_MIGRATION_GO_PATCHES.md` - Cambios línea por línea
- `assets/db/migrations/001_standardize_boolean_columns.sql` - Script de migración
- `assets/db/migrations/001_rollback_boolean_columns.sql` - Script de rollback

---

## PREREQUISITOS

### 1. Verificar estado actual

```bash
# Conectar a PostgreSQL Railway
PGPASSWORD="cGA2dBF6G33BgfefcgDb1CDa6CagFcC5" psql \
  -h metro.proxy.rlwy.net \
  -p 48204 \
  -U postgres \
  -d railway

# Verificar tipos actuales
\d account
\d account_favorite_locations

# Verificar valores actuales
SELECT is_premium, receive_notifications, is_private_profile FROM account LIMIT 5;
SELECT crime, traffic_accident, status FROM account_favorite_locations LIMIT 3;

# Contar registros (para estimar duración)
SELECT COUNT(*) FROM account;
SELECT COUNT(*) FROM account_favorite_locations;
```

**Resultado esperado:**
- `account.is_premium` = SMALLINT (valores: 0, 1)
- `account.receive_notifications` = BOOLEAN (valores: t, f)
- `account_favorite_locations.crime` = BOOLEAN (valores: t, f)

### 2. Crear backup

```bash
# Opción 1: Backup completo de Railway (recomendado)
# Ir a Railway Dashboard → Database → Create Snapshot

# Opción 2: Dump SQL manual
PGPASSWORD="cGA2dBF6G33BgfefcgDb1CDa6CagFcC5" pg_dump \
  -h metro.proxy.rlwy.net \
  -p 48204 \
  -U postgres \
  -d railway \
  --no-owner \
  --no-acl \
  > backup_before_boolean_migration_$(date +%Y%m%d_%H%M%S).sql
```

### 3. Detener servicios (opcional pero recomendado)

```bash
# Si la app está en Railway/EC2, detener temporalmente:
# - Servidor HTTP (cmd/app)
# - Cronjobs (cmd/cronjob)

# Verificar que no hay conexiones activas
SELECT COUNT(*) FROM pg_stat_activity WHERE datname = 'railway' AND state = 'active';
```

---

## PASO 1: EJECUTAR MIGRACIÓN SQL

**Duración:** 1-2 minutos
**Impacto:** Bloqueo temporal de tablas `account` y `account_favorite_locations`

```bash
# Conectar a PostgreSQL
PGPASSWORD="cGA2dBF6G33BgfefcgDb1CDa6CagFcC5" psql \
  -h metro.proxy.rlwy.net \
  -p 48204 \
  -U postgres \
  -d railway

# Ejecutar migración
\i /Users/garyeikoow/Desktop/alertly/backend/assets/db/migrations/001_standardize_boolean_columns.sql

# Verificar resultado
# Debe mostrar:
# ✅ account: 7 columns successfully migrated to BOOLEAN
# ✅ account_favorite_locations: status column successfully migrated to BOOLEAN
# ✅ No NULL values found in boolean columns
# ✅ MIGRATION COMPLETED SUCCESSFULLY
```

**Validación inmediata:**

```sql
-- Verificar tipos de columnas
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns
WHERE table_name = 'account'
  AND column_name IN (
      'is_premium', 'receive_notifications', 'is_private_profile',
      'has_finished_tutorial', 'has_watch_new_incident_tutorial',
      'can_update_email', 'can_update_nickname', 'can_update_fullname', 'can_update_birthdate'
  )
ORDER BY column_name;

-- Resultado esperado: data_type = 'boolean' para todas las columnas

-- Verificar valores actuales
SELECT is_premium, receive_notifications, is_private_profile, has_finished_tutorial
FROM account
LIMIT 5;

-- Resultado esperado: valores t/f (no 0/1)

-- Test de queries críticas
SELECT COUNT(*) FROM account WHERE is_premium = true;
SELECT COUNT(*) FROM account WHERE receive_notifications = true;
SELECT COUNT(*) FROM account_favorite_locations WHERE crime = true;

-- Todos deben ejecutar sin errores
```

**Si hay errores:**
```bash
# Rollback inmediato
\i /Users/garyeikoow/Desktop/alertly/backend/assets/db/migrations/001_rollback_boolean_columns.sql

# Investigar causa, ajustar script, reintentar
```

---

## PASO 2: APLICAR CAMBIOS EN CÓDIGO GO

**Duración:** 5-10 minutos

### Opción A: Manual (recomendado para primera vez)

Usar el archivo `BOOLEAN_MIGRATION_GO_PATCHES.md` como guía:

```bash
cd /Users/garyeikoow/Desktop/alertly/backend

# Editar archivos según documento de parches
# Verificar cada cambio manualmente
```

### Opción B: Script automatizado

```bash
cd /Users/garyeikoow/Desktop/alertly/backend

# Crear script de aplicación
cat > apply_boolean_patches.sh << 'EOF'
#!/bin/bash
set -e

echo "🔧 Aplicando parches de migración BOOLEAN..."

# Archivo 1: premium_expiration.go
sed -i '' 's/WHERE is_premium = 1/WHERE is_premium = true/g' internal/cronjob/premium_expiration.go
sed -i '' 's/SET is_premium = 0,/SET is_premium = false,/g' internal/cronjob/premium_expiration.go

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
EOF

chmod +x apply_boolean_patches.sh
./apply_boolean_patches.sh
```

### Verificación de cambios

```bash
# Verificar que no quedan referencias a = 1 o = 0 en columnas booleanas
grep -rn "is_premium = 1" internal/
grep -rn "receive_notifications = 1" internal/
# No debe devolver resultados en archivos modificados

# Verificar que no quedan BoolToInt() innecesarios
grep -rn "dbtypes.BoolToInt" internal/editprofile/repository.go
grep -rn "dbtypes.BoolToInt" internal/tutorial/repository.go
grep -rn "dbtypes.BoolToInt" internal/account/repository.go
# Solo debe devolver imports, no usos en queries de columnas BOOLEAN

# Compilación
go build ./cmd/app
go build ./cmd/cronjob

# Si hay errores, revisar manualmente los archivos modificados
```

---

## PASO 3: TESTING EXHAUSTIVO

### 3.1 Tests Unitarios

```bash
cd /Users/garyeikoow/Desktop/alertly/backend

# Tests existentes
go test ./internal/cronjobs/cjincidentexpiration/...
go test ./internal/analytics/...
go test ./internal/common/...
go test ./internal/dbtypes/...

# Todos deben pasar sin errores
```

### 3.2 Tests de Integración (Requiere DB activa)

```bash
# Iniciar servidor HTTP
go run cmd/app/main.go &
SERVER_PID=$!

# Esperar a que inicie
sleep 5

# Test 1: Login (valida is_premium, has_finished_tutorial)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}' | jq

# Test 2: Profile (valida 4 columnas booleanas)
TOKEN="<token-del-login>"
curl http://localhost:8080/profile/1 \
  -H "Authorization: Bearer $TOKEN" | jq

# Test 3: Editar perfil - Toggle privacidad
curl -X PUT http://localhost:8080/profile/update-privacy \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"is_private":true}' | jq

# Test 4: Editar perfil - Toggle notificaciones
curl -X PUT http://localhost:8080/profile/update-notifications \
  -H "Authorization: Bearer $TOKEN" | jq

# Test 5: MyPlaces (valida 12 columnas BOOLEAN)
curl http://localhost:8080/myplaces/1/favorite-locations \
  -H "Authorization: Bearer $TOKEN" | jq

# Test 6: Tutorial
curl -X POST http://localhost:8080/tutorial/mark-finished \
  -H "Authorization: Bearer $TOKEN" | jq

# Detener servidor
kill $SERVER_PID
```

### 3.3 Tests de Cronjobs

```bash
# Test cronjob de notificaciones
go run cmd/cronjob/main.go

# Verificar logs:
# - ✅ Queries ejecutadas sin errores
# - ✅ Usuarios premium detectados correctamente
# - ✅ Notificaciones enviadas según configuración receive_notifications

# Verificar en base de datos
PGPASSWORD="cGA2dBF6G33BgfefcgDb1CDa6CagFcC5" psql \
  -h metro.proxy.rlwy.net \
  -p 48204 \
  -U postgres \
  -d railway \
  -c "SELECT COUNT(*) FROM notifications WHERE created_at > NOW() - INTERVAL '5 minutes';"

# Debe mostrar notificaciones creadas por el cronjob
```

### 3.4 Tests de Queries Directas

```sql
-- Conectar a PostgreSQL
PGPASSWORD="cGA2dBF6G33BgfefcgDb1CDa6CagFcC5" psql \
  -h metro.proxy.rlwy.net \
  -p 48204 \
  -U postgres \
  -d railway

-- Test 1: Queries de premium_expiration.go
SELECT account_id, email, premium_expired_date
FROM account
WHERE is_premium = true
  AND premium_expired_date IS NOT NULL
  AND premium_expired_date <= NOW();

-- Test 2: Queries de cjnewcluster (compleja)
SELECT COUNT(*)
FROM incident_clusters ic
JOIN account_favorite_locations afl ON
  ST_DistanceSphere(
    ST_MakePoint(ic.center_longitude, ic.center_latitude),
    ST_MakePoint(afl.longitude, afl.latitude)
  ) <= afl.radius
JOIN account a ON afl.account_id = a.account_id
WHERE a.status = 'active'
  AND a.is_premium = true
  AND a.receive_notifications = true
  AND afl.status = true
  AND afl.crime = true;

-- Test 3: Update con operador NOT
UPDATE account
SET receive_notifications = NOT receive_notifications
WHERE account_id = 1
RETURNING account_id, receive_notifications;

-- Test 4: Insert con valores booleanos
INSERT INTO account_favorite_locations
  (account_id, title, crime, traffic_accident, status)
VALUES
  (1, 'Test Location', true, false, true)
RETURNING afl_id, crime, traffic_accident, status;

-- Limpiar datos de prueba
DELETE FROM account_favorite_locations WHERE title = 'Test Location';
```

---

## PASO 4: DEPLOYMENT A PRODUCCIÓN

### 4.1 Commit de Cambios

```bash
cd /Users/garyeikoow/Desktop/alertly/backend

# Verificar cambios
git status
git diff

# Crear commit
git add assets/db/migrations/001_standardize_boolean_columns.sql
git add assets/db/migrations/001_rollback_boolean_columns.sql
git add BOOLEAN_MIGRATION_ANALYSIS.md
git add BOOLEAN_MIGRATION_GO_PATCHES.md
git add BOOLEAN_MIGRATION_EXECUTION_PLAN.md

git add internal/cronjob/premium_expiration.go
git add internal/cronjobs/cjuserank/repository.go
git add internal/cronjobs/cjbadgeearn/repository.go
git add internal/cronjobs/cjcomments/repository.go
git add internal/cronjobs/cjincidentupdate/repository.go
git add internal/editprofile/repository.go
git add internal/tutorial/repository.go
git add internal/account/repository.go

git commit -m "fix: Estandarizar columnas booleanas a tipo BOOLEAN nativo de PostgreSQL

- Migra 8 columnas (SMALLINT/CHAR → BOOLEAN) en account y account_favorite_locations
- Actualiza 21 líneas en 8 archivos Go (queries SQL y conversiones)
- Elimina comparaciones inconsistentes (= 1, = 0, = true, = false)
- Incluye scripts de migración y rollback
- Mejora performance y legibilidad del código

Tablas afectadas:
- account: is_premium, can_update_*, has_finished_tutorial, has_watch_new_incident_tutorial
- account_favorite_locations: status

Archivos modificados:
- internal/cronjob/premium_expiration.go (6 líneas)
- internal/cronjobs/cjuserank/repository.go (1 línea)
- internal/cronjobs/cjbadgeearn/repository.go (1 línea)
- internal/cronjobs/cjcomments/repository.go (1 línea)
- internal/cronjobs/cjincidentupdate/repository.go (2 líneas)
- internal/editprofile/repository.go (5 líneas)
- internal/tutorial/repository.go (1 línea)
- internal/account/repository.go (4 líneas)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
"
```

### 4.2 Deploy en Railway/EC2

```bash
# Si está en Railway con autodeploy:
git push origin main

# Si está en EC2 manual:
ssh user@ec2-instance
cd /path/to/alertly/backend
git pull origin main
go build -o alertly-api cmd/app/main.go
go build -o alertly-cronjob cmd/cronjob/main.go

# Restart services
sudo systemctl restart alertly-api
sudo systemctl restart alertly-cronjob

# Verificar logs
sudo journalctl -u alertly-api -f
sudo journalctl -u alertly-cronjob -f
```

### 4.3 Monitoreo Post-Deploy

```bash
# Verificar health endpoint
curl https://api.alertly.ca/health | jq

# Verificar logs de errores en producción
# Railway: Dashboard → Logs
# EC2: journalctl -u alertly-api -f

# Verificar métricas de base de datos
PGPASSWORD="..." psql -h ... -p ... -U postgres -d railway -c "
SELECT
    COUNT(*) FILTER (WHERE is_premium = true) AS premium_users,
    COUNT(*) FILTER (WHERE receive_notifications = true) AS notif_enabled,
    COUNT(*) AS total_users
FROM account
WHERE status = 'active';
"

# Verificar que cronjobs se ejecutan correctamente
# (deben aparecer notificaciones nuevas cada intervalo configurado)
```

---

## PASO 5: VALIDACIÓN FINAL

### Checklist de Validación

- [ ] Migración SQL ejecutada sin errores
- [ ] Tipos de columnas verificados (todos BOOLEAN)
- [ ] Código Go compilado sin errores
- [ ] Tests unitarios pasan
- [ ] Tests de integración pasan
- [ ] Login funciona correctamente
- [ ] Editar perfil funciona (toggle privacidad, notificaciones)
- [ ] Cronjobs se ejecutan sin errores
- [ ] Queries de notificaciones funcionan (usuarios premium detectados)
- [ ] Sistema de lugares favoritos funciona
- [ ] Tutorial marcado como completado funciona
- [ ] No hay errores en logs de producción
- [ ] Backup de base de datos disponible

---

## ROLLBACK (SI ES NECESARIO)

### Escenario 1: Error en migración SQL

```bash
# Conectar a PostgreSQL
PGPASSWORD="cGA2dBF6G33BgfefcgDb1CDa6CagFcC5" psql \
  -h metro.proxy.rlwy.net \
  -p 48204 \
  -U postgres \
  -d railway

# Ejecutar rollback
\i /Users/garyeikoow/Desktop/alertly/backend/assets/db/migrations/001_rollback_boolean_columns.sql

# Verificar revert
\d account
# is_premium debe ser SMALLINT nuevamente
```

### Escenario 2: Errores en producción después de deploy

```bash
# 1. Revertir código Go
cd /Users/garyeikoow/Desktop/alertly/backend
git revert HEAD

# 2. Ejecutar rollback SQL
PGPASSWORD="..." psql -h ... -p ... -U postgres -d railway \
  -f assets/db/migrations/001_rollback_boolean_columns.sql

# 3. Redeploy
git push origin main
# o en EC2: rebuild y restart services

# 4. Verificar que todo vuelve a funcionar
curl https://api.alertly.ca/health
```

### Escenario 3: Restaurar desde backup

```bash
# Si el rollback SQL no funciona, restaurar snapshot de Railway
# Railway Dashboard → Database → Restore Snapshot

# O restaurar desde dump SQL manual
PGPASSWORD="..." psql -h ... -p ... -U postgres -d railway \
  < backup_before_boolean_migration_YYYYMMDD_HHMMSS.sql
```

---

## CONCLUSIÓN

Esta migración es **segura y necesaria** para:

1. **Eliminar errores de tipo** (queries con `= 1` en columnas BOOLEAN fallan)
2. **Mejorar performance** (BOOLEAN es más eficiente que SMALLINT)
3. **Simplificar código** (eliminar conversiones `dbtypes.BoolToInt()`)
4. **Seguir buenas prácticas** (usar tipos nativos de PostgreSQL)

**Confianza:** Alta - El análisis es exhaustivo y todos los cambios están mapeados.

**Próximos pasos:**
1. Ejecutar migración SQL en horario de bajo tráfico
2. Aplicar parches en código Go
3. Testing exhaustivo antes de deploy
4. Monitorear producción post-deploy

**Soporte:**
- Documentación completa en `BOOLEAN_MIGRATION_ANALYSIS.md`
- Parches detallados en `BOOLEAN_MIGRATION_GO_PATCHES.md`
- Scripts SQL en `assets/db/migrations/`

---

**Fecha de última actualización:** 2026-01-18
**Preparado por:** Claude Code (Análisis Exhaustivo)
