#!/bin/bash

#
# This script uploads/publishes binary artifacts to GitHub for a release
#
# Example:
#   GITHUB_TOKEN=... ./upload.sh
#   GITHUB_TOKEN=... ./upload.sh -r 0.1.0
#   GITHUB_TOKEN=... ./upload.sh --release 0.2.0
#

set -euo pipefail

red='\033[1;31m'
green='\033[1;32m'
yellow='\033[1;33m'
blue='\033[1;36m'
nocolor='\033[0m'

binary="gocert"
build_dir="./artifacts"


# set: release
function process_args {
  while [[ $# > 0 ]]; do
    key="$1"
    case $key in
      -r|--release)
      release="$2"
      shift
      ;;
    esac
    shift
  done

  release=${release:-"latest"}
  if [ "$release" != "latest" ] && [[ ! "$release" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    printf "${red}Release $release is not valid.${nocolor}\n"
    exit 1
  fi
}

function ensure_command {
  for cmd in $@; do
    which $cmd 1> /dev/null || (
      printf "${red}$cmd not available!${nocolor}\n"
      exit 1
    )
  done
}

function ensure_env_var {
  for var in $@; do
    if [ "${!var}" == "" ]; then
      printf "${red}$var is not set.${nocolor}\n"
      exit 1
    fi
  done
}

function ensure_repo_clean {
  status=$(git status --porcelain | tail -n 1)
  if [[ -n $status ]]; then
    printf "${red}Working direcrory is not clean.${nocolor}\n"
    exit 1
  fi
}

# set: repo
function get_repo_name {
  if [ "$(git remote -v)" == "" ]; then
    printf "${red}GitHub repo not fonud.${nocolor}\n"
    exit 1
  fi

  repo=$(
    git remote -v |
    sed -n 's/origin[[:blank:]]git@github.com://; s/.git[[:blank:]](push)// p'
  )
}

function upload_to_github {
  if [ ! -d "$build_dir" ]; then
    printf "${red}Artifacts not built.${nocolor}\n"
    exit 1
  fi

  upload_url=$(
    curl "https://api.github.com/repos/$repo/releases/$release" \
      -s \
      -X GET \
      -H "Authorization: token $GITHUB_TOKEN" \
      -H "Accept: application/vnd.github.v3+json" \
    | jq -r ".upload_url" | sed "s/{?name,label}//1"
  )

  for filepath in $build_dir/$binary-*; do
    filename=$(basename $filepath)
    mimetype=$(file -b --mime-type $filepath)

    printf "${blue}Uploading ${green}$filename${blue} to release ${yellow}$release${blue} ...${nocolor}\n"
    curl "$upload_url?name=$filename" \
      -s -o /dev/null \
      -X POST \
      -H "Authorization: token $GITHUB_TOKEN" \
      -H "Content-Type: $mimetype" \
      -H "Accept: application/vnd.github.v3+json" \
      --data-binary @$filepath
  done
}

function finish {
  if [ "$?" == "0" ]; then
    printf "${green}Done.${nocolor}\n"
  fi
}


trap finish EXIT

process_args "$@"
ensure_command "curl" "jq" "sed"
ensure_env_var "GITHUB_TOKEN"
#ensure_repo_clean

get_repo_name
upload_to_github
