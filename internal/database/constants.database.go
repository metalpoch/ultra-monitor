package database

const SQL_TABLE_USERS string = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    fullname VARCHAR(255),
    username VARCHAR(10) NOT NULL,
    password VARCHAR(255) NOT NULL,
    change_password BOOLEAN,
    is_admin BOOLEAN,
    is_disabled BOOLEAN,
    created_at TIMESTAMP DEFAULT NOW()
);`

const SQL_TABLE_REPORT string = `
CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY,
    category VARCHAR(128) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(128) NOT NULL,
    basepath VARCHAR(255) NOT NULL,
    filepath VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);`

const SQL_TABLE_OLT string = `
CREATE TABLE IF NOT EXISTS olts (
    ip VARCHAR(15) PRIMARY KEY,
    community VARCHAR(255),
    sys_name VARCHAR(255) NOT NULL UNIQUE,
    sys_location VARCHAR(255),
    is_alive BOOLEAN,
    last_check TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);`

const SQL_TABLE_PON string = `
CREATE TABLE IF NOT EXISTS pons (
    id SERIAL PRIMARY KEY,
    olt_ip VARCHAR(15) NOT NULL,
    if_index BIGINT NOT NULL,
    if_name VARCHAR(128),
    if_descr VARCHAR(255),
    if_alias VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (olt_ip, if_index),
    FOREIGN KEY (olt_ip) REFERENCES olts(ip)
);`

const SQL_TABLE_MEASUREMENT_PON string = `
CREATE TABLE IF NOT EXISTS measurement_pons (
    pon_id INTEGER UNIQUE NOT NULL,
    bandwidth NUMERIC(20,0) NOT NULL,
    bytes_in_count NUMERIC(20,0) NOT NULL,
    bytes_out_count NUMERIC(20,0) NOT NULL,
    date TIMESTAMP NOT NULL,
    FOREIGN KEY (pon_id) REFERENCES pons(id)
);`

const SQL_TABLE_TRAFFIC_PON string = `
CREATE TABLE IF NOT EXISTS traffic_pons (
    pon_id INTEGER NOT NULL,
    bps_in DOUBLE PRECISION,
    bps_out DOUBLE PRECISION,
    bandwidth_mbps_sec DOUBLE PRECISION,
    bytes_in DOUBLE PRECISION,
    bytes_out DOUBLE PRECISION,
    date TIMESTAMP NOT NULL,
    FOREIGN KEY (pon_id) REFERENCES pons(id)
);`

const SQL_TABLE_TRAFFIC_PON_SUMMARY string = `
CREATE TABLE IF NOT EXISTS traffic_pons_summary (
    day DATE NOT NULL,
    olt_ip VARCHAR(15) NOT NULL,
    mbps_in DOUBLE PRECISION,
    mbps_out DOUBLE PRECISION,
    mbytes_in DOUBLE PRECISION,
    mbytes_out DOUBLE PRECISION,
    PRIMARY KEY (day, olt_ip),
    FOREIGN KEY (olt_ip) REFERENCES olts(ip)
);`

const SQL_TABLE_ONT_SUMMARY_STATUS_COUNTS string = `
CREATE TABLE IF NOT EXISTS ont_summary_status_count (
    day DATE NOT NULL,
    olt_ip VARCHAR(15) NOT NULL,
    ports_pon INTEGER NOT NULL,
    actives INTEGER NOT NULL,
    inactives INTEGER NOT NULL,
    unknowns INTEGER NOT NULL,
    PRIMARY KEY (day, olt_ip),
    FOREIGN KEY (olt_ip) REFERENCES olts(ip)
);`

const SQL_TABLE_MEASUREMENT_ONT string = `
CREATE TABLE IF NOT EXISTS measurement_onts (
    pon_id INTEGER NOT NULL,
    idx BIGINT NOT NULL,
    despt VARCHAR(255),
    serial_number VARCHAR(64),
    line_prof_name VARCHAR(128),
    olt_distance INTEGER,
    control_mac_count SMALLINT,
    control_run_status SMALLINT,
    bytes_in_count NUMERIC(20,0) NOT NULL,
    bytes_out_count NUMERIC(20,0) NOT NULL,
    date TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (pon_id) REFERENCES pons(id)
);`

const SQL_TABLE_FAT string = `
CREATE TABLE IF NOT EXISTS fats (
    id SERIAL PRIMARY KEY,
    fat VARCHAR(128) NOT NULL,
    region VARCHAR(128) NOT NULL,
    state VARCHAR(128) NOT NULL,
    municipality VARCHAR(128) NOT NULL,
    county VARCHAR(128) NOT NULL,
    odn VARCHAR(128) NOT NULL,
    olt_ip VARCHAR(15) NOT NULL,
    pon_shell SMALLINT NOT NULL,
    pon_port SMALLINT NOT NULL,
    pon_card SMALLINT NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (fat, state, municipality, county, olt_ip, odn, pon_shell, pon_card, pon_port)
);`

const SQL_INDEX_FAT_STATE string = `CREATE INDEX IF NOT EXISTS idx_fats_state ON fats(state);`
const SQL_INDEX_FAT_COUNTY string = `CREATE INDEX IF NOT EXISTS idx_fats_county ON fats(county);`
const SQL_INDEX_FAT_MUNICIPALITY string = `CREATE INDEX IF NOT EXISTS idx_fats_municipality ON fats(municipality);`
const SQL_INDEX_FAT_ODN string = `CREATE INDEX IF NOT EXISTS idx_fats_odn ON fats(odn);`
const SQL_INDEX_FAT_OLT string = `CREATE INDEX IF NOT EXISTS idx_fats_olt_ip ON fats(olt_ip);`

const SQL_INDEX_MEASUREMENT_ONT_DATE string = `CREATE INDEX IF NOT EXISTS idx_measurement_onts_date ON measurement_onts(date);`
const SQL_INDEX_MEASUREMENT_ONT_IDX string = `CREATE INDEX IF NOT EXISTS idx_measurement_onts_pon_id ON measurement_onts(idx);`
const SQL_INDEX_MEASUREMENT_ONT_DESPT string = `CREATE INDEX IF NOT EXISTS idx_measurement_onts_despt ON measurement_onts(despt);`
const SQL_INDEX_MEASUREMENT_ONT_PON_ID string = `CREATE INDEX IF NOT EXISTS idx_measurement_onts_pon_id ON measurement_onts(pon_id);`
const SQL_INDEX_MEASUREMENT_ONT_PON_STATS string = `CREATE INDEX IF NOT EXISTS idx_measurement_onts_stats ON measurement_onts (date, pon_id, control_run_status) INCLUDE (idx);`

const SQL_INDEX_REPORT_CATEGORY string = `CREATE INDEX IF NOT EXISTS idx_reports_category ON reports(category);`
const SQL_INDEX_REPORT_USER_ID string = `CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports(user_id);`

const SQL_INDEX_USERS_USERNAME string = `CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);`
