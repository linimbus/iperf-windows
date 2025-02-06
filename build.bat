cd iperf3
go-bindata -pkg iperf3 -o static_files.go main.ico start.ico status.ico stop.ico iperf3.exe cygwin1.dll flow.ico
cd ..

cd server
rsrc -manifest exe.manifest -ico ..\iperf3\main.ico
go build -ldflags="-H windowsgui -w -s" -o ..\iperf3-server-windows.exe
cd ..

cd client
rsrc -manifest exe.manifest -ico ..\iperf3\main.ico
go build -ldflags="-H windowsgui -w -s" -o ..\iperf3-client-windows.exe
cd ..