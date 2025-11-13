# 🔗 Guía de Integración - Sistema de Referrals (Para Web App)

## 📋 Resumen Ejecutivo

El sistema de referrals para influencers ha sido **completamente implementado y desplegado en producción** en el backend de Alertly. Este documento contiene toda la información necesaria para que la web app pueda consumir los endpoints del sistema de referrals.

**Estado:** ✅ Producción - Completamente funcional
**Fecha de deployment:** 2 de Noviembre, 2025
**Base URL:** `https://api.alertly.ca`

---

## 🎯 ¿Qué se implementó?

Se creó un sistema completo de referidos para influencers que permite:

1. **Validar códigos de referral** antes del signup (endpoint público)
2. **Registrar conversiones automáticamente** cuando un usuario se registra con un código
3. **Registrar conversiones premium** cuando un usuario referido compra una suscripción
4. **Obtener métricas** individuales y agregadas de influencers
5. **Sincronizar influencers** desde la web app al backend móvil

### 💰 Comisiones Automáticas

- **Signup con código de referral:** $0.10 USD por usuario
- **Conversión a premium:** 15% del monto de la suscripción
- **Bonus:** Usuarios con código obtienen 14 días de trial (vs 7 días sin código)

---

## 🔐 Autenticación

La mayoría de los endpoints requieren autenticación con API Key.

### API Key de Producción

```bash
API_KEY: 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071
```

### Cómo usar el API Key

Incluir en el header `Authorization` de cada request:

```bash
Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071
```

**Importante:** El único endpoint que NO requiere API Key es `/api/v1/referral/validate` (validación de código).

---

## 📡 Endpoints Disponibles

### 1️⃣ Validar Código de Referral (PÚBLICO - Sin API Key)

Valida si un código de referral existe y está activo.

**Endpoint:** `POST /api/v1/referral/validate`
**Autenticación:** ❌ No requiere (público)

**Request:**
```json
{
  "referral_code": "INF-IG0001"
}
```

**Response exitoso (200):**
```json
{
  "valid": true,
  "influencer_id": 1,
  "influencer_name": "John Doe",
  "premium_trial_days": 14
}
```

**Response código inválido (200):**
```json
{
  "valid": false,
  "message": "Invalid referral code"
}
```

**Ejemplo cURL:**
```bash
curl -X POST https://api.alertly.ca/api/v1/referral/validate \
  -H "Content-Type: application/json" \
  -d '{"referral_code": "INF-IG0001"}'
```

---

### 2️⃣ Registrar Conversión (Signup)

Registra cuando un usuario se registra usando un código de referral.

**Endpoint:** `POST /api/v1/referral/conversion`
**Autenticación:** ✅ Requiere API Key

**Request:**
```json
{
  "referral_code": "INF-IG0001",
  "user_id": 123,
  "registered_at": "2025-11-02T20:20:00Z",
  "platform": "iOS"
}
```

**Campos:**
- `referral_code`: Código del influencer (string, requerido)
- `user_id`: ID del usuario en la tabla `account` (int, requerido)
- `registered_at`: Timestamp ISO 8601 (string, requerido)
- `platform`: "iOS" o "Android" (string, requerido)

**Response (201):**
```json
{
  "success": true,
  "message": "Conversion recorded",
  "influencer_earnings_added": 0.10
}
```

**Ejemplo cURL:**
```bash
curl -X POST https://api.alertly.ca/api/v1/referral/conversion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071" \
  -d '{
    "referral_code": "INF-IG0001",
    "user_id": 123,
    "registered_at": "2025-11-02T20:20:00Z",
    "platform": "iOS"
  }'
```

**⚠️ Nota Importante:** El backend móvil YA hace esto automáticamente cuando un usuario se registra con un código. Solo necesitas llamar este endpoint si quieres registrar conversiones manualmente desde la web.

---

### 3️⃣ Registrar Conversión Premium

Registra cuando un usuario referido compra una suscripción premium.

**Endpoint:** `POST /api/v1/referral/premium-conversion`
**Autenticación:** ✅ Requiere API Key

**Request:**
```json
{
  "referral_code": "INF-IG0001",
  "user_id": 123,
  "subscription_type": "monthly",
  "amount": 4.99,
  "converted_at": "2025-11-02T20:21:00Z"
}
```

