
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// Ajusta según el nombre de tu módulo en go.mod
	"github.com/Jorge10001/servA-Castro-Aragundi/internal/tickets"
)

func main() {
	reset := flag.Bool("reset", false, "borra las tablas y arranca con la base vacía")
	flag.Parse()

	dsn := "host=localhost port=5432 user=admin password=secreto dbname=parqueo_eventos sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error de conexión:", err)
	}

	if *reset {
		db.Migrator().DropTable(&tickets.Ticket{}, &tickets.Evento{})
	}

	err = db.Debug().AutoMigrate(&tickets.Evento{}, &tickets.Ticket{})
	if err != nil {
		log.Fatal("Error en AutoMigrate:", err)
	}

	tickets.Sembrar(db)

	r := chi.NewRouter()

	(&tickets.Manejador{DB: db}).Rutas(r)

	log.Println("Servidor iniciado en puerto :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}