MANAGER_IP=$(/sbin/ip route | awk '/default/ { print $3 }')
/nodes -ip $MANAGER_IP -ll $LOG_LEVEL