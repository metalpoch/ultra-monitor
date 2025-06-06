package constants

// SQL_TOTAL_TRAFFIC retrieves total traffic data aggregated by minute.
const SQL_TOTAL_TRAFFIC string = `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth_mbps_sec,
			SUM(bytes_in_sec) / 1000000 AS mbytes_in_sec,
			SUM(bytes_out_sec) / 1000000 AS mbytes_out_sec
		FROM traffic_pons
		WHERE date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`

// SQL_TRAFFIC_BY_STATE retrieves traffic data for a specific state aggregated by minute.
const SQL_TRAFFIC_BY_STATE string = `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
    FROM traffic_pons
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats ON fats.olt_ip = olts.ip
    WHERE fats.state = $1 AND traffic_pons.date BETWEEN $2 AND $3
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`

// SQL_TRAFFIC_BY_COUNTY retrieves traffic data for a specific county within a state, aggregated by minute.
const SQL_TRAFFIC_BY_COUNTY string = `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
    FROM traffic_pons
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats ON fats.olt_ip = olts.ip
    WHERE fats.state = $1 AND fats.county = $2 AND traffic_pons.date BETWEEN $3 AND $4
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`

// SQL_TRAFFIC_BY_MUNICIPALITY retrieves traffic data for a specific municipality within a county and state, aggregated by minute.
const SQL_TRAFFIC_BY_MUNICIPALITY string = `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
    FROM traffic_pons
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats ON fats.olt_ip = olts.ip
    WHERE fats.state = $1 AND fats.county = $2 AND fats.municipality = $3 AND traffic_pons.date BETWEEN $4 AND $5
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`

// SQL_TRAFFIC_BY_ODN retrieves traffic data for a specific ODN within a state, aggregated by minute.
const SQL_TRAFFIC_BY_ODN string = `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
    FROM traffic_pons
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats ON fats.olt_ip = olts.ip
    WHERE fats.state = $1 AND fats.odn = $2 AND traffic_pons.date BETWEEN $3 AND $4
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`

// SQL_TRAFFIC_BY_OLT retrieves traffic data for a specific OLT, aggregated by minute.
const SQL_TRAFFIC_BY_OLT string = `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
    FROM traffic_pons 
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    WHERE olts.sys_name = $1 AND traffic_pons.date BETWEEN $2 AND $3
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`

// SQL_TRAFFIC_BY_PON retrieves traffic data for a specific PON interface on an OLT, aggregated by minute.
const SQL_TRAFFIC_BY_PON string = `
    SELECT
        DATE_TRUNC('minute', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
        SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
    FROM traffic_pons 
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    WHERE olts.sys_name = $1 AND pons.if_name = $2 AND traffic_pons.date BETWEEN $3 AND $4
    GROUP BY DATE_TRUNC('minute', date)
    ORDER BY date;`

// SQL_PONS_BY_OLT retrieves all PON interfaces associated with a specific OLT.

// SQL_PON_BY_PORT retrieves a specific PON interface by its port name on a given OLT.

// SQL_ALL_ONT_STATUS retrieves ONT status counts for all states within a date range.
const SQL_ALL_ONT_STATUS string = `
    WITH ranked_status AS (
        SELECT
            fats.state AS state,
            DATE_TRUNC('day', measurement_onts.date) AS date,
            measurement_onts.pon_id,
            measurement_onts.idx,
            MIN(
                CASE
                    WHEN measurement_onts.control_run_status = 1 THEN 1
                    WHEN measurement_onts.control_run_status = 2 THEN 2
                    ELSE 3
                END
            ) AS status_priority
        FROM measurement_onts
        JOIN pons ON measurement_onts.pon_id = pons.id
        JOIN olts ON pons.olt_ip = olts.ip
        JOIN fats ON fats.olt_ip = olts.ip
        WHERE measurement_onts.date BETWEEN $1 AND $2
        GROUP BY fats.state, DATE_TRUNC('day', measurement_onts.date), measurement_onts.pon_id, measurement_onts.idx
    )
    SELECT
        state,
        date,
        COUNT(DISTINCT pon_id) AS ports_pon,
        SUM(CASE WHEN status_priority = 1 THEN 1 ELSE 0 END) AS actives,
        SUM(CASE WHEN status_priority = 2 THEN 1 ELSE 0 END) AS inactives,
        SUM(CASE WHEN status_priority = 3 THEN 1 ELSE 0 END) AS unknowns,
        COUNT(*) AS total
    FROM ranked_status
    GROUP BY state, date
    ORDER BY state, date;`

