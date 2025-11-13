# ✅ CHECKLIST FINAL - Sistema de Referrals

**Estado:** 🚧 Implementación completa - Pendiente testing y deployment
**Fecha:** 2 de Noviembre, 2025

---

## 📋 FASE 1: TESTING LOCAL (OBLIGATORIO)

### ✅ Paso 1: Aplicar Migración SQL a Base de Datos Local

```bash
cd /Users/garyeikoow/Documents/www/alertly/backend
mysql -u root -p alertly < assets/db/migrations/add_referral_system.sql
```

**Verifica que se crearon las tablas:**
```bash
mysql -u root -p alertly -e "SHOW TABLES LIKE '%referral%';"
mysql -u root -p alertly -e "SHOW TABLES LIKE 'influencers';"
```

**Deberías ver:**
- `influencers`
- `referral_conversions`
- `referral_premium_conversions`
- `referral_metrics_cache`

**Verifica influencers de prueba:**
```bash
mysql -u root -p alertly -e "SELECT * FROM influencers;"
```

Deberías ver 5 influencers (INF-IG0001, INF-TT0002, INF-RD0003, INF-IG0004, INF-TT0005)

---

### ✅ Paso 2: Compilar Backend (Verificar que no hay errores)

```bash
cd /Users/garyeikoow/Documents/www/alertly/backend
go build -o alertly-api cmd/app/main.go
```

**Si hay errores de compilación:**
- Ejecutar: `go mod tidy` para resolver dependencias
- Revisar logs de error y corregir

**Si compila exitosamente:**
```
✓ Backend compilado sin errores
```

---

### ✅ Paso 3: Ejecutar Backend Localmente

```bash
cd /Users/garyeikoow/Documents/www/alertly/backend
go run cmd/app/main.go
```

**Buscar en los logs:**
```
✅ Referral system endpoints registered
🚀 Starting Alertly Backend...
Alertly Backend starting on port 8080
```

**Si no arranca:**
- Verificar que MySQL está corriendo
- Verificar credenciales en `.env`
- Verificar que el puerto 8080 está libre

---

### ✅ Paso 4: Probar Endpoints (En otra terminal)

#### Test 1: Validar código válido (PÚBLICO)
```bash
curl -X POST http://localhost:8080/api/v1/referral/validate \
  -H "Content-Type: application/json" \
  -d '{"referral_code":"INF-IG0001"}' | jq
```

**Resultado esperado:**
```json
{
  "valid": true,
  "influencer_id": 1,
  "influencer_name": "John Doe",
  "premium_trial_days": 14
}
```

**✅ PASSED** / **❌ FAILED**

---

#### Test 2: Validar código inválido
```bash
curl -X POST http://localhost:8080/api/v1/referral/validate \
  -H "Content-Type: application/json" \
  -d '{"referral_code":"INVALID"}' | jq
```

**Resultado esperado:**
```json
{
  "valid": false,
  "message": "Invalid referral code"
}
```

**✅ PASSED** / **❌ FAILED**

---

#### Test 3: Registrar conversión (PROTEGIDO)
```bash
curl -X POST http://localhost:8080/api/v1/referral/conversion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071" \
  -d '{
    "referral_code":"INF-IG0001",
    "user_id":999999,
    "registered_at":"2025-11-02T12:30:00Z",
    "platform":"iOS"
  }' | jq
```

**Resultado esperado:**
```json
{
  "success": true,
  "message": "Conversion recorded",
  "influencer_earnings_added": 0.10
}
```

**Verificar en BD:**
```bash
mysql -u root -p alertly -e "SELECT * FROM referral_conversions WHERE user_id = 999999;"
```

**✅ PASSED** / **❌ FAILED**

---

