# Quick Start: Usar PostgreSQL Railway en Desarrollo

**Actualizado:** 17 de Enero, 2026
**Estado:** ✅ Base de datos corregida y lista para uso

---

## 🚀 Configuración Rápida

### 1. Actualizar archivo `.env`

Reemplaza la configuración de MySQL con PostgreSQL:

```bash
# ANTES (MySQL AWS)
# DB_USER=adminalertly
# DB_PASS=your_password
# DB_HOST=alertly-mysql-freetier.c3qmq4y86s84.us-west-2.rds.amazonaws.com
# DB_PORT=3306
# DB_NAME=alertly

# AHORA (PostgreSQL Railway)
DATABASE_URL=postgres://postgres:cGA2dBF6G33BgfefcgDb1CDa6CagFcC5@metro.proxy.rlwy.net:48204/railway
```

**Nota:** El backend Go debe estar configurado para leer `DATABASE_URL` en lugar de las variables individuales (DB_USER, DB_PASS, etc.).

---

## 🔧 Verificar Conexión

### Opción 1: Usando psql (CLI)

```bash
psql "postgres://postgres:cGA2dBF6G33BgfefcgDb1CDa6CagFcC5@metro.proxy.rlwy.net:48204/railway"
```

**Comandos útiles:**
```sql
-- Listar todas las tablas
\dt

-- Ver estructura de una tabla
\d account

-- Contar registros
SELECT COUNT(*) FROM account;

-- Salir
\q
```

### Opción 2: Usando el backend Go

```bash
cd /Users/garyeikoow/Desktop/alertly/backend
go run cmd/app/main.go
```

Deberías ver en los logs:
```
✅ Connected to PostgreSQL database
```

---

## 🔍 Diferencias importantes vs MySQL

### 1. Tipos de datos

| MySQL | PostgreSQL | Código Go |
|-------|-----------|-----------|
| `TINYINT(1)` | `SMALLINT` o `BOOLEAN` | `int` o `bool` |
| `INT UNSIGNED` | `INTEGER` | `int` |
| `ENUM('a','b')` | `VARCHAR + CHECK` | `string` |
| `TIMESTAMP` | `timestamp without time zone` | `time.Time` |
| `DECIMAL(3,1)` | `NUMERIC(3,1)` | `float64` |

**IMPORTANTE:** El backend Go debe manejar estos tipos correctamente. Si usas `database/sql` o `pgx`, no deberías necesitar cambios.

### 2. Funciones de fecha/hora

| MySQL | PostgreSQL |
|-------|-----------|
| `NOW()` | `NOW()` ✅ (compatible) |
| `CURRENT_TIMESTAMP` | `CURRENT_TIMESTAMP` ✅ (compatible) |
| `DATE_SUB(NOW(), INTERVAL 24 HOUR)` | `NOW() - INTERVAL '24 hours'` ⚠️ (diferente sintaxis) |

### 3. Geolocation queries

**MySQL (con ST_Distance_Sphere):**
```sql
WHERE ST_Distance_Sphere(
    point(longitude, latitude),
    point(?, ?)
) <= ?
```

**PostgreSQL (con PostGIS):**
```sql
WHERE ST_Distance(
    ST_MakePoint(longitude, latitude)::geography,
    ST_MakePoint(?, ?)::geography
) <= ?
```

**IMPORTANTE:** Si tu backend usa funciones geoespaciales, necesitarás actualizar las queries para PostgreSQL + PostGIS.

---

## 📊 Queries Comunes Adaptadas

### Login
```sql
-- MySQL y PostgreSQL (compatible)
SELECT * FROM account
WHERE email = $1 AND status = 'active'
LIMIT 1;
```

### Obtener clusters por ubicación
```sql
-- PostgreSQL (usando PostGIS si está instalado)
SELECT * FROM incident_clusters
WHERE ST_Distance(
    ST_MakePoint(center_longitude, center_latitude)::geography,
    ST_MakePoint($1, $2)::geography
) <= $3
AND is_active = '1'
AND created_at >= NOW() - INTERVAL '24 hours'
ORDER BY created_at DESC
LIMIT 100;
```

