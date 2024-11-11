rsrc -manifest exe.manifest -ico main.ico

go build -ldflags="-H windowsgui -w -s" -o iperf-windows.exe