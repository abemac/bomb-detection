FROM ubuntu:latest

RUN apt-get update
RUN apt-get install git -y
RUN apt-get install golang -y
RUN apt-get install vim -y
RUN echo "MANAGER_IP=\$(/sbin/ip route|awk '/default/ { print \$3 }')" >> /entry.sh
RUN echo "cd \$GOPATH/src/github.com/abemac/bomb-detection/run" >> /entry.sh
RUN echo "go run nodes.go -ip \$MANAGER_IP -ll \$LOG_LEVEL" >> entry.sh

ENTRYPOINT [ "/bin/bash","/entry.sh" ] 