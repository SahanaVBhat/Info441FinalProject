docker rm -f database

docker pull svbhat/finalprojsqldatabase

docker run \
    -d \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD="password" \
    -e MYSQL_DATABASE="mysqldatabase" \
    --name database \
	--network info441 \
    svbhat/finalprojsqldatabase