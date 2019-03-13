#! /usr/bin/env bash
# COLORS
RED="\033[1;31m"
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
CYAN="\033[1;36m"
WHITE="\033[1;37m"
# BLUE="\033[1;34m"
# PURPLE="\033[1;35m"
# RESET="\033[0m"

QUESTION_FLAG="${GREEN}?"
WARNING_FLAG="${YELLOW}!"
ERROR_FLAG="${RED}!!"
NOTICE_FLAG="${CYAN}❯"

# PATHS
PROJECT_PATH=$(dirname $(dirname $(readlink -f $0)))
APP_PATH=${PROJECT_PATH}/app
BUILD_PATH=${PROJECT_PATH}/build
NOW=$(date "+%Y%m%d-%H%M%S")

# ERROR CODE
BUMPVERSION_ERROR=-1
TESTING_FAILED=-2
BUILDING_FAILED=-3

APIDOC_VERSION_PATTERN='// @version [0-9]+\.[0-9]+\.[0-9]+'
CONST_VERSION_PATTERN='const VERSION = "[0-9]+\.[0-9]+\.[0-9]+"'

# Bump version
bumpVersion() {
    lastest_commit=$(git log --pretty=format:"%h %cd %d %s" -1)
    # check version
    current_apidoc_version_line=($(grep -oE "${APIDOC_VERSION_PATTERN}" ${APP_PATH}/apis/init.go))
    current_apidoc_version=${current_apidoc_version_line[2]}
    current_const_version_line=$(grep -oE "${CONST_VERSION_PATTERN}" ${APP_PATH}/apis/init.go)
    current_const_version=$(echo "${current_const_version_line}" | grep -o '".*"' | sed 's/"//g')
    if [ "${current_apidoc_version}" != "${current_const_version}" ]; then
        echo -e "${ERROR_FLAG} ${RED}apidoc version ${current_apidoc_version} is not match with const version ${current_const_version}"
        exit $BUMPVERSION_ERROR
    fi
    current_version=${current_apidoc_version}
    echo -e "${NOTICE_FLAG} Current version: ${WHITE}${current_version}"
    echo -e "${NOTICE_FLAG} Latest commit: ${WHITE}${lastest_commit}"

    # get new version
    num_list=($(echo ${current_version} | tr '.' ' '))
    major=${num_list[0]}
    minor=${num_list[1]}
    patch=${num_list[2]}
    patch=$((patch + 1))
    suggested_version="$major.$minor.$patch"
    echo -ne "${QUESTION_FLAG} ${CYAN}Enter a version number [${WHITE}${suggested_version}${CYAN}]: "
    read new_version
    if [ "${new_version}" = "" ]; then
        new_version=${suggested_version}
    fi

    echo -e "${NOTICE_FLAG} Will set new version to be ${WHITE}${new_version}"
    # update version in apidoc and const
    sed -i -r "s|${APIDOC_VERSION_PATTERN}|// @version ${new_version}|" ${APP_PATH}/apis/init.go && sed -i -r "s|${CONST_VERSION_PATTERN}|const VERSION = \"${new_version}\"|" ${APP_PATH}/apis/init.go
}

tests() {
    # Running tests
    echo -e "${NOTICE_FLAG} Running tests"
    cd ${APP_PATH}
    if !(go test ./...)
    then
        echo -e "${ERROR_FLAG} ${RED}Tests failed."
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
    echo -ne "${QUESTION_FLAG} ${CYAN}Enter release binary name [${WHITE}${BINARY_NAME}${CYAN}]: "
    read binary_name
    if [ "${binary_name}" != "" ]; then
        BINARY_NAME=${binary_name}
    fi

    echo -e "${NOTICE_FLAG} Will build binary name to be ${WHITE}${BINARY_NAME}"
    # Update docs
    echo "Updating swag docs"
    # check swag
    if !(swag > /dev/null 2>&1); then
        echo -e "${WARNING_FLAG} ${CYAN}No swag for generate API docs. You need to install it for auto generate the docs"
    else
        echo -e "${NOTICE_FLAG} Generating API docs..."
        cd ${APP_PATH}
        swag init -g apis/init.go
        cd -
    fi

    # Building
    echo -e "${NOTICE_FLAG} Building..."
    go build -o ${BUILD_PATH}/${BINARY_NAME} -tags=jsoniter -v ${APP_PATH}
    if [ $? -ne 0 ]
    then
        echo -e "${ERROR_FLAG} ${RED}Build failed."
        exit ${BUILDING_FAILED}
    fi
}

tarball() {
    tarname=${BINARY_NAME}-${new_version}-${NOW}
    tardir=${BUILD_PATH}/${tarname}
    mkdir ${tardir}
    cp ${BUILD_PATH}/${BINARY_NAME} ${tardir}
    if [ -e ${APP_PATH}/${BINARY_NAME}.yaml ]; then
        cp ${APP_PATH}/${BINARY_NAME}.yaml ${tardir}
    fi
    if [ -e ${APP_PATH}/config.yaml ]; then
        cp ${APP_PATH}/config.yaml ${tardir}
    fi
    tar czvf ${tardir}.tar.gz -C ${BUILD_PATH} ${tarname} && rm -rf ${tardir}
}

commit() {
    echo -ne "${QUESTION_FLAG} ${CYAN}Do you want to commit this version bump to git[${WHITE}Y/n${CYAN}]: "
    read do_commit
    if [ "${do_commit}" == "" ] || [ "${do_commit}" == "Y" ]; then
        git add ${APP_PATH}/docs ${APP_PATH}/apis/init.go
        git commit -m "Bump version ${current_version} -> ${new_version}"
        echo -ne "${QUESTION_FLAG} ${CYAN}Do you want to tag this version bump to git[${WHITE}Y/n${CYAN}]: "
        read do_tag
        if [ "${do_tag}" == "" ] || [ "${do_tag}" == "Y" ]; then
            git tag ${new_version}
            if [ $? -ne 0 ]; then
                echo -e "${WARNING_FLAG} ${CYAN}git tag failed"
            fi
        fi
    fi
}

main() {
    echo -e "${NOTICE_FLAG} This tool will help you to release your binary app.\n It will bump the version in your code, run tests then update apidocs and build the binary app and tar it as tar.gz file.\n Last do an optional commit the changed code and tag it with then version name"
    bumpVersion
    tests
    build
    tarball
    commit
}

main
