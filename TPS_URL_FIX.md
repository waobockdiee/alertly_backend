# TPS API URL Fix - Segunda Corrección

**Fecha:** 27 de Noviembre, 2025

---

## 🐛 Problema Detectado

El scraper de TPS estaba retornando **400 Bad Request** porque el formato del WHERE clause era incorrecto.

### Error Original:
```go
// ❌ INCORRECTO (retorna 400):
whereClause := fmt.Sprintf("OCCURRENCE_TIME_AGOL>=timestamp '%s'", yesterdayStr)
```

ArcGIS FeatureServer **NO acepta** el formato `timestamp 'YYYY-MM-DD HH:MM:SS'`.

---

## ✅ Solución Aplicada

### Cambio en `tps.go` (línea 70):
```go
// ✅ CORRECTO (funciona):
whereClause := fmt.Sprintf("OCCURRENCE_TIME_AGOL>=date '%s'", yesterdayStr)
```

**Diferencia clave:** Cambiar `timestamp` → `date`

---

## 📊 Archivos Modificados

| Archivo | Cambio |
|---------|--------|
| `scrapers/tps.go` | Línea 70: `timestamp` → `date` |
| `normalizer.go` | Agregados mappings para nuevos TPS call types |

---

## 🆕 Nuevos Call Types Agregados a TPS Mappings

### Tráfico:
- `IMPDR` (Impaired Driver) → `traffic_accident` / `vehicle_collision`
- `TRAOB` (Traffic Obstruct) → `traffic_accident` / `vehicle_collision`

### Crimen:
- `ASSPR` (Assault in Progress) → `crime` / `assault`
- `PERGU` (Person with Gun) → `crime` / `assault` (prioridad alta)
- `ATTBR` (Attempt Break & Enter) → `crime` / `robbery`
- `BREEN` (Break & Enter) → `crime` / `robbery`
- `FRA` (Fraud) → `crime` / `fraud`
- `THE` (Theft) → `crime` / `theft`

### Infraestructura:
- `HAZ` (Hazard) → `infrastructure_issues` / `public_utility_issues`

### Fire:
- `FIR` (Fire) → `fire_incident` / `other_fire_incident`
- `SEEFI` (See Fire Dept) → `fire_incident` / `other_fire_incident`

### Otros:
- `ANICO` (Animal Complaint) → `dangerous_wildlife_sighting` / `other_wildlife`
- `TAXAL` (Taxi Alarm) → `suspicious_activity` / `unusual_behavior`
- `DAMJU` (Damage Just Occurred) → `vandalism` / `public_property_damage`

---

## 🧪 Test de Verificación

```bash
# URL que ahora funciona correctamente:
curl 'https://services.arcgis.com/S9th0jAJ7bqgIRjw/arcgis/rest/services/C4S_Public_NoGO/FeatureServer/0/query?f=json&resultOffset=0&resultRecordCount=5&where=OCCURRENCE_TIME_AGOL%3E%3Ddate%20%272025-11-27%2020%3A00%3A00%27&orderByFields=OCCURRENCE_TIME_AGOL%20DESC&outFields=*&outSR=4326'
```

**Resultado esperado:** ✅ 200 OK con datos de incidentes

---

## 🚀 Deployment

### 1. Build
```bash
cd backend
docker build -t alertly-backend:tps-url-fix .
```

### 2. Tag y Push a ECR
```bash
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 905418451690.dkr.ecr.us-west-2.amazonaws.com

docker tag alertly-backend:tps-url-fix 905418451690.dkr.ecr.us-west-2.amazonaws.com/alertly-backend:latest

docker push 905418451690.dkr.ecr.us-west-2.amazonaws.com/alertly-backend:latest
```

### 3. Deploy en EC2
```bash
ssh -i alertly-debug.pem ec2-user@44.243.7.9

sudo docker stop alertly-api
sudo docker rm alertly-api
sudo docker pull 905418451690.dkr.ecr.us-west-2.amazonaws.com/alertly-backend:latest

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

### 4. Verificar
```bash
# Ver logs TPS
sudo docker logs alertly-api 2>&1 | grep -i tps | tail -20

# Deberías ver:
# ✅ TPS scraper found X incidents
# (en vez de "❌ TPS API returned status 400")
```

---

## 🎯 Resultado Esperado

Después del deploy:
1. ✅ TPS API retorna 200 OK (no más 400 Bad Request)
2. ✅ Incidentes de TPS se procesan correctamente
3. ✅ Todos los call types están mapeados a subcategorías válidas
4. ✅ No más errores "subcategory not found" para TPS

---

## 📝 Análisis de Causa Raíz

**Por qué el error 400:**
- ArcGIS REST API de TPS NO soporta el keyword `timestamp` en WHERE clauses
- Solo acepta `date` para comparaciones de fechas
- Documentación oficial de ArcGIS: https://developers.arcgis.com/rest/services-reference/query-feature-service-layer-.htm

**Lección aprendida:**
- Siempre probar URLs de API manualmente antes de implementar
- Revisar documentación oficial de APIs externas
- No asumir formatos de fecha sin verificar

---

✅ **LISTO PARA REDEPLOY**
