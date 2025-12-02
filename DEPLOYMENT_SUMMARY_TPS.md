# Deployment Summary - TPS Scraper + Bug Fixes
**Fecha:** 26 de Noviembre, 2025

---

## 🆕 Cambios en Este Deploy

### 1. ✅ TPS Scraper Implementado (NUEVO)
- **Archivo:** `internal/cronjobs/cjbot_creator/scrapers/tps.go`
- **Endpoint:** Toronto Police Service - Calls for Service API
- **Frecuencia:** Cada 1 hora
- **Categorías:** Traffic accidents, crime, medical emergencies

### 2. ✅ Bug Crítico Corregido - Subcategory Codes
- **Problema:** Bot usaba códigos inválidos (`crime`, `traffic_accident`, `building_fire`, etc.)
- **Solución:** Corregidos para coincidir con `Categories.tsx`
- **Archivos modificados:**
  - `internal/cronjobs/cjbot_creator/normalizer.go` (TPS + TFS mappings)
- **Base de datos:** Incidentes existentes corregidos con `sql_fixed_script.sql`

### 3. ✅ Descripciones Narrativas (Human-Readable)
- Antes: `"Call Type: PIACC | Division: D14 | Location: Yonge & Bloor"`
- Ahora: `"Officers responded to a personal injury collision at Yonge & Bloor, Toronto, ON."`
- Source attribution agregado: `"Source: Toronto Police Service"`

### 4. ✅ Scheduler Actualizado
- TFS: 15 min → **1 hora**
- TPS: **NUEVO - 1 hora** (agregado)
- Hydro: DESACTIVADO (mock data)

---

## 📋 Cambios por Archivo

| Archivo | Tipo | Descripción |
|---------|------|-------------|
| `scrapers/tps.go` | NUEVO | Scraper completo de TPS |
| `normalizer.go` | MODIFICADO | Mappings TPS/TFS corregidos |
| `service.go` | MODIFICADO | Método `RunTPS()` implementado |
| `scheduler.go` | MODIFICADO | Cronjob TPS agregado (1h) |
| `sql_fixed_script.sql` | NUEVO | Script de corrección DB |

---

## 🔧 Códigos de Subcategoría Corregidos

### Crime
- ❌ `"crime"` → ✅ `"theft"`, `"assault"`, `"robbery"`

### Traffic Accident
- ❌ `"traffic_accident"` → ✅ `"vehicle_collision"`
- ✅ `"pedestrian_nvolvement"` (typo del frontend, mantenido)
- ✅ `"hit_and_run"`

### Fire Incident
- ❌ `"building_fire"` → ✅ `"residential_fire"`
- ❌ `"fire_incident"` → ✅ `"other_fire_incident"`
- ✅ `"vehicle_fire"`

### Medical Emergency
- ❌ `"medical_emergency"` → ✅ `"other_medical_emergency"`
- ✅ `"cardiac_arrest"`, `"stroke"`, `"trauma_Injury"`, `"overdose_poisoning"`

### Infrastructure
- ❌ `"utility_issues"` → ✅ `"public_utility_issues"`

---

## 🚀 Pasos de Deployment

### 1. Build Docker Image
```bash
# Desde /backend
docker build -t alertly-backend:tps-fix .
```

### 2. Tag y Push a ECR
```bash
# Login to ECR
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 905418451690.dkr.ecr.us-west-2.amazonaws.com

# Tag
docker tag alertly-backend:tps-fix 905418451690.dkr.ecr.us-west-2.amazonaws.com/alertly-backend:latest

# Push
docker push 905418451690.dkr.ecr.us-west-2.amazonaws.com/alertly-backend:latest
```

### 3. Deploy en EC2
```bash
# SSH a EC2
ssh -i alertly-debug.pem ec2-user@44.243.7.9

# Stop container actual
sudo docker stop alertly-api
sudo docker rm alertly-api

# Pull nueva imagen
sudo docker pull 905418451690.dkr.ecr.us-west-2.amazonaws.com/alertly-backend:latest

# Run nuevo container
sudo docker run -d \
  --name alertly-api \
  -p 80:8080 \
  -e DB_USER=adminalertly \
  -e DB_PASS='Po1Ng2O3;' \
  -e DB_HOST=alertly-mysql-freetier.c3qmq4y86s84.us-west-2.rds.amazonaws.com \
  -e DB_PORT=3306 \
  -e DB_NAME=alertly \
  -e SERVER_PORT=8080 \
  -e IP_SERVER=0.0.0.0 \
  -e JWT_SECRET=your_jwt_secret \
  -e IMAGE_BASE_URL=https://api.alertly.ca \
  --restart unless-stopped \
  905418451690.dkr.ecr.us-west-2.amazonaws.com/alertly-backend:latest
```