**Campos:**
- `referral_code`: Código del influencer (string, requerido)
- `user_id`: ID del usuario en la tabla `account` (int, requerido)
- `subscription_type`: "monthly" o "yearly" (string, requerido)
- `amount`: Monto de la suscripción en USD (float, requerido)
- `converted_at`: Timestamp ISO 8601 (string, requerido)

**Response (201):**
```json
{
  "success": true,
  "message": "Premium conversion recorded",
  "influencer_commission": 0.75,
  "commission_percentage": 15.00
}
```

**Cálculo de comisión:** `commission = amount * 0.15` (15%)

**Ejemplo cURL:**
```bash
curl -X POST https://api.alertly.ca/api/v1/referral/premium-conversion \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071" \
  -d '{
    "referral_code": "INF-IG0001",
    "user_id": 123,
    "subscription_type": "monthly",
    "amount": 4.99,
    "converted_at": "2025-11-02T20:21:00Z"
  }'
```

---

### 4️⃣ Obtener Métricas Individuales de Influencer

Obtiene las métricas detalladas de un influencer específico.

**Endpoint:** `GET /api/v1/referrals/metrics?code={referral_code}`
**Autenticación:** ✅ Requiere API Key

**Parámetros Query:**
- `code`: Código de referral del influencer (requerido)

**Response (200):**
```json
{
  "referral_code": "INF-IG0001",
  "influencer_id": 1,
  "total_registrations": 5,
  "total_premium_conversions": 2,
  "total_earnings": 1.95,
  "current_month_registrations": 3,
  "current_month_premium": 1,
  "current_month_earnings": 0.85,
  "projected_month_earnings": 12.75,
  "daily_metrics": [
    {
      "date": "2025-11-02T00:00:00Z",
      "registrations": 2,
      "premium_conversions": 1,
      "earnings": 0.95
    }
  ],
  "rank": 1,
  "total_influencers": 5
}
```

**Ejemplo cURL:**
```bash
curl -X GET "https://api.alertly.ca/api/v1/referrals/metrics?code=INF-IG0001" \
  -H "Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071"
```

---

### 5️⃣ Obtener Métricas Agregadas

Obtiene las métricas globales del sistema de referrals (todos los influencers).

**Endpoint:** `GET /api/v1/referrals/aggregate`
**Autenticación:** ✅ Requiere API Key

**Response (200):**
```json
{
  "total_referrals": 15,
  "total_premium_conversions": 5,
  "total_earnings_paid": 4.25,
  "active_influencers": 4,
  "conversion_rate": 33.33,
  "top_performers": [
    {
      "influencer_id": 1,
      "referral_code": "INF-IG0001",
      "name": "John Doe",
      "total_registrations": 5,
      "total_premium_conversions": 2,
      "total_earnings": 1.95,
      "platform": "Instagram"
    }
  ],
  "monthly_trend": [
    {
      "month": "2025-11",
      "registrations": 15,
      "premium": 5,
      "earnings": 4.25
    }
  ],
  "platform_breakdown": {
    "Instagram": {
      "influencers": 2,
      "registrations": 8,
      "premium": 3,
      "earnings": 2.15
    },
    "TikTok": {
      "influencers": 1,
      "registrations": 5,
      "premium": 2,
      "earnings": 1.50
    }
  }
}
```

**Ejemplo cURL:**
```bash
curl -X GET https://api.alertly.ca/api/v1/referrals/aggregate \
  -H "Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071"
```

---

### 6️⃣ Sincronizar Influencer

Crea o actualiza un influencer en la base de datos del backend móvil.

**Endpoint:** `POST /api/v1/referral/sync-influencer`
**Autenticación:** ✅ Requiere API Key

**Request:**
```json
{
  "web_influencer_id": 42,
  "referral_code": "INF-IG0042",
  "name": "Jane Smith",
  "platform": "Instagram",
  "is_active": true
}
```

**Campos:**
- `web_influencer_id`: ID del influencer en la base de datos web (int, requerido)
- `referral_code`: Código único del influencer (string, requerido, max 20 chars)
- `name`: Nombre del influencer (string, requerido)
- `platform`: "Instagram", "TikTok", "Reddit", u "Other" (string, requerido)
- `is_active`: Estado activo/inactivo (boolean, requerido)

