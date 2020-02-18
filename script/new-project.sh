#! /usr/bin/env bash
CYAN="\033[1;36m"
GREEN="\033[0;32m"
WHITE="\033[1;37m"
NOTICE_FLAG="${CYAN}❯"
QUESTION_FLAG="${GREEN}?"

OS=`uname`
# $(replace_in_file pattern file)
function replace_in_file() {
    if [ "$OS" = 'Darwin' ]; then
        # for MacOS
        sed -i '' -e "$1" "$2"
    else
        # for Linux and Windows
        sed -i'' -e "$1" "$2"
    fi
}

main() {
    gopath=`go env GOPATH`
    if [ $? = 127 ]; then
        echo "GOPATH not exists"
        exit -1
    fi
    echo -e "${NOTICE_FLAG} New project will be create in ${WHITE}${gopath}/src/"
    echo -ne "${QUESTION_FLAG} ${CYAN}Enter your new project full name${CYAN}: "
    read projname
    echo -ne "${QUESTION_FLAG} ${CYAN}Do you want to the demo code[${WHITE}Y/n${CYAN}]: "
    read rmdemo

    # get skeleton
    echo -e "${NOTICE_FLAG} Downloading the skeleton..."
    git clone https://github.com/axiaoxin/pink-lady.git ${gopath}/src/${projname}
    # replace project name
    echo -e "${NOTICE_FLAG} Generating the project..."
    cd ${gopath}/src/${projname} && rm -rf .git && cp ${gopath}/src/${projname}/app/config.toml.example ${gopath}/src/${projname}/app/config.toml
    replace_in_file "s|github.com/axiaoxin/pink-lady|${projname}|g"  `grep "github.com/axiaoxin/pink-lady" --include ".travis.yml" --include "*.go" --include "go.*" -rl .`

    # remove demo
    if [ "${rmdemo}" == "n" ] || [ "${rmdemo}" == "N" ]; then
        rm -rf app/apis/demo
        rm -rf app/services/demo
        rm -rf app/models/demo
        replace_in_file "/demo routes start/,/demo routes end/d"  ${gopath}/src/${projname}/app/apis/routes.go
        replace_in_file '/app\/apis\/demo"$/d' ${gopath}/src/${projname}/app/apis/routes.go
    fi
    echo -e "${NOTICE_FLAG} Create project ${projname} in ${gopath}/src succeed."

    # init a git repo
    echo -ne "${QUESTION_FLAG} ${CYAN}Do you want to init a git repo[${WHITE}N/y${CYAN}]: "
    read initgit
    if [ "${initgit}" == "y" ] || [ "${rmdemo}" == "Y" ]; then
        cd ${gopath}/src/${projname} && git init && git add . && git commit -m "init project with pink-lady"
        cp ${gopath}/src/${projname}/script/pre-commit.githook ${gopath}/src/${projname}/.git/hooks/pre-commit
        chmod +x ${gopath}/src/${projname}/.git/hooks/pre-commit
    fi
}
main
