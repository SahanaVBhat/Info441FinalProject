docker rm -f rabbit
docker run -d --name rabbit -p 5672:5672 -p 15672:15672 --network info441 rabbitmq:3-management

docker rm -f customMongoContainer
docker run -d --name customMongoContainer -p 27017:27017 --network info441 mongo

docker rm -f microservices

docker pull svbhat/finalprojmicroservices

docker run -d \
    --name microservices \
    --network info441 \
    svbhat/finalprojmicroservices 