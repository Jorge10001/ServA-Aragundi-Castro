
package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/Jorge10001/servA-Castro-Aragundi/internal/respuesta" // Ajusta con el nombre de tu módulo
)

type grabador struct {
	http.ResponseWriter
	estado int
}

func (g *grabador) WriteHeader(codigo int) {
	g.estado = codigo
	g.ResponseWriter.WriteHeader(codigo)
}

// Registro escribe una línea en el log por cada petición
func Registro(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		g := &grabador{ResponseWriter: w, estado: http.StatusOK}
		
		next.ServeHTTP(g, r) // Pasa el grabador en vez de w
		
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, g.estado, time.Since(inicio))
	})
}

// Recuperacion atrapa pánicos y evita que el servidor se caiga
func Recuperacion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("PÁNICO en %s %s: %v", r.Method, r.URL.Path, p)
				respuesta.Error(w, http.StatusInternalServerError, "error_interno", "ocurrió un error inesperado")
			}
		}()
		next.ServeHTTP(w, r)
	})
}