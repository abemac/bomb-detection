if ($args.Contains("all") -or $args.Contains("nodes")){
    Write-Host "**Killing nodesim image on Server..."
    ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 docker kill nodesim1 2>$null
}
if ($args.Contains("all") -or $args.Contains("manager")){
    Write-Host "**Killing manager..."
    ssh -p 2005 -t -o LogLevel=QUIET abraham@128.4.27.184 'kill $(pidof manager)' 2>$null
}