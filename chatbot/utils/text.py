import os
import re

from dotenv import load_dotenv

load_dotenv(override=True)


def Union_string(msg: str, schema: str) -> str:
    """Join the strings to create the promt."""
    p1 = os.getenv("PROMP1")
    p2 = os.getenv("PROMP2")
    p3 = os.getenv("PROMP3")
    p4 = os.getenv("PROMP4")
    if msg == "":
        return "Error en el mensaje del cliente"

    return p1 + schema + ". " + p2 + msg + ". " + p3 + p4


def Sql_extract(msg: str):
    """Extract the SQL from the string of the AI response."""
    pattern = r"```sql(.*?)```"
    result = re.search(pattern, msg, re.DOTALL)

    if result:
        return result.group(1).strip()
    else:
        return None


def SQL_revision(sql: str) -> str:
    """Check that the SQL only has one ”;”."""
    contador = sql.count(';')
    if contador >= 2:
        i = sql.rfind(';')
        if i == -1:
            return sql
        new_sql = sql.replace(";", "")+";"
        return new_sql
    return sql

def SQL_check(sql: str)-> str:
   p1 = os.getenv("PROM5")
   p2 = os.getenv("PROM6") 
   p3 = os.getenv("PROM7")
   return p1 + sql + p2 + p3


def Delete_table(texto, nombre_tabla):
    patron_tabla = f'Tabla: {nombre_tabla}.*?(?=Tabla:|$)'
    texto_modificado = re.sub(patron_tabla, '', texto, flags=re.DOTALL)
    texto_modificado = ' '.join(texto_modificado.split())
    return texto_modificado