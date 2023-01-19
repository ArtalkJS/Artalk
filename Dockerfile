### build Artalk
FROM golang:1.19.4-alpine3.17 as builder

WORKDIR /source

# install tools
RUN set -ex \
    && apk add --no-cache make git gcc musl-dev nodejs bash npm\
    && npm install -g pnpm@7.25.0

COPY . ./Artalk

# build
RUN set -ex \
    && cd ./Artalk \
    && export VERSION=$(git describe --tags --abbrev=0) \
    && export COMMIT_SHA=$(git rev-parse --short HEAD) \
    && make all

### build final image
FROM alpine:3.15

# we set the timezone `Asia/Shanghai` by default, you can be modified
# by `docker build --build-arg="TZ=Other_Timezone ..."`
ARG TZ="Asia/Shanghai"

ENV TZ ${TZ}

COPY --from=builder /source/Artalk/bin/artalk /artalk

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
