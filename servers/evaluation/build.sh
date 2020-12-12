GOOS=linux go build -o summary
docker build -t svbhat/summary .
go clean

docker push svbhat/summary

ssh ec2-user@api.info441-deploy.me  < deploy.sh 

