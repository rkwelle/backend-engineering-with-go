package main

import (
	"log"

	"github.com/rkwelle/social-app/internal/db"
	"github.com/rkwelle/social-app/internal/env"
	"github.com/rkwelle/social-app/internal/store"
)

func main() {

	// fmt.Println("DB_ADDR = ", env.GetString("DB_ADDR", "fallback not used"))

	addr := env.GetString("DB_ADDR", "postgres://postgres:gentle@localhost:5432/social_app?sslmode=disable")

	// fmt.Printf("DB_ADDR = %q\n", addr)

	conn, err := db.New(addr, 3, 3, "15m")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	store := store.NewStorage(conn)
	db.Seed(store)
}