**Response (200):**
```json
{
  "success": true,
  "message": "Influencer INF-IG0042 synced successfully"
}
```

**Comportamiento:**
- Si `web_influencer_id` ya existe: **actualiza** el registro
- Si no existe: **crea** nuevo registro

**Ejemplo cURL:**
```bash
curl -X POST https://api.alertly.ca/api/v1/referral/sync-influencer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071" \
  -d '{
    "web_influencer_id": 42,
    "referral_code": "INF-IG0042",
    "name": "Jane Smith",
    "platform": "Instagram",
    "is_active": true
  }'
```

---

## 🗄️ Estructura de Base de Datos

El sistema creó 4 tablas en la base de datos RDS de producción:

### Tabla: `influencers`

```sql
CREATE TABLE influencers (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    web_influencer_id INT UNSIGNED NOT NULL,
    referral_code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    platform ENUM('Instagram', 'TikTok', 'Reddit', 'Other'),
    is_active TINYINT(1) DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_web_id (web_influencer_id)
);
```

### Tabla: `referral_conversions`

```sql
CREATE TABLE referral_conversions (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    referral_code VARCHAR(20) NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    registered_at DATETIME NOT NULL,
    platform ENUM('iOS', 'Android'),
    earnings DECIMAL(10,2) DEFAULT 0.10,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_referral (user_id),
    FOREIGN KEY (user_id) REFERENCES account(account_id) ON DELETE CASCADE
);
```

### Tabla: `referral_premium_conversions`

```sql
CREATE TABLE referral_premium_conversions (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    referral_code VARCHAR(20) NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    conversion_id BIGINT UNSIGNED NULL,
    subscription_type ENUM('monthly', 'yearly'),
    amount DECIMAL(10,2) NOT NULL,
    commission DECIMAL(10,2) NOT NULL,
    commission_percentage DECIMAL(5,2) DEFAULT 15.00,
    converted_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES account(account_id) ON DELETE CASCADE,
    FOREIGN KEY (conversion_id) REFERENCES referral_conversions(id)
);
```

### Tabla: `referral_metrics_cache`

```sql
CREATE TABLE referral_metrics_cache (
    referral_code VARCHAR(20) PRIMARY KEY,
    total_registrations INT UNSIGNED DEFAULT 0,
    total_premium_conversions INT UNSIGNED DEFAULT 0,
    total_earnings DECIMAL(10,2) DEFAULT 0.00,
    last_updated DATETIME NOT NULL
);
```

---

## 🔄 Flujos de Integración Recomendados

### Flujo 1: Dashboard de Influencer (Métricas)

```javascript
// 1. Obtener métricas individuales del influencer
const response = await fetch(
  `https://api.alertly.ca/api/v1/referrals/metrics?code=${influencerCode}`,
  {
    headers: {
      'Authorization': 'Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071'
    }
  }
);

const metrics = await response.json();

// Mostrar en dashboard:
// - Total earnings: metrics.total_earnings
// - Registrations: metrics.total_registrations
// - Premium conversions: metrics.total_premium_conversions
// - Rank: metrics.rank / metrics.total_influencers
// - Chart: metrics.daily_metrics
```

### Flujo 2: Panel de Admin (Métricas Globales)

```javascript
// 1. Obtener métricas agregadas
const response = await fetch(
  'https://api.alertly.ca/api/v1/referrals/aggregate',
  {
    headers: {
      'Authorization': 'Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071'
    }
  }
);

const data = await response.json();

