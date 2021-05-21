#!/bin/sh

HOST='sqre.ru'
USER='sqre'
DOC_ROOT="/home/sqre"


go build 

if [ $? -ne 0 ]
then
  exit 1;
fi;




if [ $USER ]
then
  SSH_HOST="$USER@$HOST"
else
  SSH_HOST=$HOST
fi

echo '* Создаем архив...'
tar -czf yourproject.tar.gz sqre
if [ $? -ne 0 ]
then
  exit 1;
fi;



echo '* Копируем архив на сервер...'

scp ./yourproject.tar.gz  $SSH_HOST:$DOC_ROOT
if [ $? -ne 0 ]
then
  exit 1;
fi;

echo '* Распаковываем архив на серверe...'
ssh $SSH_HOST "cd $DOC_ROOT; tar -xzf yourproject.tar.gz 2> /dev/null && rm -rf $DOC_ROOT/goapp/sqre && mv sqre $DOC_ROOT/goapp && chmod -R a+w+x $DOC_ROOT/goapp/sqre"
if [ $? -ne 0 ]
then
  exit 1;
fi;


echo '* Удаляем архив на сервере ...'
ssh $SSH_HOST "cd $DOC_ROOT; rm -rf yourproject.tar.gz"
if [ $? -ne 0 ]
then
  exit 1;
fi;

echo '* Удаляем архив локально ...'

rm -rf yourproject.tar.gz