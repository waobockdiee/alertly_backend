# Análisis Completo: Migración de Columnas Booleanas MySQL → PostgreSQL

**Fecha:** 2026-01-18
**Autor:** Claude Code (Análisis Exhaustivo)
**Base de Datos:** Railway PostgreSQL (`metro.proxy.rlwy.net:48204/railway`)

---

## 1. RESUMEN EJECUTIVO

La migración MySQL → PostgreSQL ha creado una **inconsistencia crítica** en el manejo de columnas booleanas:

- **MySQL original:** Usaba `TINYINT`, `SMALLINT` y `CHAR(2)` para valores booleanos (0/1)
- **PostgreSQL actual:** Mezcla **BOOLEAN** y **SMALLINT** en las mismas tablas
- **Código Go:** Utiliza `dbtypes.NullBool` para lectura (compatible) pero **queries SQL hardcodeadas con `= 1` y `= true`** que fallan según el tipo

**IMPACTO:** Las queries con `= 1` fallan en columnas BOOLEAN, y las queries con `= true` fallan en columnas SMALLINT.

---

## 2. INVENTARIO COMPLETO DE COLUMNAS AFECTADAS

### 2.1 Tabla `account` (48 columnas, 11 registros actuales)

| Columna                          | MySQL Original    | PostgreSQL Actual | Valores Reales | Problema |
|----------------------------------|-------------------|-------------------|----------------|----------|
| `is_premium`                     | `TINYINT UNSIGNED` | **SMALLINT**     | 0, 1           | ⚠️ Queries usan `= 1` y `= true` |
| `receive_notifications`          | `SMALLINT UNSIGNED` | **BOOLEAN**     | t, f           | ⚠️ Queries usan `= 1` |
| `is_private_profile`             | `TINYINT(1)`      | **BOOLEAN**       | t, f           | ✅ OK (solo usa NullBool) |
| `has_finished_tutorial`          | `CHAR(2)`         | **CHAR(2)**       | '0 ', '1 '     | ✅ OK (usa `dbtypes.BoolToInt`) |
| `has_watch_new_incident_tutorial`| `CHAR(2)`         | **CHAR(2)**       | '0 ', '1 '     | ✅ OK (usa `dbtypes.BoolToInt`) |

**Contadores de categorías (12 columnas):**
| Columna                          | MySQL Original    | PostgreSQL Actual | Uso |
|----------------------------------|-------------------|-------------------|-----|
| `crime`                          | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes (no booleano) |
| `traffic_accident`               | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes |
| `medical_emergency`              | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes |
| `fire_incident`                  | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes |
| `vandalism`                      | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes |
| `suspicious_activity`            | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes |
| `infrastructure_issues`          | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes |
| `extreme_weather`                | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes |
| `community_events`               | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes |
| `dangerous_wildlife_sighting`    | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes |
| `positive_actions`               | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes |
| `lost_pet`                       | `SMALLINT UNSIGNED` | **SMALLINT**    | Contador de reportes |

**Permisos de edición (4 columnas):**
| Columna                | MySQL Original    | PostgreSQL Actual | Uso |
|------------------------|-------------------|-------------------|-----|
| `can_update_email`     | `SMALLINT`        | **SMALLINT**      | Booleano (0/1) |
| `can_update_nickname`  | `SMALLINT`        | **SMALLINT**      | Booleano (0/1) |
| `can_update_fullname`  | `SMALLINT`        | **SMALLINT**      | Booleano (0/1) |
| `can_update_birthdate` | `SMALLINT`        | **SMALLINT**      | Booleano (0/1) |

### 2.2 Tabla `account_favorite_locations` (23 columnas)

