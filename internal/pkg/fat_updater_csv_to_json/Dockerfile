# Usa una imagen base con una versión de glibc más antigua, como Ubuntu 20.04
FROM ubuntu:20.04

# Instala Python y PyInstaller
RUN apt-get update && apt-get install -y python3 python3-pip

# Copia tu código a la imagen
WORKDIR /app
COPY . .

RUN pip3 install pyinstaller pandas

# Compila tu script con el comando original
CMD pyinstaller -F -n fats-csv-to-json --optimize 2 main.py
