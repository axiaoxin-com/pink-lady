#! /usr/bin/env bash
CYAN="\033[1;36m"
GREEN="\033[0;32m"
WHITE="\033[1;37m"
NOTICE_FLAG="${CYAN}❯"
QUESTION_FLAG="${GREEN}?"

main() {
    gopath=`go env GOPATH`
    echo -e "${NOTICE_FLAG} New project will be create in ${WHITE}${gopath}/src/"
    echo -e "${NOTICE_FLAG} You should enter the project full name such like <github.com/username/projectname>"
    echo -ne "${QUESTION_FLAG} ${CYAN}Enter your new project full name:${CYAN}]: "
    read projname
    projname_dir=`dirname ${projname}`
    mkdir -p ${gopath}/src/${projname_dir}
    echo -ne "${QUESTION_FLAG} ${CYAN}Do you want to the demo code[${WHITE}N/y${CYAN}]: "
    read rmdemo

    # get skeleton
    echo -e "${NOTICE_FLAG} Downloading the skeleton..."
    go get -u github.com/axiaoxin/pink-lady/app
    # replace project name
    echo -e "${NOTICE_FLAG} Generating the project..."
    cp -r ${gopath}/src/github.com/axiaoxin/pink-lady ${gopath}/src/${projname}
    cd ${gopath}/src/${projname} && rm -rf .git
    sed -i "s|github.com/axiaoxin/pink-lady|${projname}|g"  `grep "github.com/axiaoxin/pink-lady" --include *.go -rl .`

    # remove demo
    if [ "${rmdemo}" == "n" ] || [ "${rmdemo}" == "N" ]; then
        rm -rf app/docs
        rm -rf app/apis/demo
        rm -rf app/services/demo
        rm -rf app/models/demo
        sed -i "/demo routes start/,/demo routes end/d" app/apis/routes.go
    fi
    echo -e "${NOTICE_FLAG} Create ${projname} succeed."
}
main
