from datetime import datetime


def transform_date(year, month, day):
    new_day = ""
    new_month = ""
    if day < 10:
        new_day = "0" + str(day)
    else:
        new_day = str(day)
    if month < 10:
        new_month = "0" + str(month)
    else:
        new_month = str(day)
    init = f"{year}-{new_month}-{new_day}T00:00:00-04:00"
    end = f"{year}-{new_month}-{new_day}T23:59:59-04:00"

    return init, end


def new_date(year, month, day):
    new_day = ""
    new_month = ""
    if day < 10:
        new_day = "0" + str(day)
    else:
        new_day = str(day)
    if month < 10:
        new_month = "0" + str(month)
    else:
        new_month = str(day)
    init = f"{year}-{new_month}-{new_day}T00:00:00-04:00"
    return datetime.fromisoformat(init)


def day_per_month(year):
    months = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31]

    if (year % 4 == 0 and year % 100 != 0) or year % 400 == 0:
        months[1] = 29

    return months