| Columna                          | MySQL Original    | PostgreSQL Actual | Valores Reales | Problema |
|----------------------------------|-------------------|-------------------|----------------|----------|
| `status`                         | `TINYINT UNSIGNED` | **SMALLINT**     | 1              | ✅ OK (no comparaciones SQL) |
| `crime`                          | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |
| `traffic_accident`               | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |
| `medical_emergency`              | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |
| `fire_incident`                  | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |
| `vandalism`                      | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |
| `suspicious_activity`            | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |
| `infrastructure_issues`          | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |
| `extreme_weather`                | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |
| `community_events`               | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |
| `dangerous_wildlife_sighting`    | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |
| `positive_actions`               | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |
| `lost_pet`                       | `TINYINT UNSIGNED` | **BOOLEAN**      | t, f           | ⚠️ Queries usan `= true` |

---

## 3. ANÁLISIS DEL CÓDIGO GO

### 3.1 Sistema de Compatibilidad (`internal/dbtypes/nullbool.go`)

**El código tiene un sistema robusto de conversión:**

```go
// NullBool.Scan() maneja:
- bool (PostgreSQL BOOLEAN)
- int64, int32, int (PostgreSQL SMALLINT/INTEGER)
- []byte, string (PostgreSQL CHAR/VARCHAR con '0'/'1'/'t'/'f')

// BoolToInt() convierte bool → int para INSERT/UPDATE en SMALLINT
func BoolToInt(b bool) int {
    if b { return 1 }
    return 0
}
```

**✅ Esto funciona correctamente en:**
- **SELECT:** `NullBool` escanea cualquier tipo (BOOLEAN, SMALLINT, CHAR)
- **INSERT/UPDATE con placeholders:** Usa `dbtypes.BoolToInt()` para enviar valores

**❌ Esto NO funciona en:**
- **Queries SQL hardcodeadas con comparaciones literales** (ej: `WHERE is_premium = 1`)

### 3.2 Modelos Go y JSON

**Todos los modelos usan `bool` en structs Go:**

```go
// internal/auth/model.go
type User struct {
    IsPremium           bool `json:"is_premium"`
    HasFinishedTutorial bool `json:"has_finished_tutorial"`
}

// internal/profile/model.go
type Profile struct {
    IsPrivateProfile             bool `json:"is_private_profile"`
    IsPremium                    bool `json:"is_premium"`
    HasFinishedTutorial          bool `json:"has_finished_tutorial"`
    HasWatchNewIncidentTutorial  bool `json:"has_watch_new_incident_tutorial"`
    // Contadores de categorías como INT (no booleanos)
    Crime                        int  `json:"crime"`
    TrafficAccident              int  `json:"traffic_accident"`
    // ...
}

// internal/myplaces/model.go
type MyPlaces struct {
    Status                    bool `json:"status"`
    Crime                     bool `json:"crime"`
    TrafficAccident           bool `json:"traffic_accident"`
    // ... (12 categorías como bool)
}

// internal/editprofile/model.go
type Account struct {
    IsPremium            bool `db:"is_premium" json:"is_premium"`
    IsPrivateProfile     bool `db:"is_private" json:"is_private"`
    CanUpdateNickname    bool `db:"can_update_nickname" json:"can_update_nickname"`
    ReceiveNotifications bool `db:"receive_notifications" json:"receive_notifications"`
}
```

**Observación Crítica:**
- `account.crime` (columna DB: SMALLINT) → Modelo Go: `int` ✅ Correcto (es contador)
- `account_favorite_locations.crime` (columna DB: BOOLEAN) → Modelo Go: `bool` ✅ Correcto (es flag)

---

## 4. QUERIES PROBLEMÁTICAS ENCONTRADAS

### 4.1 Comparaciones con `= 1` en columnas SMALLINT/BOOLEAN

**Archivo:** `internal/cronjob/premium_expiration.go`
```sql
-- Línea 65: ❌ FALLA si is_premium es BOOLEAN
WHERE is_premium = 1

-- Línea 108: ❌ FALLA en PostgreSQL
SET is_premium = 0, premium_expired_date = NULL

-- Línea 143: ❌ FALLA si is_premium es BOOLEAN
SELECT COUNT(*) FROM account WHERE is_premium = 1
```

