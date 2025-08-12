#!/bin/bash

# Script para probar acceso a imágenes
echo "🧪 Probando acceso a imágenes..."

# Verificar que el servidor esté corriendo
if curl -s http://localhost:8080/category/get_all > /dev/null; then
    echo "✅ Servidor está corriendo"
else
    echo "❌ Servidor no está corriendo en puerto 8080"
    exit 1
fi

# Verificar que la carpeta uploads existe
if [ -d "uploads" ]; then
    echo "✅ Carpeta uploads existe"
else
    echo "❌ Carpeta uploads no existe"
    exit 1
fi

# Listar archivos en uploads
echo "📁 Archivos en carpeta uploads:"
ls -la uploads/ | head -10

# Probar acceso a un archivo específico (si existe)
if [ -f "uploads/alerty_1740606909644165000.webp" ]; then
    echo "🔗 Probando acceso a imagen existente..."
    if curl -s -I http://localhost:8080/uploads/alerty_1740606909644165000.webp | grep "200 OK" > /dev/null; then
        echo "✅ Imagen accesible via HTTP"
    else
        echo "❌ Imagen no accesible via HTTP"
    fi
else
    echo "ℹ️ No hay imágenes de prueba disponibles"
fi

echo "🎯 Para probar completamente:"
echo "1. Crea un incidente desde la app"
echo "2. Verifica que la imagen se guarde en uploads/"
echo "3. Verifica que sea accesible via http://localhost:8080/uploads/nombre_archivo"
