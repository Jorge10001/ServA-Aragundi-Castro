package respuesta

import (
	"encoding/json"
	"net/http"
)

type detalleError struct {
	Codigo  string `json:"codigo"`
	Mensaje string `json:"mensaje"`
}

type envoltura struct {
	OK    bool          `json:"ok"`
	Datos any           `json:"datos,omitempty"`
	Error *detalleError `json:"error,omitempty"`
}

func escribir(w http.ResponseWriter, estado int, cuerpo envoltura) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(estado)
	_ = json.NewEncoder(w).Encode(cuerpo)
}

// Exito responde con estado HTTP y los datos envueltos
func Exito(w http.ResponseWriter, estado int, datos any) {
	escribir(w, estado, envoltura{OK: true, Datos: datos})
}

// Error responde con estado HTTP y el detalle del error
func Error(w http.ResponseWriter, estado int, codigo, mensaje string) {
	escribir(w, estado, envoltura{
		OK: false,
		Error: &detalleError{
			Codigo:  codigo,
			Mensaje: mensaje,
		},
	})
}