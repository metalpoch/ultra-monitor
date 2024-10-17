import os

from fastapi import FastAPI, HTTPException
import uvicorn

from database import postgres
from utils import ollama, text
from models import request


app = FastAPI()


@app.post("/")
async def create_item(sql: request.Sql):
    """Main function."""
    try:
        if not sql.text:
            raise ValueError("Para continuar debe ingresar un valor")
        else:
            if not schema:
                raise ValueError("Schema empty")
            else:
                newschema = text.delete_table(schema, "users")
                promp = text.union_string(sql.text, newschema)
                ia = ollama.chatbot(promp)
                resust = text.sql_extract(ia)
                temp = postgres.SQL_exe(resust)

                # Esta funcion es la encargada de borrar una tabla en especifico del schema, es generico. se debe colocar el nombre exacto de la
                # tabla a eliminar, por ahora esta en prueba
                # temp=text.delete_table(schema,'clientes')

                return {"pregunta": resust, "respuesta": temp}
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))


if __name__ == "__main__":
    try:
        if os.getenv("URI"):
            schema = postgres.Get_schema()
            uvicorn.run(app, port=8000)
        else:
            raise Exception("No se encontró la variable URI en el .env")
    except Exception as e:
        print(e)
