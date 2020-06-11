package templates

// Dockerfile is the template for the service's Dockerfile.
const Dockerfile = `FROM golang:1.14 AS build_stage

WORKDIR /build
COPY go.mod ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0

RUN go build -a -installsuffix cgo -o executable ./cmd


FROM scratch AS runtime

WORKDIR /app
COPY --from=build_stage /build/executable cmd
EXPOSE {{.Port}}/tcp

CMD ["/app/cmd"]`
