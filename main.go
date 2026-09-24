package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mesa-ayuda/internal/tickets"
)

const claveFirma = "uleam-2026-2-mesa-de-ayuda" // firmará los tokens en la unidad 2

func main() {
	dsn := "host=localhost user=postgres password=Secreta123 dbname=mesa_ayuda_demo port=5433"
	puerto := ":8080"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("no se pudo conectar: ", err)
	}
	if err := db.AutoMigrate(&tickets.Ticket{}, &tickets.Comentario{}); err != nil {
		log.Fatal("no se pudo migrar: ", err)
	}
	if err := tickets.Sembrar(db); err != nil {
		log.Fatal("no se pudo sembrar: ", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))
	(&tickets.Manejador{DB: db}).Rutas(r)

	servidor := &http.Server{
		Addr:         puerto,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Println("escuchando en", puerto)
	log.Fatal(servidor.ListenAndServe())
}