### 4. Verificar Deployment
```bash
# Check health
curl https://api.alertly.ca/health | jq '.'

# Ver logs
sudo docker logs alertly-api --tail 50 -f

# Buscar logs de TPS
sudo docker logs alertly-api 2>&1 | grep -i tps
```

---

## ✅ Verificaciones Post-Deploy

### 1. Verificar Cronjobs Iniciados
```bash
# Deberías ver:
# 🚀 Goroutine for bot_creator_tps cronjob started
# ✅ Cronjob 'bot_creator_tps' scheduled every 1 hour
# 🔥 About to run bot_creator_tps cronjob for the first time...
```

### 2. Verificar Primera Ejecución TPS (esperar ~2-3 min)
```bash
sudo docker logs alertly-api 2>&1 | grep "TPS"

# Deberías ver:
# 🚔 [TPS] Starting bot creator job...
# ✅ [TPS] Job completed in 2.3s. Processed X/Y incidents
```

### 3. Verificar en Database
```sql
-- Ver incidentes TPS creados
SELECT
    incl_id,
    subcategory_code,
    category_code,
    description,
    address,
    created_at
FROM incident_clusters
WHERE account_id = 1
  AND category_code IN ('crime', 'traffic_accident', 'medical_emergency')
ORDER BY created_at DESC
LIMIT 10;

-- Verificar que NO hay códigos inválidos
SELECT subcategory_code, COUNT(*) as total
FROM incident_clusters
WHERE account_id = 1
  AND subcategory_code NOT IN (
    'theft', 'robbery', 'assault', 'homicide', 'fraud',
    'vehicle_collision', 'pedestrian_nvolvement', 'hit_and_run',
    'cardiac_arrest', 'stroke', 'trauma_Injury', 'overdose_poisoning', 'other_medical_emergency',
    'residential_fire', 'wildfire', 'vehicle_fire', 'other_fire_incident'
  )
GROUP BY subcategory_code;
-- Debería devolver 0 filas
```

### 4. Verificar en Frontend (App)
- Abrir mapa de Toronto
- Buscar incidentes nuevos del bot
- Verificar que NO hay error `icon_uri undefined`
- Confirmar que los íconos se muestran correctamente

---

## 🎯 Comportamiento Esperado

### Primera Hora
1. **TFS:** Ejecuta inmediatamente + cada hora
2. **TPS:** Ejecuta inmediatamente + cada hora
3. **Hydro:** Desactivado

### TPS API
- Endpoint: TPS Calls for Service (Public, No Geographic Offense)
- Retorna últimas 24h de llamadas
- Call types: PIACC, PDACC, ASSJU, ROB, THEPR, etc.
- Límite: 100 registros por query

### Deduplicación
- Hash: `SHA256(source + external_id + timestamp)`
- Logs: `"⏭️ [tps] Skipping duplicate incident: TPS-12345"`

---

## 📊 Métricas a Monitorear

### Logs Importantes
```bash
# Éxito
✅ [TPS] Job completed in X.Xs. Processed N/M incidents

# Deduplicación funcionando
⏭️ [TPS] Skipping duplicate incident: TPS-XXXXX

# Errores posibles
❌ [TPS] Scraping failed: [error]
⚠️ [TPS] Failed to normalize TPS-XXXXX: [error]
```

### CloudWatch (Opcional)
- Crear alarma si TPS falla > 3 veces consecutivas
- Monitorear tiempo de ejecución (debería ser < 10s)

---

## 🐛 Troubleshooting

### Error: "subcategory not found"
- **Causa:** Código de subcategoría no existe en DB
- **Solución:** Verificar mappings en `normalizer.go`

### Error: "icon_uri undefined" en Frontend
- **Causa:** Subcategory code inválido
- **Solución:** Ejecutar `sql_fixed_script.sql` nuevamente

### TPS no ejecuta
- **Verificar:** Logs de scheduler
- **Comando:** `sudo docker logs alertly-api 2>&1 | grep "bot_creator_tps"`

### Sin incidentes de TPS
- **Posible:** No hay calls for service activos en Toronto
- **Normal:** TPS API puede estar vacía fuera de horas pico
- **Verificar:** Ejecutar query manual a TPS API

---

## 🎉 Resultado Esperado

Después del deploy exitoso:

1. ✅ TPS scraper corriendo cada hora
2. ✅ Incidentes de TPS aparecen en mapa
3. ✅ Frontend muestra íconos correctamente (sin error icon_uri)
4. ✅ Descripciones narrativas y legibles
5. ✅ Source attribution: "Source: Toronto Police Service"
6. ✅ Deduplicación funcionando (no duplicados)

---

**🚀 LISTO PARA DEPLOY**
