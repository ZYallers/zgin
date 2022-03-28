#!/bin/bash

jsonFile="$(pwd)/app.json"
appName=""
appAddr=""
logDir=""

echoFun(){
    str=$1
    color=$2
    case ${color} in
        ok)
            echo -e "\033[32m $str \033[0m"
        ;;
        err)
            echo -e "\033[31m $str \033[0m"
        ;;
        tip)
            echo -e "\033[33m $str \033[0m"
        ;;
        title)
            echo -e "\033[42;34m $str \033[0m"
        ;;
        *)
            echo "$str"
        ;;
    esac
}

getJsonValue(){
    echo `cat ${jsonFile} | grep \"$1\"\: | sed -n '1p' | awk -F '": ' '{print $2}' | sed 's/,//g' | sed 's/"//g' | sed 's/ //g'`
}

helpFun(){
    echoFun "Operation:" title
    echoFun "    status                                  View application status" tip
    echoFun "    sync                                    Synchronization application vendor resources" tip
    echoFun "    build                                   Compile and generate application program" tip
    echoFun "    reload                                  Smooth restart application" tip
    echoFun "    quit                                    Stop application" tip
    echoFun "    help                                    View help information for the help command" tip
    echoFun "For more information about an action, use the help command to view it" tip
}

initFun(){
    echoFun "json file:" title
    if [[ ! -f "$jsonFile" ]];then
        echoFun "file [$jsonFile] is not exist" err
        exit 1
    fi

    appName=$(getJsonValue "Name")
    if [[ "$appName" == "" ]];then
        echoFun "appName is empty" err
        exit 1
    fi
    echoFun "appName: $appName" tip

    appAddr=$(getJsonValue "HttpServerAddr")
    if [[ "$appAddr" == "" ]];then
        echoFun "appAddr is empty" err
        exit 1
    fi
    echoFun "appAddr: $appAddr" tip

    logDir=$(getJsonValue "LogDir")
    if [[ "$logDir" == "" ]];then
        echoFun "LogDir is empty" err
        exit 1
    fi
    echoFun "LogDir: $logDir" tip
}

statusFun(){
    echoFun "ps process:" title
    if [[ `pgrep ${appName}|wc -l` -gt 0 ]];then
        ps -p $(pgrep ${appName}|sed ':t;N;s/\n/,/;b t'|sed -n '1h;1!H;${g;s/\n/,/g;p;}') -o user,pid,ppid,%cpu,%mem,vsz,rss,tty,stat,start,time,command
    fi

    echoFun "lsof process:" title
    port=`echoFun ${appAddr}|awk -F ':' '{print $2}'`
    lsof -i:${port}
}

syncFun(){
    echoFun "go mod vendor:" title
    if [[ ! -f "./go.mod" ]];then
        echoFun "go.mod is not exist" err
        exit 1
    fi
    go mod tidy
    go mod vendor
    echoFun "go mod vendor finished" ok
}

buildFun(){
    env=$1
    echoFun "build runner:" title#!/bin/bash

jsonFile="$(pwd)/app.json"
appName=""
appAddr=""
logDir=""

echoFun(){
    str=$1
    color=$2
    case ${color} in
        ok)
            echo -e "\033[32m $str \033[0m"
        ;;
        err)
            echo -e "\033[31m $str \033[0m"
        ;;
        tip)
            echo -e "\033[33m $str \033[0m"
        ;;
        title)
            echo -e "\033[42;34m $str \033[0m"
        ;;
        *)
            echo "$str"
        ;;
    esac
}

getJsonValue(){
    echo `cat ${jsonFile} | grep \"$1\"\: | sed -n '1p' | awk -F '": ' '{print $2}' | sed 's/,//g' | sed 's/"//g' | sed 's/ //g'`
}

helpFun(){
    echoFun "Operation:" title
    echoFun "    status                                  View application status" tip
    echoFun "    sync                                    Synchronization application vendor resources" tip
    echoFun "    build                                   Compile and generate application program" tip
    echoFun "    reload                                  Smooth restart application" tip
    echoFun "    quit                                    Stop application" tip
    echoFun "    help                                    View help information for the help command" tip
    echoFun "For more information about an action, use the help command to view it" tip
}

