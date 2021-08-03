FROM node:lts-alpine 

WORKDIR /app

COPY package.json /app/package.json

RUN npm install -- silent 

COPY app.js /app/app.js

CMD ["npm","run"]