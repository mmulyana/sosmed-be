package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	// Load file .env biar DB_URL bisa dipakai
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("❌ DB_URL tidak ditemukan di .env")
	}

	m, err := migrate.New(
		"file://cmd/migrate/migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("❌ Gagal init migrasi: %v", err)
	}

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("❌ Gagal migrate up: %v", err)
		}
		log.Println("✅ Migration UP selesai!")

	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("❌ Gagal rollback 1 langkah: %v", err)
		}
		log.Println("⬇️  Migration DOWN 1 langkah selesai!")

	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("❌ Gagal drop semua tabel: %v", err)
		}
		log.Println("🔥 Semua tabel dihapus!")

	default:
		fmt.Println("Gunakan salah satu perintah berikut:")
		fmt.Println("  go run cmd/migrate/main.go up     # Jalankan semua migration")
		fmt.Println("  go run cmd/migrate/main.go down   # Rollback 1 langkah")
		fmt.Println("  go run cmd/migrate/main.go drop   # Hapus semua tabel")
	}
}
