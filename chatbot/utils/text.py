import os
import re

from dotenv import load_dotenv

load_dotenv(override=True)


def union_string(msg: str, schema: str) -> str:
    p1 = os.getenv("PROMP1")
    p2 = os.getenv("PROMP2")
    p3 = os.getenv("PROMP3")
    p4 = os.getenv("PROMP4")
    if msg == "":
        return "Error en el mensaje del cliente"

    return p1 + schema + ". " + p2 + msg + ". " + p3 + p4


def sql_extract(msg: str):
    pattern = r"```sql(.*?)```"
    result = re.search(pattern, msg, re.DOTALL)

    if result:
        return result.group(1).strip()
    else:
        return None

def revisar_sql(sql: str)->str:
    contador= sql.count(';')
    if contador >= 2:
        i= sql.rfind(';')
        if i == -1:
            return sql
        newSql= sql.replace(";","")+";"
        return newSql
    return sql