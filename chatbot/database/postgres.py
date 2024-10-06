import utils
import os
from sqlalchemy import MetaData, create_engine, inspect
from sqlalchemy.orm import sessionmaker
from sqlalchemy.sql import text
from dotenv import load_dotenv
load_dotenv(override=True)
import utils.text


def get_database_schema():
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
    connection_string = os.getenv('URI')

    engine = create_engine(connection_string)
    session = sessionmaker(bind=engine)
    Session = session()

    temp=utils.text.revisar_sql(input)
    sql_text= text(temp)

    try:
        result=Session.execute(sql_text)
    except:
        return "No se pudo procesar su pregunta intentelo mas tarde"   
        
    column_names = result.keys()
    data = []
    i = 0
    lista=[]
    for j in result:
        lista.append(list(j))
    for column in column_names:
        data_dict = { "column": column, "values": []}
        for row in lista:
            data_dict['values'].append(row[i])
        i += 1 
        data.append(data_dict)
    return data