from datetime import datetime
from os import getenv

import click
from dotenv import load_dotenv

from src.utils import execute, date, change
from src.database import Postgres

load_dotenv()


@click.group()
def cli():
    """ """
    pass


@cli.command(help="Gets the data from the traffic by year, processed and inserted.")
@click.option("--year", required=True, help="Year to be processed.")
def year(year: str):
    try:
        year = int(year)
        if year < 2023 or int(datetime.now().strftime("%Y")) < year:
            raise Exception("Error: Year must be greater than 2023.")
    except ValueError:
        print("Error: Year and month must be integers.")
        exit(1)
    except Exception:
        print(f"Error: Invalid data entered.")
        exit(1)

    db = Postgres(getenv("URI", ""))
    devices = []

    res_devices = execute.device_id(db)
    for device in res_devices:
        devices.append(device[0])

    date_now = datetime.now()
    month_now = date_now.month
    months = date.day_per_month(year)
    day = 1
    trends = []
    peak_out_by_days = []

    for j in devices:
        print(
            f"Device {j}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
        )
        device_id = j
        for i in range(0, len(months)):
            if i == month_now - 1:
                break

            print(
                f"Mes {i+1}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
            )

            while day < months[i] + 1:
                init, end = date.transform_date(year=year, month=(i + 1), day=day)
                traffic_by_device = execute.id_device_to_traffic(
                    db, device_id, init, end
                )
                out = change.response_to_dict(traffic_by_device)
                peak_out_by_days.append(out)
                day = day + 1
            sum_in, sum_out, sum_bandwidth = change.sum_month(peak_out_by_days)
            new_trend = change.create_trend(
                device_id, date.new_date(year, i + 1, 1), sum_out, sum_in, sum_bandwidth
            )
            trends.append(new_trend)

            print(
                f"Fin del Mes {i+1}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
            )

            i = i + 1
            day = 1
            peak_out_by_days = []
        print(
            f"Fin del Device {j}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
        )

        status_response = execute.new_trend(db, trends)
        if status_response == True:
            print(f"The device was inserted correctly: {j}\n")
        else:
            print(f"The device was not inserted correctly: {j}\n")
        trends = []

    db.close()


@cli.command(help="Gets the data from the traffic by month, processed and inserted.")
@click.option("--year", required=True, help="Year to be processed.")
@click.option("--month", required=True, help="Month to be processed.")
def month(year: str, month: str):

    try:
        year = int(year)
        if year < 2023 or int(datetime.now().strftime("%Y")) < year:
            raise Exception("Error: Year must be greater than 2023.")
        month = int(month)
        if month > 12 or month < 1:
            raise Exception("Error: Month must be a number between 1 and 12.")
    except ValueError:
        print("Error: Year and month must be integers.")
        exit(1)
    except Exception:
        print(f"Error: Invalid data entered.")
        exit(1)

    db = Postgres(getenv("URI", ""))
    devices = []

    res_devices = execute.device_id(db)
    for device in res_devices:
        devices.append(device[0])

    months = date.day_per_month(year)
    day = 1
    trends = []
    peak_out_by_days = []
    i = month - 1

    for j in devices:
        print(
            f"Device {j}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
        )
        device_id = j
        while i != month:

            print(
                f"Month {i+1}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
            )

            while day < months[i] + 1:
                init, end = date.transform_date(year=year, month=month, day=day)
                traffic_by_device = execute.id_device_to_traffic(
                    db, device_id, init, end
                )
                out = change.response_to_dict(traffic_by_device)
                peak_out_by_days.append(out)
                day = day + 1
            sum_in, sum_out, sum_bandwidth = change.sum_month(peak_out_by_days)
            new_trend = change.create_trend(
                device_id, date.new_date(year, month, 1), sum_out, sum_in, sum_bandwidth
            )
            trends.append(new_trend)

            print(
                f"End month {i+1}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
            )

            i = i + 1
            day = 1
            peak_out_by_days = []
        print(
            f"End Device {j}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
        )

        status_response = execute.new_trend(db, trends)
        if status_response == True:
            print(f"The device was inserted correctly: {j}\n")
        else:
            print(f"The device was not inserted correctly: {j}\n")
        trends = []
        i = month - 1

    db.close()


@cli.command(help="Get traffic data for one device per month, processed and inserted.")
@click.option("--year", required=True, help="Year to be processed.")
@click.option("--month", required=True, help="Month to be processed.")
@click.option("--device", required=True, help="Device to be processed")
def device(year: str, month: str, sysname: str):
    try:
        year = int(year)
        if year < 2023 or int(datetime.now().strftime("%Y")) < year:
            raise Exception("Error: Year must be greater than 2023.")
        month = int(month)
        if month > 12 or month < 1:
            raise Exception("Error: Month must be a number between 1 and 12.")
    except ValueError:
        print("Error: Year and month must be integers.")
        exit(1)
    except Exception:
        print(f"Error: Invalid data entered.")
        exit(1)

    db = Postgres(getenv("URI", ""))
    devices = []
    id = execute.sys_name_to_id(db, sysname)

    res_devices = execute.device_id(db)
    for device in res_devices:
        if device[0] == id:
            devices.append(device[0])

    months = date.day_per_month(year)
    day = 1
    trends = []
    peak_out_by_days = []
    i = month - 1

    for j in devices:
        print(
            f"Device {j}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
        )
        device_id = j
        while i != month:

            print(
                f"Month {i+1}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
            )

            while day < months[i] + 1:
                init, end = date.transform_date(year=year, month=month, day=day)
                traffic_by_device = execute.id_device_to_traffic(
                    db, device_id, init, end
                )
                out = change.response_to_dict(traffic_by_device)
                peak_out_by_days.append(out)
                day = day + 1
            sum_in, sum_out, sum_bandwidth = change.sum_month(peak_out_by_days)
            new_trend = change.create_trend(
                device_id, date.new_date(year, month, 1), sum_out, sum_in, sum_bandwidth
            )
            trends.append(new_trend)

            print(
                f"End month {i+1}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
            )

            i = i + 1
            day = 1
            peak_out_by_days = []
        print(
            f"End Device {j}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
        )

        status_response = execute.new_trend(db, trends)
        if status_response == True:
            print(f"The device was inserted correctly: {j}\n")
        else:
            print(f"The device was not inserted correctly: {j}\n")
        trends = []
        i = month - 1

    db.close()


if __name__ == "__main__":
    cli()
    device(year="2025", month="1", sysname="prp-olt-00")
