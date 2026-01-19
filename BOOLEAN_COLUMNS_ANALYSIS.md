# Análisis Completo: Columnas Booleanas MySQL vs PostgreSQL

## Fecha: 2026-01-18

---

## RESUMEN EJECUTIVO

El objetivo es hacer que PostgreSQL se comporte IGUAL que MySQL en cuanto a columnas booleanas, para que las queries existentes en Go funcionen sin modificación.

**Problema Principal:**
- MySQL usa `TINYINT UNSIGNED` y `SMALLINT UNSIGNED` para valores 0/1
- PostgreSQL tiene columnas mezcladas: algunas `BOOLEAN`, otras `SMALLINT`, otras `CHAR(2)`
- El código Go usa comparaciones como `= 1`, `= 0`, que funcionan en MySQL pero son inconsistentes en PostgreSQL

**Solución:**
Convertir todas las columnas a `SMALLINT` en PostgreSQL para mantener compatibilidad total con el código Go existente.

---

## 1. TABLA: `account`

### Columnas Analizadas (Booleanas)

| Columna | MySQL (REFERENCIA) | PostgreSQL (ACTUAL) | ¿Coincide? |
|---------|-------------------|-------------------|------------|
| `is_private_profile` | `TINYINT(1) NULL DEFAULT 0` | `BOOLEAN NOT NULL DEFAULT false` | ❌ NO |
| `is_premium` | `TINYINT UNSIGNED NULL DEFAULT 1` | `SMALLINT NOT NULL DEFAULT 1` | ✅ SÍ |
| `has_finished_tutorial` | `CHAR(2) NULL DEFAULT 0` | `CHAR(2) NOT NULL DEFAULT '0'` | ⚠️ TIPO INCORRECTO |
| `has_watch_new_incident_tutorial` | `CHAR(2) NULL DEFAULT 0` | `CHAR(2) NOT NULL DEFAULT '0'` | ⚠️ TIPO INCORRECTO |
| `receive_notifications` | `SMALLINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |

### Valores Actuales en PostgreSQL (Muestra)

```
account_id | is_premium | receive_notifications | is_private_profile | has_finished_tutorial | has_watch_new_incident_tutorial
-----------+------------+-----------------------+--------------------+-----------------------+---------------------------------
         3 |          1 | t                     | f                  | 1                     | 0
         5 |          0 | t                     | f                  | 1                     | 0
         6 |          0 | t                     | f                  | 0                     | 0
```

**Observaciones:**
- `is_premium`: CORRECTO (ya es SMALLINT)
- `receive_notifications`: INCORRECTO (es BOOLEAN, debería ser SMALLINT)
- `is_private_profile`: INCORRECTO (es BOOLEAN, debería ser SMALLINT)
- `has_finished_tutorial`: INCORRECTO (es CHAR(2) pero debería ser SMALLINT)
- `has_watch_new_incident_tutorial`: INCORRECTO (es CHAR(2) pero debería ser SMALLINT)

### Categorías de Notificación en `account`

| Columna | MySQL (REFERENCIA) | PostgreSQL (ACTUAL) | ¿Coincide? |
|---------|-------------------|-------------------|------------|
| `crime` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `traffic_accident` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `medical_emergency` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `fire_incident` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `vandalism` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `suspicious_activity` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `infrastructure_issues` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `extreme_weather` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `community_events` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `dangerous_wildlife_sighting` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `positive_actions` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `lost_pet` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |
| `incident_as_update` | `SMALLINT UNSIGNED NULL DEFAULT 0` | `SMALLINT NOT NULL DEFAULT 0` | ✅ SÍ |

**Observaciones:**
- Todas las categorías de notificación están CORRECTAS (ya son SMALLINT).

---

## 2. TABLA: `account_favorite_locations`

### Columnas Analizadas (Booleanas)

| Columna | MySQL (REFERENCIA) | PostgreSQL (ACTUAL) | ¿Coincide? |
|---------|-------------------|-------------------|------------|
| `crime` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |
| `traffic_accident` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |
| `medical_emergency` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |
| `fire_incident` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |
| `vandalism` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |
| `suspicious_activity` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |
| `infrastructure_issues` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |
| `extreme_weather` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |
| `community_events` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |
| `dangerous_wildlife_sighting` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |
| `positive_actions` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |
| `lost_pet` | `TINYINT UNSIGNED NULL DEFAULT 1` | `BOOLEAN NOT NULL DEFAULT true` | ❌ NO |

### Valores Actuales en PostgreSQL (Muestra)

```
afl_id | account_id | crime | traffic_accident | medical_emergency | fire_incident
-------+------------+-------+------------------+-------------------+---------------
     1 |          1 | t     | t                | t                 | t
     3 |          3 | t     | t                | t                 | t
     4 |          5 | t     | t                | t                 | t
