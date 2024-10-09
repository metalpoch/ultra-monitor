from pydantic import BaseModel


class Sql(BaseModel):
    """Request class."""

    text: str
