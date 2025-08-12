#!/bin/bash

echo "🌐 Test de Conectividad - Frontend ↔ Backend"
echo "=============================================="

# Configuración
BACKEND_IP="192.168.1.66"
BACKEND_PORT="8080"
BASE_URL="http://${BACKEND_IP}:${BACKEND_PORT}"

echo "📋 Configuración actual:"
echo "   Backend IP: ${BACKEND_IP}"
echo "   Backend Port: ${BACKEND_PORT}"
echo "   Base URL: ${BASE_URL}"
echo ""

# 1. Verificar que el servidor esté corriendo
echo "1. Verificando servidor backend..."
if curl -s --connect-timeout 5 "${BASE_URL}/category/get_all" > /dev/null; then
    echo "   ✅ Backend responde desde ${BACKEND_IP}"
else
    echo "   ❌ Backend no responde desde ${BACKEND_IP}"
    echo "   💡 Posibles problemas:"
    echo "      - Servidor no está corriendo"
    echo "      - Firewall bloqueando conexión"
    echo "      - IP incorrecta"
fi

# 2. Probar desde localhost
echo "2. Verificando servidor desde localhost..."
if curl -s --connect-timeout 5 "http://localhost:8080/category/get_all" > /dev/null; then
    echo "   ✅ Backend responde desde localhost"
else
    echo "   ❌ Backend no responde desde localhost"
fi

# 3. Verificar IP de la máquina
echo "3. Verificando IP de la máquina..."
CURRENT_IP=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1)
echo "   📱 IP actual: ${CURRENT_IP}"

# 4. Probar acceso a imágenes
echo "4. Verificando acceso a imágenes..."
if [ -f "uploads/alerty_1740606909644165000.webp" ]; then
    echo "   📸 Probando imagen existente..."
    
    # Probar desde IP configurada
    if curl -s -I "${BASE_URL}/uploads/alerty_1740606909644165000.webp" | grep "200 OK" > /dev/null; then
        echo "   ✅ Imagen accesible desde ${BACKEND_IP}"
    else
        echo "   ❌ Imagen no accesible desde ${BACKEND_IP}"
    fi
    
    # Probar desde localhost
    if curl -s -I "http://localhost:8080/uploads/alerty_1740606909644165000.webp" | grep "200 OK" > /dev/null; then
        echo "   ✅ Imagen accesible desde localhost"
    else
        echo "   ❌ Imagen no accesible desde localhost"
    fi
else
    echo "   ℹ️ No hay imágenes de prueba"
fi

echo ""
echo "🔧 Soluciones posibles:"
echo "1. Si backend responde desde localhost pero no desde IP:"
echo "   - Cambiar IP_SERVER en .env a 0.0.0.0"
echo "   - Verificar firewall"
echo ""
echo "2. Si backend no responde:"
echo "   - Reiniciar servidor: go run cmd/app/main.go"
echo "   - Verificar logs del servidor"
echo ""
echo "3. Si imágenes no accesibles:"
echo "   - Verificar que servidor muestre 'Serving uploads from: /path/to/uploads'"
echo "   - Verificar permisos de carpeta uploads"
