CREATE TABLE IF NOT EXISTS wallets (
    "valletId" UUID PRIMARY KEY,
    "balance" NUMERIC(15, 2) NOT NULL DEFAULT 0.00,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS operations (
    "operationId" SERIAL PRIMARY KEY,
    "valletId" UUID NOT NULL REFERENCES wallets("valletId") ON DELETE CASCADE,
    "operationType" TEXT NOT NULL,
    "amount" NUMERIC(15, 2) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX IF NOT EXISTS index_operations_valletId ON operations("valletId");