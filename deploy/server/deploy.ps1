$Env:GOPATH="C:\Users\abrah\go"
$Env:GOOS="linux"
$Env:GOARCH="amd64"

if ($args.Contains("status")){
    Write-Host
    Write-Host "**Nodes:"
    ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 docker ps 2>$null
    Write-Host
    Write-Host "**Manager:"
    ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 'ps aux | grep manager | grep abraham | grep -v grep' 2>$null
    Write-Host
    exit
}
if ($args.Contains("all") -or $args.Contains("nodes")){
    Write-Host "**Building node simulator go executable..."
    go build ..\..\run\nodes.go > $null

    Write-Host "**Sending files to server..."
    scp -P 2005 Dockerfile entry.sh nodes abraham@128.4.27.184:/home/abraham/nodesim/  2>&1 > $null

    Write-Host "**Building Docker image on server..."
    ssh -p 2005 -t abraham@128.4.27.184 docker build -t nodesim ./nodesim  2>&1 > $null

    Write-Host "**Starting nodesim image on Server..."
    ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 docker run -d --rm --name nodesim1 -e LOG_LEVEL=4 nodesim 2>$null
    Remove-Item nodes
}
if ($args.Contains("all") -or $args.Contains("manager")){
    Write-Host "**Building manager go executable"
    
    go build ../../run/manager.go > $null

    Write-Host "**Sending files to server..."
    scp -P 2005 manager abraham@128.4.27.184:/home/abraham/  2>&1 > $null

    Write-Host "**Building Manager Web App"
    Set-Location ../../webapp/manager
    ng build > $null
    Set-Location ../../deploy/server
    Write-Host "**Sending Web App files to server"
    scp -r -P 2005 ../../webapp/manager/dist abraham@128.4.27.184:/home/abraham/ 2>&1 > $null

    Write-Host "**Starting manager on Server"
    Start-Job -Name manager -ScriptBlock {ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 './manager' 2>$null}
    Remove-Item manager

}