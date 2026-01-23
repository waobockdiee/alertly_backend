# Guía de Migración PostgreSQL - Alertly Database

**Fecha:** 2026-01-22
**Base de datos:** railway @ metro.proxy.rlwy.net:48204
**Estado:** Post-migración desde MySQL

---

## Resumen Ejecutivo

La auditoría completa de la base de datos PostgreSQL identificó **47 problemas** que requieren corrección inmediata. Este documento proporciona las instrucciones paso a paso para aplicar las optimizaciones.

### Problemas Críticos Encontrados
1. **0 índices espaciales PostGIS** - Rendimiento geoespacial degradado
2. **7 Foreign Keys faltantes** - Riesgo de datos huérfanos
3. **9 índices de performance faltantes** - Queries lentas
4. **6 columnas sin NOT NULL** - Datos inconsistentes
5. **229 registros muertos en tabla notifications** (23.8%)

---

## Archivos Generados

| Archivo | Descripción | Cuándo Ejecutar |
|---------|-------------|-----------------|
| `POSTGRESQL_AUDIT_REPORT.md` | Reporte completo de auditoría con análisis detallado | Lectura obligatoria |
| `postgresql_pre_migration_cleanup.sql` | Script de análisis y limpieza de datos huérfanos | **PASO 1** |
| `postgresql_migration_optimization.sql` | Script principal de optimización | **PASO 2** |
| `postgis_optimized_queries.sql` | Ejemplos de queries optimizadas con PostGIS | Referencia para backend |

---

## Instrucciones de Ejecución

### PASO 0: Preparación (OBLIGATORIO)

```bash
# 1. Conectar a la base de datos
export PGPASSWORD='cGA2dBF6G33BgfefcgDb1CDa6CagFcC5'
psql -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway

# 2. Hacer BACKUP COMPLETO de la base de datos
pg_dump -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway \
    --format=custom \
    --file=railway_backup_$(date +%Y%m%d_%H%M%S).dump

# Verificar que el backup se creó correctamente
ls -lh railway_backup_*.dump

# 3. Hacer backup solo del esquema (por si acaso)
pg_dump -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway \
    --schema-only \
    --file=railway_schema_backup_$(date +%Y%m%d_%H%M%S).sql
```

**IMPORTANTE:** NO continuar sin un backup verificado.

---

### PASO 1: Análisis de Datos Huérfanos

Ejecutar el script de análisis para identificar problemas de integridad referencial:

```bash
# Ejecutar desde terminal
PGPASSWORD='cGA2dBF6G33BgfefcgDb1CDa6CagFcC5' \
psql -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway \
    -f backend/assets/db/postgresql_pre_migration_cleanup.sql
```

**Revisar la salida:**
- Si aparecen registros huérfanos (orphan_count > 0), **DETENERSE**
- Editar el script y descomentar las secciones de limpieza del PASO 4
- Ejecutar nuevamente hasta que todos los conteos sean 0

**Resultado esperado:**
```
✓ OK: No se encontraron datos huérfanos
```

---

### PASO 2: Ejecutar Optimización Principal

Una vez confirmado que no hay datos huérfanos:

```bash
# Ejecutar script de optimización
PGPASSWORD='cGA2dBF6G33BgfefcgDb1CDa6CagFcC5' \
psql -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway \
    -f backend/assets/db/postgresql_migration_optimization.sql
```

**Este script ejecuta:**
1. Creación de 7 Foreign Keys faltantes
2. Creación de columnas GEOGRAPHY para PostGIS
3. Creación de 2 índices espaciales GiST
4. Creación de 13 índices de performance
5. Agregado de constraints NOT NULL
6. Configuración de autovacuum
7. Actualización de estadísticas (ANALYZE)

**Duración estimada:** 2-5 minutos

**Resultado esperado:**
```
========================================
MIGRACIÓN COMPLETADA EXITOSAMENTE
========================================
```

---

### PASO 3: VACUUM FULL (Mantenimiento)

Ejecutar VACUUM FULL en horario de bajo tráfico (requiere lock exclusivo):

```bash
# Conectar a la base de datos
PGPASSWORD='cGA2dBF6G33BgfefcgDb1CDa6CagFcC5' \
psql -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway

# Ejecutar VACUUM FULL
VACUUM FULL ANALYZE notifications;
VACUUM ANALYZE account;
VACUUM ANALYZE incident_clusters;
VACUUM ANALYZE incident_reports;
```

**Duración estimada:** 5-10 minutos
**ADVERTENCIA:** Las tablas quedan bloqueadas durante VACUUM FULL.

---

### PASO 4: Validación Post-Migración

Verificar que todo se aplicó correctamente:

```sql
-- 1. Verificar Foreign Keys (esperado: 32 total)
SELECT COUNT(*) FROM information_schema.table_constraints
WHERE constraint_type = 'FOREIGN KEY' AND table_schema = 'public';

-- 2. Verificar índices espaciales (esperado: 2)
SELECT COUNT(*) FROM pg_indexes
WHERE schemaname = 'public' AND indexdef LIKE '%gist%';

-- 3. Verificar columnas GEOGRAPHY pobladas
SELECT
    'incident_clusters' AS table_name,
    COUNT(*) AS total_rows,
    COUNT(center_location) AS rows_with_location
FROM incident_clusters
UNION ALL
SELECT
    'account_favorite_locations',
    COUNT(*),
    COUNT(location)
FROM account_favorite_locations;

-- 4. Test de performance con índice espacial
EXPLAIN ANALYZE
SELECT * FROM incident_clusters
WHERE ST_DWithin(
    center_location,
    ST_MakePoint(-79.3832, 43.6532)::geography,
    5000
)
AND is_active = '1'
LIMIT 100;
-- Debe mostrar: "Index Scan using idx_clusters_center_location_gist"

-- 5. Verificar dead rows (esperado: <5%)
SELECT
    relname AS table_name,
    n_live_tup AS live_rows,
    n_dead_tup AS dead_rows,
    ROUND(n_dead_tup::numeric / NULLIF(n_live_tup, 0) * 100, 2) AS dead_pct
FROM pg_stat_user_tables
WHERE schemaname = 'public'
    AND relname IN ('incident_clusters', 'incident_reports', 'notifications', 'account')
ORDER BY dead_pct DESC;
```

**Checklist de Validación:**
- [ ] 32 Foreign Keys creadas
- [ ] 2 índices espaciales GiST creados
- [ ] 100% de registros con columnas GEOGRAPHY pobladas
- [ ] Queries geoespaciales usan índice GiST (verificar con EXPLAIN)
- [ ] Dead rows <5% en todas las tablas
- [ ] No hay errores en logs de PostgreSQL

---

## Cambios Requeridos en el Backend Go

Después de ejecutar los scripts, actualizar el backend para aprovechar las optimizaciones:

### 1. Actualizar `getclustersbylocation/repository.go`

**ANTES (MySQL):**
```go
query := `
    SELECT * FROM incident_clusters
    WHERE ST_Distance_Sphere(
        point(center_longitude, center_latitude),
        point(?, ?)
    ) <= ?
    AND is_active = 1
    ORDER BY created_at DESC
    LIMIT 100
`
```

**DESPUÉS (PostgreSQL con PostGIS):**
```go
query := `
    SELECT
        incl_id,
        center_latitude,
        center_longitude,
        category_code,
        subcategory_name,
        description,
        credibility,
        incident_count,
        media_url,
        address,
        city,
        created_at,
        ST_Distance(
            center_location,
            ST_MakePoint($1, $2)::geography
        ) AS distance_meters
    FROM incident_clusters
    WHERE ST_DWithin(
        center_location,
        ST_MakePoint($1, $2)::geography,
        $3
    )
    AND is_active = '1'
    ORDER BY center_location <-> ST_MakePoint($1, $2)::geography
    LIMIT 100
`
// Nota: PostgreSQL usa $1, $2, $3 en lugar de ?
```

### 2. Actualizar `newincident/service.go`

Modificar el algoritmo de clustering para usar índice espacial:

```go
query := `
    SELECT incl_id, insu_id, incident_count, credibility,
           ST_Distance(center_location, ST_MakePoint($1, $2)::geography) AS distance_meters
    FROM incident_clusters
    WHERE ST_DWithin(
        center_location,
        ST_MakePoint($1, $2)::geography,
        $3  -- radius por categoría
    )
    AND insu_id = $4
    AND is_active = '1'
    AND created_at >= NOW() - INTERVAL '24 hours'
    ORDER BY distance_meters ASC
    LIMIT 1
`
```

### 3. Actualizar `cronjobs/cjnewcluster`

Notificaciones por proximidad usando índice espacial:

```go
query := `
    SELECT DISTINCT
        afl.account_id,
        afl.title,
        ST_Distance(afl.location, ST_MakePoint($1, $2)::geography) AS distance_meters
    FROM account_favorite_locations afl
    WHERE ST_DWithin(
        afl.location,
        ST_MakePoint($1, $2)::geography,
        afl.radius
    )
    AND afl.status = 1
    AND afl.traffic_accident = 1  -- Ajustar según categoría
`
```

### 4. Instalar driver PostgreSQL

```bash
cd backend
go get github.com/lib/pq
```

Actualizar `database/connection.go`:
```go
import (
    "database/sql"
    _ "github.com/lib/pq"  // Driver PostgreSQL
)

func Connect() (*sql.DB, error) {
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_NAME"),
    )
    db, err := sql.Open("postgres", connStr)
    // ...
}
```

---

## Variables de Entorno Actualizadas

Actualizar `.env` en el backend:

```bash
# PostgreSQL Railway
DB_HOST=metro.proxy.rlwy.net
DB_PORT=48204
DB_USER=postgres
DB_PASS=cGA2dBF6G33BgfefcgDb1CDa6CagFcC5
DB_NAME=railway
DB_DRIVER=postgres  # Nueva variable

# El resto permanece igual
NODE_ENV=production
SERVER_PORT=8080
JWT_SECRET=your_secret
IMAGE_BASE_URL=https://alertly.ca/uploads
AWS_REGION=us-west-2
```