// Mostrar:
// - Total earnings: data.total_earnings_paid
// - Top performers: data.top_performers
// - Platform breakdown: data.platform_breakdown
// - Monthly trend: data.monthly_trend
```

### Flujo 3: Crear/Actualizar Influencer

```javascript
// Cuando creas un influencer en tu base de datos web:
async function syncInfluencer(influencer) {
  const response = await fetch(
    'https://api.alertly.ca/api/v1/referral/sync-influencer',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071'
      },
      body: JSON.stringify({
        web_influencer_id: influencer.id,
        referral_code: influencer.referral_code,
        name: influencer.name,
        platform: influencer.platform,
        is_active: influencer.is_active
      })
    }
  );

  const result = await response.json();
  console.log('Synced:', result.message);
}
```

### Flujo 4: Registrar Conversión Premium (Cuando usuario paga)

```javascript
// Cuando un usuario compra una suscripción premium:
async function registerPremiumConversion(userId, referralCode, amount, subscriptionType) {
  const response = await fetch(
    'https://api.alertly.ca/api/v1/referral/premium-conversion',
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer 0560f9d085d2fe20fd7ddcd51024abb508d68f07935ad875b8769773b62d5071'
      },
      body: JSON.stringify({
        referral_code: referralCode,
        user_id: userId,
        subscription_type: subscriptionType, // 'monthly' o 'yearly'
        amount: amount, // e.g., 4.99
        converted_at: new Date().toISOString()
      })
    }
  );

  const result = await response.json();
  console.log('Commission:', result.influencer_commission);
}
```

---

## ⚠️ Puntos Importantes

### 1. El Backend Móvil Ya Maneja Conversiones de Signup

**No necesitas** llamar `/api/v1/referral/conversion` cuando un usuario se registra desde la app móvil. El backend móvil automáticamente:

1. Recibe el `referral_code` en el signup
2. Asigna 14 días de premium trial (vs 7 sin código)
3. Llama internamente al endpoint de conversión
4. Registra los $0.10 de comisión

**Solo necesitas** llamar `/api/v1/referral/conversion` si quieres registrar conversiones manualmente desde la web.

### 2. Conversiones Premium SÍ necesitan tu intervención

Cuando un usuario compra una suscripción premium (mensual o anual), **debes llamar** al endpoint `/api/v1/referral/premium-conversion` para:

- Registrar la conversión
- Calcular la comisión (15% del monto)
- Actualizar las métricas del influencer

### 3. Sincronización de Influencers

Cada vez que creas o actualizas un influencer en tu base de datos web, **debes sincronizarlo** con el backend móvil usando `/api/v1/referral/sync-influencer`.

Esto asegura que:
- Los códigos estén disponibles en la app móvil
- Las validaciones funcionen correctamente
- Las métricas sean consistentes

---

## 🧪 Datos de Prueba (Producción)

Se insertaron 5 influencers de prueba en la base de datos de producción:

| ID | Código | Nombre | Plataforma | Activo |
|----|--------|--------|------------|--------|
| 1 | INF-IG0001 | John Doe | Instagram | ✅ |
| 2 | INF-TT0002 | Jane Smith | TikTok | ✅ |
| 3 | INF-RD0003 | Bob Johnson | Reddit | ✅ |
| 4 | INF-IG0004 | Alice Williams | Instagram | ✅ |
| 5 | INF-TT0005 | Charlie Brown | TikTok | ❌ |

Puedes usar estos códigos para testing.

---

## 📊 Tests Realizados

Todos los endpoints fueron testeados en producción:

- ✅ Validación de código válido
- ✅ Validación de código inválido
- ✅ Registro de conversión
- ✅ Registro de conversión premium
- ✅ Métricas individuales
- ✅ Métricas agregadas
- ✅ Sincronización de influencer
- ✅ Signup con código (14 días trial + conversión automática)
- ✅ Signup sin código (7 días trial + sin conversión)

**Estado:** Sistema completamente funcional en producción.

---

## 🔧 Manejo de Errores

Todos los endpoints retornan errores en formato JSON:

### Errores de Autenticación (401)
```json
{
  "error": "Invalid API key"
}
```

### Errores de Validación (400)
```json
{
  "error": "Missing required field: referral_code"
}
```

### Errores Internos (500)
```json
{
  "error": "Internal server error",
  "details": "error creating conversion: ..."
}
```

---

## 📞 Soporte Técnico

Para dudas o problemas:

1. **Documentación completa:** `/backend/REFERRAL_SYSTEM_IMPLEMENTATION.md`
2. **Checklist de testing:** `/backend/FINAL_CHECKLIST.md`
3. **Código fuente:** `/backend/internal/referrals/`

---

**Creado:** 2 de Noviembre, 2025
**Versión:** 1.0
**Estado:** ✅ Producción - Listo para integración
