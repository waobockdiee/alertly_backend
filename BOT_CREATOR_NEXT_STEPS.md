# Bot Creator - Next Steps & Status

## ✅ Completado Hasta Ahora

### 1. Implementación de Código (100%)
- ✅ 2 scrapers funcionales (TFS Fire + Toronto Hydro)
- ✅ Sistema de geocoding con Nominatim + cache MySQL
- ✅ Normalización inteligente de categorías
- ✅ Deduplicación con SHA256 hashes
- ✅ Procesamiento concurrente con goroutines
- ✅ Integración con Lambda (5 casos en main.go)
- ✅ Código compila sin errores

### 2. Base de Datos (100%)
- ✅ Migration SQL creado: `backend/assets/db/bot_seeder_tables.sql`
- ✅ Tabla `bot_incident_hashes` (deduplicación)
- ✅ Tabla `geocoding_cache` (Nominatim cache)

**Para aplicar:**
```bash
mysql -u root -p alertly < backend/assets/db/bot_seeder_tables.sql
```

### 3. Imágenes Estáticas (67% - 8/12)
- ✅ Subidas a S3 correctamente
- ✅ Acceso público verificado (HTTP 200)
- ✅ Cache-Control configurado (1 año)

**URLs funcionando:**
```
https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/crime.webp
https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/traffic_accident.webp
https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/medical_emergency.webp
https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/fire_incident.webp
https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/suspicious_activity.webp
https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/extreme_weather.webp
https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/infrastructure_issues.webp
https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/dangerous_wildlife_sighting.webp
```

---

## ⏳ Tareas Pendientes

### 1. Crear Imágenes Faltantes (4/12)

Necesitas crear y subir estas imágenes (900x1200px, formato WebP):

**Faltantes:**
- ❌ `vandalism.webp`
- ❌ `community_events.webp`
- ❌ `positive_actions.webp`
- ❌ `lost_pet.webp`

**Para subir nuevas imágenes:**
```bash
cd /Users/garyeikoow/Documents/www/alertly/AWS/images

# Agregar las 4 imágenes faltantes aquí

# Subir con AWS CLI
aws s3 cp vandalism.webp s3://alertly-images-production/incidents/ --region us-west-2 --content-type "image/webp" --cache-control "max-age=31536000"
aws s3 cp community_events.webp s3://alertly-images-production/incidents/ --region us-west-2 --content-type "image/webp" --cache-control "max-age=31536000"
aws s3 cp positive_actions.webp s3://alertly-images-production/incidents/ --region us-west-2 --content-type "image/webp" --cache-control "max-age=31536000"
aws s3 cp lost_pet.webp s3://alertly-images-production/incidents/ --region us-west-2 --content-type "image/webp" --cache-control "max-age=31536000"
```

---

### 2. Investigar APIs Reales (3/5 pendientes)

#### ✅ TFS (Fire) - Implementado
- Tiene código funcional con mock data
- **Acción:** Investigar URL real de "Active Incidents"

#### ✅ Toronto Hydro - Implementado
- Tiene código funcional con mock data
- **Acción:** Reverse-engineer API del mapa de apagones

#### ❌ TPS (Toronto Police) - Pendiente
**URL a investigar:** `https://data.torontopolice.on.ca/pages/calls-for-service`

**Pasos:**
```bash
# 1. Abrir página en Chrome DevTools (Network tab)
open https://data.torontopolice.on.ca/pages/calls-for-service

# 2. Buscar peticiones XHR/Fetch de tipo JSON
# Filtrar por: XHR, Type: json

# 3. Copiar request URL y structure
# 4. Test con curl:
curl 'https://data.torontopolice.on.ca/api/...' | jq '.' > tps_sample.json

# 5. Enviarme el JSON de ejemplo
```

#### ❌ TTC (Transit) - Pendiente
**URL sugerida:** `https://www.ttc.ca/Service_Advisories/all_service_alerts.rss`

**Pasos:**
```bash
# 1. Verificar que el RSS funciona
curl 'https://www.ttc.ca/Service_Advisories/all_service_alerts.rss' | xmllint --format - > ttc_sample.xml

# 2. Si funciona, implementar parser XML
# 3. Si no funciona, buscar API alternativa en network tab de su sitio
```

#### ❌ Environment Canada - Pendiente
**URL sugerida:** `https://dd.weather.gc.ca/alerts/cap/`

**Pasos:**
```bash
# 1. Listar alertas disponibles
curl 'https://dd.weather.gc.ca/alerts/cap/' | grep -i toronto

# 2. Descargar alerta de muestra para Toronto
curl 'https://dd.weather.gc.ca/alerts/cap/...' > weather_sample.xml

# 3. Analizar formato CAP (Common Alerting Protocol)
```

---

### 3. Deployment a AWS Lambda

#### Actualizar Backend en ECR:
```bash
cd /Users/garyeikoow/Documents/www/alertly/backend

# 1. Login en ECR
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 129158986318.dkr.ecr.us-west-2.amazonaws.com

# 2. Build con Dockerfile.ecs
docker build -f Dockerfile.ecs -t alertly-api:latest .

# 3. Tag y push
docker tag alertly-api:latest 129158986318.dkr.ecr.us-west-2.amazonaws.com/alertly-api:latest
docker push 129158986318.dkr.ecr.us-west-2.amazonaws.com/alertly-api:latest
```

