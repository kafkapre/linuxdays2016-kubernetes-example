
# Build and Run SimpleCrudServer in Local Docker

Run Redis
```
docker run -p 6379:6379 --name redis -d redis
```

Obtain IPAddress value of redis docker container: docker inspect redis
```
export REDIS_IP=`docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis`
```

Build SimpleCrudServer Docker image
```
docker build -t simple-crud-server .
```
   
Run SimpleCrudServer Docker image
```
docker run -p 3000:3000 -e REDIS_IP=$REDIS_IP simple-crud-server
```