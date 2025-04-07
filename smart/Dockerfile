FROM python:3.12-alpine

WORKDIR /app

RUN apk add --no-cache g++ python3-dev py3-pip make

COPY ./smart/ .

RUN pip install --upgrade pip

RUN pip install -r requirements.txt

EXPOSE 3003

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "3003"]
