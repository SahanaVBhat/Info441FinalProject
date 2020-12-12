docker rm -f summary

docker pull svbhat/summary

export TLSCERT=/etc/letsencrypt/live/api.info441-deploy.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.info441-deploy.me/privkey.pem

docker run \
    -d \
    -p 80:80 \
    --name summary \
    --network info441 \
    -v /etc/letsencrypt:/etc/letsencrypt:ro \
    -e TLSCERT=$TLSCERT \
    -e TLSKEY=$TLSKEY \
    svbhat/summary
