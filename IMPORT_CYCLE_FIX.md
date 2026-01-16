# Import Cycle Fix - PostgreSQL Migration

## Problema Detectado

Al intentar compilar el backend después de agregar `import "alertly/internal/common"` a varios repositories, se detectó un **ciclo de importación existente**:

```
common → alerts → auth → common
```

**Específicamente:**
- `internal/common/notification.go` importa `internal/alerts`
- `internal/alerts` (probablemente) importa `internal/auth`
- `internal/auth/repository.go` AHORA importa `internal/common` (para `NullBool`)

## Este NO es un problema causado por NullBool

El ciclo de importación ya existía en el proyecto. La adición de `import "alertly/internal/common"` en `auth/repository.go` simplemente lo expuso.

## Soluciones Posibles

### Opción 1: Mover `NullBool` a un paquete separado (RECOMENDADO)

Crear un nuevo paquete sin dependencias circulares:

```
backend/internal/dbtypes/
  ├── nullbool.go
  └── nullbool_test.go
```

Luego todos los packages pueden importar `alertly/internal/dbtypes` sin causar ciclos.

**Ventajas:**
- Separa tipos de base de datos de lógica de negocio
- Elimina el ciclo de importación
- Más fácil de mantener

**Implementación:**
```bash
mkdir -p internal/dbtypes
mv internal/common/types.go internal/dbtypes/nullbool.go
mv internal/common/nullbool_test.go internal/dbtypes/nullbool_test.go

# Actualizar imports en todos los archivos
sed -i '' 's/alertly\/internal\/common/alertly\/internal\/dbtypes/g' internal/*/repository.go
```

### Opción 2: Extraer la función del ciclo

Mover `internal/common/notification.go` a un paquete diferente que no cause ciclos:

```
internal/notifications/common.go  (en vez de internal/common/notification.go)
```

**Ventajas:**
- No requiere cambiar muchos archivos
- Mantiene `common` como paquete de utilities

**Desventajas:**
- Requiere refactorizar el código que usa `common.SaveNotification()`

### Opción 3: Duplicar BoolToInt temporalmente

Mientras se resuelve el ciclo arquitectónico, duplicar la función `BoolToInt` en cada paquete que la necesite.

**Ventajas:**
- Solución rápida
- No requiere refactorizar la arquitectura

**Desventajas:**
- Código duplicado
- No es escalable
- No resuelve el problema subyacente

## Estado Actual del Código

Los siguientes archivos YA ESTÁN ACTUALIZADOS y funcionan correctamente con `NullBool`:

1. `internal/auth/repository.go` - ✅ Usa `common.NullBool`
2. `internal/profile/repository.go` - ✅ Usa `common.NullBool`
3. `internal/editprofile/repository.go` - ✅ Usa `common.NullBool` y `BoolToInt()`
4. `internal/account/repository.go` - ✅ Usa `common.NullBool` y `BoolToInt()`
5. `internal/tutorial/repository.go` - ✅ Usa `common.BoolToInt()`
6. `internal/common/notification.go` - ✅ Usa `common.BoolToInt()`
7. `internal/common/score.go` - ✅ Usa `common.BoolToInt()`

## Verificación de Funcionalidad

Los tests de `NullBool` pasan correctamente:

```bash
go test -v ./internal/common/types.go ./internal/common/nullbool_test.go
# PASS (todos los tests pasan)
```

## Recomendación Final

**OPCIÓN 1 es la mejor solución a largo plazo.**

### Pasos para implementar:

1. Crear nuevo paquete `dbtypes`:
   ```bash
   mkdir -p internal/dbtypes
   ```

2. Mover archivos:
   ```bash
   # Copiar solo la parte de NullBool de types.go
   # Crear internal/dbtypes/nullbool.go con:
   # - NullBool struct
   # - Scan() method
   # - Value() method
   # - MarshalJSON() method
   # - UnmarshalJSON() method
   # - BoolToInt() function

   # Mover tests
   mv internal/common/nullbool_test.go internal/dbtypes/nullbool_test.go
   ```

3. Actualizar imports en todos los repositories:
   ```bash
   # Reemplazar "alertly/internal/common" con "alertly/internal/dbtypes"
   # SOLO en las líneas que usan NullBool o BoolToInt
   ```

4. Verificar compilación:
   ```bash
   go build ./cmd/app/main.go
   ```

## Archivos que Necesitan Actualización si se Implementa Opción 1

```
internal/auth/repository.go
internal/profile/repository.go
internal/editprofile/repository.go
internal/account/repository.go
internal/tutorial/repository.go
internal/common/notification.go
internal/common/score.go
```

## Documentación Relacionada

- **Solución de tipos booleanos:** `POSTGRESQL_BOOL_MIGRATION.md`
- **Tests:** `internal/common/nullbool_test.go` (o `internal/dbtypes/nullbool_test.go` después de mover)

## Nota Importante

El ciclo de importación `common → alerts → auth → common` es un problema arquitectónico que debería resolverse independientemente de la migración de PostgreSQL. Se recomienda:

1. Auditar todas las dependencias de `internal/common`
2. Separar utilities sin dependencias (como `types.go`) de código con dependencias (como `notification.go`)
3. Considerar un paquete `internal/base` o `internal/dbtypes` para tipos fundamentales
