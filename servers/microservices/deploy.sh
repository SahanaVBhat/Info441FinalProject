docker rm -f messaging

docker pull svbhat/finalprojmicroservices

docker run -d \
    --name messaging \
    --network info441 \
    svbhat/finalprojmicroservices 