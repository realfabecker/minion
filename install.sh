#!/usr/bin/env bash
set -e

# install base configuration
install_dir="$HOME/.local/bin"
download_dir=$(mktemp -d)
download_file=kevin-0.0.1-linux-amd64.tar.gz
repo_url=https://api.github.com/repos/realfabecker/kevin/releases/latest

# bash message log with info format
info() {
  command printf '\033[1;32m%12s\033[0m %s\n' "$1" "$2" 1>&2
}

# bash message log with erro format
error() {
  command printf '\033[1;31mError\033[0m: %s\n\n' "$1" 1>&2
}

check_profile_config() {
  if [[ $(command -v kevin && echo "ok" || echo "no") == "no" ]]; then
    p="$PATH"
    cat <<EOF

  Please, include the following lines to your profile configuration

  export $p:$install_dir

EOF
  fi
}

download_latest_release() {
  info "Downloading" "release from github ${repo_url}"
  curl -s "$repo_url" \
    | grep "browser_download_url.*linux" \
    | sed -E 's/.*"(http.*)"/\1/g' \
    | xargs -n1 -I{} curl -sL -o "${download_dir}/${download_file}" {}

  if [[ ! -f "${download_dir}/${download_file}" ]];then
    error "unable to download release from github"
    exit 1
  fi;
}

install_latest_release() {
  info "Extracting" "binary into temporary directory"
  if [[ ! -d "${install_dir}" ]]; then
    mkdir -p "${install_dir}"
  fi
  tar -C "${install_dir}" -xf "${download_dir}/${download_file}" kevin
}

download_latest_release
install_latest_release
info "Completed" "kevin installation at $install_dir/kevin"
check_profile_config