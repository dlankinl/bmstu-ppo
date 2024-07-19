# FROM golang:1.21.1-bullseye

# WORKDIR /app

# COPY go.mod go.sum .

# RUN go mod download

# COPY . .

# RUN go build -o ./bin/server .

# # RUN echo "$PWD"
# # RUN ls -la bin

# ENTRYPOINT [./bin/server]

FROM golang:1.21.1-bullseye
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server

ENTRYPOINT [ "/server" ]