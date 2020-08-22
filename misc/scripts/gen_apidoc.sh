#! /usr/bin/env bash
# 生成swag api文档
BUMPVERSION_ERROR=-1
APIDOC_HOST_PATTERN='// @host .+'
APIDOC_VERSION_PATTERN='// @version [0-9]+\.[0-9]+\.[0-9]+'
CONST_VERSION_PATTERN='const VERSION = "[0-9]+\.[0-9]+\.[0-9]+"'


realpath() {
    [[ $1 = /* ]] && echo "$1" || echo "$PWD/${1#./}"
}

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

PROJECT_PATH=$(dirname $(dirname $(realpath $0)))
APP_PATH=${PROJECT_PATH}/app


# Bump version
bumpVersion() {
    lastest_commit=$(git log --pretty=format:"%h %cd %d %s" -1)
    # check version
    current_apidoc_version_line=($(grep -oE "${APIDOC_VERSION_PATTERN}" ${APP_PATH}/apis/apis.go))
    current_apidoc_version=${current_apidoc_version_line[2]}
    current_const_version_line=$(grep -oE "${CONST_VERSION_PATTERN}" ${APP_PATH}/apis/apis.go)
    current_const_version=$(echo "${current_const_version_line}" | grep -o '".*"' | sed 's/"//g')
    if [ "${current_apidoc_version}" != "${current_const_version}" ]; then
        echo -e "apidoc version ${current_apidoc_version} is not match with const version ${current_const_version}"
        exit $BUMPVERSION_ERROR
    fi
    current_version=${current_apidoc_version}
    echo -e "Current version: ${current_version}"
    echo -e "Latest commit: ${lastest_commit}"

    # get new version
    num_list=($(echo ${current_version} | tr '.' ' '))
    major=${num_list[0]}
    minor=${num_list[1]}
    patch=${num_list[2]}
    patch=$((patch + 1))
    suggested_version="$major.$minor.$patch"
    echo -ne "Enter a version number [${suggested_version}]: "
    read new_version
    if [ "${new_version}" = "" ]; then
        new_version=${suggested_version}
    fi

    echo -e "Will set new version to be ${new_version}"
    # update version in apidoc and const
    replace_in_file "s|${APIDOC_VERSION_PATTERN}|// @version ${new_version}|" ${APP_PATH}/apis/apis.go && replace_in_file "s|${CONST_VERSION_PATTERN}|const VERSION = \"${new_version}\"|" ${APP_PATH}/apis/apis.go
}

bumpVersion

APIDOC_HOST=127.0.0.1:4869
echo -ne "Enter apidocs host[${APIDOC_HOST}]: "
read apidoc_host
if [ "${apidoc_host}" != "" ]; then
    APIDOC_HOST=${apidoc_host}
fi

# replace apidoc host
replace_in_file "s|${APIDOC_HOST_PATTERN}|// @host ${APIDOC_HOST}|" ${APP_PATH}/apis/apis.go


# swag init必须在main.go所在的目录下执行，否则必须用--dir参数指定main.go的路径
swag init --dir ${APP_PATH}/ --generalInfo apis/apis.go --propertyStrategy camelcase --output ${APP_PATH}/apis/docs
