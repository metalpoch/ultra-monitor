from datetime import datetime

from pydantic import BaseModel


class Tendencias:
    Device_Id = int
    date = datetime.date
    Out = float
    In = float
    Bandwidth = float

    def __init__(
        self, device: int, date: str, out: float, _in: float, bandwidth: float
    ) -> None:
        self.Device_Id = device
        self.date = date
        self.Out = out
        self.In = _in
        self.Bandwidth = bandwidth