```

**Observaciones:**
- TODAS las columnas de categorías están INCORRECTAS (son BOOLEAN, deberían ser SMALLINT).

---

## 3. USO EN CÓDIGO GO

### Patrones de Query Encontrados

```go
// De: internal/cronjobs/cjincidentupdate/repository.go:140-141
AND a.is_premium = 1
AND a.receive_notifications = 1

// De: internal/cronjobs/cjuserank/repository.go:28
status = 'active' AND receive_notifications = 1

// De: internal/cronjobs/cjcomments/repository.go:140
WHERE dt.account_id IN (%s) AND a.status = 'active' AND a.is_premium = 1 AND a.receive_notifications = 1

// De: internal/cronjobs/cjbadgeearn/repository.go:41
status = 'active' AND receive_notifications = 1
```

**Análisis:**
- El código Go usa comparaciones directas `= 1` y `= 0`
- NO usa conversiones CAST o operadores booleanos específicos de PostgreSQL
- Espera que las columnas sean numéricas (0/1), no booleanas (true/false)

---

## 4. PLAN DE MIGRACIÓN

### Cambios Requeridos en PostgreSQL

#### Tabla: `account`

1. `is_private_profile`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 0`
2. `receive_notifications`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
3. `has_finished_tutorial`: `CHAR(2)` → `SMALLINT NOT NULL DEFAULT 0`
4. `has_watch_new_incident_tutorial`: `CHAR(2)` → `SMALLINT NOT NULL DEFAULT 0`

#### Tabla: `account_favorite_locations`

1. `crime`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
2. `traffic_accident`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
3. `medical_emergency`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
4. `fire_incident`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
5. `vandalism`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
6. `suspicious_activity`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
7. `infrastructure_issues`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
8. `extreme_weather`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
9. `community_events`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
10. `dangerous_wildlife_sighting`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
11. `positive_actions`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`
12. `lost_pet`: `BOOLEAN` → `SMALLINT NOT NULL DEFAULT 1`

### Lógica de Conversión de Datos

- `BOOLEAN TRUE` → `1`
- `BOOLEAN FALSE` → `0`
- `CHAR(2) '0'` → `0`
- `CHAR(2) '1'` → `1`
- Otros valores CHAR → `0` (por seguridad)

---

## 5. ESTRATEGIA DE MIGRACIÓN

### Fase 1: Backup
```sql
-- Crear tabla de respaldo
CREATE TABLE account_backup_20260118 AS SELECT * FROM account;
CREATE TABLE account_favorite_locations_backup_20260118 AS SELECT * FROM account_favorite_locations;
```

### Fase 2: Migración de account
```sql
-- Convertir is_private_profile
ALTER TABLE account ALTER COLUMN is_private_profile TYPE SMALLINT USING (CASE WHEN is_private_profile THEN 1 ELSE 0 END);
ALTER TABLE account ALTER COLUMN is_private_profile SET DEFAULT 0;
ALTER TABLE account ALTER COLUMN is_private_profile SET NOT NULL;

-- Convertir receive_notifications
ALTER TABLE account ALTER COLUMN receive_notifications TYPE SMALLINT USING (CASE WHEN receive_notifications THEN 1 ELSE 0 END);
ALTER TABLE account ALTER COLUMN receive_notifications SET DEFAULT 1;
ALTER TABLE account ALTER COLUMN receive_notifications SET NOT NULL;

-- Convertir has_finished_tutorial
ALTER TABLE account ALTER COLUMN has_finished_tutorial TYPE SMALLINT USING (
    CASE
        WHEN has_finished_tutorial = '1' THEN 1
        ELSE 0
    END
);
ALTER TABLE account ALTER COLUMN has_finished_tutorial SET DEFAULT 0;
ALTER TABLE account ALTER COLUMN has_finished_tutorial SET NOT NULL;

