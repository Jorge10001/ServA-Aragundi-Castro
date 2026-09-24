package tickets

import "time"

// Evento: Lado del Uno
type Evento struct {
	ID      uint      `gorm:"primaryKey" json:"id"`
	Nombre  string    `gorm:"not null" json:"nombre"`
	Fecha   time.Time `gorm:"not null" json:"fecha"`
	Lugar   string    `gorm:"not null" json:"lugar"`
	Tickets []Ticket  `gorm:"foreignKey:EventoID" json:"tickets,omitempty"`
}

// Ticket: Lado de los Muchos (Entidad que lleva los estados)
type Ticket struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	EventoID    uint      `gorm:"not null" json:"evento_id"`
	Placa       string    `gorm:"not null;size:10" json:"placa"`
	Espacio     string    `gorm:"not null;size:10" json:"espacio"` // Asignado previamente (ej. A-12)
	Propietario string    `gorm:"not null" json:"propietario"`
	Estado      string    `gorm:"type:varchar(20);default:'RESERVADO'" json:"estado"`
	Evento      *Evento   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"evento,omitempty"`
}

// Map de estados válidos para la Fase 2b
var estadosValidos = map[string]bool{
	"RESERVADO": true,
	"INGRESADO": true,
	"CANCELADO": true,
}