#### Actualizar Lambda (si está usando Lambda para cronjobs):
```bash
# Si tienes Lambda configurado para cronjobs
aws lambda update-function-code \
  --function-name alertly-production-oregon-alertlycronjobfunction-QIh2fMaVdqpc \
  --image-uri 129158986318.dkr.ecr.us-west-2.amazonaws.com/alertly-api:latest \
  --region us-west-2
```

#### Configurar EventBridge Schedules:
```bash
# Ejemplo para TFS (cada 10 minutos)
aws events put-rule \
  --name BotCreator-TFS \
  --schedule-expression "rate(10 minutes)" \
  --region us-west-2

aws events put-targets \
  --rule BotCreator-TFS \
  --targets "Id"="1","Arn"="arn:aws:lambda:us-west-2:129158986318:function:alertly-cronjob","Input"='{"task":"bot_creator_tfs"}' \
  --region us-west-2

# Ejemplo para Hydro (cada 30 minutos)
aws events put-rule \
  --name BotCreator-Hydro \
  --schedule-expression "rate(30 minutes)" \
  --region us-west-2

aws events put-targets \
  --rule BotCreator-Hydro \
  --targets "Id"="1","Arn"="arn:aws:lambda:us-west-2:129158986318:function:alertly-cronjob","Input"='{"task":"bot_creator_hydro"}' \
  --region us-west-2
```

---

### 4. Testing en Producción

#### Test Manual de Scrapers:
```bash
# Si tienes acceso a invocar Lambda directamente
aws lambda invoke \
  --function-name alertly-production-oregon-alertlycronjobfunction-QIh2fMaVdqpc \
  --payload '{"task":"bot_creator_tfs"}' \
  --region us-west-2 \
  response.json

cat response.json
```

#### Verificar Logs en CloudWatch:
```bash
# Ver logs del cronjob
aws logs tail /aws/lambda/alertly-production-oregon-alertlycronjobfunction-QIh2fMaVdqpc \
  --follow \
  --region us-west-2
```

#### Verificar Incidentes en Base de Datos:
```bash
# Conectar a RDS
mysql -h alertly-mysql-freetier.c3qmq4y86s84.us-west-2.rds.amazonaws.com \
  -u adminalertly -p alertly

# Query para ver incidentes del bot
SELECT
  ir.inre_id,
  ir.title,
  ir.created_at,
  a.username
FROM incident_reports ir
JOIN account a ON ir.user_id = a.account_id
WHERE a.username = 'System Bot' OR ir.user_id = 1
ORDER BY ir.created_at DESC
LIMIT 20;

# Ver hashes guardados (deduplicación)
SELECT * FROM bot_incident_hashes ORDER BY created_at DESC LIMIT 20;

# Ver cache de geocoding
SELECT * FROM geocoding_cache ORDER BY last_used_at DESC LIMIT 20;
```

---

## 📊 Progreso General

| Componente               | Status      | Completado |
|--------------------------|-------------|------------|
| Código Backend           | ✅ Done     | 100%       |
| Database Migration       | ⏳ Pending  | 0%         |
| Imágenes S3              | ⏳ Partial  | 67%        |
| Scrapers Implementados   | ⏳ Partial  | 40% (2/5)  |
| APIs Research            | ⏳ Pending  | 0%         |
| Lambda Deployment        | ⏳ Pending  | 0%         |
| EventBridge Config       | ⏳ Pending  | 0%         |
| Production Testing       | ⏳ Pending  | 0%         |

**Overall:** ~30% completado

---

## 🎯 Roadmap Recomendado

### Fase 1: Preparación (Esta Semana)
1. ✅ Código implementado
2. ⏳ Aplicar migration de database
3. ⏳ Crear 4 imágenes faltantes y subirlas
4. ⏳ Investigar URLs reales de TPS/TTC/Weather

### Fase 2: Deployment (Próxima Semana)
1. Build y push de nueva imagen Docker
2. Configurar EventBridge schedules
3. Testear con mock data en Lambda
4. Verificar logs y funcionamiento

### Fase 3: APIs Reales (Siguiente Semana)
1. Implementar scrapers restantes (TPS, TTC, Weather)
2. Testear con datos reales
3. Ajustar intervalos de ejecución
4. Monitorear métricas

### Fase 4: Optimización (Mes 2)
1. Analizar performance de geocoding
2. Ajustar TTLs de deduplicación
3. Monitorear costos de Lambda
4. Configurar alertas de errores

---

## 💰 Costos Estimados

| Servicio              | Costo/Mes Estimado |
|-----------------------|-------------------|
| Lambda Executions     | ~$2-5             |
| RDS Storage (hashes)  | ~$1               |
| S3 Images (12 files)  | ~$0.01            |
| CloudWatch Logs       | ~$1               |
| **Total Bot Creator** | **~$4-7/mes**     |

**Nota:** Los scrapers usan las mismas tablas MySQL del Free Tier existente, sin costo adicional.

---

## 📞 Ayuda Necesaria

**De tu parte:**
1. ✅ Crear 4 imágenes faltantes (vandalism, community_events, positive_actions, lost_pet)
2. ✅ Investigar URLs reales de APIs (TPS, TTC, Hydro, Weather)
3. ✅ Aplicar database migration
4. ✅ Testear funcionamiento en producción

**De mi parte (cuando me des las URLs):**
1. Implementar scrapers restantes
2. Ajustar parsers según estructura real de datos
3. Optimizar lógica de normalización
4. Documentar configuración de EventBridge

---

**📅 Creado:** 23 de Noviembre, 2025
**📊 Estado Actual:** 30% completado (código + imágenes parciales)
**🎯 Próximo Milestone:** Aplicar migration + completar imágenes + research APIs
