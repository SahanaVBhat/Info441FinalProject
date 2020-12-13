GOOS=linux go build
docker build -t svbhat/finalprojclient .
go clean

docker push svbhat/finalprojclient

ssh ec2-user@courseeval.info441-deploy.me  < deploy.sh 

