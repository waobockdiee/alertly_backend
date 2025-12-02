# Bot Creator - Tipos de Incidentes por Fuente

Este documento detalla qué categorías y subcategorías de incidentes guardará cada cronjob del sistema bot_creator.

---

## 🚔 bot_creator_tps (Toronto Police Service)

**Source:** Toronto Police "Calls for Service"
**Frecuencia Sugerida:** Cada 5-10 minutos
**Estado:** ⏳ Placeholder (necesita URL real)

### Categorías que Guardará:

#### 1. Crime (Crímenes)
| Pattern Detectado           | Categoría Alertly | Subcategoría        | Prioridad |
|-----------------------------|-------------------|---------------------|-----------|
| assault, attack, battery    | `crime`           | `assault`           | 10        |
| robbery, theft, burglary    | `crime`           | `theft_burglary`    | 10        |
| shooting, gun, shots fired  | `crime`           | `shooting`          | 20        |
| stabbing, knife, weapon     | `crime`           | `stabbing`          | 15        |
| break enter, b&e            | `crime`           | `theft_burglary`    | 10        |
| homicide, murder, death     | `crime`           | `homicide`          | 30        |
| sexual, rape                | `crime`           | `sexual_assault`    | 25        |
| Generic crime               | `crime`           | `other`             | 1         |

#### 2. Suspicious Activity
| Pattern Detectado              | Categoría Alertly      | Subcategoría       |
|--------------------------------|------------------------|--------------------|
| suspicious, prowler, trespassing | `suspicious_activity` | `person_vehicle`   |

#### 3. Traffic Accidents
| Pattern Detectado           | Categoría Alertly    | Subcategoría           |
|-----------------------------|----------------------|------------------------|
| collision, accident, crash  | `traffic_accident`   | `car_accident`         |
| pedestrian struck, hit run  | `traffic_accident`   | `pedestrian_involved`  |

**Ejemplo de Mapeo:**
```
TPS Input: "ASSAULT - DOMESTIC"
→ Normalizado: crime / assault
→ Imagen: crime.webp
→ Título: "Police Call: ASSAULT - DOMESTIC"
```

---

## 🔥 bot_creator_tfs (Toronto Fire Services)

**Source:** TFS Active Incidents Page
**Frecuencia Sugerida:** Cada 10-15 minutos
**Estado:** ✅ Implementado (con mock data)

### Categorías que Guardará:

#### 1. Fire Incidents
| Pattern Detectado       | Categoría Alertly | Subcategoría      | Prioridad |
|-------------------------|-------------------|-------------------|-----------|
| fire, burning, smoke    | `fire_incident`   | `building_fire`   | 10        |
| vehicle fire            | `fire_incident`   | `vehicle_fire`    | 15        |
| alarm, detector         | `fire_incident`   | `alarm_activation`| 5         |
| hazmat, chemical, leak  | `fire_incident`   | `hazmat`          | 20        |
| rescue, trapped         | `fire_incident`   | `rescue`          | 12        |
| Generic fire            | `fire_incident`   | `other`           | 1         |

#### 2. Medical Emergencies
| Pattern Detectado              | Categoría Alertly      | Subcategoría |
|--------------------------------|------------------------|--------------|
| medical, emergency, trauma     | `medical_emergency`    | `trauma`     |
| overdose, poisoning            | `medical_emergency`    | `overdose`   |

**Ejemplo de Mapeo:**
```
TFS Input: "STRUCTURE FIRE - 100 Queen St W"
→ Normalizado: fire_incident / building_fire
→ Imagen: fire_incident.webp
→ Título: "Fire Service Call: STRUCTURE FIRE"
→ Geocoding: "100 Queen St W, Toronto" → (43.6529, -79.3849)
```

**Mock Data Actual:**
- Structure Fire @ 100 Queen St W
- Medical Call @ 200 Yonge St
- Vehicle Fire @ Gardiner Expressway & Spadina

---

## 🚇 bot_creator_ttc (TTC Transit)

**Source:** TTC Service Alerts RSS
**Frecuencia Sugerida:** Cada 15 minutos
**Estado:** ⏳ Placeholder (necesita verificar RSS)

### Categorías que Guardará:

#### 1. Infrastructure Issues (Principal)
| Pattern Detectado                 | Categoría Alertly        | Subcategoría      | Prioridad |
|-----------------------------------|--------------------------|-------------------|-----------|
| delay, service, suspended         | `infrastructure_issues`  | `transit`         | 10        |
| signal, mechanical, technical     | `infrastructure_issues`  | `transit`         | 8         |
| power, electrical, outage         | `infrastructure_issues`  | `utility_issues`  | 12        |
| Generic transit issue             | `infrastructure_issues`  | `transit`         | 1         |

#### 2. Medical Emergencies
| Pattern Detectado                | Categoría Alertly     | Subcategoría |
|----------------------------------|-----------------------|--------------|
| medical, emergency, passenger injury | `medical_emergency` | `trauma`     |

