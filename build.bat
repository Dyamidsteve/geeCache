go build -o server.exe &
./server 
@REM ./server -port=8002 
@REM ./server -port=8003 -api=1 

echo "start test"

curl "http://localhost:9999/api?key=Tom" 
@REM curl "http://localhost:9999/api?key=Tom" 
@REM curl "http://localhost:9999/api?key=Tom" 

wait