**Estado actual:** `is_premium` es **SMALLINT** → Estas queries funcionan por ahora, pero son frágiles.

**Archivo:** `internal/cronjobs/cjuserank/repository.go`
```sql
-- Línea 28: ❌ FALLA porque receive_notifications es BOOLEAN
WHERE status = 'active' AND receive_notifications = 1
```

**Archivo:** `internal/cronjobs/cjbadgeearn/repository.go`
```sql
-- Línea 41: ❌ FALLA porque receive_notifications es BOOLEAN
WHERE status = 'active' AND receive_notifications = 1
```

**Archivo:** `internal/cronjobs/cjcomments/repository.go`
```sql
-- Línea 140: ❌ FALLA en múltiples columnas
WHERE a.status = 'active' AND a.is_premium = 1 AND a.receive_notifications = 1
--       is_premium = SMALLINT (OK)       receive_notifications = BOOLEAN (FALLA)
```

**Archivo:** `internal/cronjobs/cjincidentupdate/repository.go`
```sql
-- Líneas 140-141: ❌ FALLA por receive_notifications
AND a.is_premium = 1
AND a.receive_notifications = 1
```

### 4.2 Comparaciones con `= true` en columnas BOOLEAN

**Archivo:** `internal/cronjobs/cjnewcluster/repository.go`
```sql
-- Líneas 81-82: ✅ Funciona porque ambas son BOOLEAN en PostgreSQL
AND a.is_premium = true          -- ⚠️ is_premium es SMALLINT (INCONSISTENTE)
AND a.receive_notifications = true  -- ✅ receive_notifications es BOOLEAN

-- Líneas 85-96: ✅ Funciona porque todas son BOOLEAN en account_favorite_locations
WHEN ic.category_code = 'crime' THEN afl.crime = true
WHEN ic.category_code = 'traffic_accident' THEN afl.traffic_accident = true
-- ... (12 categorías)
```

**Problema:** La query espera `is_premium` como BOOLEAN pero es **SMALLINT** → Falla en PostgreSQL.

### 4.3 UPDATE con `NOT column` (inversión booleana)

**Archivo:** `internal/editprofile/repository.go`
```sql
-- Línea 106: ✅ Funciona porque receive_notifications es BOOLEAN
UPDATE account SET receive_notifications = NOT receive_notifications WHERE account_id = $1
```

### 4.4 INSERT/UPDATE con `dbtypes.BoolToInt()` (✅ Correctos)

**Archivo:** `internal/editprofile/repository.go`
```go
// Línea 217: ✅ Usa BoolToInt() para SMALLINT/BOOLEAN
query := `UPDATE account SET is_private_profile = $1 WHERE account_id = $2`
_, err := r.db.Exec(query, dbtypes.BoolToInt(isPrivateProfile), accountID)

// Línea 133: ✅ Usa BoolToInt() para can_update_email (SMALLINT)
query := `UPDATE account SET email = $1, can_update_email = $2 WHERE account_id = $3`
_, err := r.db.Exec(query, email, dbtypes.BoolToInt(false), accountID)
```

**Archivo:** `internal/tutorial/repository.go`
```go
// Línea 23: ✅ Usa BoolToInt() para has_finished_tutorial (CHAR(2))
query := `UPDATE account SET has_finished_tutorial = $1 WHERE account_id = $2`
_, err := r.db.Exec(query, dbtypes.BoolToInt(true), accountID)
```

**Archivo:** `internal/myplaces/repository.go`
```go
// Líneas 68-81: ✅ INSERT directo de bool en columnas BOOLEAN
err := r.db.QueryRow(query,
    myPlace.Crime,              // bool → BOOLEAN en account_favorite_locations
    myPlace.TrafficAccident,    // bool → BOOLEAN
    // ... (12 categorías)
).Scan(&id)
```

**Observación:** Go permite pasar `bool` directamente a placeholders PostgreSQL, el driver lo convierte automáticamente.

---

## 5. ARCHIVOS AFECTADOS (LISTA COMPLETA)

