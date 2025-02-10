CREATE TABLE scans (
    ip TEXT,
    port INTEGER,
    service TEXT,
    last_scan INTEGER,
    response TEXT,
    UNIQUE(ip, port, service) ON CONFLICT REPLACE
);
