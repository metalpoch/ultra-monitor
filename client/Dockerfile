FROM node:lts-alpine3.20

WORKDIR /app

COPY client/ .

RUN npm install
RUN npm run build

ENV PORT=4321

EXPOSE 4321
CMD node ./dist/server/entry.mjs
