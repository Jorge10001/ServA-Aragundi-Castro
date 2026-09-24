
#!/usr/bin/env bash
# Uso: ./pruebas.sh http://localhost:8081

B=${1:-http://localhost:8081}

c() { echo; echo "### $1 [esperado: $2]"; shift 2; curl -s -i "$@" | sed -n '1p;/^{/p'; }

c "POST válido" "201" -X POST $B/tickets -d '{"titulo": "Impresora sin tóner", "prioridad":"alta"}'
c "GET lista" "200" $B/tickets
c "GET por id" "200" $B/tickets/1
c "GET id inexistente" "404" $B/tickets/99
c "GET id no numérico" "400" $B/tickets/abc
c "POST JSON roto" "400" -X POST $B/tickets -d '{"titulo": "x'
c "POST título corto" "422" -X POST $B/tickets -d '{"titulo": "abc", "prioridad":"alta"}'
c "POST prioridad mala" "422" -X POST $B/tickets -d '{"titulo": "Teclado dañado", "prioridad":"urgente"}'
c "DELETE no soportado" "405" -X DELETE $B/tickets/1
c "ruta inexistente" "404" $B/usuarios
c "pánico controlado" "500" $B/explotar
c "campo desconocido" "400" -X POST $B/tickets -d '{"titulo": "Pantalla rota", "prioridad":"alta", "estado": "cerrado"}'
c "el servidor sigue vivo" "200" $B/tickets