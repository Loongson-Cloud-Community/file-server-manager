FROM cr.loongnix.cn/library/debian:buster-slim

COPY dist/file-server-manager /usr/bin/fsm

RUN apt update && apt install -y git

VOLUME ["/data"]

CMD ["/usr/bin/fsm", "-data=/data"]
