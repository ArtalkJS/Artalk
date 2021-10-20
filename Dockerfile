FROM golang:alpine as builder

ENV ARTALK_VERSION v2.0

RUN apk --no-cache add make git gcc musl-dev
WORKDIR /source
COPY . .
RUN go mod tidy \
    && go install github.com/markbates/pkger/cmd/pkger \
    && pkger -include /frontend -include /email-tpl -o pkged \
    && go build -ldflags "-X github.com/ArtalkJS/ArtalkGo.Version=${ARTALK_VERSION}" -o bin/artalk-go github.com/ArtalkJS/ArtalkGo

FROM alpine

COPY --from=builder /source/bin/artalk-go /artalk-go

VOLUME ["/conf.yml", "/local"]

ENTRYPOINT ["/artalk-go", "serve"]

CMD ["--config", "/conf.yml"]