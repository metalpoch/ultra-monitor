FROM python:alpine

WORKDIR /app

COPY ./smart/requirements.txt .

RUN python -m venv ./venv
RUN source ./venv/bin/activate
RUN pip install --no-cache-dir -r requirements.txt

COPY  . .
EXPOSE 8000

CMD ["fastapi, run", "./smart/main.py"]