#### 3. Security
| Pattern Detectado              | Categoría Alertly      | Subcategoría     |
|--------------------------------|------------------------|------------------|
| security, police, investigation| `suspicious_activity`  | `person_vehicle` |

**Ejemplo de Mapeo:**
```
TTC Input: "Line 1: Service suspended - signal problems at Bloor"
→ Normalizado: infrastructure_issues / transit
→ Imagen: infrastructure_issues.webp
→ Título: "Transit Alert: Service suspended"
→ Coordenadas: Bloor Station (geocoded)
```

---

## ⚡ bot_creator_hydro (Toronto Hydro)

**Source:** Toronto Hydro Outage Map API
**Frecuencia Sugerida:** Cada 30 minutos
**Estado:** ✅ Implementado (con mock data)

### Categorías que Guardará:

#### 1. Infrastructure Issues (Único)
| Pattern Detectado           | Categoría Alertly       | Subcategoría      | Prioridad |
|-----------------------------|-------------------------|-------------------|-----------|
| outage, power out           | `infrastructure_issues` | `utility_issues`  | 10        |
| planned, maintenance        | `infrastructure_issues` | `utility_issues`  | 8         |
| unplanned, emergency, fault | `infrastructure_issues` | `utility_issues`  | 12        |
| Generic hydro               | `infrastructure_issues` | `utility_issues`  | 1         |

**Características Especiales:**
- **Polígonos:** Calcula centroide del área afectada
- **ETR:** Incluye Estimated Time of Restoration
- **Customers Affected:** Muestra número en descripción

**Ejemplo de Mapeo:**
```
Hydro Input: {
  "outageId": "OUT-12345",
  "status": "active",
  "cause": "Equipment failure",
  "customersAffected": 450,
  "polygon": [[lng, lat], ...]
}
→ Normalizado: infrastructure_issues / utility_issues
→ Imagen: infrastructure_issues.webp
→ Título: "Power Outage - 450 customers affected"
→ Coordenadas: Centroide del polígono calculado
→ Descripción: "Cause: Equipment failure | Status: active | Source: HYDRO"
```

**Mock Data Actual:**
- Downtown Toronto: 450 customers, equipment failure, ETR 2 hours
- Scarborough: 1200 customers, storm damage, ETR 4 hours

---

## 🌦️ bot_creator_weather (Environment Canada)

**Source:** Environment Canada CAP Alerts
**Frecuencia Sugerida:** Cada 1 hora
**Estado:** ⏳ Placeholder (necesita URL CAP)

### Categorías que Guardará:

#### 1. Extreme Weather (Único)
| Pattern Detectado                | Categoría Alerty    | Subcategoría    | Prioridad |
|----------------------------------|---------------------|-----------------|-----------|
| tornado, funnel cloud            | `extreme_weather`   | `tornado`       | 30        |
| snow storm, blizzard             | `extreme_weather`   | `snow_storms`   | 20        |
| thunderstorm, lightning          | `extreme_weather`   | `thunderstorm`  | 15        |
| flood, flooding, heavy rain      | `extreme_weather`   | `flood`         | 18        |
| heat, extreme temp               | `extreme_weather`   | `heat_wave`     | 12        |
| cold, freeze, wind chill         | `extreme_weather`   | `cold_wave`     | 12        |
| wind, gale, storm wind           | `extreme_weather`   | `high_winds`    | 10        |
| fog, visibility                  | `extreme_weather`   | `fog`           | 5         |
| Generic weather                  | `extreme_weather`   | `other`         | 1         |

**Características Especiales:**
- **Regional:** Usa coordenadas centrales de Toronto (43.6532, -79.3832) si la alerta es para toda la ciudad
- **Duration:** Incluye inicio y fin de la alerta

**Ejemplo de Mapeo:**
```
Weather Input: "Severe Thunderstorm Warning - City of Toronto"
→ Normalizado: extreme_weather / thunderstorm
→ Imagen: extreme_weather.webp
→ Título: "Weather Alert: Severe Thunderstorm Warning"
→ Coordenadas: (43.6532, -79.3832) [Toronto center]
→ Descripción: "[Alert details] | Source: WEATHER"
```

---

## 📊 Resumen de Categorías por Cronjob

| Cronjob              | Categorías Principales                                    | Subcategorías Totales |
|----------------------|-----------------------------------------------------------|-----------------------|
| `bot_creator_tps`    | Crime, Suspicious Activity, Traffic Accident              | ~10 subcategorías     |
| `bot_creator_tfs`    | Fire Incident, Medical Emergency                          | ~7 subcategorías      |
| `bot_creator_ttc`    | Infrastructure Issues, Medical Emergency, Suspicious Act. | ~4 subcategorías      |
| `bot_creator_hydro`  | Infrastructure Issues (solo utility)                      | 1 subcategoría        |
| `bot_creator_weather`| Extreme Weather                                           | ~9 subcategorías      |

---

## 🎨 Mapeo de Imágenes por Categoría

