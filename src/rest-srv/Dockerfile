FROM golang
WORKDIR /app
RUN apt-get update && apt-get install -y wrk && rm -rf /var/lib/apt/lists/*
RUN go install github.com/air-verse/air@latest
COPY go.mod ./

RUN go mod download


# COPY . .
CMD [ "air", "-c", ".air.toml" ]
# CMD [ "air" ]