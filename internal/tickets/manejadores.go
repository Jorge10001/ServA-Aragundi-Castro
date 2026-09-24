
package tickets

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	// Reemplaza 'parqueo-eventos' por el nombre de tu módulo en go.mod
	"github.com/Jorge10001/servA-Castro-Aragundi/internal/respuesta"
)

type Manejador struct {
	DB *gorm.DB
}

func (m *Manejador) Rutas(r chi.Router) {
	r.Post("/tickets", m.crear)
	r.Get("/tickets", m.listar)
	r.Get("/tickets/{id}", m.verUno)
	r.Put("/tickets/{id}", m.actualizar)
	r.Delete("/tickets/{id}", m.borrar)

	// Fase 2c: Listado de la entidad del lado del uno con Preload
	r.Get("/eventos", m.listarEventos)
}

func leerID(w http.ResponseWriter, r *http.Request) (uint, bool) {
	n, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || n <= 0 {
		respuesta.Error(w, http.StatusBadRequest, "id_invalido", "El id debe ser un número positivo")
		return 0, false
	}
	return uint(n), true
}

// 1. CREAR (POST /tickets)
func (m *Manejador) crear(w http.ResponseWriter, r *http.Request) {
	var t Ticket
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "json_invalido", "El cuerpo no es un JSON válido")
		return
	}
	t.ID = 0

	// Validar estado (422)
	if !estadosValidos[t.Estado] {
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado_invalido", "Estado no válido")
		return
	}

	// Regla extra de negocio (422): Placa, Espacio y EventoID requeridos
	if t.Placa == "" || t.Espacio == "" || t.EventoID == 0 {
		respuesta.Error(w, http.StatusUnprocessableEntity, "datos_incompletos", "La placa, el espacio y el EventoID son obligatorios")
		return
	}

	if err := m.DB.Debug().Create(&t).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "No se pudo guardar el ticket")
		return
	}

	respuesta.Exito(w, http.StatusCreated, t)
}

// 2. LISTAR Y FILTRAR SEGURO (GET /tickets?estado=...)
func (m *Manejador) listar(w http.ResponseWriter, r *http.Request) {
	var list []Ticket
	query := m.DB.Debug().Preload("Evento")

	estado := r.URL.Query().Get("estado")
	if estado != "" {
		// Filtrado seguro usando ? contra SQL Injection (Fase 2d)
		query = query.Where("estado = ?", estado)
	}

	if err := query.Find(&list).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "Error al consultar los tickets")
		return
	}

	respuesta.Exito(w, http.StatusOK, list)
}

// 3. VER UNO (GET /tickets/{id})
func (m *Manejador) verUno(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}

	var t Ticket
	err := m.DB.Debug().Preload("Evento").First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respuesta.Error(w, http.StatusNotFound, "no_encontrado", "El ticket no existe")
		return
	} else if err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "Error al buscar")
		return
	}

	respuesta.Exito(w, http.StatusOK, t)
}

// 4. ACTUALIZAR (PUT /tickets/{id})
func (m *Manejador) actualizar(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}

	var existente Ticket
	err := m.DB.First(&existente, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respuesta.Error(w, http.StatusNotFound, "no_encontrado", "El ticket no existe")
		return
	}

	var entrada Ticket
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		respuesta.Error(w, http.StatusBadRequest, "json_invalido", "El cuerpo no es un JSON válido")
		return
	}

	if !estadosValidos[entrada.Estado] {
		respuesta.Error(w, http.StatusUnprocessableEntity, "estado_invalido", "Estado no válido")
		return
	}

	if entrada.Placa == "" || entrada.Espacio == "" {
		respuesta.Error(w, http.StatusUnprocessableEntity, "datos_incompletos", "La placa y espacio son obligatorios")
		return
	}

	existente.Estado = entrada.Estado
	existente.Placa = entrada.Placa
	existente.Espacio = entrada.Espacio
	existente.Propietario = entrada.Propietario

	if err := m.DB.Debug().Save(&existente).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "Error al actualizar")
		return
	}

	respuesta.Exito(w, http.StatusOK, existente)
}

// 5. BORRAR (DELETE /tickets/{id})
func (m *Manejador) borrar(w http.ResponseWriter, r *http.Request) {
	id, ok := leerID(w, r)
	if !ok {
		return
	}

	res := m.DB.Debug().Delete(&Ticket{}, id)
	if res.Error != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "Error al eliminar")
		return
	}

	if res.RowsAffected == 0 {
		respuesta.Error(w, http.StatusNotFound, "no_encontrado", "El ticket no existe")
		return
	}

	respuesta.Exito(w, http.StatusOK, map[string]string{"mensaje": "Registro eliminado correctamente"})
}

// 6. LISTAR LADO DEL UNO CON PRELOAD (GET /eventos) - Cumple Fase 2c sin problema N+1
func (m *Manejador) listarEventos(w http.ResponseWriter, r *http.Request) {
	var eventos []Evento
	if err := m.DB.Debug().Preload("Tickets").Find(&eventos).Error; err != nil {
		respuesta.Error(w, http.StatusInternalServerError, "error_base", "Error al consultar eventos")
		return
	}
	respuesta.Exito(w, http.StatusOK, eventos)
}