#!/bin/bash

gex_alias()
{
  exec 4>&1
  script="$(\./gex "$@" 3>&1 >&4 )"
  if [ -n "$script" ]; then
    source <(echo "$script")
  fi
}

alias gex="gex_alias"