| Categoría Alertly            | Imagen S3                                  | Fuentes que la Usan     |
|------------------------------|--------------------------------------------|-------------------------|
| `crime`                      | crime.webp                                 | TPS                     |
| `traffic_accident`           | traffic_accident.webp                      | TPS                     |
| `medical_emergency`          | medical_emergency.webp                     | TFS, TTC                |
| `fire_incident`              | fire_incident.webp                         | TFS                     |
| `suspicious_activity`        | suspicious_activity.webp                   | TPS, TTC                |
| `infrastructure_issues`      | infrastructure_issues.webp                 | TTC, Hydro              |
| `extreme_weather`            | extreme_weather.webp                       | Weather                 |
| `vandalism` ⚠️               | vandalism.webp *(falta)*                   | -                       |
| `community_events` ⚠️        | community_events.webp *(falta)*            | -                       |
| `positive_actions` ⚠️        | positive_actions.webp *(falta)*            | -                       |
| `lost_pet` ⚠️                | lost_pet.webp *(falta)*                    | -                       |
| `dangerous_wildlife_sighting`| dangerous_wildlife_sighting.webp           | -                       |

**Nota:** Las 4 categorías marcadas con ⚠️ no son usadas por ningún cronjob actual, pero están en el sistema de Alertly para reportes de usuarios.

---

## 🔍 Ejemplos de Incidentes Completos

### Ejemplo 1: TPS Crime
```json
{
  "source": "tps",
  "externalID": "tps_call_12345",
  "rawTitle": "ROBBERY - ARMED",
  "rawCategory": "ROBBERY",
  "address": "456 Yonge St, Toronto",
  "timestamp": "2025-11-23T14:30:00Z",

  // Después de normalización:
  "categoryCode": "crime",
  "subcategoryCode": "theft_burglary",
  "imageURL": "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/crime.webp",
  "title": "Police Call: ROBBERY - ARMED",
  "latitude": 43.6560,
  "longitude": -79.3802,
  "eventType": "public"
}
```

### Ejemplo 2: TFS Fire
```json
{
  "source": "tfs",
  "externalID": "tfs_100_queen_st_1732373400",
  "rawTitle": "Structure Fire",
  "rawCategory": "STRUCTURE FIRE",
  "address": "100 Queen St W, Toronto",
  "timestamp": "2025-11-23T14:30:00Z",

  // Después de normalización:
  "categoryCode": "fire_incident",
  "subcategoryCode": "building_fire",
  "imageURL": "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/fire_incident.webp",
  "title": "Fire Service Call: Structure Fire",
  "latitude": 43.6529,
  "longitude": -79.3849,
  "eventType": "public"
}
```

### Ejemplo 3: Hydro Outage
```json
{
  "source": "hydro",
  "externalID": "hydro_outage_001",
  "rawTitle": "Power Outage - 450 customers affected",
  "rawCategory": "unplanned outage",
  "polygon": [
    {"lat": 43.6500, "lng": -79.3800},
    {"lat": 43.6520, "lng": -79.3750}
  ],
  "ETR": "2025-11-23T16:30:00Z",

  // Después de normalización:
  "categoryCode": "infrastructure_issues",
  "subcategoryCode": "utility_issues",
  "imageURL": "https://alertly-images-production.s3.us-west-2.amazonaws.com/incidents/infrastructure_issues.webp",
  "title": "Power Outage - 450 customers affected",
  "latitude": 43.6510,  // Centroide calculado
  "longitude": -79.3775,
  "description": "Cause: Equipment failure | Location: Downtown Toronto | ETR: 4:30 PM | Source: HYDRO",
  "eventType": "public"
}
```

---

## 💡 Notas Importantes

### Sobre las Categorías Faltantes

Las siguientes categorías de Alertly **NO** son usadas por el sistema de bot_creator:
- `vandalism` - Solo reportes de usuarios
- `community_events` - Solo reportes de usuarios
- `positive_actions` - Solo reportes de usuarios
- `lost_pet` - Solo reportes de usuarios

**¿Deberías crear las imágenes?**
- ✅ **Sí**, porque los usuarios pueden reportar estos tipos de incidentes manualmente
- ❌ No las necesitas para que funcione el bot_creator

### Sobre Duplicados

El sistema previene duplicados usando:
```go
SHA256(source + external_id + timestamp)
```

**Ejemplo:**
- TPS reporta "ASSAULT" a las 2:00 PM → Hash: `a1b2c3...`
- TPS reporta el mismo "ASSAULT" a las 2:05 PM → Hash: `d4e5f6...` (diferente timestamp)
- ✅ Se guardan ambos (son eventos separados en el tiempo)

### Sobre TTL (Expiración)

Cada tipo de incidente tiene diferente tiempo de vida:
- **Fire incidents:** 6 horas
- **Power outages:** 24 horas (mientras estén activos)
- **Weather alerts:** 12 horas
- **Crime/Traffic:** 6 horas (default)

---

**📅 Creado:** 23 de Noviembre, 2025
**🔄 Basado en:** `normalizer.go` implementación actual
**✅ Estado:** Mappings completos para 5 fuentes
