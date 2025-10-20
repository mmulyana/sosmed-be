# 🧩 Migration Commands

## 🔹 Jalankan semua migration

```
go run cmd/migrate/main.go up
```

## 🔹 Rollback 1 langkah terakhir

```
go run cmd/migrate/main.go down
```

## 🔹 Hapus semua tabel (drop database)

```
go run cmd/migrate/main.go drop
```

## 🔹 Buat migration baru

```
go run cmd/migrate/main.go new nama_migration
```

### Contoh:

```
go run cmd/migrate/main.go new create_users_table
```

Akan membuat 2 file:

```
cmd/migrate/migrations/20251020123000_create_users_table.up.sql
cmd/migrate/migrations/20251020123000_create_users_table.down.sql
```

📝 Isi file `.up.sql` dan `.down.sql` sesuai kebutuhan.
