LOG_LEVEL=4
docker run -d --rm --name nodesim1 -v $GOPATH:/go -e GOPATH=/go -e LOG_LEVEL=$LOG_LEVEL nodesim
docker run -d --rm --name nodesim2 -v $GOPATH:/go -e GOPATH=/go -e LOG_LEVEL=$LOG_LEVEL nodesim
docker run -d --rm --name nodesim3 -v $GOPATH:/go -e GOPATH=/go -e LOG_LEVEL=$LOG_LEVEL nodesim
docker run -d --rm --name nodesim4 -v $GOPATH:/go -e GOPATH=/go -e LOG_LEVEL=$LOG_LEVEL nodesim
docker run -d --rm --name nodesim5 -v $GOPATH:/go -e GOPATH=/go -e LOG_LEVEL=$LOG_LEVEL nodesim
docker run -d --rm --name nodesim6 -v $GOPATH:/go -e GOPATH=/go -e LOG_LEVEL=$LOG_LEVEL nodesim
docker run -d --rm --name nodesim7 -v $GOPATH:/go -e GOPATH=/go -e LOG_LEVEL=$LOG_LEVEL nodesim
docker run -d --rm --name nodesim8 -v $GOPATH:/go -e GOPATH=/go -e LOG_LEVEL=$LOG_LEVEL nodesim
docker run -d --rm --name nodesim9 -v $GOPATH:/go -e GOPATH=/go -e LOG_LEVEL=$LOG_LEVEL nodesim
docker run -d --rm --name nodesim10 -v $GOPATH:/go -e GOPATH=/go -e LOG_LEVEL=$LOG_LEVEL nodesim