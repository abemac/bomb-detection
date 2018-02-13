export GOPATH="/mnt/c/Users/abrah/go"

args="$@"

if [[ "$args" = *"nodes"* ]] || [[ "$args" = *"all"* ]]; then
    echo "**Building node simulator go executable..."
    go build ../../run/nodes.go > /dev/null

    echo "**Sending files to server..."
    scp -P 2005 -o "LogLevel=QUIET" Dockerfile entry.sh nodes abraham@128.4.27.184:/home/abraham/nodesim/  > /dev/null

    echo "**Building Docker image on server..."
    ssh -p 2005 -t -o "LogLevel=QUIET" abraham@128.4.27.184 docker build -t nodesim ./nodesim  > /dev/null

    echo "**Starting nodesim image on Server..."
    ssh -p 2005 -t -o "LogLevel=QUIET" abraham@128.4.27.184 docker run -d --rm --name nodesim1 -e LOG_LEVEL=4 nodesim
    rm nodes
fi
if [[ "$args" = *"manager"* ]] || [[ "$args" = *"all"* ]]; then
    echo "**Building manager go executable"
    go build ../../run/manager.go > /dev/null

    echo "**Sending files to server..."
    scp -P 2005 -o "LogLevel=QUIET" manager abraham@128.4.27.184:/home/abraham/  > /dev/null

    echo "**Starting manager on Server"
    ssh -p 2005 -t -o "LogLevel=QUIET" abraham@128.4.27.184 './manager &' > /dev/null &
    rm manager
fi