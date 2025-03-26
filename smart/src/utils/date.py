from datetime import datetime


def Transform_date(year, month, day):
    newday = ""
    newmonth = ""
    if day < 10:
        newday = "0" + str(day)
    else:
        newday = str(day)
    if month < 10:
        newmonth = "0" + str(month)
    else:
        newmonth = str(day)
    init = f"{year}-{newmonth}-{newday}T00:00:00-04:00"
    end = f"{year}-{newmonth}-{newday}T23:59:59-04:00"

    return init, end


def New_date(year, month, day):
    newday = ""
    newmonth = ""
    if day < 10:
        newday = "0" + str(day)
    else:
        newday = str(day)
    if month < 10:
        newmonth = "0" + str(month)
    else:
        newmonth = str(day)
    init = f"{year}-{newmonth}-{newday}T00:00:00-04:00"
    return datetime.fromisoformat(init)


def dias_por_mes(año):
    dias_mes = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31]

    if (año % 4 == 0 and año % 100 != 0) or año % 400 == 0:
        dias_mes[1] = 29  # Febrero tiene 29 días en año bisiesto

    return dias_mes
