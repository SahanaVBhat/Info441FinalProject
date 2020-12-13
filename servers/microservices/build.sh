GOOS=linux go build
docker build -t svbhat/finalprojmicroservices .
go clean

docker push svbhat/finalprojmicroservices

ssh ec2-user@apicourseeval.info441-deploy.me  < deploy.sh 

