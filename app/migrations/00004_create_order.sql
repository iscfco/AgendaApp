-- +goose Up

-- Order table
CREATE TABLE "order" (
    id                      INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author_id               INTEGER NOT NULL,
    client_name             VARCHAR(255) NOT NULL,
    client_phone            VARCHAR(20),
    client_address          TEXT,
    total_price             NUMERIC(10, 2) DEFAULT 0.00,
    down_payment            NUMERIC(10, 2) DEFAULT 0.00,
    created_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by              INTEGER NOT NULL,
    delivery_date           TIMESTAMP,
    description             TEXT,
    status                  VARCHAR(50) NOT NULL DEFAULT 'pending',
    change_log              JSONB,
    stored_in_change_log_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_order_status     CHECK (status IN ('pending', 'delivered')),
    CONSTRAINT fk_order_author      FOREIGN KEY (author_id) REFERENCES public.user(id),
    CONSTRAINT fk_order_updated_by  FOREIGN KEY (updated_by) REFERENCES public.user(id)
);

-- +goose Down

DROP TABLE "order";
