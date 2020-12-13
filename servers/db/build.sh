GOOS=linux go build
docker build -t svbhat/finalprojsqldatabase .
go clean

docker push svbhat/finalprojsqldatabase

ssh ec2-user@apicourseeval.info441-deploy.me  < deploy.sh 