#### Test 4: Registrar conversión premium
```bash
curl -X POST http://localhost:8080/api/v1/referral/premium-conversion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071" \
  -d '{
    "user_id":999999,
    "referral_code":"INF-IG0001",
    "subscription_type":"monthly",
    "amount":7.99,
    "converted_at":"2025-11-05T18:45:00Z"
  }' | jq
```

**Resultado esperado:**
```json
{
  "success": true,
  "message": "Premium conversion recorded",
  "influencer_commission": 1.20,
  "commission_percentage": 15.00
}
```

**Verificar en BD:**
```bash
mysql -u root -p alertly -e "SELECT * FROM referral_premium_conversions WHERE user_id = 999999;"
```

**✅ PASSED** / **❌ FAILED**

---

#### Test 5: Obtener métricas individuales
```bash
curl -X GET "http://localhost:8080/api/v1/referrals/metrics?code=INF-IG0001" \
  -H "Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071" | jq
```

**Resultado esperado:** JSON con métricas del influencer

**✅ PASSED** / **❌ FAILED**

---

#### Test 6: Obtener métricas agregadas
```bash
curl -X GET http://localhost:8080/api/v1/referrals/aggregate \
  -H "Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071" | jq
```

**Resultado esperado:** JSON con métricas globales

**✅ PASSED** / **❌ FAILED**

---

#### Test 7: Sincronizar influencer
```bash
curl -X POST http://localhost:8080/api/v1/referral/sync-influencer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071" \
  -d '{
    "web_influencer_id":999,
    "referral_code":"INF-TEST001",
    "name":"Test Influencer",
    "platform":"Instagram",
    "is_active":true
  }' | jq
```

**Resultado esperado:**
```json
{
  "success": true,
  "message": "Influencer INF-TEST001 synced successfully"
}
```

**Verificar en BD:**
```bash
mysql -u root -p alertly -e "SELECT * FROM influencers WHERE referral_code = 'INF-TEST001';"
```

**✅ PASSED** / **❌ FAILED**

---

#### Test 8: Signup con referral code
```bash
curl -X POST http://localhost:8080/account/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email":"testref@example.com",
    "first_name":"Test",
    "last_name":"Referral",
    "password":"password123",
    "birth_year":"1990",
    "birth_month":"05",
    "birth_day":"15",
    "referral_code":"INF-IG0001"
  }' | jq
```

**Verificar:**
1. Usuario creado exitosamente
2. **En logs del backend:** `✅ Referral conversion registered successfully`
3. **En BD:**
```bash
mysql -u root -p alertly -e "SELECT account_id, email, is_premium, premium_expired_date FROM account WHERE email = 'testref@example.com';"
```
4. `premium_expired_date` debe ser ~14 días desde NOW()

**Verificar conversión automática:**
```bash
mysql -u root -p alertly -e "SELECT * FROM referral_conversions WHERE user_id = (SELECT account_id FROM account WHERE email = 'testref@example.com');"
```

**✅ PASSED** / **❌ FAILED**

---

#### Test 9: Signup SIN referral code (verificar 7 días)
```bash
curl -X POST http://localhost:8080/account/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email":"testnoref@example.com",
    "first_name":"Test",
    "last_name":"NoReferral",
    "password":"password123",
    "birth_year":"1990",
    "birth_month":"05",
    "birth_day":"15"
  }' | jq
```

**Verificar en BD:**
```bash
mysql -u root -p alertly -e "SELECT account_id, email, is_premium, premium_expired_date FROM account WHERE email = 'testnoref@example.com';"
```

`premium_expired_date` debe ser ~7 días desde NOW()

**✅ PASSED** / **❌ FAILED**

---

### 📊 Resumen Testing Local

**Total tests:** 9
**Passed:** _____ / 9
**Failed:** _____ / 9

**Si todos pasan:** ✅ Listo para deployment a producción
**Si alguno falla:** ❌ Revisar logs y corregir antes de continuar

---

## 📋 FASE 2: DEPLOYMENT A PRODUCCIÓN

