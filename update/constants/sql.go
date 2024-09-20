package constants

const SQL_DEVICE_SEQUENCES = "CREATE SEQUENCE IF NOT EXISTS device_sequence START 1;"

const SQL_TEMPLATE_TABLE string = `
CREATE TABLE IF NOT EXISTS template (
    id 			INTEGER PRIMARY KEY DEFAULT nextval('device_sequence'),
    name		VARCHAR UNIQUE NOT NULL,
    oid_bw		VARCHAR NOT NULL,
    oid_in		VARCHAR NOT NULL,
    oid_out		VARCHAR NOT NULL,
	prefix_bw	VARCHAR NOT NULL,
    prefix_in	VARCHAR NOT NULL,
    prefix_out	VARCHAR NOT NULL,
    created_at	TIMESTAMP,
    updated_at	TIMESTAMP
);`

const SQL_DEVICE_TABLE string = `
CREATE TABLE IF NOT EXISTS device (
    id 			INTEGER PRIMARY KEY DEFAULT nextval('device_sequence'),
    ip			VARCHAR NOT NULL,
    sysname		VARCHAR NOT NULL,
    community	VARCHAR NOT NULL,
	template_id	INTEGER REFERENCES template(id),
    is_alive    BOOLEAN NOT NULL DEFAULT false,
	last_check	TIMESTAMP,
	created_at	TIMESTAMP,
    updated_at	TIMESTAMP,
    UNIQUE (ip, community)
);`

const SQL_INTERFACE_TABLE string = `
CREATE TABLE IF NOT EXISTS interface (
    id 			INTEGER PRIMARY KEY, -- IfIndex
    name		VARCHAR,
    descr		VARCHAR,
	device_id	INTEGER REFERENCES device(id),
	created_at	TIMESTAMP,
    updated_at	TIMESTAMP
);`

const SQL_MEASUREMENT_TABLE string = `
CREATE TABLE IF NOT EXISTS measurement (
    date			TIME,
    bw				INTEGER,
    in				INTEGER,
    out				INTEGER,
    interface_id	INTEGER REFERENCES interface(id)
);`

const SQL_TMP_MEASUREMENT_TABLE string = `
CREATE TABLE IF NOT EXISTS tmp_measurement (
    date			TIME,
    bw				INTEGER,
    in				INTEGER,
    out				INTEGER,
    interface_id	INTEGER REFERENCES interface(id)
);`
