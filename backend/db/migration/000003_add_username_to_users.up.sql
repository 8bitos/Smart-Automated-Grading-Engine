-- Menambahkan kolom username yang unik dan tidak boleh null
-- Dibuat nullable terlebih dahulu untuk kompatibilitas dengan data yang ada
ALTER TABLE users ADD COLUMN username VARCHAR(50) UNIQUE;
