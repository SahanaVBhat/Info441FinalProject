GOOS=linux go build -o gateway
docker build -t svbhat/finalprojgateway .
go clean

docker push svbhat/finalprojgateway

ssh ec2-user@apicourseeval.info441-deploy.me  < deploy.sh 

