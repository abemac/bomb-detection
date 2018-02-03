FROM ubuntu:latest
RUN apt-get update
COPY ./nodes /nodes
COPY ./entry.sh /entry.sh
ENTRYPOINT [ "/bin/bash","/entry.sh" ]
