-- Metricraft logs database (DATABASE_LOGS)
-- Paste into pgAdmin Query Tool and run against your logs database.
--
-- Replace 'my-app' with your APPNAME env var before running.

-- ---------------------------------------------------------------------------
-- Schema (skip this section if the worker has already created the tables)
-- ---------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS settings (
    realtime BOOL,
    enabled TEXT,
    retention INTEGER,
    appName TEXT
);

CREATE TABLE IF NOT EXISTS logs (
    date TIMESTAMP,
    responseTime INTEGER,
    url TEXT NOT NULL,
    "user" TEXT NOT NULL,
    country TEXT NOT NULL,
    payload TEXT,
    headers TEXT NOT NULL,
    method TEXT NOT NULL,
    status INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_date ON logs (date);
CREATE INDEX IF NOT EXISTS idx_url ON logs (url);

-- ---------------------------------------------------------------------------
-- Settings (log retention + dashboard config)
-- Matches the default row inserted by worker/db/db.go InitDB()
-- ---------------------------------------------------------------------------

INSERT INTO settings (realtime, enabled, retention, appname)
SELECT
    true,
    '{"Geographical traffic":true,"P95 Latency":true,"Traffic congestion trends":false,"Uptime Score":true,"Geographic performance":false,"Status code distribution":false,"Median response time":true,"Throughput":true}',
    30,
    'my-app'
WHERE NOT EXISTS (SELECT 1 FROM settings);

-- ---------------------------------------------------------------------------
-- Sample log rows (optional — useful for testing metrics and retention)
-- "user" stores the client IP (X-Real-IP), same as worker/db/db.go Insert()
-- ---------------------------------------------------------------------------

INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status) VALUES
(
    NOW() - INTERVAL '1 hour',
    42,
    '/api/users',
    '203.0.113.10',
    'United States',
    '{"page":1}',
    '{"X-Real-IP":"203.0.113.10","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '6 hours',
    128,
    '/api/dashboard',
    '198.51.100.22',
    'Germany',
    '{}',
    '{"X-Real-IP":"198.51.100.22","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '2 days',
    890,
    '/api/reports',
    '192.0.2.55',
    'Poland',
    '{"range":"7d"}',
    '{"X-Real-IP":"192.0.2.55","User-Agent":"curl/8.5.0","Accept":"application/json"}',
    'POST',
    200
),
(
    NOW() - INTERVAL '15 days',
    2100,
    '/api/legacy/export',
    '203.0.113.88',
    'United Kingdom',
    '{"format":"csv"}',
    '{"X-Real-IP":"203.0.113.88","User-Agent":"PostmanRuntime/7.36.0","Accept":"*/*"}',
    'GET',
    504
),
(
    NOW() - INTERVAL '45 days',
    75,
    '/health',
    '198.51.100.99',
    'Unknown',
    '{}',
    '{"User-Agent":"kube-probe/1.0","Accept":"*/*"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '5 minutes',
    18,
    '/api/users',
    '203.0.113.10',
    'United States',
    '{"page":2}',
    '{"X-Real-IP":"203.0.113.10","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '20 minutes',
    31,
    '/api/users',
    '203.0.113.11',
    'United States',
    '{"page":1}',
    '{"X-Real-IP":"203.0.113.11","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '45 minutes',
    55,
    '/api/users/42',
    '203.0.113.12',
    'Canada',
    '{}',
    '{"X-Real-IP":"203.0.113.12","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '3 hours',
    102,
    '/api/dashboard',
    '198.51.100.22',
    'Germany',
    '{}',
    '{"X-Real-IP":"198.51.100.22","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '4 hours',
    167,
    '/api/dashboard',
    '198.51.100.23',
    'Germany',
    '{}',
    '{"X-Real-IP":"198.51.100.23","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '8 hours',
    2400,
    '/api/search',
    '192.0.2.10',
    'France',
    '{"q":"metrics"}',
    '{"X-Real-IP":"192.0.2.10","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '12 hours',
    95,
    '/api/auth/login',
    '192.0.2.20',
    'Spain',
    '{"email":"user@example.com"}',
    '{"X-Real-IP":"192.0.2.20","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}',
    'POST',
    200
),
(
    NOW() - INTERVAL '18 hours',
    412,
    '/api/auth/login',
    '192.0.2.21',
    'Italy',
    '{"email":"bad@example.com"}',
    '{"X-Real-IP":"192.0.2.21","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}',
    'POST',
    401
),
(
    NOW() - INTERVAL '1 day',
    63,
    '/api/settings/realtime',
    '203.0.113.20',
    'Netherlands',
    '{"enabled":true}',
    '{"X-Real-IP":"203.0.113.20","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}',
    'POST',
    200
),
(
    NOW() - INTERVAL '1 day 4 hours',
    88,
    '/api/settings/retention',
    '203.0.113.21',
    'Belgium',
    '{"retention":30}',
    '{"X-Real-IP":"203.0.113.21","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}',
    'POST',
    200
),
(
    NOW() - INTERVAL '3 days',
    156,
    '/api/reports',
    '192.0.2.55',
    'Poland',
    '{"range":"30d"}',
    '{"X-Real-IP":"192.0.2.55","User-Agent":"curl/8.5.0","Accept":"application/json"}',
    'POST',
    200
),
(
    NOW() - INTERVAL '3 days 6 hours',
    1340,
    '/api/reports',
    '192.0.2.56',
    'Czech Republic',
    '{"range":"90d"}',
    '{"X-Real-IP":"192.0.2.56","User-Agent":"curl/8.5.0","Accept":"application/json"}',
    'POST',
    200
),
(
    NOW() - INTERVAL '5 days',
    220,
    '/api/teams/invite',
    '198.51.100.40',
    'Sweden',
    '{"email":"teammate@example.com"}',
    '{"X-Real-IP":"198.51.100.40","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}',
    'POST',
    201
),
(
    NOW() - INTERVAL '7 days',
    49,
    '/api/users',
    '203.0.113.30',
    'Norway',
    '{"page":1}',
    '{"X-Real-IP":"203.0.113.30","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '7 days 2 hours',
    52,
    '/api/users',
    '203.0.113.31',
    'Denmark',
    '{"page":3}',
    '{"X-Real-IP":"203.0.113.31","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '10 days',
    3800,
    '/api/analytics/latency',
    '192.0.2.70',
    'Japan',
    '{"percentile":95}',
    '{"X-Real-IP":"192.0.2.70","User-Agent":"PostmanRuntime/7.36.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '12 days',
    190,
    '/api/analytics/throughput',
    '192.0.2.71',
    'South Korea',
    '{"window":"1h"}',
    '{"X-Real-IP":"192.0.2.71","User-Agent":"PostmanRuntime/7.36.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '14 days',
    112,
    '/api/webhooks/stripe',
    '198.51.100.60',
    'Ireland',
    '{"type":"invoice.paid"}',
    '{"X-Real-IP":"198.51.100.60","User-Agent":"Stripe/1.0","Accept":"application/json","Content-Type":"application/json"}',
    'POST',
    200
),
(
    NOW() - INTERVAL '18 days',
    404,
    '/api/old-endpoint',
    '203.0.113.50',
    'Portugal',
    '{}',
    '{"X-Real-IP":"203.0.113.50","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    404
),
(
    NOW() - INTERVAL '20 days',
    980,
    '/api/export/csv',
    '203.0.113.51',
    'Austria',
    '{"from":"2025-01-01","to":"2025-01-31"}',
    '{"X-Real-IP":"203.0.113.51","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    500
),
(
    NOW() - INTERVAL '22 days',
    145,
    '/api/users/99',
    '192.0.2.80',
    'Switzerland',
    '{"name":"Updated Name"}',
    '{"X-Real-IP":"192.0.2.80","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}',
    'PATCH',
    200
),
(
    NOW() - INTERVAL '25 days',
    72,
    '/api/users/99',
    '192.0.2.81',
    'Switzerland',
    '{}',
    '{"X-Real-IP":"192.0.2.81","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'DELETE',
    204
),
(
    NOW() - INTERVAL '28 days',
    310,
    '/api/dashboard',
    '198.51.100.22',
    'Germany',
    '{}',
    '{"X-Real-IP":"198.51.100.22","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '29 days',
    5200,
    '/api/batch/import',
    '203.0.113.60',
    'Brazil',
    '{"records":5000}',
    '{"X-Real-IP":"203.0.113.60","User-Agent":"python-requests/2.31.0","Accept":"application/json","Content-Type":"application/json"}',
    'PUT',
    502
),
(
    NOW() - INTERVAL '32 days',
    91,
    '/health',
    '198.51.100.99',
    'Unknown',
    '{}',
    '{"User-Agent":"kube-probe/1.0","Accept":"*/*"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '35 days',
    68,
    '/api/users',
    '203.0.113.70',
    'Mexico',
    '{"page":1}',
    '{"X-Real-IP":"203.0.113.70","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '40 days',
    203,
    '/api/dashboard',
    '198.51.100.80',
    'Argentina',
    '{}',
    '{"X-Real-IP":"198.51.100.80","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '50 days',
    44,
    '/health',
    '198.51.100.99',
    'Unknown',
    '{}',
    '{"User-Agent":"kube-probe/1.0","Accept":"*/*"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '60 days',
    110,
    '/api/reports',
    '192.0.2.90',
    'Australia',
    '{"range":"7d"}',
    '{"X-Real-IP":"192.0.2.90","User-Agent":"curl/8.5.0","Accept":"application/json"}',
    'POST',
    200
),
(
    NOW() - INTERVAL '75 days',
    8700,
    '/api/legacy/export',
    '203.0.113.88',
    'United Kingdom',
    '{"format":"xlsx"}',
    '{"X-Real-IP":"203.0.113.88","User-Agent":"PostmanRuntime/7.36.0","Accept":"*/*"}',
    'GET',
    503
),
(
    NOW() - INTERVAL '90 days',
    33,
    '/api/users',
    '203.0.113.90',
    'India',
    '{"page":1}',
    '{"X-Real-IP":"203.0.113.90","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '2 hours',
    403,
    '/api/admin',
    '192.0.2.30',
    'United States',
    '{}',
    '{"X-Real-IP":"192.0.2.30","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    403
),
(
    NOW() - INTERVAL '9 hours',
    256,
    '/api/search',
    '192.0.2.11',
    'France',
    '{"q":"dashboard"}',
    '{"X-Real-IP":"192.0.2.11","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '9 hours 30 minutes',
    278,
    '/api/search',
    '192.0.2.12',
    'France',
    '{"q":"retention"}',
    '{"X-Real-IP":"192.0.2.12","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '4 days',
    119,
    '/api/metrics/uptime',
    '198.51.100.45',
    'Finland',
    '{}',
    '{"X-Real-IP":"198.51.100.45","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '6 days',
    142,
    '/api/metrics/status-codes',
    '198.51.100.46',
    'Finland',
    '{"window":"24h"}',
    '{"X-Real-IP":"198.51.100.46","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '8 days',
    600,
    '/api/metrics/geo',
    '203.0.113.40',
    'Singapore',
    '{}',
    '{"X-Real-IP":"203.0.113.40","User-Agent":"Mozilla/5.0","Accept":"application/json"}',
    'GET',
    200
),
(
    NOW() - INTERVAL '11 days',
    175,
    '/api/ws/connect',
    '192.0.2.75',
    'New Zealand',
    '{}',
    '{"X-Real-IP":"192.0.2.75","User-Agent":"Mozilla/5.0","Upgrade":"websocket"}',
    'GET',
    101
),
(
    NOW() - INTERVAL '16 days',
    299,
    '/api/verify/token',
    '203.0.113.55',
    'Romania',
    '{"token":"abc123"}',
    '{"X-Real-IP":"203.0.113.55","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}',
    'POST',
    400
),
(
    NOW() - INTERVAL '24 days',
    188,
    '/api/sign',
    '192.0.2.85',
    'Greece',
    '{"mail":"signup@example.com"}',
    '{"X-Real-IP":"192.0.2.85","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}',
    'POST',
    200
);
