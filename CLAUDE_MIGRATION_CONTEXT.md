# CLAUDE MIGRATION CONTEXT - MySQL to PostgreSQL

> **IMPORTANTE**: Este archivo es el punto de entrada para Claude en futuras sesiones. Contiene todo el contexto de la migración de AWS (MySQL) a Railway (PostgreSQL).

---

## USO DE AGENTES ESPECIALIZADOS

**REGLA CRÍTICA**: Siempre usar el agente adecuado para cada tarea:

| Tarea | Agente a usar |
|-------|---------------|
| Análisis/cambios en PostgreSQL | `postgresql-guru` |
| Análisis/cambios en código Go | `golang-guru` |
| Exploración del codebase | `Explore` |
| Tareas full-stack Go + Frontend | `golang-svelte-pro` |

**NO hacer cambios reactivos sin análisis completo**. Siempre verificar primero:
1. Cómo estaba en MySQL original (AWS)
2. Cómo está actualmente en PostgreSQL (Railway)
3. Qué espera el código Go
4. Qué espera el frontend

---

## Conexiones a Bases de Datos

### AWS MySQL (BASE DE DATOS ORIGINAL - REFERENCIA)

> **IMPORTANTE**: Esta es la DB original que funcionaba correctamente. Cuando tengas dudas sobre tipos de datos, estructura o valores, CONSULTA ESTA DB como referencia.

```
Host: alertly-mysql-freetier.c3qmq4y86s84.us-west-2.rds.amazonaws.com
Puerto: 3306
Usuario: adminalertly
Password: Po1Ng2O3;
Database: alertly
```

**Comando para conectar**:
```bash
mysql -h alertly-mysql-freetier.c3qmq4y86s84.us-west-2.rds.amazonaws.com -P 3306 -u adminalertly -p'Po1Ng2O3;' alertly
```

**Documentación de AWS**:
- [AWS_DEPLOYMENT_GUIDE.md](../AWS_DEPLOYMENT_GUIDE.md) - Guía completa de deployment en AWS
- [DEPLOYMENT_GUIDE.md](../DEPLOYMENT_GUIDE.md) - Guía general de deployment

### Railway PostgreSQL (BASE DE DATOS ACTUAL - PRODUCCIÓN)

```
Host: metro.proxy.rlwy.net
Puerto: 48204
Usuario: postgres
Password: cGA2dBF6G33BgfefcgDb1CDa6CagFcC5
Database: railway
Connection string: postgres://postgres:cGA2dBF6G33BgfefcgDb1CDa6CagFcC5@metro.proxy.rlwy.net:48204/railway
```

**Comando para conectar**:
```bash
PGPASSWORD='cGA2dBF6G33BgfefcgDb1CDa6CagFcC5' psql -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway
```

---

## Estado Actual del Proyecto

**Fecha última actualización**: 2026-01-18

**Situación**: La app Alertly migró de AWS (MySQL) a Railway (PostgreSQL/PostGIS). El frontend (React Native) NO ha cambiado. Todos los problemas son del backend.

---

## Problema Principal: Incompatibilidad de Tipos Booleanos

### El problema
MySQL usaba `TINYINT`/`SMALLINT` para valores booleanos (0/1). PostgreSQL es más estricto con tipos:
- `boolean = 1` → ERROR en PostgreSQL
- `NOT smallint` → ERROR en PostgreSQL

### Estado actual de las columnas (INCONSISTENTE)

| Tabla | Columna | Tipo Actual | Problema |
|-------|---------|-------------|----------|
| `account` | `is_premium` | SMALLINT | Queries usan `= true` |
| `account` | `receive_notifications` | BOOLEAN | Queries usan `= 1` |
| `account` | `has_finished_tutorial` | CHAR(2) | Conversiones manuales |
| `account` | `is_private_profile` | CHAR(2) | Inconsistente |
| `account_favorite_locations` | `crime`, `traffic_accident`, etc. | BOOLEAN | Fueron cambiadas |

### Solución pendiente
Estandarizar TODAS las columnas booleanas a tipo `BOOLEAN` nativo de PostgreSQL y actualizar las queries en Go.