// SQL_ONT_STATUS_BY_STATE retrieves ONT status counts for a specific state within a date range.

// SQL_ONT_STATUS_BY_ODN retrieves ONT status counts for a specific ODN within a state and date range.

// SQL_ONT_STATUS_BY_SYSNAME retrieves ONT status counts for a specific ip within a date range.

// SQL_ONT_STATUS_BY_SYSNAME retrieves ONT status counts for a specific sysname within a date range.

// SQL_TRAFFIC_ONT retrieves traffic data for a specific ONT, including calculated Mbps and Mbytes per second.

// SQL_TRAFFIC_ONT_BY_DESPT retrieves  traffic data for a specific ONT by Despt

// SQL_DELETE_OLT deletes an OLT from the database by its ID.

// SQL_GET_ALL_OLTS retrieves all OLTs from the database.

// SQL_GET_OLTS_BY_STATE retrieves OLTs by their state, with pagination.

// SQL_GET_OLTS_BY_COUNTY retrieves OLTs by their state and county, with pagination.

// SQL_GET_OLTS_BY_MUNICIPALITY retrieves OLTs by their state, county, and municipality, with pagination.

// SQL_GET_OLTS_IPS retrieves all unique OLT IP addresses from the database.

// SQL_UPSERT_OLT updates an existing OLT or inserts a new one if it does not exist.
const SQL_UPSERT_OLT string = `
    UPDATE olts SET
        sys_name = :sys_name,
        sys_location = :sys_location,
        is_alive = :is_alive,
        last_check = :last_check,
    WHERE ip = :ip`

// SQL_UPSERT_PON inserts a new PON interface or updates an existing one if it already exists.
const SQL_UPSERT_PON string = `
    INSERT INTO pons (olt_ip, if_index, if_name, if_descr, if_alias)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (olt_ip, if_index) DO UPDATE SET
        if_name = EXCLUDED.if_name,
        if_descr = EXCLUDED.if_descr,
        if_alias = EXCLUDED.if_alias
    RETURNING id
`

// SQL_GET_TEMPORAL_MEASUREMENT_PON retrieves a temporal measurement for a specific PON.
const SQL_GET_TEMPORAL_MEASUREMENT_PON string = `SELECT * FROM measurement_pons WHERE pon_id = $1`

// SQL_UPSERT_TEMPORAL_MEASUREMENT_PON inserts or updates a temporal measurement for a PON.
const SQL_UPSERT_TEMPORAL_MEASUREMENT_PON string = `
    INSERT INTO measurement_pons (pon_id, bandwidth, bytes_in_count, bytes_out_count, date)
    VALUES (:pon_id, :bandwidth, :bytes_in_count, :bytes_out_count, :date)
    ON CONFLICT (pon_id) DO UPDATE SET
        bandwidth = EXCLUDED.bandwidth,
        bytes_in_count = EXCLUDED.bytes_in_count,
        bytes_out_count = EXCLUDED.bytes_out_count,
        date = EXCLUDED.date`

// SQL_INSERT_TRAFFIC_PON inserts traffic data for a PON.
const SQL_INSERT_TRAFFIC_PON string = `
    INSERT INTO traffic_pons (pon_id, date, bps_in, bps_out, bandwidth_mbps_sec, bytes_in_sec, bytes_out_sec)
    VALUES (:pon_id, :date, :bps_in, :bps_out, :bandwidth_mbps_sec, :bytes_in_sec, :bytes_out_sec)`

// SQL_INSERT_MANY_MEASUREMENT_ONT_PREFIX is the prefix for inserting multiple ONT measurements.
const SQL_INSERT_MANY_MEASUREMENT_ONT_PREFIX string = `
    INSERT INTO measurement_onts (
            pon_id, idx, despt, serial_number, line_prof_name, olt_distance,
            control_mac_count, control_run_status, bytes_in_count, bytes_out_count, date
        ) VALUES `

// SQL_DELETE_FAT_BY_ID is the SQL statement to delete a FAT record by its ID.

// SQL_SELECT_ALL_FATS retrieves all FAT records with pagination.

// SQL_SELECT_FAT_BY_ID retrieves a FAT record by its ID.

// SQL_SELECT_FAT_BY_FAT retrieves a FAT record by its FAT identifier.
const SQL_SELECT_FAT_BY_FAT = `SELECT * FROM fats WHERE fat = $1`

// SQL_SELECT_FATS_BY_ODN retrieves all FAT records for a specific ODN in a given state, ordered by FAT.
const SQL_SELECT_FATS_BY_ODN = `SELECT * FROM fats WHERE state = $1 AND odn = $2 ORDER BY fat`

