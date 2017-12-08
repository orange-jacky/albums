pid=`ps -ef | grep albums_server | grep -v grep | awk '{print $2}'`
kill -9 $pid