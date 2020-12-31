package templates

// Dockerfile is the template for the service's Dockerfile.
const Dockerfile = `FROM golang:1.14 AS build_stage

ARG GITHUB_URL
RUN if [ "$(echo ${GITHUB_URL} | wc -w | xargs)" != "0" ]; then \
        git config --global url."${GITHUB_URL}".insteadOf "https://github.com/"; \
    fi

ENV GONOSUMDB="{{.RepositoryURL | domainUserRepos}}"
ENV CGO_ENABLED=0

WORKDIR /build
COPY go.mod ./
RUN go mod download all

COPY . .
RUN go build -a -installsuffix cgo -o executable ./cmd


FROM scratch AS runtime

WORKDIR /app
COPY --from=build_stage /build/executable cmd
EXPOSE {{.Port}}/tcp

CMD ["/app/cmd"]`
