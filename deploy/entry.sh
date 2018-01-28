MANAGER_IP=$(/sbin/ip route | awk '/default/ { print $3 }')
cd $GOPATH/src/github.com/abemac/bomb-detection/run
go run nodes.go -ip $MANAGER_IP -ll $LOG_LEVEL