echo "**Building node simulator go executable..."
GOPATH=/mnt/c/Users/abrah/go
go build ../../run/nodes.go > /dev/null

echo "**Sending files to server..."
scp -P 2005 Dockerfile entry.sh nodes abraham@128.4.27.184:/home/abraham/nodesim/  > /dev/null

echo "**Building Docker image on server..."
ssh -p 2005 -t abraham@128.4.27.184 docker build -t nodesim ./nodesim  > /dev/null

echo "**Starting nodesim image on Server..."
ssh -p 2005 -t abraham@128.4.27.184 docker run -d --rm --name nodesim1 -e LOG_LEVEL=4 nodesim

rm nodes