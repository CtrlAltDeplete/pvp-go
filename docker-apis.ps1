$env:GOOS = "linux"
go build PvP-Go
$env:GOOS = "windows"
docker build .
docker run -d -p 8080:8080 --name pvpgo-apis