#! /usr/bin/env bash

realpath() {
    [[ $1 = /* ]] && echo "$1" || echo "$PWD/${1#./}"
}

# PATHS
PROJECT_PATH=$(dirname $(dirname $(realpath $0)))
APP_PATH=${PROJECT_PATH}/app
BUILD_PATH=${PROJECT_PATH}/build
NOW=$(date "+%Y%m%d-%H%M%S")

# ERROR CODE
TESTING_FAILED=-2
BUILDING_FAILED=-3



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


tests() {
    # Running tests
    echo -e "Running tests"
    cd ${APP_PATH}
    if !(go test ./...)
    then
        echo -e "Tests failed."
        cd -
        exit ${TESTING_FAILED}
    fi
    cd -
}

build() {
    if [ ! -d ${BUILD_PATH} ]; then
        mkdir ${BUILD_PATH}
    fi
    BINARY_NAME=app
    echo -ne "Enter release binary name [${BINARY_NAME}]: "
    read binary_name
    if [ "${binary_name}" != "" ]; then
        BINARY_NAME=${binary_name}
    fi

    echo -e "Will build binary name to be ${BINARY_NAME}"

    # Update docs
    echo "Updating swag docs"
    # check swag
    if !(swag > /dev/null 2>&1); then
        echo -e "No swag for generate API docs. You need to install it for auto generate the docs"
    else
        echo -e "Generating API docs..."
        bash ${PROJECT_PATH}/script/gen_apidoc.sh
    fi

    # Building
    echo -e "Building..."
    go build -o ${BUILD_PATH}/${BINARY_NAME} -tags=jsoniter -v ${APP_PATH}
    if [ $? -ne 0 ]
    then
        echo -e "Build failed."
        exit ${BUILDING_FAILED}
    fi
}

tarball() {
    configfile=config.toml
    echo -ne "Enter your configfile[${configfile}]: "
    read cf
    if [ "${cf}" != "" ]; then
        configfile=${cf}
    fi

    tarname=${BINARY_NAME}-${new_version}-${NOW}
    tardir=${BUILD_PATH}/${tarname}
    mkdir ${tardir}
    cp ${BUILD_PATH}/${BINARY_NAME} ${tardir}
    cp ${APP_PATH}/${configfile} ${tardir}
    tar czvf ${tardir}.tar.gz -C ${BUILD_PATH} ${tarname} && rm -rf ${tardir}
}

commit() {
    echo -ne "Do you want to commit this version bump to git[Y/n]: "
    read do_commit
    if [ "${do_commit}" == "" ] || [ "${do_commit}" == "Y" ]; then
        git add ${APP_PATH}/docs ${APP_PATH}/apis/apis.go
        git commit -m "Bump version ${current_version} -> ${new_version}"
        echo -ne "Do you want to tag this version bump to git[Y/n]: "
        read do_tag
        if [ "${do_tag}" == "" ] || [ "${do_tag}" == "Y" ]; then
            git tag ${new_version}
            if [ $? -ne 0 ]; then
                echo -e "git tag failed"
            fi
        fi
    fi
}

main() {
    echo -e "This tool will help you to release your binary app.\n It will bump the version in your code, run tests then update apidocs and build the binary app and tar it as tar.gz file.\n Last do an optional commit the changed code and tag it with then version name"
    tests
    build
    tarball
    commit
}

main