initFun(){
    echoFun "json file:" title
    if [[ ! -f "$jsonFile" ]];then
        echoFun "file [$jsonFile] is not exist" err
        exit 1
    fi

    appName=$(getJsonValue "Name")
    if [[ "$appName" == "" ]];then
        echoFun "appName is empty" err
        exit 1
    fi
    echoFun "appName: $appName" tip

    appAddr=$(getJsonValue "HttpServerAddr")
    if [[ "$appAddr" == "" ]];then
        echoFun "appAddr is empty" err
        exit 1
    fi
    echoFun "appAddr: $appAddr" tip

    logDir=$(getJsonValue "LogDir")
    if [[ "$logDir" == "" ]];then
        echoFun "LogDir is empty" err
        exit 1
    fi
    echoFun "LogDir: $logDir" tip
}

statusFun(){
    echoFun "ps process:" title
    if [[ `pgrep ${appName}|wc -l` -gt 0 ]];then
        ps -p $(pgrep ${appName}|sed ':t;N;s/\n/,/;b t'|sed -n '1h;1!H;${g;s/\n/,/g;p;}') -o user,pid,ppid,%cpu,%mem,vsz,rss,tty,stat,start,time,command
    fi

    echoFun "lsof process:" title
    port=`echoFun ${appAddr}|awk -F ':' '{print $2}'`
    lsof -i:${port}
}

syncFun(){
    echoFun "go mod vendor:" title
    if [[ ! -f "./go.mod" ]];then
        echoFun "go.mod is not exist" err
        exit 1
    fi
    go mod tidy
    go mod vendor
    echoFun "go mod vendor finished" ok
}

buildFun(){
    env=$1
    echoFun "build runner:" title
    tmpName="${serviceName}_$(date +'%Y%m%d%H%M%S')"
    if [[ "$env" == "debug" ]];then
        echoFun '>>>>>>>>>> build for debug mode <<<<<<<<<<' tip
        # 配合 delve 使用, @see http://wiki.sys.hxsapp.net/pages/viewpage.action?pageId=21349181
        CGO_ENABLED=0 go build -v -installsuffix cgo -gcflags 'all=-N -l' -i -o ./bin/${tmpName} -tags=jsoniter ./main.go
    elif [[ "$env" == "dev" ]];then
        echoFun '>>>>>>>>>> build for development mode <<<<<<<<<<' tip
        CGO_ENABLED=0 go build -v -installsuffix cgo -ldflags '-w' -i -o ./bin/${tmpName} -tags=jsoniter ./main.go
    else
        echoFun '>>>>>>>>>> build for production mode <<<<<<<<<<' tip
        # Build compilation parameter reference:
        # Dependency free compilation：https://blog.csdn.net/weixin_42506905/article/details/93135684
        # Detailed explanation of build parameters：https://blog.csdn.net/zl1zl2zl3/article/details/83374131
        # Ldflags parameter：https://blog.csdn.net/javaxflinux/article/details/89177863
        CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w' -i -o ./bin/${tmpName} -tags=jsoniter ./main.go
    fi

    if [[ ! -f "./bin/${tmpName}" ]];then
        echoFun "build tmp runner ($(pwd)/bin/${tmpName}) failed" err
        exit 1
    fi

    mv -f ./bin/${tmpName} ./bin/${appName}
    if [[ ! -f "./bin/${appName}" ]];then
        echoFun "mv tmp runner failed" err
        exit 1
    fi
    echoFun "build runner ($(pwd)/bin/${appName}) finished" ok
}

