#!/bin/bash

echo "🔍 Diagnóstico de Imágenes - Alertly Backend"
echo "=============================================="

# 1. Verificar que el servidor esté corriendo
echo "1. Verificando servidor..."
if curl -s http://localhost:8080/category/get_all > /dev/null; then
    echo "   ✅ Servidor corriendo en puerto 8080"
else
    echo "   ❌ Servidor no está corriendo"
    echo "   💡 Ejecuta: go run cmd/app/main.go"
    exit 1
fi

# 2. Verificar carpeta uploads
echo "2. Verificando carpeta uploads..."
if [ -d "uploads" ]; then
    echo "   ✅ Carpeta uploads existe"
    echo "   📁 Contenido:"
    ls -la uploads/ | head -5
else
    echo "   ❌ Carpeta uploads no existe"
    echo "   💡 Creando carpeta..."
    mkdir -p uploads
    echo "   ✅ Carpeta creada"
fi

# 3. Verificar configuración estática
echo "3. Verificando configuración estática..."
if curl -s -I http://localhost:8080/uploads/ | grep "200 OK\|404 Not Found" > /dev/null; then
    echo "   ✅ Endpoint /uploads/ responde"
else
    echo "   ❌ Endpoint /uploads/ no responde"
fi

# 4. Probar con una imagen existente
echo "4. Probando acceso a imagen existente..."
if [ -f "uploads/alerty_1740606909644165000.webp" ]; then
    echo "   📸 Imagen de prueba encontrada"
    if curl -s -I http://localhost:8080/uploads/alerty_1740606909644165000.webp | grep "200 OK" > /dev/null; then
        echo "   ✅ Imagen accesible via HTTP"
    else
        echo "   ❌ Imagen no accesible via HTTP"
        echo "   🔍 Headers de respuesta:"
        curl -s -I http://localhost:8080/uploads/alerty_1740606909644165000.webp
    fi
else
    echo "   ℹ️ No hay imágenes de prueba"
fi

# 5. Verificar logs del servidor
echo "5. Verificando configuración del servidor..."
echo "   📋 Busca en los logs del servidor:"
echo "   'Serving uploads from: /path/to/uploads'"

echo ""
echo "🎯 Pasos para probar:"
echo "1. Reinicia el servidor: go run cmd/app/main.go"
echo "2. Crea un incidente desde la app"
echo "3. Verifica que la imagen se guarde en uploads/"
echo "4. Verifica que sea accesible via: http://localhost:8080/uploads/nombre_archivo"
echo ""
echo "🔧 Si no funciona:"
echo "- Verifica que el servidor muestre 'Serving uploads from: /path/to/uploads'"
echo "- Verifica que las imágenes se guarden en la carpeta correcta"
echo "- Verifica que las rutas en la BD sean '/uploads/filename'"
