### build Artalk
FROM golang:1.21.5-alpine3.19 as builder

WORKDIR /source

# install tools
RUN set -ex \
    && apk add --no-cache make git gcc musl-dev bash

# download go deps
# (cache by separating the downloading of deps)
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

## build UI
ARG SKIP_UI_BUILD=false

# install ui build toolchain
RUN set -ex \
    && if [ "$SKIP_UI_BUILD" = "false" ]; then \
        apk add --no-cache nodejs npm \
        && npm install -g pnpm@8.12.1 \
    ;fi

RUN set -ex \
    && if [ "$SKIP_UI_BUILD" = "false" ]; then \
        make build-frontend \
    ;fi

## build App
ARG APP_VERSION=""
ARG APP_COMMIT_HASH=""

RUN set -ex \
    && if [[ ! -z "$APP_VERSION" ]]; then export VERSION=$APP_VERSION ;fi \
    && if [[ ! -z "$APP_COMMIT_HASH" ]]; then export COMMIT_HASH=$APP_COMMIT_HASH ;fi \
    && make build

### build final image
FROM alpine:3.19

# we set the timezone `Asia/Shanghai` by default, you can be modified
# by `docker build --build-arg="TZ=Other_Timezone ..."`
ARG TZ="Asia/Shanghai"

ENV TZ ${TZ}

COPY --from=builder /source/bin/artalk /artalk

RUN apk add --no-cache bash tzdata \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

# move runner script to `/usr/bin/` and create alias
COPY scripts/docker-artalk-runner.sh /usr/bin/artalk
RUN chmod +x /usr/bin/artalk \
    && ln -s /usr/bin/artalk /usr/bin/artalk-go

VOLUME ["/data"]

COPY docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

# expose Artalk default port
EXPOSE 23366

CMD ["server", "--host", "0.0.0.0", "--port", "23366"]
