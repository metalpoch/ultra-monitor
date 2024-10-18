from os import getenv

from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException

from src import model
from src.database import Postgres
from src.libs.chatbox import AI
from src.libs.osm import Openstreetmap
from src.libs.tracking import Telegram

load_dotenv()

app = FastAPI()
ai = AI(model="gemma2:2b")
db = Postgres(getenv("URI", ""))
telegram = Telegram(getenv("TELEGRAM_BOT_ID", ""), getenv("TELEGRAM_CHAT_ID", ""))


@app.post("/trend")
async def linear_regression():
    return {"msg": "fooziman"}


@app.post("/chatbox")
async def chatbox(request: model.QueryAI) -> dict:
    try:
        schemas = db.csv_schemas(conn=db.connect())
        msg = ai.query(schemas=schemas, body=request.message)
        # print(msg)  # for develop
        sql = ai.sql_extract(msg)

    except BaseException as e:
        raise HTTPException(status_code=400, detail=str(e))

    else:
        return {"response": sql}


@app.post("/telegram")
async def tracking(req: model.Telegram) -> dict:
    try:
        res = telegram.send_message(
            module=req.module, event=req.event, category=req.category, msg=req.message
        )

    except BaseException as e:
        raise HTTPException(status_code=400, detail=str(e))

    else:
        return {"response": res}


@app.get("/location")
async def location(latitude: float, longitude: float) -> dict:
    try:
        res = Openstreetmap(latitude, longitude).location()

    except BaseException as e:
        raise HTTPException(status_code=400, detail=str(e))

    else:
        return res
