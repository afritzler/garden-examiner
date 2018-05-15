#!/bin/bash


setup() 
{
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
