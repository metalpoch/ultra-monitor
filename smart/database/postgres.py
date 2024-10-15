import os

from sqlalchemy import MetaData, create_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.sql import text
from fastapi import HTTPException
from dotenv import load_dotenv

import utils.text
import utils.ollama

load_dotenv(override=True)


def Get_schema():
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


def SQL_exe(input: str):
    """Execute the SQL statement received by parameter in the database."""
    connection_string = os.getenv('URI')

    engine = create_engine(connection_string)
    session_maker = sessionmaker(bind=engine)
    session = session_maker()

    temp = utils.text.SQL_revision(input)
    check = utils.text.SQL_check(input)
    ia_check = utils.ollama.Chatbot(check)
    try:
        if ia_check == "False \n":
            sql_text = text(temp)
            try:
                result = session.execute(sql_text)
            except Exception:
                return HTTPException(status_code=503,
                                     detail="No se puede procesar por el momento su solicitud")
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
        elif ia_check == "False \n":
            raise ValueError("Esta intentando realizar una acción prohibida.")
    except ValueError as e:
        raise HTTPException(status_code=451, detail=str(e))
