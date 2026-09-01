-- +goose Up


-- Create Accent-Insensitive Trigram index for client_name and desscription
-- How to use: WHERE immutable_unaccent(clientname) ILIKE immutable_unaccent('%albe%')
CREATE INDEX idx_order_client_name_trgm
ON "order"
USING GIN (
    immutable_unaccent(client_name) gin_trgm_ops 
);

CREATE INDEX idx_order_description_trgm
ON "order"
USING GIN (
    immutable_unaccent(description) gin_trgm_ops 
);

-- Create B-Tree index for order status
CREATE INDEX idx_order_status
ON "order" (status);

-- Create B-Tree index for order.created_at
CREATE INDEX idx_order_created_at
ON "order" (created_at);

-- +goose Down
-- ----------------------------------------------------------------------------------------------

DROP INDEX IF EXISTS idx_order_client_name_trgm;
DROP INDEX IF EXISTS idx_order_description_trgm;
DROP INDEX IF EXISTS idx_order_status;
DROP INDEX IF EXISTS idx_order_created_at;