FROM golang:1.18-alpine as build-stage

ENV CGO_ENABLED 0
ENV PROJECT_PACKAGE github.com/viniscoelho/stocksbot
ENV OBJ_NAME stocksbot
ARG GIT_COMMIT

COPY . /go/src/${PROJECT_PACKAGE}/

RUN cd /go/src/${PROJECT_PACKAGE} && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod vendor -o ${OBJ_NAME} -ldflags "-X main.gitCommit=${GIT_COMMIT}"

# --------------------------------------------------------------------------------
FROM alpine:3.11

ENV PROJECT_PACKAGE github.com/viniscoelho/stocksbot
ENV OBJ_NAME stocksbot

RUN apk update && apk add tzdata
ENV TZ America/Sao_Paulo
RUN ln -s /usr/share/zoneinfo/$TZ /etc/localtime && echo "$TZ" > /etc/timezone

RUN adduser -D ${OBJ_NAME}
USER ${OBJ_NAME}

COPY --from=build-stage /go/src/${PROJECT_PACKAGE}/${OBJ_NAME} /usr/local/bin/${OBJ_NAME}
CMD /usr/local/bin/${OBJ_NAME}