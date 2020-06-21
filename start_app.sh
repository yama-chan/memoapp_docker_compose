#!/bin/sh
echo "starting...!"
# MySQLサーバーが起動するまでGoアプリのバイナリファイルを実行せずにループで待機する
until mysqladmin ping -h mysql --silent; do
# until mysqladmin ping -h mysql -u root -p admin -P 3306 --silent; do #細かく書きすぎると上手くいかない...
  echo 'waiting for mysqld to be connectable...'
  sleep 2
done

echo "app is starting...!"
# exec go get bitbucket.org/liamstask/goose/cmd/goose
pwd
ls
exec /server
# exec ls
# exec go run main.go