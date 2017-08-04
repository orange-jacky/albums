pid=`ps -ef | grep albums | grep -v grep | awk '{print $2}'`
kill -9 $pid