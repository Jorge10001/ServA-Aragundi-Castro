package tickets

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"mesa-ayuda/internal/respuesta"
)

// Manejador agrupa las rutas de tickets y lleva la conexión adentro.
type Manejador struct {
	DB *gorm.DB
}

// Rutas registra las rutas de tickets en el enrutador.
func (m *Manejador) Rutas(r chi.Router) {
	r.Get("/tickets", m.listar)
	r.Post("/tickets", m.crear)
	r.Get("/tickets/{id}", m.verUno)
}

// leerID convierte el {id} de la ruta; si no es un número devuelve false.
func leerID(r *http.Request) (uint, bool) {
	n, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || n <= 0 {
		return 0, false
	}
	return uint(n), true
}

// crear: 400 si el JSON está roto · 422 si un dato rompe una regla · 201 si se guardó.
// Las dos primeras salidas ocurren ANTES de tocar la base de datos.
func (m *Manejador) crear(w http.ResponseWriter, r *http.Request) {
	var t Ticket
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "el cuerpo no es JSON válido")
		return
	}
	if strings.TrimSpace(t.Asunto) == "" {
		respuesta.Error(w, http.StatusUnprocessableEntity, "el asunto es obligatorio")
		return
	}
	if !EstadosValidos[t.Estado] {
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado no permitido: "+t.Estado)
		return
	}
	// ---- desde aquí hace falta la base de datos ----
	if err := m.DB.Create(&t).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "no se pudo guardar")
		return
	}
	respuesta.JSON(w, http.StatusCreated, t)
}

// verUno: 400 si el id no es un número · 404 si no existe · 200 si existe.
func (m *Manejador) verUno(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(r)
	if !ok {
		respuesta.Error(w, http.StatusBadRequest, "el id debe ser un número entero positivo")
		return
	}
	// ---- desde aquí hace falta la base de datos ----
	var t Ticket
	err := m.DB.Preload("Comentarios").First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respuesta.Error(w, http.StatusNotFound, "no existe el ticket")
		return
	}
	if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "no se pudo leer")
		return
	}
	respuesta.JSON(w, http.StatusOK, t)
}

// listar: ?estado= filtra con parámetro (nunca pegando texto en la consulta).
func (m *Manejador) listar(w http.ResponseWriter, r *http.Request) {
	var lista []Ticket
	consulta := m.DB.Preload("Comentarios")
	if estado := r.URL.Query().Get("estado"); estado != "" {
		consulta = consulta.Where("estado = ?", estado)
	}
	if err := consulta.Find(&lista).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "no se pudo listar")
		return
	}
	respuesta.JSON(w, http.StatusOK, lista)
}
