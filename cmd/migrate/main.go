package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("❌ DB_URL tidak ditemukan di .env")
	}

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	migrationsDir := "cmd/migrate/migrations"

	switch cmd {
	case "up":
		m, err := migrate.New("file://"+migrationsDir, dbURL)
		if err != nil {
			log.Fatalf("❌ Gagal init migrasi: %v", err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("❌ Gagal migrate up: %v", err)
		}
		log.Println("✅ Migration UP selesai!")

	case "down":
		m, err := migrate.New("file://"+migrationsDir, dbURL)
		if err != nil {
			log.Fatalf("❌ Gagal init migrasi: %v", err)
		}
		if err := m.Steps(-1); err != nil {
			log.Fatalf("❌ Gagal rollback 1 langkah: %v", err)
		}
		log.Println("⬇️  Migration DOWN 1 langkah selesai!")

	case "drop":
		m, err := migrate.New("file://"+migrationsDir, dbURL)
		if err != nil {
			log.Fatalf("❌ Gagal init migrasi: %v", err)
		}
		if err := m.Drop(); err != nil {
			log.Fatalf("❌ Gagal drop semua tabel: %v", err)
		}
		log.Println("🔥 Semua tabel dihapus!")

	case "new":
		if len(os.Args) < 3 {
			log.Fatal("⚠️  Gunakan: go run cmd/migrate/main.go new nama_migration")
		}
		name := os.Args[2]
		createMigrationFiles(migrationsDir, name)

	default:
		fmt.Println("Gunakan salah satu perintah berikut:")
		fmt.Println("  go run cmd/migrate/main.go up             # Jalankan semua migration")
		fmt.Println("  go run cmd/migrate/main.go down           # Rollback 1 langkah")
		fmt.Println("  go run cmd/migrate/main.go drop           # Hapus semua tabel")
		fmt.Println("  go run cmd/migrate/main.go new nama_file  # Buat migration baru")
	}
}

func createMigrationFiles(dir, name string) {
	timestamp := time.Now().Format("20060102150405")
	safeName := strings.ToLower(strings.ReplaceAll(name, " ", "_"))

	upFile := filepath.Join(dir, fmt.Sprintf("%s_%s.up.sql", timestamp, safeName))
	downFile := filepath.Join(dir, fmt.Sprintf("%s_%s.down.sql", timestamp, safeName))

	// pastikan folder ada
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("❌ Gagal membuat folder migrations: %v", err)
	}

	// buat file kosong
	if err := os.WriteFile(upFile, []byte("-- tulis SQL untuk UP migration di sini\n"), 0644); err != nil {
		log.Fatalf("❌ Gagal membuat file UP: %v", err)
	}
	if err := os.WriteFile(downFile, []byte("-- tulis SQL untuk DOWN migration di sini\n"), 0644); err != nil {
		log.Fatalf("❌ Gagal membuat file DOWN: %v", err)
	}

	log.Printf("✅ Migration baru dibuat:\n- %s\n- %s\n", upFile, downFile)
}
