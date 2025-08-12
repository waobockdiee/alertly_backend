# Configuración de Producción - Alertly Backend

## 🚀 **Configuración de Imágenes para Producción**

### **Problema**:
En producción, el backend es una API y no debe servir archivos estáticos directamente. Las imágenes deben ser servidas desde un CDN o servidor web separado.

### **Solución Implementada**:
- ✅ **Configuración centralizada**: `common.GetImageURL()`
- ✅ **Variables de entorno**: `IMAGE_BASE_URL`
- ✅ **Fallback para desarrollo**: `http://192.168.1.66:8080`
- ✅ **Producción**: `https://cdn.alertly.ca`

## 📋 **Configuración de Variables de Entorno**

### **Desarrollo (.env)**:
```bash
NODE_ENV=development
IMAGE_BASE_URL=http://192.168.1.66:8080
```

### **Producción (.env.production)**:
```bash
NODE_ENV=production
IMAGE_BASE_URL=https://cdn.alertly.ca
```

## 🔧 **Opciones de CDN para Producción**

### **1. AWS S3 + CloudFront**:
```bash
IMAGE_BASE_URL=https://d1234567890.cloudfront.net
```

### **2. Cloudflare**:
```bash
IMAGE_BASE_URL=https://cdn.alertly.ca
```

### **3. Google Cloud Storage**:
```bash
IMAGE_BASE_URL=https://storage.googleapis.com/alertly-images
```

### **4. Azure Blob Storage**:
```bash
IMAGE_BASE_URL=https://alertly.blob.core.windows.net/images
```

## 📁 **Estructura de Archivos en CDN**

```
https://cdn.alertly.ca/
├── uploads/
│   ├── alerty_1754681013820625000.webp
│   ├── alerty_1754681013820625001.jpg
│   └── ...
```

## 🔄 **Flujo de Subida de Imágenes**

### **Desarrollo**:
1. Usuario sube imagen → API la guarda en `uploads/`
2. API sirve imagen desde `http://192.168.1.66:8080/uploads/`
3. Frontend accede directamente

### **Producción**:
1. Usuario sube imagen → API la guarda en `uploads/`
2. API sube imagen a CDN (S3, CloudFront, etc.)
3. API guarda URL del CDN en BD
4. Frontend accede desde CDN

## 🛠️ **Implementación de CDN**

### **Opción 1: AWS S3 + CloudFront**
```bash
# Instalar AWS CLI
aws s3 sync uploads/ s3://alertly-images/
aws cloudfront create-invalidation --distribution-id E1234567890 --paths "/*"
```

### **Opción 2: Script de Subida Automática**
```bash
#!/bin/bash
# Subir imágenes a CDN automáticamente
aws s3 sync uploads/ s3://alertly-images/ --delete
```

## 📊 **Beneficios de CDN**

- ✅ **Performance**: Imágenes servidas desde edge locations
- ✅ **Escalabilidad**: Manejo de tráfico global
- ✅ **Seguridad**: API separada de archivos estáticos
- ✅ **Costo**: Reducción de carga en servidor principal

## 🔍 **Monitoreo**

### **Logs Importantes**:
```
✅ Image uploaded to CDN: https://cdn.alertly.ca/uploads/alerty_123.jpg
⚠️ CDN upload failed: timeout
✅ Image URL updated in database
```

### **Métricas a Monitorear**:
- Tiempo de subida a CDN
- Tasa de éxito de subidas
- Tiempo de respuesta del CDN
- Uso de ancho de banda

## 🚨 **Consideraciones de Seguridad**

1. **CORS**: Configurar CORS en CDN para permitir acceso desde app
2. **Autenticación**: Considerar autenticación para imágenes privadas
3. **Rate Limiting**: Limitar subidas de imágenes
4. **Validación**: Validar tipos y tamaños de archivo

## 📞 **Soporte**

Para configurar CDN en producción:
1. Elegir proveedor (AWS, Cloudflare, etc.)
2. Configurar bucket/dominio
3. Actualizar `IMAGE_BASE_URL` en `.env.production`
4. Implementar script de subida automática
5. Probar acceso desde app móvil
