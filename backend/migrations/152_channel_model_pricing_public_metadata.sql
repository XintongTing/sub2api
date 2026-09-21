-- Store optional public marketplace metadata on each channel model pricing row.
-- These fields are operator-facing overrides for display only; routing and
-- billing still use platform/models/billing/price fields from the same row.

ALTER TABLE channel_model_pricing
    ADD COLUMN IF NOT EXISTS provider TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS endpoint_types JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS tags JSONB NOT NULL DEFAULT '[]'::jsonb;

COMMENT ON COLUMN channel_model_pricing.provider IS 'Optional public provider label override for model marketplace.';
COMMENT ON COLUMN channel_model_pricing.endpoint_types IS 'Optional public endpoint type labels, JSON array.';
COMMENT ON COLUMN channel_model_pricing.description IS 'Optional public model description override.';
COMMENT ON COLUMN channel_model_pricing.tags IS 'Optional public tag labels, JSON array.';
