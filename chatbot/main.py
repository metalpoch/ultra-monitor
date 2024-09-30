from functions.ollama import *
from functions.text import *
from database.postgre import *
from fastapi import FastAPI, Request, HTTPException
from pydantic import BaseModel

try:
    schema = Get_database_schema()
except Exception as e:
    print("Error al conectar a la base de datos:", e)

class sql(BaseModel):
    text: str

app= FastAPI()


@app.post("/RodolfIA/")
async def create_item(sql:sql ):
    promp= UnionString(sql.text,schema)
    ia=Chatbot(promp)
    resust=extraer_sql(ia)
    temp= Ejecute_SQL(resust)
    return {"pregunta": resust, "respuesta":temp}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)