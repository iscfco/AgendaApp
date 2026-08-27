-- +goose Up

-- Function to make the index accent insensitive
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION immutable_unaccent(text)
RETURNS text
LANGUAGE sql
IMMUTABLE
PARALLEL SAFE
AS $$
    SELECT public.unaccent('public.unaccent'::regdictionary, $1);
$$;
-- +goose StatementEnd

-- Create Accent-Insensitive Trigram index for user_full_name
-- How to use: WHERE immutable_unaccent(user_full_name) ILIKE immutable_unaccent('%albe%')
CREATE INDEX idx_user_full_name_trgm
ON "user"
USING GIN (
    immutable_unaccent(user_full_name) gin_trgm_ops 
);

-- Create B-Tree index for user role
CREATE INDEX idx_user_role
ON "user" (role);

-- Create B-Tree index for user status
CREATE INDEX idx_user_status
ON "user" (status);

-- +goose Down

DROP INDEX IF EXISTS idx_user_full_name_trgm;
DROP INDEX IF EXISTS idx_user_role;
DROP INDEX IF EXISTS idx_user_status;

DROP FUNCTION IF EXISTS immutable_unaccent(text);