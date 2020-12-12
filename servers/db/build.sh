GOOS=linux go build
docker build -t svbhat/sqldatabase .
go clean

docker push svbhat/sqldatabase

ssh ec2-user@api.info441-deploy.me  < deploy.sh 

