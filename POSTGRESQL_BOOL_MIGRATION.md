# PostgreSQL Boolean Type Migration Guide

Este documento describe la solución implementada para manejar tipos booleanos en PostgreSQL cuando se migra desde MySQL.

## Problema

Al migrar de MySQL a PostgreSQL, las columnas booleanas pueden tener tipos inconsistentes:

- **MySQL**: `TINYINT(1)` (valores: 0, 1)
- **PostgreSQL**: Puede ser `BOOLEAN`, `SMALLINT`, `CHAR(1)`, o strings con espacios (`"1 "`)

### Errores Comunes

```
sql: Scan error on column index 8, name "has_finished_tutorial":
sql/driver: couldn't convert "1 " into type bool

ERROR: invalid input syntax for type smallint: "true"
```

## Solución Implementada

### 1. Tipo Personalizado: `dbtypes.NullBool`

Se creó un tipo que implementa `sql.Scanner` y `driver.Valuer` para manejar todos los formatos posibles:

**Archivo:** `/Users/garyeikoow/Desktop/alertly/backend/internal/common/types.go`

```go
type NullBool struct {
    Bool  bool
    Valid bool // Valid es true si Bool no es NULL
}
```

**Características:**
- Escanea: `bool`, `int64`, `int32`, `int`, `string`, `[]byte`
- Maneja espacios en blanco: `"1 "`, `" t "`
- Soporta: `"1"`, `"0"`, `"true"`, `"false"`, `"t"`, `"f"`, `"TRUE"`, `"FALSE"`
- NULL-safe: Si el valor es NULL, `Valid = false`

### 2. Función Helper: `dbtypes.BoolToInt()`

Para INSERT/UPDATE en columnas `SMALLINT`:

```go
func BoolToInt(b bool) int {
    if b {
        return 1
    }
    return 0
}
```

## Uso en Repositories

### Escanear Resultados (SELECT)

**ANTES (Error):**
```go
var user User
err := row.Scan(&user.IsPremium, &user.HasFinishedTutorial)
// Error: can't convert "1 " to bool
```

**DESPUÉS (Correcto):**
```go
var isPremium, hasFinishedTutorial dbtypes.NullBool

err := row.Scan(&isPremium, &hasFinishedTutorial)

// Convertir a bool
user.IsPremium = isPremium.Valid && isPremium.Bool
user.HasFinishedTutorial = hasFinishedTutorial.Valid && hasFinishedTutorial.Bool
```

### Insertar/Actualizar (INSERT/UPDATE)

**ANTES (Error):**
```go
query := `UPDATE account SET has_finished_tutorial = 1 WHERE account_id = $1`
_, err := db.Exec(query, accountID)
// Error en PostgreSQL si la columna es SMALLINT y espera int, no string "1"
```

**DESPUÉS (Correcto):**
```go
import "alertly/internal/dbtypes"

query := `UPDATE account SET has_finished_tutorial = $1 WHERE account_id = $2`
_, err := db.Exec(query, dbtypes.BoolToInt(true), accountID)
```

**Otro Ejemplo:**
```go
query := `UPDATE account SET can_update_nickname = $1 WHERE account_id = $2 AND can_update_nickname = $3`
_, err := db.Exec(query, dbtypes.BoolToInt(false), accountID, dbtypes.BoolToInt(true))
```

## ⚠️ Evitando Ciclos de Importación

El tipo `NullBool` y la función `BoolToInt` están en el paquete `internal/dbtypes` (no en `internal/common`) para evitar ciclos de importación.

**Ciclo Anterior:**
```
common → alerts → auth → common (ERROR!)
```

**Solución:**
```
dbtypes (sin dependencias)
   ↑
   |-- auth
   |-- profile
   |-- editprofile
   └-- account
```

El paquete `dbtypes` solo importa paquetes estándar de Go, por lo que puede ser importado por cualquier otro paquete sin causar ciclos.

## Archivos Actualizados

Los siguientes archivos fueron modificados para usar `dbtypes.NullBool` y `dbtypes.BoolToInt()`:

1. `/Users/garyeikoow/Desktop/alertly/backend/internal/auth/repository.go`
   - `GetUserByEmail()`: Escanea `is_premium`, `has_finished_tutorial`

2. `/Users/garyeikoow/Desktop/alertly/backend/internal/profile/repository.go`
   - `GetById()`: Escanea múltiples campos booleanos

