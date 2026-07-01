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
-- Last 24 hours — endpoints from data-1782848361682.csv
-- Quick seed: 96 rows every 15 minutes. Run this block alone for recent metrics.
-- "user" stores the client IP (X-Real-IP), same as worker/db/db.go Insert()
-- ---------------------------------------------------------------------------

INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status)
SELECT
    NOW() - (step * INTERVAL '15 minutes'),
    (18 + (step * 23 + (step / 4) * 11) % 420)::integer,
    endpoints.url,
    ips.ip,
    ips.country,
    endpoints.payload,
    format(
        '{"X-Real-IP":"%s","User-Agent":"Mozilla/5.0","Accept":"application/json","Content-Type":"application/json"}',
        ips.ip
    ),
    endpoints.method,
    endpoints.status
FROM generate_series(0, 95) AS step
CROSS JOIN LATERAL (
    SELECT *
    FROM (VALUES
        (0,  'https://34.116.244.111:8000/cleanup',              '{}',                                  'POST', 200),
        (1,  '/api/users/99',                                     '{"name":"Updated Name"}',             'PATCH', 200),
        (2,  '/api/old-endpoint',                                 '{}',                                  'GET',  404),
        (3,  'https://34.116.244.111:8000/receive/config-v2',   '{}',                                  'POST', 200),
        (4,  '/api/batch/import',                                 '{"records":5000}',                    'PUT',  502),
        (5,  'https://34.116.244.111:8000/drop',                  '{}',                                  'POST', 200),
        (6,  'https://codrawapp.com/api/sync',                    '{}',                                  'POST', 200),
        (7,  '/api/settings/realtime',                            '{"enabled":true}',                    'POST', 200),
        (8,  '/api/export/csv',                                   '{"from":"2025-01-01","to":"2025-01-31"}', 'GET', 500),
        (9,  '/api/analytics/throughput',                         '{"window":"1h"}',                     'GET',  200),
        (10, '/api/users/42',                                       '{}',                                  'GET',  200),
        (11, '/api/settings/retention',                           '{"retention":30}',                    'POST', 200),
        (12, '/api/ws/connect',                                   '{}',                                  'GET',  101),
        (13, '/api/sign',                                           '{"mail":"signup@example.com"}',       'POST', 200),
        (14, '/api/verify/token',                                 '{"token":"abc123"}',                  'POST', 400),
        (15, '/api/teams/invite',                                 '{"email":"teammate@example.com"}',    'POST', 201),
        (16, '/api/admin',                                          '{}',                                  'GET',  403),
        (17, '/api/reports',                                        '{"range":"7d"}',                      'POST', 200),
        (18, '/api/dashboard',                                      '{}',                                  'GET',  200),
        (19, '/api/metrics/geo',                                    '{}',                                  'GET',  200),
        (20, '/api/webhooks/stripe',                                '{"type":"invoice.paid"}',             'POST', 200),
        (21, '/api/users',                                          '{"page":1}',                          'GET',  200),
        (22, '/api/analytics/latency',                              '{"percentile":95}',                   'GET',  200),
        (23, '/api/auth/login',                                     '{"email":"user@example.com"}',        'POST', 200),
        (24, 'https://metricraft.io/api/v1/ingest',                 '{}',                                  'POST', 200),
        (25, '/api/metrics/uptime',                                 '{}',                                  'GET',  200),
        (26, '/health',                                             '{}',                                  'GET',  200),
        (27, '/api/metrics/status-codes',                           '{"window":"24h"}',                    'GET',  200),
        (28, '/api/search',                                         '{"q":"metrics"}',                     'GET',  200),
        (29, '/api/legacy/export',                                  '{"format":"csv"}',                    'GET',  504)
    ) AS t(idx, url, payload, method, status)
    WHERE t.idx = (step % 30)
) AS endpoints
CROSS JOIN LATERAL (
    SELECT *
    FROM (VALUES
        ('203.0.113.10', 'United States'),
        ('198.51.100.22', 'Germany'),
        ('192.0.2.55', 'Poland'),
        ('203.0.113.11', 'Canada'),
        ('198.51.100.23', 'Netherlands'),
        ('192.0.2.20', 'Spain')
    ) AS t(ip, country)
    WHERE t.ip = (ARRAY['203.0.113.10', '198.51.100.22', '192.0.2.55', '203.0.113.11', '198.51.100.23', '192.0.2.20'])[1 + (step % 6)]
) AS ips;

-- ---------------------------------------------------------------------------
-- Older retention test rows (optional — useful for testing retention policy)
-- ---------------------------------------------------------------------------

INSERT INTO logs (date, responseTime, url, "user", country, payload, headers, method, status) VALUES
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
);
