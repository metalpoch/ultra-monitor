import os

from sqlalchemy import MetaData, create_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.sql import text
from dotenv import load_dotenv

import utils.text

load_dotenv(override=True)


def get_database_schema():
    """Return the database schema."""
    connection_string = os.getenv('URI')

    engine = create_engine(connection_string)
    metadata = MetaData()
    metadata.reflect(bind=engine)

    schema_info = ""
    for table_name, table in metadata.tables.items():
        schema_info += f"Tabla: {table_name} "
        schema_info += "Columnas: "
        for column in table.columns:
            schema_info += f" {column.name}: {column.type}"
        schema_info += " "

    return schema_info.strip()


def execute_sql(input: str):
    """Execute the SQL statement received by parameter in the database."""
    connection_string = os.getenv('URI')

    engine = create_engine(connection_string)
    session_maker = sessionmaker(bind=engine)
    session = session_maker()

    temp = utils.text.revisar_sql(input)
    sql_text = text(temp)

    try:
        result = session.execute(sql_text)
    except Exception as e:
        return "No se pudo procesar su pregunta intentelo mas tarde.", e
    column_names = result.keys()
    data = []
    i = 0
    lista = []
    for j in result:
        lista.append(list(j))
    for column in column_names:
        data_dict = {"column": column, "values": []}
        for row in lista:
            data_dict['values'].append(row[i])
        i += 1
        data.append(data_dict)
    return data
