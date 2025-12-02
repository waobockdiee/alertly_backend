# Scraper Testing Guide

Este documento explica cómo probar los scrapers localmente para verificar si están trayendo datos de las fuentes externas.

## 🚀 Comandos de Prueba

### Probar un scraper individual:

```bash
# Desde el directorio backend/
cd backend

# Probar Toronto Fire Services
go run test_scrapers.go tfs

# Probar Toronto Hydro
go run test_scrapers.go hydro

# Probar todos los scrapers
go run test_scrapers.go all
```

## 📊 Interpretación de Resultados

### ✅ **Scraping exitoso con datos**
```
✅ Real scraping successful!
📊 RESULTS: Found 5 incidents
```
- El scraper se conectó correctamente al endpoint
- Encontró y parseó datos exitosamente
- Los datos se muestran en pantalla

### ⚠️ **Scraping exitoso pero sin datos**
```
✅ Real scraping successful!
📊 RESULTS: Found 0 incidents
```
- La conexión fue exitosa
- El HTML/JSON se parseó correctamente
- No hay incidentes activos en este momento (normal)
- O los selectores HTML necesitan ajustes

### ❌ **Scraping fallido (usando MOCK data)**
```
⚠️  Real scraping failed: Hydro API returned status 403
📦 Falling back to MOCK data...
```
- El endpoint real no está accesible (403, 404, timeout, etc.)
- El scraper usa datos de prueba automáticamente
- Necesitas:
  - Verificar la URL del endpoint
  - Verificar si requiere autenticación
  - Verificar headers (User-Agent, API keys, etc.)

## 🔍 Estado Actual de los Scrapers

### **TFS (Toronto Fire Services)**
- **Status:** ✅ Conecta correctamente
- **URL:** `https://www.toronto.ca/community-people/public-safety-alerts/alerts-notifications/`
- **Problema:** No encuentra incidentes (selectores HTML necesitan ajuste)
- **Acción requerida:** Inspeccionar el HTML real de la página y actualizar selectores en `scrapers/tfs.go`

### **Hydro (Toronto Hydro)**
- **Status:** ❌ API no accesible (403 Forbidden)
- **URL:** `https://api.torontohydro.com/outages/current` (placeholder)
- **Problema:** URL es hipotética, necesita reverse-engineering
- **Acción requerida:**
  1. Visitar https://www.torontohydro.com/outage-map
  2. Abrir DevTools → Network tab
  3. Filtrar por XHR/Fetch
  4. Encontrar el endpoint real de la API
  5. Actualizar `HYDRO_API_URL` en `scrapers/hydro.go`

### **TPS, TTC, Weather**
- **Status:** ⏳ No implementados aún
- Retornan mensajes de "not yet implemented"

## 🛠️ Próximos Pasos

1. **Para TFS:**
   - Visitar la página en el navegador
   - Inspeccionar el HTML (F12 → Elements)
   - Encontrar la tabla/lista de incidentes activos
   - Actualizar los selectores CSS en `scrapers/tfs.go:61`

2. **Para Hydro:**
   - Reverse-engineer el endpoint real de la API
   - Verificar estructura del JSON response
   - Actualizar `HYDRO_API_URL` y structs en `scrapers/hydro.go`

3. **Para nuevos scrapers:**
   - Crear archivo en `backend/internal/cronjobs/cjbot_creator/scrapers/`
   - Implementar interfaz con métodos `Scrape()` y `ScrapeMockData()`
   - Agregar test case en `test_scrapers.go`
   - Agregar mapeos de categorías en `normalizer.go`

## 📝 Ejemplo de Output Exitoso

```
[1] 🔥 TFS Incident
    External ID:  tfs_100_queen_st_w_1732558800
    Title:        Fire Service Call: STRUCTURE FIRE
    Category:     STRUCTURE FIRE
    Address:      100 Queen St W, Toronto
    Timestamp:    2025-11-25 14:00:00
    Status:       active
    Coordinates:  ⚠️  Not available (needs geocoding)
    Description:  Fire crews responding to structure fire
```

## 🔗 Recursos

- **Nominatim (Geocoding):** https://nominatim.openstreetmap.org/
- **goquery (HTML parsing):** https://github.com/PuerkitoBio/goquery
- **Rate limiting:** El sistema respeta los límites de Nominatim (1 req/sec)

## ⚡ Tips de Debugging

### Ver logs detallados:
```bash
# Los logs del scraper incluyen emojis para fácil identificación
🔥 = TFS (Fire)
⚡ = Hydro (Power)
👮 = TPS (Police)
🚇 = TTC (Transit)
🌦️ = Weather
```

### Verificar estructura HTML de TFS:
```bash
curl -s "https://www.toronto.ca/community-people/public-safety-alerts/alerts-notifications/" | grep -i "incident\|fire\|active"
```

### Verificar endpoint de Hydro:
```bash
# En el navegador, abre DevTools y busca llamadas a APIs que contengan "outage" o "power"
```

## 🚨 Importante

- **NUNCA** hacer scraping agresivo (respetar rate limits)
- **SIEMPRE** incluir User-Agent apropiado
- **VERIFICAR** términos de servicio de cada fuente
- **USAR** mock data para desarrollo cuando los endpoints reales no estén disponibles
