#! /usr/bin/env bash

OS=`uname`
# $(replace_in_file pattern file)
function replace_in_file() {
    if [ "$OS" = 'Darwin' ]; then
        # for MacOS
        sed -i '' -E "$1" "$2"
    else
        # for Linux and Windows
        sed -i'' "$1" "$2"
    fi
}

main() {
    gopath=`go env GOPATH`
    if [ $? = 127 ]; then
        echo "GOPATH not exists"
        exit -1
    fi
    echo -e "New project will be create in ${gopath}/src/"
    echo -ne "Enter your new project full name: "
    read projname
    echo -ne "Do you want to save the demo code[Y/n]: "
    read rmdemo

    # get skeleton
    echo -e "Downloading the skeleton..."
    git clone https://github.com/axiaoxin/pink-lady.git ${gopath}/src/${projname}
    # replace project name
    echo -e "Generating the project..."
    cd ${gopath}/src/${projname} && rm -rf .git && cp ${gopath}/src/${projname}/app/config.toml.example ${gopath}/src/${projname}/app/config.toml
    if [ "$OS" = 'Darwin' ]; then
        sed -i '' -e "s|pink-lady|${projname}|g" `grep "pink-lady" --include "swagger.*" --include ".travis.yml" --include "*.go" --include "go.*" -rl .`
    else
        sed -i "s|pink-lady|${projname}|g" `grep "pink-lady" --include "swagger.*" --include ".travis.yml" --include "*.go" --include "go.*" -rl .`
    fi

    # remove demo
    if [ "${rmdemo}" == "n" ] || [ "${rmdemo}" == "N" ]; then
        rm -rf app/apis/demo*
        rm -rf app/handlers/demohdl
        rm -rf app/models/demomod
        replace_in_file "/demo routes start/,/demo routes end/d"  ${gopath}/src/${projname}/app/apis/routes.go
    fi
    echo -e "Create project ${projname} in ${gopath}/src succeed."

    # init a git repo
    echo -ne "Do you want to init a git repo[N/y]: "
    read initgit
    if [ "${initgit}" == "y" ] || [ "${rmdemo}" == "Y" ]; then
        cd ${gopath}/src/${projname} && git init && git add . && git commit -m "init project with pink-lady"
        cp ${gopath}/src/${projname}/script/pre-commit.githook ${gopath}/src/${projname}/.git/hooks/pre-commit
        chmod +x ${gopath}/src/${projname}/.git/hooks/pre-commit
    fi
}
main