sendMsg(){
    app="App: $appName"
    listen="Listen: $appAddr"
    hostName="HostName: $(hostname)"
    time="Time: $(date "+%Y/%m/%d %H:%M:%S")"
    sip="SystemIP: $(ifconfig -a |grep inet |grep -v 127.0.0.1 |grep -v inet6|awk '{print $2}' |tr -d "addr:")"

    token=$(getJsonValue "GracefulRobotToken")
    url="https://oapi.dingtalk.com/robot/send?access_token=$token"
    content="$1\n---------------------------\n$app\n$listen\n$hostName\n$time\n$sip"
    cnt=$(echo ${content//\"/\\\"})
    header="Content-Type: application/json"
    curl -o /dev/null -m 3 -s "$url" -H "$header" -d "{\"msgtype\":\"text\",\"text\":{\"content\":\"$cnt\"}}"
}

reloadFun(){
    echoFun "runner reloading:" title

    if [[ ! -f "./bin/$appName" ]];then
        echoFun "runner [`pwd`/bin/$appName] is not exist" err
        exit 1
    fi

    if [[ ! -d "$logDir" ]];then
        mkdir -p ${logDir}
    fi
    if [[ ! -d "$logDir" ]];then
        echoFun "logDir [$logDir] is not exist" err
        exit 1
    fi

    logfile=${logDir}/${appName}.log
    if [[ ! -f "$logfile" ]];then
        touch ${logfile}
    fi
    echoFun "logfile: $logfile" tip

    if [[ ! -x "./bin/$appName" ]];then
        chmod u+x ./bin/${appName}
    fi

    quitFun

    # Prevent Jenkins from killing all derived processes after the end of build by default
    export BUILD_ID=dontKillMe

    nohup ./bin/${appName} >> ${logfile} 2>&1 &
    echoFun "app $appName($appAddr) is reloaded, pid: `echo $!`" ok

    # Check whether the health interface is accessed normally
    sleep 3s
    resp=`curl -m 3 -s "http://$appAddr/health" | xargs echo -n`
    echoFun "curl \"http://$appAddr\" health: $resp" tip
    sendMsg "curl \"http://$appAddr\" health: $resp"
}

quitFun(){
    port=`echo ${appAddr}|awk -F ':' '{print $2}'`
    counter=0
    while true;
    do
        pid=`lsof -i tcp:${port}|grep LISTEN|awk '{print $2}'`
        if [[ ${pid} -gt 0 ]];then
            if [[ ${counter} -ge 30 ]];then
                kill -9 ${pid}
                echoFun "app($appName) has been killed for 30s and is ready to be forcibly killed" tip
                sendMsg "app($appName) has been killed for 30s and is ready to be forcibly killed"
                break
            else
                kill ${pid}
                counter=$(($counter+1))
                echoFun "killing app $appName($port), pid($pid), $counter tried" tip
                sleep 1s
            fi
        else
            echoFun "app $appName($port) is stopped" ok
            break
        fi
    done
}


cmd=$1
arg1=$2
arg2=$3

initFun
case $1 in
    status)
        statusFun
    ;;
    sync)
        syncFun
    ;;
    build)
        buildFun $2
    ;;
    quit)
        quitFun
    ;;
    reload)
        reloadFun
    ;;
    *)
        helpFun
    ;;
esac

    tmpName="${serviceName}_$(date +'%Y%m%d%H%M%S')"
    if [[ "$env" == "debug" ]];then
        echoFun '>>>>>>>>>> build for debug mode <<<<<<<<<<' tip
        # 配合 delve 使用, @see http://wiki.sys.hxsapp.net/pages/viewpage.action?pageId=21349181
        CGO_ENABLED=0 go build -v -installsuffix cgo -gcflags 'all=-N -l' -i -o ./bin/${tmpName} -tags=jsoniter ./main.go
    elif [[ "$env" == "dev" ]];then
        echoFun '>>>>>>>>>> build for development mode <<<<<<<<<<' tip
        CGO_ENABLED=0 go build -v -installsuffix cgo -ldflags '-w' -i -o ./bin/${tmpName} -tags=jsoniter ./main.go
    else
        echoFun '>>>>>>>>>> build for production mode <<<<<<<<<<' tip
        # Build compilation parameter reference:
        # Dependency free compilation：https://blog.csdn.net/weixin_42506905/article/details/93135684
        # Detailed explanation of build parameters：https://blog.csdn.net/zl1zl2zl3/article/details/83374131
        # Ldflags parameter：https://blog.csdn.net/javaxflinux/article/details/89177863
        CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w' -i -o ./bin/${tmpName} -tags=jsoniter ./main.go
    fi

    if [[ ! -f "./bin/${tmpName}" ]];then
        echoFun "build tmp runner ($(pwd)/bin/${tmpName}) failed" err
        exit 1
    fi

    mv -f ./bin/${tmpName} ./bin/${appName}
    if [[ ! -f "./bin/${appName}" ]];then
        echoFun "mv tmp runner failed" err
        exit 1
    fi
    echoFun "build runner ($(pwd)/bin/${appName}) finished" ok
}

