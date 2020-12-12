docker rm -f apiserver

docker pull svbhat/gateway

export SESSIONKEY="key"
export MYSQL_ROOT_PASSWORD="password"
export MYSQL_DATABASE="mysqldatabase"
export DSN="root:password@tcp(database:3306)/mysqldatabase"
export REDISADDR=redis:6379
export MESSAGESADDR="http://messaging:80"
export SUMMARYADDR="http://summary:80"

docker run -d \
    -p 443:443  \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    -e TLSCERT=/etc/letsencrypt/live/api.info441-deploy.me/fullchain.pem \
    -e TLSKEY=/etc/letsencrypt/live/api.info441-deploy.me/privkey.pem \
    -e ADDR=:443 \
    -e SESSIONKEY=$SESSIONKEY -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD -e MYSQL_DATABASE=$MYSQL_DATABASE -e DSN=$DSN -e REDISADDR=$REDISADDR -e MESSAGESADDR=$MESSAGESADDR -e SUMMARYADDR=$SUMMARYADDR \
    --name apiserver \
    --network info441 \
    svbhat/gateway

