FROM golang:1.26-alpine

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /gen/sse .


FROM scratch

COPY --from=0 /gen/sse /sse
ENTRYPOINT [ "/sse" ]