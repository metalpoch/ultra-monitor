FROM node:lts-alpine3.20

WORKDIR /workspace

COPY client/ .

# Instalar las dependencias
RUN npm install

# Exponer el puerto 4321
EXPOSE 4321

# Comando para iniciar la aplicación
CMD ["npm", "run", "start", "--","--host"]
