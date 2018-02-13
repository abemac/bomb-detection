echo
echo "**Nodes:"
ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 docker ps
echo
echo "**Manager:"
ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 'ps aux | grep manager | grep abraham | grep -v grep'
echo