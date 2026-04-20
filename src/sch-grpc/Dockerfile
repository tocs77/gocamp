FROM golang

WORKDIR /workspace/sch-grpc

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum /workspace/sch-grpc/
COPY cert/ /workspace/cert/
RUN go mod download

CMD ["air", "-c", ".air.toml"]
