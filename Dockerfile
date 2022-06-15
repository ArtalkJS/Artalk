### build ArtalkGo
FROM golang:1.18.1-alpine3.15 as builder

WORKDIR /source

# install tools
RUN set -ex \
    && apk upgrade \
    && apk add make git gcc musl-dev nodejs bash npm\
    && npm install -g pnpm@7.2.1

COPY . ./ArtalkGo

# build
RUN set -ex \
    && cd ./ArtalkGo \
    && git fetch --tags -f \
    && export VERSION=$(git describe --tags --abbrev=0) \
    && export COMMIT_SHA=$(git rev-parse --short HEAD) \
    && make all

### build final image
FROM alpine:3.15

# we set the timezone `Asia/Shanghai` by default, you can be modified
# by `docker build --build-arg="TZ=Other_Timezone ..."`
ARG TZ="Asia/Shanghai"

ENV TZ ${TZ}

COPY --from=builder /source/ArtalkGo/bin/artalk-go /artalk-go

RUN apk upgrade \
    && apk add bash tzdata \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

# add alias
RUN echo -e '#!/bin/bash\n/artalk-go -w / -c /data/artalk-go.yml "$@"' > /usr/bin/artalk-go \
    && chmod +x /usr/bin/artalk-go \
    && cp -p /usr/bin/artalk-go /usr/bin/artalk

VOLUME ["/data"]

COPY docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

# expose ArtalkGo default port
EXPOSE 23366

CMD ["server", "--host", "0.0.0.0", "--port", "23366"]
