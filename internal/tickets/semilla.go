package tickets

import "gorm.io/gorm"

// Sembrar deja tres tickets si la tabla está vacía.
func Sembrar(db *gorm.DB) error {
	var n int64
	db.Model(&Ticket{}).Count(&n)
	if n > 0 {
		return nil
	}
	return db.Create(&[]Ticket{
		{Asunto: "No enciende el proyector del aula 309", Estado: "abierto",
			Comentarios: []Comentario{{Texto: "Ya se avisó a mantenimiento"}}},
		{Asunto: "Sin internet en el laboratorio 2", Estado: "en_proceso"},
		{Asunto: "Contraseña del aula virtual vencida", Estado: "cerrado"},
	}).Error
}