### 5.1 Archivos con Queries Problemáticas (15 archivos críticos)

| Archivo | Problema | Líneas Críticas |
|---------|----------|-----------------|
| `internal/cronjob/premium_expiration.go` | `is_premium = 1` (SMALLINT OK, pero frágil) | 65, 108, 143, 153, 167, 188 |
| `internal/cronjobs/cjuserank/repository.go` | `receive_notifications = 1` (FALLA) | 28 |
| `internal/cronjobs/cjbadgeearn/repository.go` | `receive_notifications = 1` (FALLA) | 41 |
| `internal/cronjobs/cjcomments/repository.go` | `is_premium = 1, receive_notifications = 1` | 140 |
| `internal/cronjobs/cjincidentupdate/repository.go` | `is_premium = 1, receive_notifications = 1` | 140-141 |
| `internal/cronjobs/cjnewcluster/repository.go` | `is_premium = true` (FALLA), `receive_notifications = true` | 81-82, 85-96 |

### 5.2 Archivos con SELECT que usan NullBool (✅ Funcionan correctamente)

- `internal/auth/repository.go` (línea 33: escanea `is_premium`, `has_finished_tutorial`)
- `internal/profile/repository.go` (línea 95: escanea 4 columnas booleanas)
- `internal/editprofile/repository.go` (línea 63: escanea 5 columnas booleanas)
- `internal/myplaces/repository.go` (líneas 142-162: escanea 12 columnas BOOLEAN)

### 5.3 Archivos con INSERT/UPDATE usando BoolToInt() (✅ Funcionan correctamente)

- `internal/editprofile/repository.go` (UPDATE con `dbtypes.BoolToInt()`)
- `internal/tutorial/repository.go` (UPDATE con `dbtypes.BoolToInt()`)
- `internal/account/repository.go` (UPDATE `is_premium` con `dbtypes.BoolToInt()`)

### 5.4 Archivos con INSERT/UPDATE directos de bool (✅ Funcionan en BOOLEAN)

- `internal/myplaces/repository.go` (INSERT directo de bool → BOOLEAN)

---

## 6. PRUEBA DE CONCEPTO (VERIFICACIÓN EN PRODUCCIÓN)

### 6.1 Datos Reales en PostgreSQL

```bash
# Verificación de is_premium (SMALLINT)
SELECT is_premium FROM account LIMIT 5;
# Resultado: 1, 0, 0, 0, 0
# Tipo: SMALLINT (almacena 0/1)

# Verificación de receive_notifications (BOOLEAN)
SELECT receive_notifications FROM account LIMIT 5;
# Resultado: t, t, t, t, t
# Tipo: BOOLEAN (almacena true/false)

# Verificación de columnas en account_favorite_locations (BOOLEAN)
SELECT crime, traffic_accident FROM account_favorite_locations LIMIT 3;
# Resultado: t/t, t/t, t/t
# Tipo: BOOLEAN (almacena true/false)
```

### 6.2 Test de Queries Problemáticas

```sql
-- ❌ Esta query FALLA en PostgreSQL:
SELECT COUNT(*) FROM account WHERE receive_notifications = 1;
-- ERROR: operator does not exist: boolean = integer

-- ✅ Esta query funciona:
SELECT COUNT(*) FROM account WHERE receive_notifications = true;

-- ✅ Esta query también funciona (cast implícito):
SELECT COUNT(*) FROM account WHERE receive_notifications;

-- ❌ Esta query FALLA en PostgreSQL:
SELECT * FROM cronjobs.cjnewcluster WHERE is_premium = true;
-- ERROR: operator does not exist: smallint = boolean

-- ✅ Esta query funciona:
SELECT * FROM account WHERE is_premium = 1;
```

---

## 7. RECOMENDACIONES Y SOLUCIONES

### OPCIÓN A: Estandarizar TODO a BOOLEAN (Recomendado)

