package constants

const SQL_TABLE_OLT string = `
CREATE TABLE olt (
    id BIGINT UNSIGNED PRIMARY KEY,
    ip VARCHAR(15) NOT NULL,
    community VARCHAR(255),
    sys_name VARCHAR(255),
    sys_location VARCHAR(255),
    is_alive BOOLEAN,
    template_id INT UNSIGNED,
    last_check DATETIME,
    created_at DATETIME
);`

const SQL_TABLE_MEASUREMENT_OLT string = `
CREATE TABLE measurement_olt (
    pon_id BIGINT UNSIGNED NOT NULL,
    bandwidth BIGINT UNSIGNED,
    bytes_in_count BIGINT UNSIGNED,
    bytes_out_count BIGINT UNSIGNED,
    date DATETIME NOT NULL
    FOREIGN KEY (pon_id) REFERENCES pon(id)

);`
const SQL_TABLE_MEASUREMENT_ONT string = `
CREATE TABLE measurement_ont (
    pon_id BIGINT UNSIGNED,
    idx BIGINT UNSIGNED,
    despt VARCHAR(255),
    serial_number VARCHAR(255),
    line_prof_name VARCHAR(255),
    olt_distance BIGINT,
    control_mac_count BIGINT,
    control_run_status TINYINT,
    bytes_in BIGINT UNSIGNED,
    bytes_out BIGINT UNSIGNED,
    date DATETIME
    FOREIGN KEY (pon_id) REFERENCES pon(id)

);`

const SQL_TABLE_TRAFFIC_OLT string = `
CREATE TABLE traffic_olt (
    date DATETIME NOT NULL,
    mbps_in DOUBLE,
    mbps_out DOUBLE,
    bandwidth_mbps_sec DOUBLE,
    mbytes_in_sec DOUBLE,
    mbytes_out_sec DOUBLE
);`

const SQL_TABLE_PON string = `
CREATE TABLE pon (
    id BIGINT UNSIGNED PRIMARY KEY,
    olt_id BIGINT UNSIGNED,
    if_index BIGINT UNSIGNED,
    if_name VARCHAR(128),
    if_descr VARCHAR(255),
    if_alias VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    FOREIGN KEY (olt_id) REFERENCES olt(id)
);`

const SQL_TABLE_LOCATION string = `
CREATE TABLE location (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    state VARCHAR(128) NOT NULL,
    county VARCHAR(128) NOT NULL,
    municipality VARCHAR(128) NOT NULL,
    UNIQUE KEY unique_location (state, county, municipality)
);`

const SQL_TABLE_FAT string = `
CREATE TABLE fat (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    location_id INT UNSIGNED NOT NULL,
    splitter TINYINT UNSIGNED,
    latitude DOUBLE,
    longitude DOUBLE,
    address VARCHAR(255),
    fat VARCHAR(64),
    odn VARCHAR(64),
    created_at DATETIME,
    FOREIGN KEY (location_id) REFERENCES location(id)
);`

const SQL_TABLE_FAT_PON string = `
CREATE TABLE fat_pon (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    fat_id INT UNSIGNED NOT NULL,
    pon_id BIGINT UNSIGNED NOT NULL,
    UNIQUE KEY unique_fat_pon (fat_id, pon_id),
    FOREIGN KEY (fat_id) REFERENCES fat(id),
    FOREIGN KEY (pon_id) REFERENCES pon(id)
);`
