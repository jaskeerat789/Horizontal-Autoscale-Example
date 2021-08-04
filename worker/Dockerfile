FROM node:lts-alpine 

WORKDIR /app

COPY package.json package.json
COPY package-lock.json package-lock.json

RUN npm install -- silent 

COPY app.js app.js

CMD ["node","app.js"]