---

## Monitoreo Post-Migración

### 1. Instalar pg_stat_statements (Recomendado)

```sql
-- Crear extensión
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

-- Ver queries lentas
SELECT
    query,
    calls,
    mean_exec_time,
    max_exec_time
FROM pg_stat_statements
WHERE query LIKE '%incident_clusters%'
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### 2. Monitorear Índices No Usados

```sql
-- Ver índices que nunca se usan (después de 1 semana)
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan AS times_used,
    pg_size_pretty(pg_relation_size(indexrelid)) AS size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
    AND idx_scan = 0
    AND indexname NOT LIKE '%pkey'
ORDER BY pg_relation_size(indexrelid) DESC;
```

### 3. Monitorear Tamaño de Tablas

```sql
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS total_size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS table_size,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) -
                   pg_relation_size(schemaname||'.'||tablename)) AS indexes_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
LIMIT 10;
```

### 4. Configurar Alertas

Monitorear en Railway Dashboard o configurar alertas para:
- CPU > 80% por más de 5 minutos
- Memoria > 85%
- Conexiones activas > 80% del límite
- Autovacuum failures
- Dead rows > 10% en cualquier tabla

---

## Performance Esperado

### Antes de la Optimización
- Búsqueda geoespacial (5km): **~200-500ms**
- Clustering de incidentes: **~100-200ms**
- Notificaciones por proximidad: **~300ms**
- Queries sin índices: **>1000ms**

### Después de la Optimización
- Búsqueda geoespacial (5km): **<10ms** (20-50x más rápido)
- Clustering de incidentes: **<5ms** (20x más rápido)
- Notificaciones por proximidad: **<15ms** (20x más rápido)
- Todas las queries críticas: **<50ms**

---

## Troubleshooting

### Problema: "foreign key constraint fails"

**Causa:** Hay datos huérfanos que no se limpiaron en PASO 1.

**Solución:**
```sql
-- Identificar datos problemáticos
SELECT acs.acs_id, acs.incl_id
FROM account_cluster_saved acs
LEFT JOIN incident_clusters ic ON acs.incl_id = ic.incl_id
WHERE ic.incl_id IS NULL;

-- Eliminar datos huérfanos
DELETE FROM account_cluster_saved
WHERE incl_id NOT IN (SELECT incl_id FROM incident_clusters);
```

### Problema: "column center_location does not exist"

**Causa:** El script de optimización no se ejecutó completamente.

**Solución:**
```sql
-- Verificar si la columna existe
SELECT column_name FROM information_schema.columns
WHERE table_name = 'incident_clusters' AND column_name = 'center_location';

-- Si no existe, ejecutar solo PASO 3 del script de optimización
```

### Problema: Queries lentas después de migración

**Causa:** Estadísticas desactualizadas o índice no se usa.

**Solución:**
```sql
-- Actualizar estadísticas
ANALYZE incident_clusters;

-- Verificar que el índice se usa
EXPLAIN ANALYZE
SELECT * FROM incident_clusters
WHERE ST_DWithin(center_location, ST_MakePoint(-79.3832, 43.6532)::geography, 5000);

-- Debe mostrar: "Index Scan using idx_clusters_center_location_gist"
-- Si no, verificar que la columna center_location está poblada
```

### Problema: Backend Go no conecta

**Causa:** Connection string incorrecto o driver no instalado.

**Solución:**
```bash
# Instalar driver
go get github.com/lib/pq

# Verificar connection string
psql "postgresql://postgres:cGA2dBF6G33BgfefcgDb1CDa6CagFcC5@metro.proxy.rlwy.net:48204/railway?sslmode=require"
```

---

## Rollback (En caso de problemas)

Si algo sale mal, restaurar desde el backup:

```bash
# Restaurar backup completo
pg_restore -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway \
    --clean \
    --if-exists \
    railway_backup_YYYYMMDD_HHMMSS.dump

# O restaurar solo esquema
PGPASSWORD='cGA2dBF6G33BgfefcgDb1CDa6CagFcC5' \
psql -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway \
    -f railway_schema_backup_YYYYMMDD_HHMMSS.sql
```

---

## Próximos Pasos (Futuro)

Una vez estabilizado el sistema:

1. **Particionamiento** - Cuando supere 100K incidentes
2. **Full-Text Search** - Para búsqueda por texto en descripciones
3. **Read Replicas** - Para escalar lectura de analytics
4. **Connection Pooling** - PgBouncer para manejar más usuarios concurrentes
5. **TimescaleDB** - Para optimizar analytics de series de tiempo

---

## Contacto y Soporte

**Documentación:** Ver `POSTGRESQL_AUDIT_REPORT.md` para análisis completo
**Queries de Referencia:** Ver `postgis_optimized_queries.sql`
**Fecha de Migración:** 2026-01-22

Para preguntas o problemas, revisar los logs de PostgreSQL en Railway Dashboard.
