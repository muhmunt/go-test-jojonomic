-- Tabel Harga
CREATE TABLE harga (
    admin_id VARCHAR(255) PRIMARY KEY,
    harga_topup INT,
    harga_buyback INT
);

-- Tabel Rekening
CREATE TABLE rekening (
    norek VARCHAR(255) PRIMARY KEY,
    saldo FLOAT
);

-- Tabel Transaksi
CREATE TABLE transaksi (
    id VARCHAR(255) PRIMARY KEY,
    norek VARCHAR(255),
    type VARCHAR(10) CHECK (type IN ('TOPUP', 'BUYBACK')),
    gram FLOAT,
    harga_topup INT,
    harga_buyback INT,
    saldo_terakhir FLOAT,
    date INT
);
