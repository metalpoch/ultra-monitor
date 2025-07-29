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
    shell SMALLINT NOT NULL,
    card SMALLINT NOT NULL,
    port SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (fat, municipality, county, ip, odn, shell, card, port)
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
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (ip, idx, shell, card, port)
);`

const SQL_INDEX_REPORT_CATEGORY string = `CREATE INDEX IF NOT EXISTS idx_reports_category ON reports(category);`
const SQL_INDEX_REPORT_USER_ID string = `CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports(user_id);`

const SQL_INDEX_USERS_USERNAME string = `CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);`