**Nota:** Si PostGIS no está instalado en Railway, necesitarás usar cálculo de distancia manual o instalar la extensión.

### Verificar si PostGIS está instalado:
```sql
SELECT PostGIS_version();
```

Si no está instalado:
```sql
CREATE EXTENSION postgis;
```

---

## 🐛 Troubleshooting Común

### Error: "pq: SSL is not enabled on the server"

**Solución:** Agregar `?sslmode=disable` a la DATABASE_URL:
```bash
DATABASE_URL=postgres://postgres:***@metro.proxy.rlwy.net:48204/railway?sslmode=disable
```

### Error: "column does not exist"

**Causa:** Nombre de columna incorrecto o tabla no encontrada.
**Solución:** Verificar estructura con `\d nombre_tabla` en psql.

### Error: "null value in column violates not-null constraint"

**Causa:** Intentando insertar NULL en columna NOT NULL.
**Solución:** Verificar que todos los campos requeridos tengan valores.

### Error: "function st_distance_sphere does not exist"

**Causa:** Función de MySQL no existe en PostgreSQL.
**Solución:** Instalar PostGIS y usar `ST_Distance()` en su lugar.

---

## 🧪 Testing Rápido

### 1. Probar conexión
```bash
psql "postgres://postgres:cGA2dBF6G33BgfefcgDb1CDa6CagFcC5@metro.proxy.rlwy.net:48204/railway" -c "SELECT version();"
```

### 2. Probar query simple
```bash
psql "postgres://postgres:cGA2dBF6G33BgfefcgDb1CDa6CagFcC5@metro.proxy.rlwy.net:48204/railway" -c "SELECT COUNT(*) FROM account;"
```

### 3. Probar INSERT
```bash
psql "postgres://postgres:cGA2dBF6G33BgfefcgDb1CDa6CagFcC5@metro.proxy.rlwy.net:48204/railway" -c "
INSERT INTO account (email, password, nickname, role, status)
VALUES ('test_quick@example.com', 'hashed', 'QuickTest', 'citizen', 'pending_activation')
RETURNING account_id;
"
```

### 4. Limpiar test
```bash
psql "postgres://postgres:cGA2dBF6G33BgfefcgDb1CDa6CagFcC5@metro.proxy.rlwy.net:48204/railway" -c "
DELETE FROM account WHERE email = 'test_quick@example.com';
"
```

---

## 📝 Checklist Pre-Deploy

Antes de cambiar a PostgreSQL en producción, verifica:

- [ ] DATABASE_URL configurado en `.env`
- [ ] Backend se conecta exitosamente a Railway PostgreSQL
- [ ] Queries de autenticación funcionan
- [ ] Queries de incidents funcionan
- [ ] PostGIS instalado si usas funciones geoespaciales
- [ ] Queries geoespaciales adaptadas a PostgreSQL
- [ ] Cronjobs pueden conectarse y ejecutarse
- [ ] Testing end-to-end del frontend completo
- [ ] Sin errores en logs después de 1 hora de uso

---

## 🔐 Seguridad

**IMPORTANTE:** El connection string contiene credenciales. NO lo subas a Git.

**Para producción:**
1. Usa variables de entorno en el servidor
2. Considera usar Railway's built-in DATABASE_URL
3. Habilita SSL si es posible (Railway lo soporta)

---

## 📞 Soporte

**Problemas con la base de datos:**
- Verificar estado en Railway Dashboard
- Revisar logs de conexión en backend
- Consultar `/Users/garyeikoow/Desktop/alertly/backend/POSTGRESQL_MIGRATION_FIX_REPORT.md`

**Documentación oficial:**
- PostgreSQL: https://www.postgresql.org/docs/
- Railway: https://docs.railway.app/
- PostGIS: https://postgis.net/documentation/

---

## ✅ Resumen

- ✅ Base de datos corregida con todos los constraints NOT NULL
- ✅ 119 columnas verificadas
- ✅ 985 incident clusters + 2,177 incident reports migrados
- ✅ Todas las tablas del sistema de referidos funcionando
- ✅ Lista para testing del backend Go

**Próximo paso:** Actualizar `.env` y ejecutar `go run cmd/app/main.go`
