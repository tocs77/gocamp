FROM golang
WORKDIR /workspace/grpc-client
RUN go install github.com/air-verse/air@latest
COPY go.work /workspace/go.work
COPY grpc-client/go.mod grpc-client/go.sum /workspace/grpc-client/
COPY grpc-srv/go.mod grpc-srv/go.sum /workspace/grpc-srv/

RUN go mod download

CMD [ "air", "-c", ".air.toml" ]
