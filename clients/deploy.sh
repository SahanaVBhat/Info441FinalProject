docker rm -f client

docker pull svbhat/finalprojclient

export TLSCERT=/etc/letsencrypt/live/courseeval.info441-deploy.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/courseeval.info441-deploy.me/privkey.pem

docker run \
    -d \
    -e ADDR=:443 \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    -e TLSCERT=$TLSCERT \
    -e TLSKEY=$TLSKEY \
    -p 443:443  \
    -p 80:80 \
    --name client \
    svbhat/finalprojclient
