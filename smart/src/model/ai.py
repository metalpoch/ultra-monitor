from pydantic import BaseModel


class QueryAI(BaseModel):
    message: str
