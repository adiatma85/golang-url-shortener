# Build tiny docker image
FROM alpine:3.12

WORKDIR /project

# hadolint ignore=DL3018
RUN apk add --no-cache \
        libstdc++ \
        tzdata \
        libx11 \
        libxrender \
        libxext \
        libssl1.1 \
        ca-certificates;

COPY ./build/app ./app
COPY ./docs docs
COPY ./etc/cfg etc/cfg
COPY ./etc/tpl etc/tpl

CMD [ "./app" ]
