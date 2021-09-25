FROM namely/protoc-all:1.37_2

## Add git to be able to download repositories.
RUN set -ex && apt-get update && apt-get install -y --no-install-recommends \
    git

ADD gpm /usr/local/bin
RUN chmod +x /usr/local/bin/gpm

WORKDIR /defs
ENTRYPOINT [ "gpm", "generate", "--protoGenerator=dockerized", "." ]
