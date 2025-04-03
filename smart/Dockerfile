FROM python:3.12-alpine

WORKDIR /app

COPY ./smart/ .
RUN pip install --no-cache-dir -r requirements.txt

EXPOSE 3003

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "3003"]
