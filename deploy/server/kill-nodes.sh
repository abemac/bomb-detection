echo "**Killing nodesim image on Server..."
ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 docker kill nodesim1