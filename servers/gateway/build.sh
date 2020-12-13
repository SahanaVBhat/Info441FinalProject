GOOS=linux go build
docker build -t svbhat/finalprojgateway .
go clean

docker push svbhat/finalprojgateway

ssh ec2-user@apicourseeval.info441-deploy.me  < deploy.sh 

