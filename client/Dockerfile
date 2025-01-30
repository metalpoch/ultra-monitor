FROM node:lts-alpine3.20

# Crear un directorio para la aplicación
RUN mkdir /app
WORKDIR /app

# Copiar los archivos de la aplicación
COPY . .

# Instalar las dependencias
RUN npm install

# Exponer el puerto 4321
EXPOSE 4321

# Comando para iniciar la aplicación
CMD ["npm", "run", "start", "--","--host"]