// SQL_SELECT_DISTINCT_ODN_BY_STATE retrieves distinct ODNs for a specific state, ordered by ODN.
const SQL_SELECT_DISTINCT_ODN_BY_STATE = `SELECT DISTINCT odn FROM fats WHERE state = $1 ORDER BY odn`

// SQL_SELECT_DISTINCT_ODN_BY_COUNTY retrieves distinct ODNs for a specific state and county, ordered by ODN.
const SQL_SELECT_DISTINCT_ODN_BY_COUNTY = `SELECT DISTINCT odn FROM fats WHERE state = $1 AND municipality = $2 AND county = $3 ORDER BY odn`

// SQL_SELECT_DISTINCT_ODN_BY_MUNICIPALITY retrieves distinct ODNs for a specific state, county, and municipality, ordered by ODN.
const SQL_SELECT_DISTINCT_ODN_BY_MUNICIPALITY = `SELECT DISTINCT odn FROM fats WHERE state = $1 AND municipality = $3 ORDER BY odn`

// SQL_SELECT_DISTINCT_ODN_BY_OLT retrieves distinct ODNs for a specific OLT, ordered by ODN.
const SQL_SELECT_DISTINCT_ODN_BY_OLT = `SELECT DISTINCT odn FROM fats WHERE olt_ip = $1 ORDER BY odn`

// SQL_SELECT_DISTINCT_ODN_BY_OLT_PORT retrieves distinct ODNs for a specific OLT and PON port, ordered by ODN.
const SQL_SELECT_DISTINCT_ODN_BY_OLT_PORT = `SELECT DISTINCT odn FROM fats WHERE olt_ip = $1 AND pon_shell = $2 AND pon_card = $3 AND pon_port = $4 ORDER BY odn`

// SQL_SELECT_DISTINCT_ALL_ODN retrieves all distinct ODNs from the fats table.
const SQL_SELECT_DISTINCT_ALL_ODN = `SELECT DISTINCT odn FROM fats`

// SQL_INSERT_REPORT inserts a new report into the reports table.
const SQL_INSERT_REPORT = `
    INSERT INTO reports (
        id, category, original_filename, content_type, basepath, filepath, user_id
    ) VALUES (
        :id, :category, :original_filename, :content_type, :basepath, :filepath, :user_id
    )`

// SQL_SELECT_REPORT_BY_ID retrieves a report by its ID.
const SQL_SELECT_REPORT_BY_ID = `SELECT * FROM reports WHERE id = $1 ORDER BY created_at`

// SQL_SELECT_DISTINCT_REPORT_CATEGORIES retrieves all distinct report categories.
const SQL_SELECT_DISTINCT_REPORT_CATEGORIES = `SELECT DISTINCT category FROM reports ORDER BY category`

// SQL_SELECT_REPORTS_BY_USER retrieves paginated reports for a specific user.
const SQL_SELECT_REPORTS_BY_USER = `SELECT * FROM reports WHERE user_id = $1 LIMIT $2 OFFSET $3 ORDER BY created_at`

// SQL_SELECT_REPORTS_BY_CATEGORY retrieves paginated reports for a specific category.
const SQL_SELECT_REPORTS_BY_CATEGORY = `SELECT * FROM reports WHERE category = $1 LIMIT $2 OFFSET $3 ORDER BY created_at`

// SQL_DELETE_REPORT_BY_ID deletes a report by its ID.
const SQL_DELETE_REPORT_BY_ID = `DELETE FROM reports WHERE id = $1`

// SQL_DAILY_AVERAGED_HOURLY_TRAFFIC_TREND retieves traffic averaged hourly traffic
const SQL_DAILY_AVERAGED_HOURLY_TRAFFIC_TREND string = `
    WITH hourly_max AS (
        SELECT
            DATE(date) AS day,
            EXTRACT(HOUR FROM date) AS hour,
            MAX(bps_in) AS max_bps_in,
            MAX(bps_out) AS max_bps_out,
            MAX(bytes_in_sec) AS max_bytes_in_sec,
            MAX(bytes_out_sec) AS max_bytes_out_sec
        FROM traffic_pons
        GROUP BY day, hour
    )
    SELECT
        day,
        AVG(max_bps_in) / 1e6 AS mbps_in,
        AVG(max_bps_out) / 1e6 AS mbps_out,
        AVG(max_bytes_in_sec) / 1e6 AS mbytes_in_sec,
        AVG(max_bytes_out_sec) / 1e6 AS mbytes_out_sec
    FROM hourly_max
    WHERE date BETWEEN $1 AND $2
    GROUP BY day
    ORDER BY day;`
