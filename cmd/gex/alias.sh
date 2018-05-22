#!/bin/bash


setup()
{
  os=`uname`
  if [[ "$os" == 'Darwin' ]]; then
    if brew ls --versions coreutils > /dev/null; then
      export PATH="/usr/local/opt/coreutils/libexec/gnubin:$PATH"
    else
      echo 'run: "brew install coreutils" first'
      exit 1
    fi
  fi
  local dir="$(readlink -f "$(dirname "${BASH_SOURCE[0]}")")"
  echo $dir
  alias gex="gex_alias \"$dir\""
}

gex_alias()
{
  exec 4>&1
  local script="$("$1/gex" "${@:2}" 3>&1 >&4 || { exec 4>&-; exit 1; })"
  exec 4>&-
  if [ -n "$script" ]; then
    source <(echo "$script")
  fi
}

setup
unset setup
