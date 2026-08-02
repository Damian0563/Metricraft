-- Migrates logs.payload and logs.headers from TEXT to JSONB and adds the GIN
-- indexes the custom metric jsonpath predicates rely on.
--
-- Run once per existing logs database, before starting a worker built from this
-- revision. Databases created by this revision already have JSONB columns, so
-- running it there is a no-op. Safe to run more than once.
--
--     psql "$LOGS_DATABASE_URL" -f others/migrate_logs_jsonb.sql
--
-- ALTER TABLE ... TYPE rewrites the whole table under an ACCESS EXCLUSIVE lock,
-- so expect writes to logs to block for the duration on a large table.

BEGIN;

-- A plain payload::jsonb aborts the entire ALTER on the first row that is not
-- valid JSON. A regex gate such as payload ~ '^\s*[{\[]' is not enough either:
-- a value like '{oops' passes it and still fails the cast. Catching the cast
-- exception per row is the only way to guarantee the migration completes.
CREATE OR REPLACE FUNCTION pg_temp.safe_jsonb(value TEXT) RETURNS JSONB AS $$
BEGIN
	RETURN value::jsonb;
EXCEPTION WHEN others THEN
	RETURN NULL;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Rows that never held valid JSON become NULL, or an empty object for the NOT
-- NULL headers column. Inspect what would be dropped before committing with:
--
--     SELECT count(*) FROM logs WHERE payload IS NOT NULL
--       AND pg_temp.safe_jsonb(payload) IS NULL;

DO $$
BEGIN
	IF EXISTS (
		SELECT 1 FROM information_schema.columns
		WHERE table_schema = current_schema() AND table_name = 'logs'
			AND column_name = 'payload' AND data_type <> 'jsonb'
	) THEN
		ALTER TABLE logs ALTER COLUMN payload TYPE JSONB
			USING pg_temp.safe_jsonb(payload);
	END IF;

	IF EXISTS (
		SELECT 1 FROM information_schema.columns
		WHERE table_schema = current_schema() AND table_name = 'logs'
			AND column_name = 'headers' AND data_type <> 'jsonb'
	) THEN
		ALTER TABLE logs ALTER COLUMN headers TYPE JSONB
			USING COALESCE(pg_temp.safe_jsonb(headers), '{}'::jsonb);
	END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_payload ON logs USING GIN (payload);
CREATE INDEX IF NOT EXISTS idx_headers ON logs USING GIN (headers);

COMMIT;
