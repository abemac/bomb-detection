FROM ubuntu:latest

RUN apt-get update
RUN apt-get install git -y
RUN apt-get install golang -y
RUN apt-get install vim -y
ENTRYPOINT [ "/bin/bash","/go/src/github.com/abemac/bomb-detection/deploy/entry.sh" ] 