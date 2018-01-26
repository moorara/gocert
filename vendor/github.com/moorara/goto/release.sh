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
purple='\033[1;35m'
blue='\033[1;36m'
nocolor='\033[0m'


function whitelist_variable {
  if [[ ! $2 =~ (^|[[:space:]])$3($|[[:space:]]) ]]; then
    printf "${red}Invalid $1 $3${nocolor}\n"
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


# set: component, comment
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
  whitelist_variable "version component" "patch minor major" "$component"
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

# set: release_version
function release_current_version {
  printf "${blue}Releasing current version ...${nocolor}\n"

  version=$(cat VERSION)
  components=(${version//./ })
  major="${components[0]}"
  minor="${components[1]}"
  patch="${components[2]/-0/}"

  case "$component" in
    patch)  release_version="$major.$minor.$patch"       next_version="$major.$minor.$(( patch + 1 ))-0"  ;;
    minor)  release_version="$major.$(( minor + 1 )).0"  next_version="$major.$(( minor + 1 )).1-0"       ;;
    major)  release_version="$(( major + 1 )).0.0"       next_version="$(( major + 1 )).0.1-0"            ;;
  esac

  echo "$release_version" > VERSION
}

function generate_changelog {
  printf "${blue}Generating changelog ...${nocolor}\n"
  CHANGELOG_GITHUB_TOKEN=$GITHUB_TOKEN \
  github_changelog_generator \
    --no-filter-by-milestone \
    --exclude-labels question,duplicate,invalid,wontfix \
    --future-release v$release_version \
  &> /dev/null
}

function create_github_release {
  printf "${blue}Creating github release ${purple}v$release_version${blue} ...${nocolor}\n"

  {
    git add .
    git commit -m "Releasing v$release_version"
    git tag "v$release_version"
    git push
    git push --tags
  } &> /dev/null

  changelog=$(
    git diff HEAD~1 CHANGELOG.md |
    sed '/^+/!d; /^+++\|+##/d; s/^+//; s/\\/\\\\/g;' |
    sed -E ':a; N; $!ba; s/\r{0,1}\n/\\n/g'
  )

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

function prepare_next_version {
  printf "${blue}Preparing next version v$next_version ...${nocolor}\n"
  echo "$next_version" > VERSION
  {
    git add .
    git commit -m "Beginning v$next_version [skip ci]"
    git push
  } &> /dev/null
}

function finish {
  status=$?

  disable_master_push

  if [ "$status" == "0" ]; then
    printf "${green}Done.${nocolor}\n"
  fi
}


trap finish EXIT

process_args "$@"
ensure_command "sed" "curl" "jq" "git" "github_changelog_generator"
ensure_env_var "GITHUB_TOKEN"
get_repo_name
ensure_repo_clean

enable_master_push
release_current_version
generate_changelog
create_github_release
prepare_next_version