3. `/Users/garyeikoow/Desktop/alertly/backend/internal/editprofile/repository.go`
   - `GetAccountByID()`: Escanea campos booleanos
   - `UpdateEmail()`: Usa `BoolToInt()` para `can_update_email`
   - `UpdateNickname()`: Usa `BoolToInt()` para `can_update_nickname`
   - `UpdateFullName()`: Usa `BoolToInt()` para `can_update_fullname`
   - `UpdateBirthDate()`: Usa `BoolToInt()` para `can_update_birthdate`

4. `/Users/garyeikoow/Desktop/alertly/backend/internal/account/repository.go`
   - `GetMyInfo()`: Escanea `is_premium`, `has_finished_tutorial`
   - `SetHasFinishedTutorial()`: Usa `BoolToInt(true)`

5. `/Users/garyeikoow/Desktop/alertly/backend/internal/tutorial/repository.go`
   - `MarkTutorialAsFinished()`: Usa `BoolToInt(true)`

## Tests

Se crearon tests exhaustivos en `/Users/garyeikoow/Desktop/alertly/backend/internal/common/nullbool_test.go`:

```bash
go test -v ./internal/common/types.go ./internal/common/nullbool_test.go
```

**Resultados:**
- 18 casos de prueba para `Scan()`
- 3 casos de prueba para `Value()`
- 3 casos de prueba para `MarshalJSON()`
- 4 casos de prueba para `UnmarshalJSON()`
- 2 casos de prueba para `BoolToInt()`

**Todos los tests pasan correctamente.**

## Patrón de Migración

Si encuentras más errores de tipo booleano en otros archivos:

### Para SELECT (escanear resultados):

1. Importar `"alertly/internal/dbtypes"`
2. Declarar variable temporal:
   ```go
   var myBoolField dbtypes.NullBool
   ```
3. Escanear en la variable temporal:
   ```go
   err := row.Scan(&myBoolField)
   ```
4. Convertir a bool:
   ```go
   myStruct.MyBoolField = myBoolField.Valid && myBoolField.Bool
   ```

### Para INSERT/UPDATE (insertar valores):

1. Importar `"alertly/internal/dbtypes"`
2. Reemplazar valores literales (`0`, `1`, `true`, `false`) con placeholders (`$1`, `$2`, etc.)
3. Usar `dbtypes.BoolToInt()`:
   ```go
   _, err := db.Exec(query, dbtypes.BoolToInt(myBool), otherParams...)
   ```

## Búsqueda de Archivos Afectados

Para encontrar más archivos que puedan necesitar actualización:

```bash
# Buscar queries con campos booleanos en WHERE
grep -rn "is_premium\|has_finished_tutorial\|can_update\|receive_notifications" internal/ --include="*.go"

# Buscar queries con valores literales 0/1 en UPDATE/INSERT
grep -rn "UPDATE.*= [01]" internal/ --include="*.go"
grep -rn "INSERT.*VALUES.*[01]" internal/ --include="*.go"
```

## Columnas Booleanas Conocidas

Lista de columnas booleanas en la tabla `account` que pueden requerir este patrón:

- `is_premium`
- `has_finished_tutorial`
- `has_watch_new_incident_tutorial`
- `is_private_profile`
- `can_update_nickname`
- `can_update_fullname`
- `can_update_birthdate`
- `can_update_email`
- `receive_notifications`

## Ventajas de esta Solución

1. **Type-Safe**: Implementa interfaces estándar de Go (`sql.Scanner`, `driver.Valuer`)
2. **Reutilizable**: Un solo tipo para todo el proyecto
3. **NULL-Safe**: Maneja valores NULL correctamente
4. **Compatible**: Funciona con cualquier formato de PostgreSQL (BOOLEAN, SMALLINT, CHAR, VARCHAR)
5. **Testeado**: Tests completos garantizan el funcionamiento correcto
6. **Idiomático**: Sigue el patrón de `sql.NullString`, `sql.NullInt64`, etc.

## Referencias

- **Tipo implementado:** `internal/dbtypes/nullbool.go`
- **Tests:** `internal/dbtypes/nullbool_test.go`
- **Driver PostgreSQL:** `github.com/lib/pq`
- **Go SQL Package:** `database/sql`, `database/sql/driver`
- **Ciclos de importación:** Ver `IMPORT_CYCLE_FIX.md` para detalles

## Soporte

Si encuentras un nuevo caso de uso o error relacionado con booleanos:

1. Verifica que el archivo tenga `import "alertly/internal/dbtypes"`
2. Aplica el patrón descrito arriba
3. Ejecuta los tests para asegurar que funciona
4. Documenta cualquier nuevo caso de uso en este archivo
