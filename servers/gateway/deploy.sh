docker rm -f redis
docker run -d --name redis --network info441 redis

docker rm -f apiserver
docker pull svbhat/finalprojgateway

export SESSIONKEY="key"
export MYSQL_ROOT_PASSWORD="password"
export MYSQL_DATABASE="mysqldatabase"
export DSN="root:password@tcp(database:3306)/mysqldatabase"
export REDISADDR=redis:6379
export MICROSERVICEADDR="http://microservices:80"

docker run -d \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    -p 443:443  \
    -e TLSCERT=/etc/letsencrypt/live/apicourseeval.info441-deploy.me/fullchain.pem \
    -e TLSKEY=/etc/letsencrypt/live/apicourseeval.info441-deploy.me/privkey.pem \
    -e ADDR=:443 \
    -e SESSIONKEY=$SESSIONKEY -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD -e MYSQL_DATABASE=$MYSQL_DATABASE \
    -e DSN=$DSN -e REDISADDR=$REDISADDR -e MICROSERVICEADDR=$MICROSERVICEADDR \
    --name apiserver \
    --network info441 \
    svbhat/finalprojgateway

