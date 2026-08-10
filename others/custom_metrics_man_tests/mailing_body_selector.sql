-- Metricraft logs database (DATABASE_LOGS)
-- Smoke test seed for the custom metric:
--   {"name":"mailing","method":"POST","path":"/api/auth/login","source":"body",
--    "selector":"users.gmail","aggregation":"count","timeframe":"Last 365 days",
--    "valueType":"string","applyRules":false,"chartType":"line"}
--
-- With source=body the worker turns the selector into the jsonpath $."users"."gmail"
-- and keeps a row only when payload @? that jsonpath is true.
-- applyRules=false means the url has to match exactly (no prefix/grouping rules)
-- and blacklisting is skipped.
--
-- "Last 365 days" resolves to a 365 day window bucketed 30 days at a time, so the
-- chart should end up with 13 points.
--
-- Paste into pgAdmin Query Tool and run against your logs database.

-- ---------------------------------------------------------------------------
-- Clean previous runs of this seed only
-- ---------------------------------------------------------------------------

DELETE FROM logs WHERE "user" LIKE '198.51.100.20%';

-- ---------------------------------------------------------------------------
-- MATCHING ROWS
-- POST /api/auth/login with users.gmail present in the body.
-- One row every 3 days over the last 360 days, and the number of rows per day
-- grows with the 30 day bucket index so the line chart is not flat.
-- ---------------------------------------------------------------------------

INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status)
SELECT
    (NOW() AT TIME ZONE 'UTC') - (day_offset * INTERVAL '1 day'),
    (45 + (day_offset * 7) % 260)::integer,
    '/api/auth/login',
    '198.51.100.200',
    'Poland',
    format('{"users":{"gmail":"user%s@gmail.com","id":%s},"remember":true}', day_offset, day_offset)::jsonb,
    '{"X-Real-IP":"198.51.100.200","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}'::jsonb,
    'POST',
    200
FROM generate_series(0, 360, 3) AS day_offset
CROSS JOIN generate_series(1, 1 + ((day_offset / 30)::int % 4)) AS rep;

-- Nested siblings and extra keys must not change the outcome: still a match.
INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status) VALUES
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '2 hours',
    120,
    '/api/auth/login',
    '198.51.100.201',
    'Germany',
    '{"users":{"gmail":"fresh@gmail.com","outlook":"fresh@outlook.com"},"meta":{"source":"web"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.201","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '1 day',
    98,
    '/api/auth/login',
    '198.51.100.201',
    'Germany',
    '{"users":{"gmail":""}}'::jsonb,
    '{"X-Real-IP":"198.51.100.201","User-Agent":"curl/8.5.0","Accept":"*/*"}'::jsonb,
    'POST',
    401
),
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '200 days',
    150,
    '/api/auth/login',
    '198.51.100.201',
    'Germany',
    '{"users":{"gmail":null}}'::jsonb,
    '{"X-Real-IP":"198.51.100.201","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
);

-- ---------------------------------------------------------------------------
-- NON MATCHING ROWS
-- Every row below hits the same endpoint or looks close enough to be tempting,
-- but must be excluded from the metric.
-- ---------------------------------------------------------------------------

INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status) VALUES
-- selector missing: sibling key only
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '3 days',
    130,
    '/api/auth/login',
    '198.51.100.202',
    'Spain',
    '{"users":{"email":"nogmail@example.com"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.202","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- selector at the wrong depth: top level gmail, not users.gmail
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '4 days',
    140,
    '/api/auth/login',
    '198.51.100.202',
    'Spain',
    '{"gmail":"toplevel@gmail.com"}'::jsonb,
    '{"X-Real-IP":"198.51.100.202","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- users is an array, so users.gmail does not resolve
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '5 days',
    160,
    '/api/auth/login',
    '198.51.100.202',
    'Spain',
    '{"users":[{"gmail":"inarray@gmail.com"}]}'::jsonb,
    '{"X-Real-IP":"198.51.100.202","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- key casing differs: json paths are case sensitive
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '6 days',
    110,
    '/api/auth/login',
    '198.51.100.202',
    'Spain',
    '{"users":{"Gmail":"casing@gmail.com"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.202","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- empty body and NULL body: must yield no match instead of an error
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '7 days',
    75,
    '/api/auth/login',
    '198.51.100.203',
    'France',
    NULL,
    '{"X-Real-IP":"198.51.100.203","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    400
),
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '8 days',
    80,
    '/api/auth/login',
    '198.51.100.203',
    'France',
    NULL,
    '{"X-Real-IP":"198.51.100.203","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    400
),
-- right body, wrong method
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '9 days',
    60,
    '/api/auth/login',
    '198.51.100.204',
    'Italy',
    '{"users":{"gmail":"wrongmethod@gmail.com"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.204","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'GET',
    200
),
-- right body, sub path: applyRules=false requires an exact url match
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '10 days',
    95,
    '/api/auth/login/refresh',
    '198.51.100.204',
    'Italy',
    '{"users":{"gmail":"subpath@gmail.com"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.204","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- right body, url carries a query string so it is not an exact match either
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '11 days',
    105,
    '/api/auth/login?redirect=/dashboard',
    '198.51.100.204',
    'Italy',
    '{"users":{"gmail":"query@gmail.com"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.204","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
),
-- outside the 365 day window
(
    (NOW() AT TIME ZONE 'UTC') - INTERVAL '400 days',
    88,
    '/api/auth/login',
    '198.51.100.205',
    'Norway',
    '{"users":{"gmail":"tooold@gmail.com"}}'::jsonb,
    '{"X-Real-IP":"198.51.100.205","User-Agent":"Mozilla/5.0","Accept":"application/json"}'::jsonb,
    'POST',
    200
);

-- ---------------------------------------------------------------------------
-- Expected result: mirrors the predicate the worker builds for this metric.
-- The count here should equal the sum of the values plotted on the chart.
-- ---------------------------------------------------------------------------

SELECT COUNT(*) AS expected_total
FROM logs
WHERE date >= (NOW() AT TIME ZONE 'UTC') - INTERVAL '365 days'
  AND date < (NOW() AT TIME ZONE 'UTC')
  AND url = '/api/auth/login'
  AND method = 'POST'
  AND payload @? '$."users"."gmail"'::jsonpath;

-- Per bucket breakdown, roughly what the 13 chart points should look like.
SELECT
    FLOOR(EXTRACT(EPOCH FROM date - ((NOW() AT TIME ZONE 'UTC') - INTERVAL '365 days')) / (30 * 86400))::int AS bucket,
    COUNT(*) AS value
FROM logs
WHERE date >= (NOW() AT TIME ZONE 'UTC') - INTERVAL '365 days'
  AND date < (NOW() AT TIME ZONE 'UTC')
  AND url = '/api/auth/login'
  AND method = 'POST'
  AND payload @? '$."users"."gmail"'::jsonpath
GROUP BY bucket
ORDER BY bucket;

-- ---------------------------------------------------------------------------
-- Optional: malformed body on the tracked endpoint.
-- payload is JSONB, so invalid JSON cannot be inserted; the worker still fails if
-- legacy TEXT rows contain non-JSON text. Run only when testing that path.
-- ---------------------------------------------------------------------------

-- INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status) VALUES
-- (
--     (NOW() AT TIME ZONE 'UTC') - INTERVAL '12 days',
--     70,
--     '/api/auth/login',
--     '198.51.100.206',
--     'Unknown',
--     '{"broken":}'::jsonb,
--     '{"X-Real-IP":"198.51.100.206","User-Agent":"curl/8.5.0","Accept":"*/*"}'::jsonb,
--     'POST',
--     500
-- );
