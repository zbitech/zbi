FROM --platform=linux/amd64 node:18

WORKDIR /usr/src/app

COPY package*.json ./ 
COPY tsconfig.json ./ 

RUN npm install 

COPY . .

RUN npm run build

EXPOSE 4000

RUN ls -al /usr/src/app
CMD ["node", "./dist/index.js"]
