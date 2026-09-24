
package tickets

import (
	"time"
	"gorm.io/gorm"
)

func Sembrar(db *gorm.DB) {
	var total int64
	db.Model(&Ticket{}).Count(&total)
	if total > 0 {
		return
	}

	eventoSample := Evento{
		Nombre: "Concierto Rock 2026",
		Fecha:  time.Now().Add(24 * time.Hour),
		Lugar:  "Estadio Central",
		Tickets: []Ticket{
			{
				Placa:       "PBA-1234",
				Espacio:     "A-01",
				Propietario: "Carlos Pérez",
				Estado:      "RESERVADO",
			},
			{
				Placa:       "GBA-5678",
				Espacio:     "A-02",
				Propietario: "María López",
				Estado:      "INGRESADO",
			},
		},
	}

	db.Create(&eventoSample)
}