#!/usr/bin/env sh

#
# This script will install the latest version of gocert
#
# USAGE:
#   ./install.sh
#

# -o pipefail will cause the grep and head pipes exit with 141 code.
# https://stackoverflow.com/questions/19120263/why-exit-code-141-with-grep-q
set -eu


get_latest_release() {
  os="$1"
  arch="$2"
  release_url="https://github.com/moorara/gocert/releases"
  bin_pattern="/moorara/gocert/releases/download/v.*/gocert-$os-$arch"

  if hash curl 2>/dev/null; then
    content=$(curl -sL $release_url)
  elif hash wget 2>/dev/null; then
    content=$(wget -qO- $release_url)
  else
    printf "No command available to get %s\n" "$release_url"
    exit 1
  fi

  bin_path=$(echo "$content" | grep -o "$bin_pattern" | head -n 1)
  download_url="https://github.com$bin_path"
  latest_version=$(echo "$bin_path" | cut -d '/' -f6 | cut -d 'v' -f 2)
}

install_gocert() {
  download_url="$1"
  bin_path="/usr/local/bin/gocert"

  if hash curl 2>/dev/null; then
    curl -fsSL -o "$bin_path" "$download_url"
  elif hash wget 2>/dev/null; then
    wget -qO "$bin_path" "$download_url"
  else
    printf "No command available to download %s\n" "$download_url"
    exit 1
  fi

  chmod 755 $bin_path
}

main() {
  printf "Installing the latest release of gocert ...\n"

  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  arch=$(uname -m)

  if [ "$arch" = "i386" ]; then
    arch="386"
  elif [ "$arch" = "x86_64" ]; then
    arch="amd64"
  fi

  get_latest_release "$os" "$arch"
  install_gocert "$download_url"

  printf "gocert %s installed successfully.\n" "$latest_version"
}


main