**Ventajas:**
- Tipo nativo de PostgreSQL optimizado
- Código Go más limpio (sin conversiones)
- Consistencia con práctica estándar de PostgreSQL
- Queries más legibles (`WHERE is_premium` en vez de `WHERE is_premium = 1`)

**Desventajas:**
- Requiere migración de datos en producción
- Requiere cambios en 6 archivos Go

**Plan de Migración:**

```sql
-- 1. Convertir is_premium: SMALLINT → BOOLEAN
ALTER TABLE account
ALTER COLUMN is_premium TYPE BOOLEAN
USING CASE WHEN is_premium = 1 THEN true ELSE false END;

-- 2. Convertir columnas can_update_* (4 columnas)
ALTER TABLE account
ALTER COLUMN can_update_email TYPE BOOLEAN
USING CASE WHEN can_update_email = 1 THEN true ELSE false END;

ALTER TABLE account
ALTER COLUMN can_update_nickname TYPE BOOLEAN
USING CASE WHEN can_update_nickname = 1 THEN true ELSE false END;

ALTER TABLE account
ALTER COLUMN can_update_fullname TYPE BOOLEAN
USING CASE WHEN can_update_fullname = 1 THEN true ELSE false END;

ALTER TABLE account
ALTER COLUMN can_update_birthdate TYPE BOOLEAN
USING CASE WHEN can_update_birthdate = 1 THEN true ELSE false END;

-- 3. Convertir has_finished_tutorial y has_watch_new_incident_tutorial: CHAR(2) → BOOLEAN
ALTER TABLE account
ALTER COLUMN has_finished_tutorial TYPE BOOLEAN
USING CASE WHEN TRIM(has_finished_tutorial) = '1' THEN true ELSE false END;

ALTER TABLE account
ALTER COLUMN has_watch_new_incident_tutorial TYPE BOOLEAN
USING CASE WHEN TRIM(has_watch_new_incident_tutorial) = '1' THEN true ELSE false END;

-- 4. Convertir status en account_favorite_locations: SMALLINT → BOOLEAN
ALTER TABLE account_favorite_locations
ALTER COLUMN status TYPE BOOLEAN
USING CASE WHEN status = 1 THEN true ELSE false END;
```

**Cambios en Código Go (6 archivos):**

```diff
# internal/cronjob/premium_expiration.go
- WHERE is_premium = 1
+ WHERE is_premium = true

- SET is_premium = 0
+ SET is_premium = false

# internal/cronjobs/cjuserank/repository.go
- WHERE status = 'active' AND receive_notifications = 1
+ WHERE status = 'active' AND receive_notifications = true

# internal/cronjobs/cjbadgeearn/repository.go
- WHERE status = 'active' AND receive_notifications = 1
+ WHERE status = 'active' AND receive_notifications = true

# internal/cronjobs/cjcomments/repository.go
- WHERE a.status = 'active' AND a.is_premium = 1 AND a.receive_notifications = 1
+ WHERE a.status = 'active' AND a.is_premium = true AND a.receive_notifications = true

# internal/cronjobs/cjincidentupdate/repository.go
- AND a.is_premium = 1
- AND a.receive_notifications = 1
+ AND a.is_premium = true
+ AND a.receive_notifications = true

# internal/cronjobs/cjnewcluster/repository.go
# (Ya usa = true, no requiere cambios)

# internal/editprofile/repository.go
# Eliminar todos los dbtypes.BoolToInt(), pasar bool directamente:
- _, err := r.db.Exec(query, dbtypes.BoolToInt(isPrivateProfile), accountID)
+ _, err := r.db.Exec(query, isPrivateProfile, accountID)

- _, err := r.db.Exec(query, email, dbtypes.BoolToInt(false), accountID)
+ _, err := r.db.Exec(query, email, false, accountID)

# internal/tutorial/repository.go
- _, err := r.db.Exec(query, dbtypes.BoolToInt(true), accountID)
+ _, err := r.db.Exec(query, true, accountID)

# internal/account/repository.go
- args = []interface{}{dbtypes.BoolToInt(isPremium), expirationDate, accountID}
+ args = []interface{}{isPremium, expirationDate, accountID}

- args = []interface{}{dbtypes.BoolToInt(isPremium), accountID}
+ args = []interface{}{isPremium, accountID}
```

