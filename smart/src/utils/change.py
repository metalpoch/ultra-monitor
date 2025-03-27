import math

import pytz as tz

from src.model.trends import Trend


def bits_to_gbits(value: float):
    new_value = value / 1000000000
    return new_value

    
def gbits_to_bits(value: float):
    new_value = value * 1000000000
    return new_value


def round_up(n: float, decimals: int = 0) -> float:
    expoN = n * 10**decimals
    if abs(expoN) - abs(math.floor(expoN)) < 0.5:
        return math.floor(expoN) / 10**decimals
    return math.ceil(expoN) / 10**decimals


def response_to_dict(list_date):
    sumOut = {"Fecha": "", "In": 0, "Out": 0, "Bandwidth": 0}
    flaOut = 0
    for i in list_date:
        if i[2] > flaOut:
            sumOut["Fecha"] = i[0].astimezone(tz.timezone("Etc/GMT+4"))
            sumOut["In"] = i[1]
            sumOut["Out"] = i[2]
            sumOut["Bandwidth"] = i[3]
            flaOut = i[2]
    return sumOut


def sum_month(listdict: list):
    sumIn = 0
    sumOut = 0
    ban = 0

    for i in listdict:
        sumIn = sumIn + i["In"]
        sumOut = sumOut + i["Out"]
        if i["Bandwidth"] != 0:
            ban = i["Bandwidth"]

    return sumIn / len(listdict), sumOut / len(listdict), ban / len(listdict)


def create_trend(device, date, out, in_, bandwidth):
    trend = Trend(device, date, float(out), float(in_), float(bandwidth))
    return trend
