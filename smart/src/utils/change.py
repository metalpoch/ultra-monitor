import math

import pytz as tz

from src.model.trends import Tendencias as te


def Bits_to__GBits(valor: float):
    conversion = valor / 1000000000
    if conversion < 0.01:
        redondeo = Redondear(conversion, 5)
        return redondeo
    else:
        redondeo = Redondear(conversion, 2)
        return redondeo


def Redondear(n: float, decimals: int = 0) -> float:
    expoN = n * 10**decimals
    if abs(expoN) - abs(math.floor(expoN)) < 0.5:
        return math.floor(expoN) / 10**decimals
    return math.ceil(expoN) / 10**decimals


def Respone_to_dict(list_date):
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


def Sum_Month(listdict: list):
    sumIn = 0
    sumOut = 0
    ban = 0

    for i in listdict:
        sumIn = sumIn + i["In"]
        sumOut = sumOut + i["Out"]
        if i["Bandwidth"] != 0:
            ban = i["Bandwidth"]

    return sumIn / len(listdict), sumOut / len(listdict), ban / len(listdict)


def Create_tendencia(device, date, out, in_, bandwidth):
    tendencia = te(device, date, float(out), float(in_), float(bandwidth))
    return tendencia
