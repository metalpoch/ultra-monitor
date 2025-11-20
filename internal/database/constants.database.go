package database

const SQL_TABLE_USERS string = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    fullname VARCHAR(255),
    username VARCHAR(10) UNIQUE NOT NULL,
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

const SQL_TABLE_FAT string = `
CREATE TABLE IF NOT EXISTS fats (
    id SERIAL PRIMARY KEY,
    ip VARCHAR(15) NOT NULL,
    region VARCHAR(128) NOT NULL,
    state VARCHAR(128) NOT NULL,
    municipality VARCHAR(128) NOT NULL,
    county VARCHAR(128) NOT NULL,
    odn VARCHAR(128) NOT NULL,
    fat VARCHAR(128) NOT NULL,
    bras VARCHAR(15) NOT NULL,
    shell SMALLINT NOT NULL,
    card SMALLINT NOT NULL,
    port SMALLINT NOT NULL,
		plans TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (ip, region, state, municipality, county, odn, fat, bras, shell, card, port)
);`

const SQL_TABLE_FAT_STATUS string = `
CREATE TABLE IF NOT EXISTS fat_status (
    fats_id SERIAL NOT NULL,
    date DATE NOT NULL,
    actives INTEGER NOT NULL,
    provisioned_offline INTEGER NOT NULL,
    cut_off INTEGER NOT NULL,
    in_progress INTEGER NOT NULL,
    FOREIGN KEY (fats_id) REFERENCES fats(id),
    UNIQUE (fats_id, date)
);`

const SQL_TABLE_PROMETHEUS_DEVICES string = `
CREATE TABLE IF NOT EXISTS prometheus_devices (
    region VARCHAR(128) NOT NULL,
    state VARCHAR(128) NOT NULL,
    ip VARCHAR(15) NOT NULL,
    idx BIGINT NOT NULL,
    shell SMALLINT NOT NULL,
    card SMALLINT NOT NULL,
    port SMALLINT NOT NULL,
    status SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (ip, idx, shell, card, port)
);`

const SQL_TABLE_SUMMARY_TRAFFIC string = `
CREATE TABLE IF NOT EXISTS summary_traffic (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    ip VARCHAR(15) NOT NULL,
    state VARCHAR(128) NOT NULL,
    region VARCHAR(128) NOT NULL,
    sysname VARCHAR(255),
    bps_in DOUBLE PRECISION NOT NULL,
    bps_out DOUBLE PRECISION NOT NULL,
    bytes_in DOUBLE PRECISION NOT NULL,
    bytes_out DOUBLE PRECISION NOT NULL,
    PRIMARY KEY (ip, time)
);`

const SQL_TABLE_ONT string = `
CREATE TABLE IF NOT EXISTS onts (
    id SERIAL PRIMARY KEY,
    ip VARCHAR(15) NOT NULL,
		ont_idx VARCHAR(15) NOT NULL,
    serial VARCHAR(19) UNIQUE NOT NULL,
    despt VARCHAR(50) UNIQUE NOT NULL,
		line_prof VARCHAR(50) NOT NULL,
		description VARCHAR(255) NOT NULL,
		enabled BOOLEAN DEFAULT true,
	  status BOOLEAN,
		olt_distance INTEGER,
		last_check TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);`

const SQL_TABLE_ONT_TRAFFIC string = `
CREATE TABLE IF NOT EXISTS onts_traffic (
    ont_id INTEGER NOT NULL,
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    bps_in DOUBLE PRECISION NOT NULL,
    bps_out DOUBLE PRECISION NOT NULL,
    bytes_in DOUBLE PRECISION NOT NULL,
    bytes_out DOUBLE PRECISION NOT NULL,
		temperature INTEGER NOT NULL,
		rx INTEGER NOT NULL,
		tx INTEGER NOT NULL,
		FOREIGN KEY (ont_id) REFERENCES onts(id),
		PRIMARY KEY (ont_id, time)
);`

const SQL_TABLE_INTERFACE_BANDWIDTH string = `
CREATE TABLE IF NOT EXISTS interfaces_bandwidth (
    olt_verbose VARCHAR(50) NOT NULL,
    interface VARCHAR(255) NOT NULL,
    bandwidth FLOAT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (olt_verbose, interface)
);`

const SQL_TABLE_INTERFACE_OLT string = `
CREATE TABLE IF NOT EXISTS interfaces_olt (
    olt_verbose VARCHAR(50) UNIQUE NOT NULL,
    olt_ip VARCHAR(15) UNIQUE NOT NULL,
    PRIMARY KEY (olt_verbose, olt_ip)
);`

const SQL_INDEX_REPORT_CATEGORY string = `CREATE INDEX IF NOT EXISTS idx_reports_category ON reports(category);`
const SQL_INDEX_REPORT_USER_ID string = `CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports(user_id);`
const SQL_INDEX_USERS_USERNAME string = `CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);`
