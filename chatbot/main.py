import os

from fastapi import FastAPI
from pydantic import BaseModel
import uvicorn

from database import postgres
from utils import ollama, text


class Sql(BaseModel):
    """Request class."""

    text: str


app = FastAPI()


@app.post("/")
async def create_item(sql: Sql):
    """Main function."""
    promp = text.union_string(sql.text, schema)
    ia = ollama.chatbot(promp)
    resust = text.sql_extract(ia)
    temp = postgres.execute_sql(resust)
    return {"pregunta": resust, "respuesta": temp}


if __name__ == "__main__":
    try:
        if os.getenv('URI'):
            schema = postgres.get_database_schema()
            uvicorn.run(app, port=8000)
        else:
            raise Exception("No se encontró la variable URI en el .env")
    except Exception as e:
        print(e)