---

## Fixes Ya Aplicados (NO revertir)

### 1. ST_DistanceSphere operador corregido
**Archivo**: `internal/newincident/repository.go`
```go
// ANTES (incorrecto): < $6
// AHORA (correcto): <= $6
ST_DistanceSphere(...) <= $6
```
**Razón**: Clusters a distancia exacta no matcheaban.

### 2. is_active comparación directa
**Archivo**: `internal/newincident/repository.go`
```go
// ANTES: CAST(is_active AS TEXT) = '1'
// AHORA: is_active = '1'
```
**Razón**: is_active es CHAR(1), comparación directa funciona.

### 3. vote permite NULL
**Tabla**: `incident_reports`
```sql
ALTER TABLE incident_reports ALTER COLUMN vote DROP NOT NULL;
```
**Razón**: Al crear incidente nuevo, vote es NULL (solo tiene valor al votar).

### 4. receive_notifications es BOOLEAN
**Tabla**: `account`
```sql
ALTER TABLE account ALTER COLUMN receive_notifications TYPE BOOLEAN;
```
**Razón**: PostgreSQL no permite `NOT smallint`.

### 5. Categorías en account_favorite_locations son BOOLEAN
**Tabla**: `account_favorite_locations`
```sql
-- crime, traffic_accident, medical_emergency, etc. son BOOLEAN
```
**Razón**: El código Go enviaba `true/false`, no `1/0`.

---

## Fixes Pendientes (REQUIEREN ACCIÓN)

### 1. Estandarizar columnas booleanas en `account`
Las siguientes columnas deben ser BOOLEAN:
- `is_premium` (actualmente SMALLINT)
- `is_private_profile` (actualmente CHAR)
- `has_finished_tutorial` (actualmente CHAR)
- `has_watch_new_incident_tutorial` (actualmente CHAR)

### 2. Actualizar queries en Go
Cambiar todas las comparaciones `= 1` a `= true` para columnas BOOLEAN.

**Archivos afectados**:
- `internal/cronjobs/cjnewcluster/repository.go`
- `internal/cronjobs/cjuserank/repository.go`
- `internal/cronjobs/cjcomments/repository.go`
- Otros (ver documentación detallada)

---

## Documentación Detallada

### Análisis completo
- [BOOLEAN_MIGRATION_README.md](./BOOLEAN_MIGRATION_README.md) - Resumen ejecutivo
- [BOOLEAN_MIGRATION_ANALYSIS.md](./BOOLEAN_MIGRATION_ANALYSIS.md) - Análisis técnico exhaustivo
- [BOOLEAN_MIGRATION_GO_PATCHES.md](./BOOLEAN_MIGRATION_GO_PATCHES.md) - Cambios exactos en código Go
- [BOOLEAN_MIGRATION_EXECUTION_PLAN.md](./BOOLEAN_MIGRATION_EXECUTION_PLAN.md) - Plan paso a paso

### Scripts SQL
- [001_standardize_boolean_columns.sql](./assets/db/migrations/001_standardize_boolean_columns.sql) - Migración
- [001_rollback_boolean_columns.sql](./assets/db/migrations/001_rollback_boolean_columns.sql) - Rollback

### Otros documentos relevantes
- [POSTGRESQL_MIGRATION_FIX_REPORT.md](./POSTGRESQL_MIGRATION_FIX_REPORT.md) - Reporte de fixes iniciales
- [CLAUDE.md](../CLAUDE.md) - Instrucciones generales del proyecto

---

## Esquema Original MySQL

**Ubicación**: `assets/db/db.sql`

Referencia para ver cómo estaban definidas las columnas originalmente:
```sql
-- account
is_premium TINYINT UNSIGNED NULL DEFAULT 1
receive_notifications SMALLINT UNSIGNED NULL DEFAULT 1
is_private_profile CHAR(1) NULL DEFAULT '0'
has_finished_tutorial CHAR(2) NULL DEFAULT '0'

-- account_favorite_locations
crime SMALLINT NULL DEFAULT 1
traffic_accident SMALLINT NULL DEFAULT 1
-- etc.
```