### ⚠️ SOLO PROCEDER SI TODOS LOS TESTS LOCALES PASARON

---

### ✅ Paso 5: Aplicar Migración a RDS (Producción)

```bash
mysql -h alertly-mysql-freetier.c3qmq4y86s84.us-west-2.rds.amazonaws.com \
  -u adminalertly \
  -p \
  alertly < /Users/garyeikoow/Documents/www/alertly/backend/assets/db/migrations/add_referral_system.sql
```

**Password:** `Po1Ng2O3;`

**Verificar:**
```bash
mysql -h alertly-mysql-freetier.c3qmq4y86s84.us-west-2.rds.amazonaws.com \
  -u adminalertly \
  -p \
  alertly -e "SHOW TABLES LIKE '%referral%';"
```

**✅ MIGRACIÓN APLICADA** / **❌ ERROR**

---

### ✅ Paso 6: Build y Push Docker Image

```bash
cd /Users/garyeikoow/Documents/www/alertly/backend

# 1. Login en ECR
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 129158986318.dkr.ecr.us-west-2.amazonaws.com

# 2. Build imagen
docker build -f Dockerfile.ecs -t alertly-api:latest .

# 3. Tag imagen
docker tag alertly-api:latest 129158986318.dkr.ecr.us-west-2.amazonaws.com/alertly-api:latest

# 4. Push a ECR
docker push 129158986318.dkr.ecr.us-west-2.amazonaws.com/alertly-api:latest
```

**✅ IMAGEN PUSHED** / **❌ ERROR**

---

### ✅ Paso 7: Deploy a EC2

```bash
# 1. SSH a instancia
ssh -i /Users/garyeikoow/Documents/www/alertly/alertly-debug.pem ec2-user@44.243.7.9

# 2. Dentro de EC2: Stop y remove container
sudo docker stop alertly-api && sudo docker rm alertly-api

# 3. Login en ECR
aws ecr get-login-password --region us-west-2 | sudo docker login --username AWS --password-stdin 129158986318.dkr.ecr.us-west-2.amazonaws.com

# 4. Pull nueva imagen
sudo docker pull 129158986318.dkr.ecr.us-west-2.amazonaws.com/alertly-api:latest

# 5. Run nuevo container (COPIAR EXACTAMENTE - incluye REFERRAL_API_KEY)
sudo docker run -d --name alertly-api -p 80:8080 --restart unless-stopped \
  -e DB_HOST=alertly-mysql-freetier.c3qmq4y86s84.us-west-2.rds.amazonaws.com \
  -e DB_PORT=3306 \
  -e DB_USER=adminalertly \
  -e DB_PASS="Po1Ng2O3;" \
  -e DB_NAME=alertly \
  -e GIN_MODE=release \
  -e PORT=8080 \
  -e JWT_SECRET="AlertlySecretKey2024!ProductionJWT" \
  -e IMAGE_BASE_URL="https://cdn.alertly.ca" \
  -e AWS_REGION=us-west-2 \
  -e REFERRAL_API_KEY="0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071" \
  129158986318.dkr.ecr.us-west-2.amazonaws.com/alertly-api:latest

# 6. Verificar logs
sudo docker logs alertly-api

# Buscar: "✅ Referral system endpoints registered"

# 7. Exit SSH
exit
```

**✅ DEPLOYED** / **❌ ERROR**

---

### ✅ Paso 8: Probar en Producción

#### Health Check
```bash
curl -w "\n⏱️  Response: %{time_total}s - Status: %{http_code}\n" https://api.alertly.ca/health
```

**Resultado esperado:** Status 200, "healthy"

---

#### Test Endpoint Público en Producción
```bash
curl -X POST https://api.alertly.ca/api/v1/referral/validate \
  -H "Content-Type: application/json" \
  -d '{"referral_code":"INF-IG0001"}' | jq
```

