import time

import requests


class Openstreetmap:
    def __init__(self, latitude: float, longitude: float) -> None:
        self.latitude = latitude
        self.longitude = longitude

    def location(self) -> dict[str, str]:
        c = 0
        while c < 3:
            c += 1
            res = requests.get(
                f"https://nominatim.openstreetmap.org/reverse?lat={self.latitude}&lon={self.longitude}&format=json",
                headers={
                    "User-Agent": "CANTV Analizator/1.0 (https://cantv.com.ve)",
                },
            )

            if res.status_code == 200 and "application/json" in res.headers.get(
                "Content-Type", ""
            ):
                data = res.json()
                return data["address"]

            time.sleep(1.5)

        return {}
