# -*- coding: utf-8 -*-
# vim: ft=Dockerfile

### container - builder
FROM golang:1.26.0-bookworm AS build
LABEL maintainer="mindhunter86 <mindhunter86@vkom.cc>"

ARG GOAPP_MAIN_VERSION="devel"
ARG GOAPP_MAIN_BUILDTIME="N/A"

ENV MAIN_VERSION=$GOAPP_MAIN_VERSION
ENV MAIN_BUILDTIME=$GOAPP_MAIN_BUILDTIME

ENV NODE_ENV=production
ENV NODE_VER=20.19

ENV DEBIAN_FRONTEND=noninteractive

# hadolint/hadolint - DL4006
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /usr/sources/eyesonly
COPY . .

# skipcq: DOK-DL3008 pinning version for upx is not required
RUN echo "ready" \
  && echo "" \
  && echo ">>> Some preparations before building... <<<" \
  && echo "deb https://deb.debian.org/debian bookworm-backports main contrib non-free-firmware" > /etc/apt/sources.list.d/debian-12-backports.list \
  && apt-get update && apt-get install --no-install-recommends -t bookworm-backports -y upx-ucl pkg-config libssl-dev build-essential \
  && rm -rf /var/lib/apt/lists/* \
  && ls -la . && ls -la ./internal/web/dist \
  && find ./internal/web/dist -type f -name "*.js.map" -delete \
  && echo "" \
  && echo ">>> Buildtools installation completed, executing go build... <<<" \
  && go mod download \
  && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags netgo,osusergo -trimpath -ldflags="-s -w -X 'main.version=$MAIN_VERSION' -X 'main.buildtime=$MAIN_BUILDTIME'" -o eyesonly cmd/eyesonly/*.go \
  && upx -9 -k eyesonly \
  && ls -lah eyesonly

  # disabled PGO
  # && go tool pprof -proto extras/pgo/*.pprof > merged.pprof \
  # && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -pgo=merged.pprof -tags netgo,osusergo -trimpath -ldflags="-s -w -X 'main.version=$MAIN_VERSION' -X 'main.buildtime=$MAIN_BUILDTIME'" -o eyesonly cmd/eyesonly/*.go \


### container - runner
###   for image debugging use tag :debug
FROM gcr.io/distroless/cc-debian12:latest-amd64
LABEL maintainer="mindhunter86 <mindhunter86@vkom.cc>"

WORKDIR /usr/local/bin/
COPY --from=build --chmod=0555 /usr/sources/eyesonly/eyesonly eyesonly

USER nobody
ENTRYPOINT ["/usr/local/bin/eyesonly"]
CMD ["--help"]
