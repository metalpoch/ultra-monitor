from datetime import datetime


class Trend:
    device_id = int
    date = datetime.date
    out = float
    in_ = float
    bandwidth = float

    def __init__(
        self, device: int, date: str, out: float, _in: float, bandwidth: float
    ) -> None:
        self.device_id = device
        self.date = date
        self.out = out
        self.in_ = _in
        self.bandwidth = bandwidth
