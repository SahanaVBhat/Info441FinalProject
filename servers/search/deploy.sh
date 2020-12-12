docker rm -f messaging

docker pull svbhat/messaging

docker run -d \
    --name messaging \
    --network info441 \
    svbhat/messaging 