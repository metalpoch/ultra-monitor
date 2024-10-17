from os import getenv

import uvicorn
from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException

from src import model
from src.database import Postgres
from src.libs.chatbox import AI

load_dotenv()

app = FastAPI()
ai = AI(model="gemma2:2b")
db = Postgres(getenv("URI", ""))


@app.post("/trend")
async def linear_regression():
    return {"msg": "fooziman"}


@app.post("/chatbox")
async def chatbox(request: model.QueryAI):
    try:
        schemas = db.csv_schemas(conn=db.connect())
        msg = ai.query(schemas=schemas, body=request.message)
        # print(msg)  # for develop
        sql = ai.sql_extract(msg)

    except BaseException as e:
        raise HTTPException(status_code=400, detail=str(e))

    else:
        return {"response": sql}