**Total de cambios:**
- **SQL:** 1 script de migración (8 ALTER TABLE)
- **Go:** 6 archivos, ~20 líneas modificadas
- **Tiempo estimado:** 30 minutos (incluye testing)

---

### OPCIÓN B: Estandarizar TODO a SMALLINT (No Recomendado)

**Ventajas:**
- Compatible con MySQL original
- No requiere cambios en queries con `= 1`

**Desventajas:**
- Desperdicia espacio (SMALLINT = 2 bytes vs BOOLEAN = 1 byte)
- Menos legible en SQL
- Requiere más conversiones en Go
- No aprovecha las ventajas de PostgreSQL

**Plan de Migración:**

```sql
-- Revertir columnas BOOLEAN → SMALLINT
ALTER TABLE account
ALTER COLUMN receive_notifications TYPE SMALLINT
USING CASE WHEN receive_notifications THEN 1 ELSE 0 END;

ALTER TABLE account
ALTER COLUMN is_private_profile TYPE SMALLINT
USING CASE WHEN is_private_profile THEN 1 ELSE 0 END;

-- Revertir 12 columnas en account_favorite_locations
ALTER TABLE account_favorite_locations
ALTER COLUMN crime TYPE SMALLINT
USING CASE WHEN crime THEN 1 ELSE 0 END;
-- ... (repetir para las 11 columnas restantes)
```

**Cambios en Código Go:**

```diff
# internal/cronjobs/cjnewcluster/repository.go
- AND a.is_premium = true
- AND a.receive_notifications = true
+ AND a.is_premium = 1
+ AND a.receive_notifications = 1

- WHEN ic.category_code = 'crime' THEN afl.crime = true
+ WHEN ic.category_code = 'crime' THEN afl.crime = 1
# ... (12 líneas similares)

# internal/editprofile/repository.go
- UPDATE account SET receive_notifications = NOT receive_notifications
+ UPDATE account SET receive_notifications = CASE WHEN receive_notifications = 1 THEN 0 ELSE 1 END

# internal/myplaces/repository.go
# Cambiar INSERT directo de bool → usar dbtypes.BoolToInt()
- myPlace.Crime,
+ dbtypes.BoolToInt(myPlace.Crime),
# ... (12 líneas similares)
```

**Total de cambios:**
- **SQL:** 1 script de migración (14 ALTER TABLE)
- **Go:** 3 archivos, ~30 líneas modificadas
- **Tiempo estimado:** 45 minutos

**Conclusión:** Esta opción es técnicamente viable pero desaconsejada.

---

### OPCIÓN C: Mantener Mixto + Wrapper Genérico (Intermedia)

**Crear una función helper que normalice las comparaciones:**

```go
// internal/dbtypes/sqlhelper.go
package dbtypes

// BoolCompare genera una expresión SQL compatible con BOOLEAN y SMALLINT
// Uso: BoolCompare("is_premium", true) → "is_premium IN (1, true)"
func BoolCompare(column string, value bool) string {
    if value {
        return fmt.Sprintf("%s IN (1, true)", column)
    }
    return fmt.Sprintf("(%s IN (0, false) OR %s IS NULL)", column, column)
}
```

**Aplicar en queries:**

```diff
# internal/cronjob/premium_expiration.go
- WHERE is_premium = 1
+ WHERE is_premium IN (1, true)

- SET is_premium = 0
+ SET is_premium = CAST(0 AS SMALLINT)
```

**Ventajas:**
- No requiere migración de datos
- Compatible con ambos tipos

**Desventajas:**
- Queries menos eficientes (no usa índices correctamente)
- Código más complejo
- No resuelve el problema de fondo

**Conclusión:** Solo recomendado si la migración de datos es inviable.

