#!/bin/bash

appConfigFile="$(pwd)/src/config/app/base.go"
name=""
httpServerAddr=""
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

helpFun(){
    echoFun "操作:" title
    echoFun "    status [httpAddr(监听IP:Port)]          查看服务状态" tip
    echoFun "    sync                                    同步服务vendor资源" tip
    echoFun "    build [branch(分支)]                    编译生成服务程序" tip
    echoFun "    reload [httpAddr(监听IP:Port)]          平滑重启服务" tip
    echoFun "    quit [httpAddr(监听IP:Port)]            停止服务" tip
    echoFun "    help                                    查看命令的帮助信息" tip
    echoFun "有关某个操作的详细信息，请使用 help 命令查看" tip
    exit 0
}

initFun(){
    if [[ ! -f "$appConfigFile" ]];then
        echoFun "file [$appConfigFile] is not exist" err
        exit 1
    fi

    name=`cat ${appConfigFile}|grep "Name"|awk -F '"' '{print $2}'`
    if [[ "$name" == "" ]];then
        echoFun "name is null" err
        exit 1
    fi
    echoFun "name: $name" tip

    if [[ "$1" == "" ]];then
        httpServerAddr=`cat ${appConfigFile}|grep "HttpServerDefaultAddr"|awk -F '"' '{print $2}'`
        echoFun "httpServerAddr[httpServerDefaultAddr]: $httpServerAddr" tip
    else
        httpServerAddr=$1
        echoFun "httpServerAddr[shellArgs]: $httpServerAddr" tip
    fi

    if [[ "$httpServerAddr" == "" ]];then
        echoFun "httpServerAddr is empty" err
        exit 1
    fi

    logDir=`cat ${appConfigFile}|grep "LogDir"|awk -F '"' '{print $2}'`
    if [[ "$logDir" == "" ]];then
        echoFun "logDir is null" err
        exit 1
    fi

    echoFun "logDir: $logDir" tip
}

statusFun(){
    initFun $1

    echoFun "ps process:" title
    if [[ `pgrep ${name}|wc -l` -gt 0 ]];then
        ps -p $(pgrep ${name}|sed ':t;N;s/\n/,/;b t') -o user,pid,ppid,%cpu,%mem,vsz,rss,tty,stat,start,time,command
    fi

    echoFun "lsof process:" title
    port=`echo ${httpServerAddr}|awk -F ':' '{print $2}'`
    lsof -i:${port}
}

syncFun(){
    initFun
    cd ./src

    echoFun "go mod vendor:" title
    if [[ ! -f "./go.mod" ]];then
        go mod init src
    fi
    go mod tidy
    rm -rf ./vendor
    go mod vendor
    echoFun "go mod vendor finished" ok
}

buildFun(){
    echoFun "git pull:" title
    branch=$1
    env=$2
    if [[ "$branch" == "" ]];then
        echoFun "branch of the build is empty" err
        exit 1
    fi
    if [[ "$branch" == "local" ]];then
        echoFun "ignore git pull, direct build by local" tip
    else
        git remote update origin --prune # 更新远程分支列表
        git checkout ${branch} # 切换分支
        git pull # 拉取最新版本
        echoFun "git pull [$branch] finish" ok
    fi

    initFun

    echoFun "build runner:" title
    cd ./src
    tmpName="${name}_tmp_$(date +'%Y-%m-%d-%H-%M-%S')"
    if [[ "$dlv" ]];then
        echoFun "build with -gcflags 'all=-N -l'" tip
        if [[ "$env" == "dev" ]];then
            echoFun 'build in develop environment' tip
            CGO_ENABLED=0 go build -v -installsuffix cgo -gcflags 'all=-N -l' -i -o ../bin/${tmpName} -tags=jsoniter ./main.go
        else
            CGO_ENABLED=0 go build -a -installsuffix cgo -gcflags 'all=-N -l' -i -o ../bin/${tmpName} -tags=jsoniter ./main.go
        fi
    else
        echoFun "no extra build options" tip
        if [[ "$env" == "dev" ]];then
            echoFun 'build in develop environment' tip
            CGO_ENABLED=0 go build -v -installsuffix cgo -ldflags '-w' -i -o ../bin/${tmpName} -tags=jsoniter ./main.go
        else
            ### build编译参数参考资料：
            # 无依赖编译：https://blog.csdn.net/weixin_42506905/article/details/93135684
            # build参数详解：https://blog.csdn.net/zl1zl2zl3/article/details/83374131
            # ldflags参数：https://blog.csdn.net/javaxflinux/article/details/89177863
            CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w' -i -o ../bin/${tmpName} -tags=jsoniter ./main.go
        fi
    fi

    cd ../
    if [[ ! -f "./bin/${tmpName}" ]];then
        echoFun "build tmp runner ($(pwd)/bin/${tmpName}) failed" err
        exit 1
    fi

    mv -f ./bin/${tmpName} ./bin/${name}
    if [[ ! -f "./bin/${name}" ]];then
        echoFun "mv tmp runner failed" err
        exit 1
    fi

    echoFun "build runner ($(pwd)/bin/${name}) finished" ok
}

