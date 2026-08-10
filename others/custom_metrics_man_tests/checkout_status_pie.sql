-- Metricraft logs database (DATABASE_LOGS)
-- Smoke test seed for the custom metric:
--   {"name":"checkout-status","method":"POST","path":"/api/orders/checkout","source":"body",
--    "selector":"order.status","aggregation":"count","timeframe":"Last 7 days",
--    "valueType":"string","applyRules":false,"chartType":"pie"}
--
-- Pie charts use the cumulative query path: rows are grouped by the extracted
-- selector value (not by time bucket). Each slice is one distinct order.status
-- and the slice size is COUNT(*) for that value inside the timeframe window.
--
-- With source=body the worker turns the selector into the jsonpath $."order"."status"
-- and keeps a row only when payload @? that jsonpath is true.
-- applyRules=false means the url has to match exactly (no prefix/grouping rules)
-- and blacklisting is skipped.
--
-- "Last 7 days" filters date >= start AND date < now; the pie should show six
-- slices: completed (16), pending (8), failed (3), refunded (2), plus one row
-- each for an empty-string and a null order.status (edge cases that still match
-- the jsonpath existence check). expected_total = 31.
--
-- Paste into pgAdmin Query Tool and run against your logs database.

-- ---------------------------------------------------------------------------
-- Clean previous runs of this seed only
-- ---------------------------------------------------------------------------

DELETE FROM logs WHERE "user" LIKE '198.51.100.30%';

-- ---------------------------------------------------------------------------
-- MATCHING ROWS
-- POST /api/orders/checkout with order.status present in the body.
-- Distinct status values with fixed counts so pie slices are easy to verify.
-- ---------------------------------------------------------------------------

INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status)
SELECT
    (NOW() AT TIME ZONE 'UTC') - ((row_num % 6) * INTERVAL '1 day') - ((row_num % 12) * INTERVAL '1 hour'),
    (80 + (row_num * 11) % 220)::integer,
    '/api/orders/checkout',
    '198.51.100.300',
    CASE (row_num % 4)
        WHEN 0 THEN 'Poland'
        WHEN 1 THEN 'Germany'
        WHEN 2 THEN 'France'
        ELSE 'Spain'
    END,
    format('{"order":{"status":"completed","id":%s},"currency":"EUR"}', row_num)::jsonb,
    '{"X-Real-IP":"198.51.100.300","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}'::jsonb,
    'POST',
    200
FROM generate_series(1, 15) AS row_num;

INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status)
SELECT
    (NOW() AT TIME ZONE 'UTC') - ((row_num % 5) * INTERVAL '1 day'),
    (90 + (row_num * 13) % 180)::integer,
    '/api/orders/checkout',
    '198.51.100.301',
    'Italy',
    format('{"order":{"status":"pending","id":%s},"retry":true}', 100 + row_num)::jsonb,
    '{"X-Real-IP":"198.51.100.301","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}'::jsonb,
    'POST',
    202
FROM generate_series(1, 8) AS row_num;

INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status)
SELECT
    (NOW() AT TIME ZONE 'UTC') - (row_num * INTERVAL '18 hours'),
    (110 + row_num * 17)::integer,
    '/api/orders/checkout',
    '198.51.100.302',
    'Netherlands',
    format('{"order":{"status":"failed","id":%s},"error":"card_declined"}', 200 + row_num)::jsonb,
    '{"X-Real-IP":"198.51.100.302","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}'::jsonb,
    'POST',
    402
FROM generate_series(1, 3) AS row_num;

INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status)
SELECT
    (NOW() AT TIME ZONE 'UTC') - (row_num * INTERVAL '2 days'),
    (95 + row_num * 21)::integer,
    '/api/orders/checkout',
    '198.51.100.303',
    'Belgium',
    format('{"order":{"status":"refunded","id":%s},"refundReason":"customer_request"}', 300 + row_num)::jsonb,
    '{"X-Real-IP":"198.51.100.303","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}'::jsonb,
    'POST',
    200
FROM generate_series(1, 2) AS row_num;

-- Extra keys and sibling fields must not change the outcome: still a match.
INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status) VALUES
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '3 hours',
    140,
    '/api/orders/checkout',
    '198.51.100.304',
    'Austria',
    '{"order":{"status":"completed","total":49.99,"items":[{"sku":"A1"}]},"meta":{"channel":"web"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.304","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '1 day',
    102,
    '/api/orders/checkout',
    '198.51.100.304',
    'Austria',
    '{"order":{"status":""}}'::jsonb,
    '{"X-Real-IP":"198.51.100.304","User-Agent":"curl/8.5.0","Accept":"*/*"}'::jsonb,
    'POST',
    400
),
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '4 days',
    118,
    '/api/orders/checkout',
    '198.51.100.304',
    'Austria',
    '{"order":{"status":null}}'::jsonb,
    '{"X-Real-IP":"198.51.100.304","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
);

