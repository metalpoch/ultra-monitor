FROM python:3.12-alpine

RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    g++ \
    python3-dev \
    libatlas-base-dev \
    && rm -rf /var/lib/apt/lists/*


WORKDIR /app

COPY ./smart/ .

RUN pip install --upgrade pip

RUN pip install --no-cache-dir -r requirements.txt

EXPOSE 3003

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "3003"]
