### build ArtalkGo
FROM golang:1.17.2-alpine3.13 as builder

WORKDIR /source

# install tools
RUN set -ex \
    && apk upgrade \
    && apk add make git gcc musl-dev yarn

# 1. build frontend (using latest tag)
RUN git clone https://github.com/ArtalkJS/Artalk.git Artalk \
    && cd Artalk \
    && git fetch --tags \
    && git checkout $(git describe --tags --abbrev=0)

RUN set -ex \
    && cd ./Artalk \
    && yarn install --network-timeout 1000000 \
    && yarn build

# 2. build backend
COPY . ./ArtalkGo

# copy frontend dist folder
RUN cp -r ./Artalk/dist ./ArtalkGo/frontend

RUN set -ex \
    && cd ./ArtalkGo \
    && export VERSION=$(git describe --tags --abbrev=0) \
    && export COMMIT_SHA=$(git rev-parse --short HEAD) \
    && go mod tidy \
    && go install github.com/markbates/pkger/cmd/pkger \
    && pkger -include /frontend -include /email-tpl -o pkged \
    && go build -ldflags "-X github.com/ArtalkJS/ArtalkGo/lib.Version=${VERSION} \
                          -X github.com/ArtalkJS/ArtalkGo/lib.LastCommit=${COMMIT_SHA}" \
                          -o bin/artalk-go \
                          github.com/ArtalkJS/ArtalkGo

### build final image
FROM alpine:3.13

# we set the timezone `Asia/Shanghai` by default, you can be modified
# by `docker build --build-arg="TZ=Other_Timezone ..."`
ARG TZ="Asia/Shanghai"

ENV TZ ${TZ}

COPY --from=builder /source/ArtalkGo/bin/artalk-go /artalk-go

RUN apk upgrade \
    && apk add bash tzdata \
    && ln -s /artalk-go /usr/bin/artalk-go \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

VOLUME ["/conf.yml", "/data"]

ENTRYPOINT ["/artalk-go"]

# expose ArtalkGo default port
EXPOSE 23366

CMD ["serve", "--config", "/conf.yml", "--host", "0.0.0.0", "--port", "23366"]
