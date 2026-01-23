# PostgreSQL Quick Reference - Alertly

**Database:** railway @ metro.proxy.rlwy.net:48204
**Status:** Post-migración desde MySQL
**Date:** 2026-01-22

---

## Conexión Rápida

```bash
# Terminal
export PGPASSWORD='cGA2dBF6G33BgfefcgDb1CDa6CagFcC5'
psql -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway

# Connection String
postgresql://postgres:cGA2dBF6G33BgfefcgDb1CDa6CagFcC5@metro.proxy.rlwy.net:48204/railway?sslmode=require
```

---

## Estado Actual (Post-Auditoría)

| Métrica | Valor |
|---------|-------|
| Versión PostgreSQL | 16.9 |
| PostGIS | 3.7.0dev |
| Tablas | 33 |
| Índices Total | 104 |
| Índices Espaciales | 0 ❌ (DEBE SER 2) |
| Foreign Keys | 25 ❌ (DEBE SER 32) |
| Registros en incident_clusters | 1,058 |
| Registros en incident_reports | 2,359 |
| Tamaño DB Total | ~5 MB |

---

## Problemas Críticos

| # | Problema | Impacto | Solución |
|---|----------|---------|----------|
| 1 | Sin índices espaciales PostGIS | Queries geoespaciales 50x más lentas | Ejecutar migration script |
| 2 | 7 Foreign Keys faltantes | Datos huérfanos posibles | Agregar FKs |
| 3 | 9 índices de performance faltantes | Queries lentas | Crear índices |
| 4 | 229 dead rows en notifications (23.8%) | Espacio desperdiciado | VACUUM FULL |
| 5 | Columnas críticas permiten NULL | Datos inconsistentes | Agregar NOT NULL |

---

## Ejecución Rápida (3 Comandos)

```bash
# 1. BACKUP (OBLIGATORIO)
pg_dump -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway \
    --format=custom --file=railway_backup_$(date +%Y%m%d).dump

# 2. ANÁLISIS + LIMPIEZA
PGPASSWORD='cGA2dBF6G33BgfefcgDb1CDa6CagFcC5' \
psql -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway \
    -f backend/assets/db/postgresql_pre_migration_cleanup.sql

# 3. OPTIMIZACIÓN PRINCIPAL
PGPASSWORD='cGA2dBF6G33BgfefcgDb1CDa6CagFcC5' \
psql -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway \
    -f backend/assets/db/postgresql_migration_optimization.sql
```

**Duración total:** ~5-10 minutos

---

## Queries Críticas a Actualizar

### 1. Búsqueda por Radio (getclustersbylocation)

```sql
-- ANTES (MySQL - LENTO)
SELECT * FROM incident_clusters
WHERE ST_Distance_Sphere(point(center_longitude, center_latitude), point(?, ?)) <= ?
AND is_active = 1
ORDER BY created_at DESC LIMIT 100;

-- DESPUÉS (PostgreSQL - RÁPIDO)
SELECT * FROM incident_clusters
WHERE ST_DWithin(center_location, ST_MakePoint($1, $2)::geography, $3)
AND is_active = '1'
ORDER BY center_location <-> ST_MakePoint($1, $2)::geography LIMIT 100;
```

### 2. Clustering de Incidentes (newincident)

```sql
-- DESPUÉS (PostgreSQL con índice espacial)
SELECT incl_id, ST_Distance(center_location, ST_MakePoint($1, $2)::geography) AS dist
FROM incident_clusters
WHERE ST_DWithin(center_location, ST_MakePoint($1, $2)::geography, $3)
AND insu_id = $4 AND is_active = '1'
AND created_at >= NOW() - INTERVAL '24 hours'
ORDER BY dist ASC LIMIT 1;
```

### 3. Notificaciones por Proximidad (cjnewcluster)

```sql
-- DESPUÉS (PostgreSQL)
SELECT DISTINCT account_id
FROM account_favorite_locations
WHERE ST_DWithin(location, ST_MakePoint($1, $2)::geography, radius)
AND status = 1 AND traffic_accident = 1;
```

---

## Backend Go - Cambios Requeridos

### 1. Driver PostgreSQL

```bash
go get github.com/lib/pq
```

### 2. Connection String

```go
import _ "github.com/lib/pq"

connStr := fmt.Sprintf(
    "host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASS"),
    os.Getenv("DB_NAME"),
)
db, err := sql.Open("postgres", connStr)
```

### 3. Placeholders

MySQL usa `?` → PostgreSQL usa `$1, $2, $3`

```go
// ANTES (MySQL)
db.Query("SELECT * FROM table WHERE id = ?", id)

// DESPUÉS (PostgreSQL)
db.Query("SELECT * FROM table WHERE id = $1", id)
```

---

## Variables de Entorno

```bash
# .env (backend)
DB_HOST=metro.proxy.rlwy.net
DB_PORT=48204
DB_USER=postgres
DB_PASS=cGA2dBF6G33BgfefcgDb1CDa6CagFcC5
DB_NAME=railway
DB_DRIVER=postgres

NODE_ENV=production
SERVER_PORT=8080
IP_SERVER=0.0.0.0
JWT_SECRET=your_secret
IMAGE_BASE_URL=https://alertly.ca/uploads
AWS_REGION=us-west-2
REFERRAL_API_KEY=your_api_key

# APNs (Push Notifications)
APNS_ENV=production
APNS_P12_BASE64=<base64_encoded_cert>
APNS_P12_PASS=alertly123
APNS_TOPIC=com.anonymous.Alertly
```

