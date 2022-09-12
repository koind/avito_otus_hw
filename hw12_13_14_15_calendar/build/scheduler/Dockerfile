FROM golang:1.17.3 as build

ENV BIN_FILE /opt/calendar/sheduler
ENV CODE_DIR /go/src/
ENV CGO_ENABLED=0

RUN go install github.com/antelman107/net-wait-go@latest

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/scheduler/*

FROM alpine:3.15

ENV BIN_FILE /opt/calendar/sheduler
COPY --from=build ${BIN_FILE} ${BIN_FILE}

COPY --from=build /go/bin/net-wait-go /bin/net-wait-go

CMD ${BIN_FILE}