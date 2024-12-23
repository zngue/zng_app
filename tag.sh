#!/bin/bash


git add .

git commit -m "update"
git push
sleep 3
echo "推送完成 开始打标签"
version=${1:-"v0.0.25"}
git tag -d "${version}"
git push origin :refs/tags/"${version}"
msg=${2:-"Release ${version} 新增日志项目"}
git tag -a "${version}" -m "${msg}"
git push origin "${version}"
echo  "推送标签"