-- Convertir has_watch_new_incident_tutorial
ALTER TABLE account ALTER COLUMN has_watch_new_incident_tutorial TYPE SMALLINT USING (
    CASE
        WHEN has_watch_new_incident_tutorial = '1' THEN 1
        ELSE 0
    END
);
ALTER TABLE account ALTER COLUMN has_watch_new_incident_tutorial SET DEFAULT 0;
ALTER TABLE account ALTER COLUMN has_watch_new_incident_tutorial SET NOT NULL;
```

### Fase 3: Migración de account_favorite_locations
```sql
-- Lista de columnas a convertir
ALTER TABLE account_favorite_locations
  ALTER COLUMN crime TYPE SMALLINT USING (CASE WHEN crime THEN 1 ELSE 0 END),
  ALTER COLUMN crime SET DEFAULT 1,
  ALTER COLUMN crime SET NOT NULL;

-- Repetir para cada categoría...
```

### Fase 4: Verificación
```sql
-- Verificar tipos de datos
SELECT column_name, data_type, column_default, is_nullable
FROM information_schema.columns
WHERE table_name IN ('account', 'account_favorite_locations')
  AND column_name IN ('is_private_profile', 'receive_notifications', 'has_finished_tutorial', 'crime', 'traffic_accident');

-- Verificar valores
SELECT account_id, is_premium, receive_notifications, is_private_profile, has_finished_tutorial
FROM account LIMIT 5;

SELECT afl_id, account_id, crime, traffic_accident, medical_emergency
FROM account_favorite_locations LIMIT 5;
```

---

## 6. VALIDACIÓN POST-MIGRACIÓN

### Tests a Ejecutar

1. **Test de Query Numérica:**
```sql
SELECT COUNT(*) FROM account WHERE is_premium = 1;
SELECT COUNT(*) FROM account WHERE receive_notifications = 1;
SELECT COUNT(*) FROM account_favorite_locations WHERE crime = 1;
```

2. **Test de Operadores:**
```sql
SELECT COUNT(*) FROM account WHERE is_private_profile = 0;
SELECT COUNT(*) FROM account WHERE NOT receive_notifications;
```

3. **Test de Inserción:**
```sql
INSERT INTO account (email, password, nickname, is_private_profile, receive_notifications)
VALUES ('test@test.com', 'hash', 'testuser', 0, 1);

SELECT is_private_profile, receive_notifications FROM account WHERE email = 'test@test.com';
```

---

## 7. IMPACTO EN CÓDIGO GO

### Código Afectado (Confirmado Compatible)

✅ **No requiere cambios** - Las queries ya usan `= 1` y `= 0`:
- `internal/cronjobs/cjincidentupdate/repository.go`
- `internal/cronjobs/cjuserank/repository.go`
- `internal/cronjobs/cjcomments/repository.go`
- `internal/cronjobs/cjbadgeearn/repository.go`

### Ventajas de SMALLINT vs BOOLEAN

1. **Compatibilidad Total:** `= 1` y `= 0` funcionan sin CAST
2. **Consistencia:** Mismo tipo en MySQL y PostgreSQL
3. **Sin Cambios en Go:** El código existente funciona sin modificación
4. **Espacio Eficiente:** SMALLINT (2 bytes) vs INT (4 bytes)

---

## 8. RIESGOS Y MITIGACIÓN

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|--------------|---------|------------|
| Pérdida de datos durante conversión | Baja | Alto | Backup completo antes de migrar |
| Queries en producción fallan | Media | Alto | Ejecutar en horario de baja demanda |
| Valores inconsistentes (NULL, otros) | Media | Medio | USING clauses con CASE para manejar todos los casos |
| Rollback necesario | Baja | Medio | Mantener backups y script de rollback |

### Script de Rollback

```sql
-- Restaurar desde backup
DROP TABLE account;
DROP TABLE account_favorite_locations;

CREATE TABLE account AS SELECT * FROM account_backup_20260118;
CREATE TABLE account_favorite_locations AS SELECT * FROM account_favorite_locations_backup_20260118;

-- Recrear índices y constraints
```

---

## 9. CONCLUSIONES

1. **Total de Columnas a Migrar:** 16
   - `account`: 4 columnas
   - `account_favorite_locations`: 12 columnas

2. **Tipo Objetivo:** SMALLINT NOT NULL
   - Compatible con MySQL TINYINT UNSIGNED
   - Soporta valores 0 y 1
   - Funciona con operadores numéricos en Go

3. **Impacto en Código:** CERO cambios requeridos en Go

4. **Tiempo Estimado:** 5-10 minutos (depende del volumen de datos)

5. **Recomendación:** Ejecutar INMEDIATAMENTE para evitar bugs futuros con queries mixtas.

---

## ANEXO: Script Completo de Migración

Ver archivo: `fix_boolean_columns_postgresql_final.sql`
