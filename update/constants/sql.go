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
    created_at	TIMESTAMP,
    updated_at	TIMESTAMP
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
	created_at	TIMESTAMP,
    updated_at	TIMESTAMP,
    UNIQUE (ip, community)
);`

// The file db of this interfaces will be named with the device.ip
const SQL_CREATE_TABLE_INTERFACE string = `
CREATE TABLE IF NOT EXISTS interface (
    id 			INTEGER PRIMARY KEY, -- IfIndex
    name		VARCHAR,
    descr		VARCHAR,
    alias		VARCHAR,
	created_at	TIMESTAMP,
    updated_at	TIMESTAMP
);`

const SQL_CREATE_TABLE_MEASUREMENT string = `
CREATE TABLE IF NOT EXISTS measurement (
    date			TIME,
    bps_bw			BIGINT,
    bps_in			BIGINT,
    bps_out			BIGINT,
    interface_id	INTEGER REFERENCES interface(id)
);`

const SQL_CREATE_TABLE_TMP_MEASUREMENT string = `
CREATE TABLE IF NOT EXISTS tmp_measurement (
    date			TIME,
    snmp_bw			BIGINT,
    snmp_in			BIGINT,
    snmp_out		BIGINT,
    interface_id	INTEGER REFERENCES interface(id)
);`

const (
	SQL_SELECT_ALL_DEVICES          string = "SELECT id, ip, sysname, community, template_id, is_alive, last_check, created_at, updated_at FROM device;"
	SQL_SELECT_ALL_DEVICES_AND_OIDS string = `
    SELECT 
        d.id,
        d.ip,
        d.sysname,
        d.community,
        d.template_id,
        d.is_alive,
        d.last_check,
        d.created_at,
        d.updated_at,
        t.oid_bw,
        t.oid_in,
        t.oid_out
    FROM 
        device d
    JOIN 
        template t ON d.template_id = t.id;`
)

const (
	SQL_SELECT_ALL_INTERFACES string = "SELECT id, name, descr, alias, created_at, updated_at FROM interface;"
	SQL_SELECT_INTERFACE      string = "SELECT id, name, descr, alias, created_at, updated_at FROM interface WHERE id=$1;"
	SQL_UPSERT_INTERFACE      string = `
	INSERT INTO interface (id, name, descr, alias, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		descr = EXCLUDED.descr,
		alias = EXCLUDED.alias,
		updated_at = EXCLUDED.updated_at
	RETURNING id;`
)

const (
	TABLE_MEASUREMENT          string = "measurement"
	TABLE_TMP_MEASUREMENT      string = "tmp_measurement"
	SQL_INSERT_MEASUREMENT     string = "INSERT INTO measurement (date, bw, in, out, interface_id) VALUES ($1, $2, $3, $4, $5)"
	SQL_INSERT_TMP_MEASUREMENT string = "INSERT INTO tmp_measurement (date, bw, in, out, interface_id) VALUES ($1, $2, $3, $4, $5)"
	SQL_SELECT_TMP_MEASUREMENT string = "SELECT date, bw, in, out, interface_id FROM tmp_measurement;"
)
