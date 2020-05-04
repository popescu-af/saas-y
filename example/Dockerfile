FROM golang:1.14 AS build_stage

WORKDIR /build
COPY . .
ENV CGO_ENABLED=0

RUN go build -a -installsuffix cgo -o executable ./cmd


FROM scratch AS runtime

WORKDIR /app
COPY --from=build_stage /build/executable cmd
EXPOSE 8080/tcp

CMD ["/app/cmd"]
