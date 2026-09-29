-- Add operator-facing visibility and API availability flags to model pricing.

SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '10min';

ALTER TABLE channel_model_pricing
    ADD COLUMN IF NOT EXISTS public_visible BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS api_enabled BOOLEAN NOT NULL DEFAULT TRUE;

COMMENT ON COLUMN channel_model_pricing.public_visible IS 'Whether the model is shown in public/user-facing model catalogues.';
COMMENT ON COLUMN channel_model_pricing.api_enabled IS 'Whether the model can be called through the gateway.';