**Resultado esperado:**
```json
{
  "valid": true,
  "influencer_id": 1,
  "influencer_name": "John Doe",
  "premium_trial_days": 14
}
```

**✅ PRODUCCIÓN FUNCIONANDO** / **❌ ERROR**

---

#### Test Endpoint Protegido en Producción
```bash
curl -X GET "https://api.alertly.ca/api/v1/referrals/aggregate" \
  -H "Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071" | jq
```

**Resultado esperado:** JSON con métricas agregadas

**✅ AUTENTICACIÓN FUNCIONANDO** / **❌ ERROR**

---

## 📋 FASE 3: INTEGRACIÓN CON BACKEND WEB

### ✅ Paso 9: Compartir Credenciales con Equipo Web

**Enviar al equipo del backend web (Symfony/PHP):**

#### API Key:
```
0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071
```

#### Base URL:
```
https://api.alertly.ca/api/v1
```

#### Endpoints disponibles:
1. `POST /referral/validate` (PÚBLICO)
2. `POST /referral/conversion` (Requiere API Key)
3. `POST /referral/premium-conversion` (Requiere API Key)
4. `GET /referrals/metrics?code={CODE}` (Requiere API Key)
5. `GET /referrals/aggregate` (Requiere API Key)
6. `POST /referral/sync-influencer` (Requiere API Key)

#### Documentación:
Compartir el archivo: `backend/REFERRAL_SYSTEM_IMPLEMENTATION.md`

**✅ COMPARTIDO** / **❌ PENDIENTE**

---

### ✅ Paso 10: Probar Integración End-to-End

**Una vez que el backend web esté configurado:**

1. Backend web crea un nuevo influencer
2. Backend web llama a `POST /api/v1/referral/sync-influencer`
3. Usuario usa código en app móvil
4. App llama a `POST /api/v1/referral/validate`
5. Usuario completa signup
6. Backend nativo llama automáticamente a `POST /api/v1/referral/conversion`
7. Usuario se suscribe a premium
8. Backend nativo llama a `POST /api/v1/referral/premium-conversion`
9. Backend web consulta métricas con `GET /api/v1/referrals/metrics`

**✅ INTEGRACIÓN COMPLETA** / **❌ PENDIENTE**

---

## 🎯 CRITERIOS DE ÉXITO FINAL

### Testing Local:
- ✅ Migración SQL aplicada
- ✅ Backend compila sin errores
- ✅ 9/9 tests locales pasan
- ✅ Signup con código funciona (14 días premium)
- ✅ Signup sin código funciona (7 días premium)
- ✅ Conversión automática se registra

### Producción:
- ✅ Migración aplicada a RDS
- ✅ Docker image deployed
- ✅ Endpoints públicos funcionan con HTTPS
- ✅ Endpoints protegidos validan API Key
- ✅ Health check retorna "healthy"

### Integración:
- ✅ API Key compartido con backend web
- ✅ Documentación compartida
- ✅ Backend web puede consumir endpoints
- ✅ Flujo end-to-end funciona

---

## 📞 SIGUIENTE ACCIÓN

**SI ESTÁS LISTO PARA TESTING LOCAL:**
```bash
# Ejecuta estos comandos en orden:
cd /Users/garyeikoow/Documents/www/alertly/backend
mysql -u root -p alertly < assets/db/migrations/add_referral_system.sql
go run cmd/app/main.go
```

**SI ALGO FALLA:**
- Revisar logs del servidor
- Verificar conexión a MySQL
- Verificar que `.env` tiene todas las variables
- Consultar `backend/REFERRAL_SYSTEM_IMPLEMENTATION.md`

---

## ✅ ESTADO ACTUAL

**Implementación:** ✅ 100% COMPLETA
**Testing Local:** ⏳ PENDIENTE
**Deployment Producción:** ⏳ PENDIENTE
**Integración Web:** ⏳ PENDIENTE

**🎯 Próximo paso:** Ejecutar Fase 1 - Testing Local