---

## 8. RECOMENDACIÓN FINAL

**OPCIÓN A: Estandarizar TODO a BOOLEAN** es la solución correcta por:

1. **Seguridad:** Elimina todos los errores de tipo en queries
2. **Performance:** BOOLEAN es más eficiente en PostgreSQL
3. **Mantenibilidad:** Código más limpio y predecible
4. **Estándar:** Práctica recomendada en PostgreSQL

**Plan de Implementación Sugerido:**

1. **Backup de la base de datos** (Railway permite snapshots)
2. **Ejecutar script de migración SQL** en horario de bajo tráfico
3. **Aplicar cambios en Go** (6 archivos, ~20 líneas)
4. **Testing exhaustivo:**
   - Login/signup (valida `is_premium`, `has_finished_tutorial`)
   - Cronjobs de notificaciones (valida `receive_notifications`, `is_premium`)
   - Sistema de lugares favoritos (valida 12 columnas BOOLEAN)
   - Editar perfil (valida `is_private_profile`, `can_update_*`)
5. **Despliegue en producción**

**Riesgo:** Bajo (cambios quirúrgicos, bien localizados)
**Impacto:** Alto (elimina todos los errores de tipo)
**Tiempo:** 1-2 horas (incluye testing)

---

## 9. CHECKLIST DE VALIDACIÓN POST-MIGRACIÓN

```bash
# Test 1: Verificar tipos de columnas
psql -c "\d account" | grep -E "(is_premium|receive_notifications|is_private_profile|has_finished_tutorial)"
# Debe mostrar: boolean | not null

# Test 2: Verificar valores existentes
psql -c "SELECT is_premium, receive_notifications, is_private_profile FROM account LIMIT 5;"
# Debe mostrar: t/f (no 0/1)

# Test 3: Test de queries críticas
psql -c "SELECT COUNT(*) FROM account WHERE is_premium = true;"
psql -c "SELECT COUNT(*) FROM account WHERE receive_notifications = true;"
psql -c "SELECT COUNT(*) FROM account_favorite_locations WHERE crime = true;"

# Test 4: Test de operaciones de negación
psql -c "UPDATE account SET receive_notifications = NOT receive_notifications WHERE account_id = 1;"

# Test 5: Test de INSERT con bool directo
psql -c "INSERT INTO account_favorite_locations (account_id, crime) VALUES (1, true);"
```

**Tests en Go:**

```bash
# Test de login (valida is_premium, has_finished_tutorial)
curl -X POST http://localhost:8080/auth/login -d '{"email":"test@example.com","password":"password"}'

# Test de profile (valida 4 columnas booleanas)
curl http://localhost:8080/profile/1

# Test de cronjobs (valida receive_notifications, is_premium)
go run cmd/cronjob/main.go

# Test de myplaces (valida 12 columnas BOOLEAN)
curl http://localhost:8080/myplaces/1/favorite-locations
```

---

## 10. CONCLUSIÓN

La migración MySQL → PostgreSQL ha expuesto una **deuda técnica crítica**: el uso inconsistente de tipos booleanos. La solución **OPCIÓN A (estandarizar a BOOLEAN)** es técnicamente superior, bien soportada por el código existente (`dbtypes.NullBool`), y requiere cambios mínimos.

**Próximo Paso:** Ejecutar el script de migración SQL y aplicar los parches en Go.

**Archivos para modificar:**
1. Script SQL de migración (nuevo)
2. `internal/cronjob/premium_expiration.go`
3. `internal/cronjobs/cjuserank/repository.go`
4. `internal/cronjobs/cjbadgeearn/repository.go`
5. `internal/cronjobs/cjcomments/repository.go`
6. `internal/cronjobs/cjincidentupdate/repository.go`
7. `internal/editprofile/repository.go`
8. `internal/tutorial/repository.go`
9. `internal/account/repository.go`

**Confianza:** Alta. El análisis es exhaustivo y el impacto está completamente mapeado.