sendMsg(){
    app="App: $appName"
    listen="Listen: $appAddr"
    hostName="HostName: $(hostname)"
    time="Time: $(date "+%Y/%m/%d %H:%M:%S")"
    sip="SystemIP: $(ifconfig -a |grep inet |grep -v 127.0.0.1 |grep -v inet6|awk '{print $2}' |tr -d "addr:")"

    token=$(getJsonValue "GracefulRobotToken")
    url="https://oapi.dingtalk.com/robot/send?access_token=$token"
    content="$1\n---------------------------\n$app\n$listen\n$hostName\n$time\n$sip"
    cnt=$(echo ${content//\"/\\\"})
    header="Content-Type: application/json"
    curl -o /dev/null -m 3 -s "$url" -H "$header" -d "{\"msgtype\":\"text\",\"text\":{\"content\":\"$cnt\"}}"
}

reloadFun(){
    echoFun "runner reloading:" title

    if [[ ! -f "./bin/$appName" ]];then
        echoFun "runner [`pwd`/bin/$appName] is not exist" err
        exit 1
    fi

    if [[ ! -d "$logDir" ]];then
        mkdir -p ${logDir}
    fi
    if [[ ! -d "$logDir" ]];then
        echoFun "logDir [$logDir] is not exist" err
        exit 1
    fi

    logfile=${logDir}/${appName}.log
    if [[ ! -f "$logfile" ]];then
        touch ${logfile}
    fi
    echoFun "logfile: $logfile" tip

    if [[ ! -x "./bin/$appName" ]];then
        chmod u+x ./bin/${appName}
    fi

    quitFun

    # Prevent Jenkins from killing all derived processes after the end of build by default
    export BUILD_ID=dontKillMe

    nohup ./bin/${appName} >> ${logfile} 2>&1 &
    echoFun "app $appName($appAddr) is reloaded, pid: `echo $!`" ok

    # Check whether the health interface is accessed normally
    sleep 3s
    resp=`curl -m 3 -s "http://$appAddr/health" | xargs echo -n`
    echoFun "curl \"http://$appAddr\" health: $resp" tip
    sendMsg "curl \"http://$appAddr\" health: $resp"
}

quitFun(){
    port=`echo ${appAddr}|awk -F ':' '{print $2}'`
    counter=0
    while true;
    do
        pid=`lsof -i tcp:${port}|grep LISTEN|awk '{print $2}'`
        if [[ ${pid} -gt 0 ]];then
            if [[ ${counter} -ge 30 ]];then
                kill -9 ${pid}
                echoFun "app($appName) has been killed for 30s and is ready to be forcibly killed" tip
                sendMsg "app($appName) has been killed for 30s and is ready to be forcibly killed"
                break
            else
                kill ${pid}
                counter=$(($counter+1))
                echoFun "killing app $appName($port), pid($pid), $counter tried" tip
                sleep 1s
            fi
        else
            echoFun "app $appName($port) is stopped" ok
            break
        fi
    done
}


cmd=$1
arg1=$2
arg2=$3

initFun
case $1 in
    status)
        statusFun
    ;;
    sync)
        syncFun
    ;;
    build)
        buildFun $2
    ;;
    quit)
        quitFun
    ;;
    reload)
        reloadFun
    ;;
    *)
        helpFun
    ;;
esac