reloadFun(){
    initFun $1

    echoFun "runner reloading:" title

    if [[ ! -f "./bin/$name" ]];then
        echoFun "runner [`pwd`/bin/$name] is not exist" err
        exit 1
    fi

    # 日志目录
    if [[ ! -d "$logDir" ]];then
        mkdir -p ${logDir}
    fi
    if [[ ! -d "$logDir" ]];then
        echoFun "logDir [$logDir] is not exist" err
        exit 1
    fi

    # 日志文件
    logfile=${logDir}/${name}.log
    if [[ ! -f "$logfile" ]];then
        touch ${logfile}
    fi
    echoFun "logfile: $logfile" tip

    if [[ ! -x "./bin/$name" ]];then
        chmod u+x ./bin/${name}
    fi

    quitFun ${httpServerAddr}

    # 防止Jenkins默认会在Build结束后Kill掉所有的衍生进程
    export BUILD_ID=dontKillMe

    nohup ./bin/${name} -http.addr=${httpServerAddr} >> ${logfile} 2>&1 &
    echoFun "$name($httpServerAddr) is reloaded, pid: `echo $!`" ok

    # 检查对应服务端口是否已经正常启动
    sleep 3s
    port=`echo ${httpServerAddr}|awk -F ':' '{print $2}'`
    if [[ `lsof -i tcp:${port}|grep LISTEN|wc -l` -le 0 ]];then
        echoFun "$name($httpServerAddr) service is not running" err
        exit 1
    else
        echoFun "$name($httpServerAddr) service is running..." ok
        # 检查健康接口是否访问正常
        respHttpCode=`curl -m 3 -s -w "%{http_code}" "http://{$httpServerAddr}/health"`
        if [[ "${respHttpCode}" != 'ok200' ]];then
            echoFun "curl 'http://$httpServerAddr/health' error: ($respHttpCode)" err
            exit 1
        else
            echoFun "curl 'http://$httpServerAddr/health' succeed" ok
        fi
    fi
}

quitFun(){
    addr=$1
    port=`echo ${addr}|awk -F ':' '{print $2}'`

    counter=1
    while true;
    do
        pid=`lsof -i tcp:${port}|grep LISTEN|awk '{print $2}'`
        if [[ ${pid} -gt 0 ]];then
            echoFun "killing service $name($port), pid($pid), $counter tried" tip
            kill ${pid}
            counter=$(($counter+1))
            sleep 1s
        else
            echoFun "$name($port) service is stopped" ok
            break
        fi
    done
}

while getopts ':d' OPT; do
    case ${OPT} in
        d)
            # 配合 delve 使用, http://wiki.sys.hxsapp.net/pages/viewpage.action?pageId=21349181
            dlv=d
            shift 1;
        ;;
        ?)  #当有不认识的选项的时候arg为?
            echo "unkonw argument"
            exit 1
        ;;
    esac
done

cmd=$1
arg1=$2
arg2=$3

case ${cmd} in
        status)
            statusFun ${arg1}
        ;;
        sync)
            syncFun
        ;;
        build)
            buildFun ${arg1} ${arg2}
        ;;
        quit)
            initFun ${arg1}
            quitFun ${httpServerAddr}
        ;;
        reload)
            reloadFun ${arg1}
        ;;
        *)
            helpFun
        ;;
esac
