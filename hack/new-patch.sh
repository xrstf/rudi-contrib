#!/usr/bin/env bash

set -euo pipefail

cd $(dirname $0)/..

if (( $# < 1 )); then
  echo "Usage: hack/new-patch.sh module[ module module module]"
  exit 2
fi

for module in "$@"; do
  if [ ! -d "$module" ]; then
    echo "Error: No such module: $module"
    exit 1
  fi
done

for module in "$@"; do
  # remove trailing slashes
  module="${module%/}"

  # list all tags for this module,
  # turn "module/vA.B.C" into "vA.B.C",
  # then sort version,
  # then take the last one, the most recent.
  latest="$(git tag --list "$module/*" | xargs -n 1 basename | sort --version-sort | tail -n 1)"

  # trim leading v
  latest="${latest#v}"

  # get next patch release
  next="$(semver next patch "$latest")"

  echo "Bumping $module module from $latest to $next…"
  git tag --message "$module version $next" "$module/v$next"
done
