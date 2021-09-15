# Create the building image for compiling
FROM node:alpine as build

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run build

# Create the final image with the http server
FROM nginx:alpine as exec

COPY  --from=build /app/dist /usr/share/nginx/html
