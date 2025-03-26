from datetime import datetime
from os import getenv

from dotenv import load_dotenv

import src.utils.change as change
import src.utils.date as d
import src.utils.execute as ex
from src.database import Postgres

load_dotenv()

db = Postgres(getenv("URI", ""))

lisdevice = []
g = ex.Device_id(db)
for p in g:
    lisdevice.append(p[0])
lisdevice.remove(18)  # para dev

datenow = datetime.now()
monthnow = datenow.month

year = 2025
meses = d.dias_por_mes(year)
day = 1
listten = []
list_dayout = []
i = 0


for j in lisdevice:
    print(
        f"Device {j}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
    )
    deviceId = j
    for i in range(len(meses)):
        if i == monthnow - 1:
            break

        print(
            f"Mes {i+1}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
        )

        while day < meses[i] + 1:
            init, end = d.Transform_date(year, i + 1, day)
            p = ex.Id_Device_to_traffic(db, deviceId, init, end)
            out = change.Respone_to_dict(p)
            list_dayout.append(out)
            day = day + 1
        flag1, flag2, flag3 = change.Sum_Month(list_dayout)
        o = change.Create_tendencia(
            deviceId, d.New_date(year, i + 1, 1), flag2, flag1, flag3
        )
        listten.append(o)

        print(
            f"Fin del Mes {i+1}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
        )

        i = i + 1
        day = 1
        list_dayout = []
    print(
        f"Fin del Device {j}= {datetime.now().hour}:{datetime.now().minute}:{datetime.now().second}"
    )
    
    
    op=ex.NewTendecia(db,listten)
    if op==True:
        print(f"Insertado: {j}")
        print()
    else:
        print(f"No insertado: {j}")
        print()
    listten=[]
db.close()
    
    