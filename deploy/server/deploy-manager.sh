echo "**Building manager go executable"
GOPATH=/mnt/c/Users/abrah/go
go build ../../run/manager.go > /dev/null

echo "**Sending files to server..."
scp -P 2005 manager abraham@128.4.27.184:/home/abraham/  > /dev/null

echo "**Starting manager on Server"
ssh -p 2005 -t abraham@128.4.27.184 ./manager &