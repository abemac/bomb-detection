args="$@"

if [[ "$args" = *"nodes"* ]] || [[ "$args" = *"all"* ]]; then
    echo "**Killing nodesim image on Server..."
    ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 docker kill nodesim1
fi
if [[ "$args" = *"manager"* ]] || [[ "$args" = *"all"* ]]; then
    echo "**Killing manager..."
    ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 'kill $(pidof manager)'
fi