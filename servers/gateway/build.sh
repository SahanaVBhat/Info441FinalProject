GOOS=linux go build -o gateway
docker build -t svbhat/gateway .
go clean

docker push svbhat/gateway

ssh ec2-user@api.info441-deploy.me  < deploy.sh 