---

## Errores Comunes y Soluciones

### Error: `operator does not exist: boolean = integer`
**Causa**: Comparando columna BOOLEAN con número (1 o 0)
**Solución**: Cambiar `= 1` a `= true` o `= 0` a `= false`

### Error: `argument of NOT must be type boolean, not type smallint`
**Causa**: Usando `NOT` con columna SMALLINT
**Solución**: Cambiar columna a BOOLEAN o usar `CASE WHEN col = 1 THEN 0 ELSE 1 END`

### Error: `null value in column "X" violates not-null constraint`
**Causa**: El agente de migración agregó NOT NULL a columnas que deberían permitir NULL
**Solución**: `ALTER TABLE X ALTER COLUMN Y DROP NOT NULL;`

### Error: `invalid input syntax for type smallint: "true"`
**Causa**: Código Go envía boolean pero columna es SMALLINT
**Solución**: Cambiar columna a BOOLEAN o código a enviar 1/0

---

## Comandos Útiles

### Conectar a PostgreSQL Railway
```bash
PGPASSWORD='cGA2dBF6G33BgfefcgDb1CDa6CagFcC5' psql -h metro.proxy.rlwy.net -p 48204 -U postgres -d railway
```

### Ver tipos de columnas de una tabla
```sql
SELECT column_name, data_type, column_default, is_nullable
FROM information_schema.columns
WHERE table_name = 'account';
```

### Cambiar columna a BOOLEAN
```sql
ALTER TABLE tabla ALTER COLUMN columna DROP DEFAULT;
ALTER TABLE tabla ALTER COLUMN columna TYPE BOOLEAN USING CASE WHEN columna = 1 THEN true ELSE false END;
ALTER TABLE tabla ALTER COLUMN columna SET DEFAULT true;
```

### Build y test del backend
```bash
cd /Users/garyeikoow/Desktop/alertly/backend
go build ./cmd/app/...
go build ./internal/...
```

---

## Historial de Cambios

### 2026-01-18
- Migración inicial de AWS a Railway
- Fix: ST_DistanceSphere `<` → `<=`
- Fix: is_active comparación directa
- Fix: vote permite NULL
- Fix: receive_notifications → BOOLEAN
- Fix: categorías en account_favorite_locations → BOOLEAN
- Análisis completo de columnas booleanas pendientes
- Creación de documentación y scripts de migración

---

## Para Claude en Futuras Sesiones

### Proceso obligatorio cuando el usuario reporte un error:

1. **Lee este archivo primero** para entender el contexto
2. **Revisa los logs** en `/Users/garyeikoow/Desktop/alertly/md/LOGS.md`
3. **ANTES de hacer cualquier cambio**:
   - Conecta a AWS MySQL para ver cómo estaba originalmente
   - Conecta a Railway PostgreSQL para ver el estado actual
   - Compara ambos
4. **Usa el agente adecuado**:
   - `postgresql-guru` → Cambios en la DB
   - `golang-guru` → Análisis de código Go
   - `Explore` → Buscar en el codebase
5. **NO hagas cambios reactivos** - Analiza primero, cambia después

### Agentes disponibles y cuándo usarlos:

| Situación | Agente | Ejemplo |
|-----------|--------|---------|
| Error de tipo en PostgreSQL | `postgresql-guru` | "boolean = integer" |
| Error en código Go | `golang-guru` | "undefined: sql.NullString" |
| Buscar queries que usan X columna | `golang-guru` o `Explore` | "dónde se usa is_premium" |
| Cambiar tipo de columna en DB | `postgresql-guru` | "ALTER TABLE" |
| Entender flujo de datos | `Explore` | "cómo funciona el login" |

### Principios clave:

1. **El frontend NO ha cambiado** - Si algo falla, el problema es la DB o el código Go
2. **AWS MySQL es la referencia** - Cuando tengas dudas, consulta cómo estaba originalmente
3. **PostgreSQL debe comportarse igual que MySQL** - Mismos tipos, mismos valores
4. **Analiza antes de actuar** - No hagas cambios sin entender el impacto completo
