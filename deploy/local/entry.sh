MANAGER_IP=$(/sbin/ip route | awk '/default/ { print $3 }')
cd $GOPATH/src/github.com/abemac/bomb-detection/run
go run nodes.go -ip $MANAGER_IP -ll 4 -nodeConfigFile "$GOPATH/src/github.com/abemac/bomb-detection/run/config-parts/$NODE_CONFIG_FILE"