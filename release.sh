#!/bin/bash

#
# This script
#   - releases current version,
#   - increases the version number for next release, and
#   - commits/pushes the changes with version tag to GitHub
#
# Example:
#   GITHUB_TOKEN=... release.sh
#   GITHUB_TOKEN=... release.sh minor
#   GITHUB_TOKEN=... release.sh major
#

set -euo pipefail


function process_args {
  component=${1:-patch}
  # component=${component,,}
  if [ "$component" != "patch" ] && [ "$component" != "minor" ] && [ "$component" != "major" ]; then
    echo "Version component $component is not valid."
    exit 1
  fi
}

function ensure_command {
  for cmd in $@; do
    which $cmd 1> /dev/null || (
      echo "$cmd not available!"
      exit 1
    )
  done
}

function ensure_env_var {
  for var in $@; do
    if [ "${!var}" == "" ]; then
      echo "$var is not set."
      exit 1
    fi
  done
}

function ensure_repo_ok {
  status=$(git status --porcelain | tail -n 1)
  if [[ -n $status ]]; then
    echo "Working direcrory is not clean."
    exit 1
  fi
}

function get_repo_name {
  if [ "$(git remote -v)" == "" ]; then
    echo "No GitHub repo is fonud."
    exit 1
  fi

  repo=$(git remote -v| grep push | grep -oE 'github\.com:.*\.git')
  repo=${repo/github.com:/}
  repo=${repo/.git/}
}

function enable_master_push {
  echo "Temporarily enabling push to master branch ..."

  repo=$1

  github_state_original=$(
    curl "https://api.github.com/repos/$repo/branches/master" \
      -s -X GET \
      -H "Authorization: token $GITHUB_TOKEN" \
      -H "Accept: application/vnd.github.loki-preview" \
    | jq '{ protection: .protection }'
  )
  github_state_disabled=$(
    echo $github_state_original \
    | jq '.protection.required_status_checks.enforcement_level="off"'
  )
  curl "https://api.github.com/repos/$repo/branches/master" \
    -s -o /dev/null -X PATCH \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.loki-preview" \
    -d "$github_state_disabled"
}

function disable_master_push {
  echo "Re-disabling push to master branch ..."

  repo=$1
  github_state_original=$2

  curl "https://api.github.com/repos/$repo/branches/master" \
    -s -o /dev/null -X PATCH \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.loki-preview" \
    -d "$github_state_original"
}

function release_version {
  component=$1

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

  echo "Releasing current version ..."
  echo "$release_version" > version/VERSION
  {
    git add .
    git commit -m "Releasing v$release_version"
    git tag "v$release_version"
    git push
    git push --tags
  } &> /dev/null

  echo "Preparing next version ..."
  echo "$next_version" > version/VERSION
  {
    git add .
    git commit -m "Beginning v$next_version [skip ci]"
    git push
  } &> /dev/null
}

function generate_changelog {
  echo "Generating changelog ..."

  CHANGELOG_GITHUB_TOKEN=$GITHUB_TOKEN \
  github_changelog_generator \
    --no-filter-by-milestone \
    --exclude-labels question,duplicate,invalid,wontfix \
  &> /dev/null

  {
    git add .
    git commit -m "Change Log [skip ci]"
    git push
  } &> /dev/null
}


process_args $@
ensure_command curl git github_changelog_generator
ensure_env_var GITHUB_TOKEN
ensure_repo_ok

get_repo_name
enable_master_push $repo
release_version $component
generate_changelog
disable_master_push $repo $github_state_original