---

## Comandos Útiles

### Monitoreo

```sql
-- Ver tamaño de tablas
SELECT tablename, pg_size_pretty(pg_total_relation_size('public.'||tablename))
FROM pg_tables WHERE schemaname = 'public' ORDER BY pg_total_relation_size('public.'||tablename) DESC;

-- Ver dead rows
SELECT relname, n_live_tup, n_dead_tup,
       ROUND(n_dead_tup::numeric / NULLIF(n_live_tup, 0) * 100, 2) AS dead_pct
FROM pg_stat_user_tables WHERE schemaname = 'public' AND n_dead_tup > 0;

-- Ver índices no usados
SELECT tablename, indexname, idx_scan
FROM pg_stat_user_indexes WHERE schemaname = 'public' AND idx_scan = 0;

-- Ver queries lentas (requiere pg_stat_statements)
SELECT query, calls, mean_exec_time, max_exec_time
FROM pg_stat_statements WHERE query LIKE '%incident_clusters%'
ORDER BY mean_exec_time DESC LIMIT 10;
```

### Mantenimiento

```sql
-- VACUUM tablas críticas
VACUUM ANALYZE incident_clusters;
VACUUM ANALYZE incident_reports;
VACUUM ANALYZE notifications;
VACUUM FULL ANALYZE notifications;  -- Requiere lock exclusivo

-- Actualizar estadísticas
ANALYZE;

-- Ver procesos activos
SELECT pid, usename, query, state FROM pg_stat_activity WHERE datname = 'railway';

-- Matar query lenta
SELECT pg_cancel_backend(pid);  -- Cancelar
SELECT pg_terminate_backend(pid);  -- Forzar término
```

---

## Performance Esperado

| Query | Antes | Después | Mejora |
|-------|-------|---------|--------|
| Búsqueda 5km | 200-500ms | <10ms | 20-50x |
| Clustering | 100-200ms | <5ms | 20-40x |
| Notificaciones | 300ms | <15ms | 20x |
| Bounding box | 150ms | <5ms | 30x |

---

## Validación Post-Migración

```sql
-- ✓ 32 Foreign Keys
SELECT COUNT(*) FROM information_schema.table_constraints
WHERE constraint_type = 'FOREIGN KEY' AND table_schema = 'public';

-- ✓ 2 Índices Espaciales
SELECT COUNT(*) FROM pg_indexes
WHERE schemaname = 'public' AND indexdef LIKE '%gist%';

-- ✓ center_location poblado
SELECT COUNT(*), COUNT(center_location) FROM incident_clusters;

-- ✓ Índice se usa
EXPLAIN SELECT * FROM incident_clusters
WHERE ST_DWithin(center_location, ST_MakePoint(-79.3832, 43.6532)::geography, 5000);
-- Debe mostrar: "Index Scan using idx_clusters_center_location_gist"
```

---

## Troubleshooting Rápido

| Error | Causa | Fix |
|-------|-------|-----|
| `foreign key constraint fails` | Datos huérfanos | Ejecutar cleanup script |
| `column center_location does not exist` | Script no completado | Re-ejecutar migration script |
| `Index Scan not used` | Estadísticas desactualizadas | `ANALYZE incident_clusters;` |
| `connection refused` | Credenciales incorrectas | Verificar .env |
| Queries lentas | Sin índices espaciales | Ejecutar migration script |

---

## Rollback

```bash
# Restaurar desde backup
pg_restore -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway \
    --clean --if-exists railway_backup_YYYYMMDD.dump
```

---

## Archivos Generados

1. `POSTGRESQL_AUDIT_REPORT.md` - Análisis completo (47 problemas)
2. `POSTGRESQL_MIGRATION_GUIDE.md` - Guía paso a paso
3. `postgresql_pre_migration_cleanup.sql` - Limpieza de datos
4. `postgresql_migration_optimization.sql` - Script principal
5. `postgis_optimized_queries.sql` - 17 ejemplos de queries optimizadas
6. `POSTGRESQL_QUICK_REFERENCE.md` - Esta guía

---

## Checklist de Migración

- [ ] Backup completo creado y verificado
- [ ] Pre-migration cleanup ejecutado (sin huérfanos)
- [ ] Migration optimization ejecutado exitosamente
- [ ] VACUUM FULL ejecutado
- [ ] 32 Foreign Keys verificadas
- [ ] 2 Índices espaciales creados
- [ ] Columnas GEOGRAPHY pobladas 100%
- [ ] Backend Go actualizado (driver + queries)
- [ ] Variables de entorno actualizadas
- [ ] Performance validado (<10ms búsquedas)
- [ ] Tests end-to-end pasando
- [ ] Monitoring configurado

---

## Contacto

**Auditoría:** Claude Sonnet 4.5 (PostgreSQL Expert)
**Fecha:** 2026-01-22
**Documentación completa:** Ver `POSTGRESQL_AUDIT_REPORT.md`

Para más detalles, revisar los 5 archivos generados en `/backend/assets/db/` y `/backend/`.