-- ---------------------------------------------------------------------------
-- NON MATCHING ROWS
-- Every row below hits a nearby endpoint or looks close enough to be tempting,
-- but must be excluded from the metric.
-- ---------------------------------------------------------------------------

INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status) VALUES
-- selector missing: sibling key only
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '2 days',
    130,
    '/api/orders/checkout',
    '198.51.100.305',
    'Portugal',
    '{"order":{"state":"shipped"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.305","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- selector at the wrong depth: top level status, not order.status
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '2 days 6 hours',
    125,
    '/api/orders/checkout',
    '198.51.100.305',
    'Portugal',
    '{"status":"completed"}'::jsonb,
    '{"X-Real-IP":"198.51.100.305","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- order is an array, so order.status does not resolve
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '3 days',
    160,
    '/api/orders/checkout',
    '198.51.100.305',
    'Portugal',
    '{"order":[{"status":"completed"}]}'::jsonb,
    '{"X-Real-IP":"198.51.100.305","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- key casing differs: json paths are case sensitive
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '3 days 12 hours',
    111,
    '/api/orders/checkout',
    '198.51.100.305',
    'Portugal',
    '{"order":{"Status":"completed"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.305","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- empty body and NULL body
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '4 days',
    75,
    '/api/orders/checkout',
    '198.51.100.306',
    'Sweden',
    NULL,
    '{"X-Real-IP":"198.51.100.306","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    400
),
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '4 days 6 hours',
    80,
    '/api/orders/checkout',
    '198.51.100.306',
    'Sweden',
    NULL,
    '{"X-Real-IP":"198.51.100.306","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    400
),
-- right body, wrong method
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '5 days',
    60,
    '/api/orders/checkout',
    '198.51.100.307',
    'Denmark',
    '{"order":{"status":"completed"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.307","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'GET',
    200
),
-- right body, sub path: applyRules=false requires an exact url match
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '5 days 12 hours',
    95,
    '/api/orders/checkout/confirm',
    '198.51.100.307',
    'Denmark',
    '{"order":{"status":"completed"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.307","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- right body, url carries a query string so it is not an exact match either
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '6 days',
    105,
    '/api/orders/checkout?source=mobile',
    '198.51.100.307',
    'Denmark',
    '{"order":{"status":"completed"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.307","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- wrong route entirely: login endpoint from the line-chart smoke test
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '6 days 12 hours',
    88,
    '/api/auth/login',
    '198.51.100.308',
    'Norway',
    '{"order":{"status":"completed"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.308","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- outside the 7 day window
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '10 days',
    88,
    '/api/orders/checkout',
    '198.51.100.308',
    'Norway',
    '{"order":{"status":"completed"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.308","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
);

-- ---------------------------------------------------------------------------
-- Expected result: mirrors the predicate the worker builds for this metric.
-- The sum of slice values should equal expected_total (31 rows).
-- ---------------------------------------------------------------------------

SELECT COUNT(*) AS expected_total
FROM logs
WHERE date >= (NOW() AT TIME ZONE 'UTC') - INTERVAL '7 days'
  AND date < (NOW() AT TIME ZONE 'UTC')
  AND url = '/api/orders/checkout'
  AND method = 'POST'
  AND payload @? '$."order"."status"'::jsonpath;

-- Per-slice breakdown: this is what the four pie segments should look like.
SELECT
    jsonb_path_query_first(payload, '$."order"."status"'::jsonpath)#>>'{}' AS slice,
    COUNT(*) AS value
FROM logs
WHERE date >= (NOW() AT TIME ZONE 'UTC') - INTERVAL '7 days'
  AND date < (NOW() AT TIME ZONE 'UTC')
  AND url = '/api/orders/checkout'
  AND method = 'POST'
  AND payload @? '$."order"."status"'::jsonpath
GROUP BY slice
ORDER BY value DESC, slice;

-- ---------------------------------------------------------------------------
-- Optional: malformed body on the tracked endpoint.
-- payload is JSONB, so invalid JSON cannot be inserted; the worker still fails if
-- legacy TEXT rows contain non-JSON text. Run only when testing that path.
-- ---------------------------------------------------------------------------

-- INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status) VALUES
-- (
--     (NOW() AT TIME ZONE 'UTC') - INTERVAL '1 day',
--     70,
--     '/api/orders/checkout',
--     '198.51.100.309',
--     'Unknown',
--     '{"broken":}'::jsonb,
--     '{"X-Real-IP":"198.51.100.309","User-Agent":"curl/8.5.0","Accept":"*/*"}'::jsonb,
--     'POST',
--     500
-- );
