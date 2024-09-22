package constants

const SQL_DEVICE_SEQUENCES = "CREATE SEQUENCE IF NOT EXISTS device_sequence START 1;"

const SQL_CREATE_TABLE_TEMPLATE string = `
CREATE TABLE IF NOT EXISTS template (
    id 			INTEGER PRIMARY KEY DEFAULT nextval('device_sequence'),
    name		VARCHAR UNIQUE NOT NULL,
    oid_bw		VARCHAR NOT NULL,
    oid_in		VARCHAR NOT NULL,
    oid_out		VARCHAR NOT NULL,
	prefix_bw	VARCHAR NOT NULL,
    prefix_in	VARCHAR NOT NULL,
    prefix_out	VARCHAR NOT NULL,
    created_at	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at	TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

const SQL_CREATE_TABLE_DEVICE string = `
CREATE TABLE IF NOT EXISTS device (
    id 			INTEGER PRIMARY KEY DEFAULT nextval('device_sequence'),
    ip			VARCHAR NOT NULL,
    sysname		VARCHAR NOT NULL,
    community	VARCHAR NOT NULL,
	template_id	INTEGER REFERENCES template(id),
    is_alive    BOOLEAN NOT NULL DEFAULT false,
	last_check	TIMESTAMP,
	created_at	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (ip, community)
);`

// The file db of this interfaces will be named with the device.id
const SQL_CREATE_TABLE_INTERFACE string = `
CREATE TABLE IF NOT EXISTS interface (
    id 			INTEGER PRIMARY KEY, -- IfIndex
    name		VARCHAR,
    descr		VARCHAR,
    alias		VARCHAR,
	created_at	TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at	TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

// The file db of this interfaces will be named with the device.id
const SQL_CREATE_TABLE_MEASUREMENT string = `
CREATE TABLE IF NOT EXISTS measurement (
    date			TIME,
    bw				INTEGER,
    in				INTEGER,
    out				INTEGER,
    interface_id	INTEGER REFERENCES interface(id)
);`

// The file db of this interfaces will be named with the device.id
const SQL_CREATE_TABLE_TMP_MEASUREMENT string = `
CREATE TABLE IF NOT EXISTS tmp_measurement (
    date			TIME,
    bw			    INTEGER,
    in			    INTEGER,
    out		        INTEGER,
    interface_id	INTEGER REFERENCES interface(id)
);`

const (
	SQL_SELECT_ALL_INTERFACES string = "SELECT id, name, descr, alias, created_at, updated_at FROM interface;"
	SQL_SELECT_INTERFACE      string = "SELECT id, name, descr, alias, created_at, updated_at FROM interface WHERE id=$1;"
	SQL_UPSERT_INTERFACE      string = `
	INSERT INTO interface (id, name, descr, alias) VALUES ($1, $2, $3, $4)
	ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		descr = EXCLUDED.descr,
		alias = EXCLUDED.alias,
		updated_at = CURRENT_TIMESTAMP
	RETURNING id;`
)

const (
	SQL_INSERT_MEASUREMENT     string = "INSERT INTO measurement (date, bw, in, out, interface_id) VALUES ($1, $2, $3, $4, $5)"
	SQL_INSERT_TMP_MEASUREMENT string = "INSERT INTO tmp_measurement (date, bw, in, out, interface_id) VALUES ($1, $2, $3, $4, $5)"
	SQL_SELECT_TMP_MEASUREMENT string = "SELECT date, bw, in, out, interface_id FROM tmp_measurement;"
)
