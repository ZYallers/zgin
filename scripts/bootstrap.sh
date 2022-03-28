#!/bin/bash
# Initialize the script, pull the data, and generate the compilation and deployment script

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

echoFun "current path: $(pwd)" title
echoFun "download produce.sh" title
src='http://github.com/ZYallers/zgin/raw/master/scripts/produce.sh'
des=`dirname $0`"/produce.sh"
curl -o ${des} ${src}
if [[ ! -f "$des" ]];then
    echoFun "download produce.sh($(pwd)/$des) failed" err
    exit 1
fi
chmod u+x ${des}
echoFun "download produce.sh($(pwd)/$des) finished" ok
$des help

