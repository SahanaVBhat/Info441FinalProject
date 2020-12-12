GOOS=linux go build
docker build -t svbhat/client .
go clean

docker push svbhat/client

ssh ec2-user@info441-deploy.me  < deploy.sh 

