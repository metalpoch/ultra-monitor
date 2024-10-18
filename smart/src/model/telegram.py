from pydantic import BaseModel


class Telegram(BaseModel):
    event: str
    module: str
    message: str
    category: str
