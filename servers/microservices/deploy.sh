docker rm -f microservices

docker pull svbhat/finalprojmicroservices

docker run -d \
    --name microservices \
    --network info441 \
    svbhat/finalprojmicroservices 