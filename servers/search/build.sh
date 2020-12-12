GOOS=linux go build
docker build -t svbhat/messaging .
go clean

docker push svbhat/messaging

ssh ec2-user@api.info441-deploy.me  < deploy.sh 

