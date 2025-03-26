import traceback


def Device_id(db):
    query = """
        SELECT id
        FROM devices;
    """
    params = ()
    r = db.execute_fet(query, params)
    return r


def Sys_Name_to_id(db, name: str):
    query = """
        SELECT id
        FROM devices
        WHERE sys_name = %s;
    """
    params = (name,)
    r = db.execute_fet(query, params)
    return r[0][0]


def Id_Device_to_traffic(db, id, init, end):
    query = """
    SELECT
        traffics.date,
        SUM(traffics."in") AS "in",
        SUM(traffics.out) AS out,
        SUM(traffics.bandwidth) AS bandwidth
    FROM
        traffics
    JOIN
        interfaces ON interfaces.id = traffics.interface_id
    WHERE
        interfaces.device_id = %s
        AND traffics.date BETWEEN %s AND %s
    GROUP BY
        traffics.date
    ORDER BY
        traffics.date;
    """
    params = (id, init, end)
    r = db.execute_fet(query, params)
    return r


def NewTendecia(db, data):
    query = """
    INSERT INTO trends (device_id, date, bandwidth, "in", out)
    VALUES (%s, %s, %s, %s, %s);
    """
    try:
        if not data:  # Verificar si la lista está vacía
            print("Error: La lista 'data' está vacía.")
            return False

        for i in data:
            params = (i.Device_Id, i.date, i.Bandwidth, i.In, i.Out)
            db.execute_commit(query, params)
        return True
    except Exception as e:
        print(f"Error al insertar datos: {e}")
        traceback.print_exc()  # Imprime el traceback completo para depuración
        return False


def Trends(db, id):
    query = """
    SELECT * FROM trends WHERE device_id=%s
    """
    params = (id,)
    r = db.execute_fet(query, params)
    return r
