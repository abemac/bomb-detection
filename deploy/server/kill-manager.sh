echo "**Killing manager..."
ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 'kill $(pidof manager)'