#!/bin/bash

#
# This script
#   - releases current version,
#   - increases the version number for next release, and
#   - commits/pushes the changes with version tag to GitHub
#
# Example:
#   GITHUB_TOKEN=... ./release.sh
#   GITHUB_TOKEN=... ./release.sh minor -c "..."
#   GITHUB_TOKEN=... ./release.sh major --comment "..."
#

set -euo pipefail

red='\033[1;31m'
green='\033[1;32m'
yellow='\033[1;33m'
blue='\033[1;36m'
nocolor='\033[0m'


# set: comment, component
function process_args {
  while [[ $# > 0 ]]; do
    key="$1"
    case $key in
      patch)
      component="$1"
      ;;
      minor)
      component="$1"
      ;;
      major)
      component="$1"
      ;;
      -c|--comment)
      comment="$2"
      shift
      ;;
    esac
    shift
  done

  comment=${comment:-""}
  component=${component:-patch}
  # component=${component,,}
  if [ "$component" != "patch" ] && [ "$component" != "minor" ] && [ "$component" != "major" ]; then
    printf "${red}Version component $component is not valid.${nocolor}\n"
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

function ensure_repo_ok {
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

function enable_master_push {
  printf "${yellow}Temporarily enabling push to master branch ...${nocolor}\n"
  curl "https://api.github.com/repos/$repo/branches/master/protection/enforce_admins" \
    -s -o /dev/null \
    -X DELETE \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json"
}

function disable_master_push {
  printf "${yellow}Re-disabling push to master branch ...${nocolor}\n"
  curl "https://api.github.com/repos/$repo/branches/master/protection/enforce_admins" \
    -s -o /dev/null \
    -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json"
}

function generate_changelog {
  printf "${blue}Generating changelog v$release_version ...${nocolor}\n"
  CHANGELOG_GITHUB_TOKEN=$GITHUB_TOKEN \
  github_changelog_generator \
    --no-filter-by-milestone \
    --exclude-labels question,duplicate,invalid,wontfix \
  &> /dev/null

  {
    git add .
    git commit -m "Changelog v$release_version [skip ci]"
    git push
  } &> /dev/null
}

function create_github_release {
  changelog=$(
    git diff v$release_version CHANGELOG.md |
    sed '/^+/!d; /^+++\|+##/d; s/^+//; s/\\/\\\\/g;' |
    sed -E ':a; N; $!ba; s/\r{0,1}\n/\\n/g'
  )

  printf "${blue}Creating github release v$release_version ...${nocolor}\n"
  curl "https://api.github.com/repos/$repo/releases" \
    -s -o /dev/null \
    -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    -d @- << EOF
    {
      "tag_name": "v$release_version",
      "target_commitish": "master",
      "name": "$release_version",
      "body": "$comment\n\n$changelog",
      "draft": false,
      "prerelease": false
    }
EOF
}

# set: release_version
function release_version {
  version=$(cat version/VERSION)
  components=(${version//./ })
  major="${components[0]}"
  minor="${components[1]}"
  patch="${components[2]/-0/}"

  case "$component" in
    patch)  release_version="$major.$minor.$patch"       next_version="$major.$minor.$(( patch + 1 ))-0"  ;;
    minor)  release_version="$major.$(( minor + 1 )).0"  next_version="$major.$(( minor + 1 )).1-0"       ;;
    major)  release_version="$(( major + 1 )).0.0"       next_version="$(( major + 1 )).0.1-0"            ;;
  esac

  printf "${blue}Releasing current version v$release_version ...${nocolor}\n"
  echo "$release_version" > version/VERSION
  {
    git add .
    git commit -m "Releasing v$release_version"
    git tag "v$release_version"
    git push
    git push --tags
  } &> /dev/null

  generate_changelog
  create_github_release

  printf "${blue}Preparing next version v$next_version ...${nocolor}\n"
  echo "$next_version" > version/VERSION
  {
    git add .
    git commit -m "Beginning v$next_version [skip ci]"
    git push
  } &> /dev/null
}

function finish {
  disable_master_push

  printf "${green}Done.${nocolor}\n"
}


trap finish EXIT

process_args "$@"
ensure_command "sed" "curl" "git" "github_changelog_generator"
ensure_env_var "GITHUB_TOKEN"
ensure_repo_ok

get_repo_name
enable_master_push
release_version
