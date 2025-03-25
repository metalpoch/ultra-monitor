from os import getenv

from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware

from src import model
from src.database import Postgres
from src.libs.chatbox import AI
from src.libs.osm import Openstreetmap
from src.libs.tracking import Telegram

load_dotenv()

telegram = Telegram(getenv("TELEGRAM_BOT_ID", ""), getenv("TELEGRAM_CHAT_ID", ""))
app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Permite todos los orígenes
    allow_credentials=True,
    allow_methods=["*"],  # Permite todos los métodos HTTP
    allow_headers=["*"],  # Permite todos los headers
)
db = Postgres(getenv("URI", ""))


@app.post("/trend")
async def linear_regression():
    return {"msg": "fooziman"}


@app.post("/chatbox")
async def chatbox(request: model.QueryAI) -> dict:
    count = 0
    ai = AI(model="gemma2:2b", schemas=db.csv_schemas())

    try:
        msg = ai.query(request.message)
        sql = ai.sql_extract(msg)

    except BaseException as e:
        raise HTTPException(status_code=400, detail=str(e))

    while True:
        res, err = db.execute_sql(sql)
        if err is not None:
            return {"response": res, "sql": sql}

        prompt = f"""
        acabo de recibir el siguiente error
        {type(err).__name__}: {str(err)}
        puedes darme la sentencia sql correcta? responde solo con la sentencia sql
        """

        msg = ai.query(prompt)
        sql = ai.sql_extract(msg)

        if count == 3:
            raise HTTPException(status_code=400, detail=msg)

        count += 1


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
