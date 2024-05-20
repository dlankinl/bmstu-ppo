# FROM node:lts-alpine as build-stage

# RUN npm install -g http-server

# WORKDIR /app

# COPY package*.json ./

# RUN npm install

# COPY . .

# ENV HOST 0.0.0.0
# EXPOSE 80

# RUN npm run build

# CMD [ "npm", "run", "dev"]



FROM node:latest as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY ./ .
RUN npm run build

FROM nginx as production-stage
RUN mkdir /app
COPY --from=build-stage /app/dist /app
COPY nginx.conf /etc/nginx/nginx.conf