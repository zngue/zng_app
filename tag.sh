#!/bin/bash

version=${1:-"v0.0.1"}
git tag -d "${version}"
git push origin :refs/tags/"${version}"
msg=${2:-"Release ${version}"}
git tag -a "${version}" -m "${msg}"
git push origin "${version}"




