cd iperf3
rice embed-go
cd ..


cd server
rsrc -manifest exe.manifest -ico ..\iperf3\static\main.ico
go build -ldflags="-H windowsgui -w -s" -o ..\iperf3-server-windows.exe
cd ..

cd client
rsrc -manifest exe.manifest -ico ..\iperf3\static\main.ico
go build -ldflags="-H windowsgui -w -s" -o ..\iperf3-client-windows.exe
cd ..