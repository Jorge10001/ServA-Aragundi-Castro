package tickets

// Ticket es la entidad con estados de la mesa de ayuda.
type Ticket struct {
	ID          uint         `json:"id"`
	Asunto      string       `json:"asunto"`
	Estado      string       `json:"estado"`
	Comentarios []Comentario `json:"comentarios,omitempty"`
}

// Comentario pertenece a un ticket (uno a muchos).
type Comentario struct {
	ID       uint   `json:"id"`
	TicketID uint   `json:"ticket_id"`
	Texto    string `json:"texto"`
}

// EstadosValidos es la barrera: un estado que no está aquí responde 422.
var EstadosValidos = map[string]bool{
	"abierto":    true,
	"en_proceso": true,
	"cerrado":    true,
}
