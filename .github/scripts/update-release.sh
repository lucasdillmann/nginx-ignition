#!/bin/bash
set -e

CHANGELOG_FILE="CHANGELOG.md"
DOCKER_IMAGE_NAME="dillmann/nginx-ignition"
DOCKER_IMAGE_HASH=$1
PRERELEASE=$2
VERSION=$3
DRAFT=${4:-false}

if ! command -v gh &> /dev/null; then
  echo "Error: gh CLI is not installed" >&2
  exit 1
fi

if [[ -z "$DOCKER_IMAGE_HASH" ]]; then
  echo "Error: Docker image hash is required." >&2
  exit 1
fi

if [[ ! "$DOCKER_IMAGE_HASH" =~ ^sha256- ]]; then
  echo "Error: Docker image hash must start with 'sha256-'" >&2
  exit 1
fi

if [[ -z "$PRERELEASE" ]]; then
  PRERELEASE="true"
fi

if [[ -z "$VERSION" ]]; then
  echo "Error: Version argument is required." >&2
  exit 1
fi

VERSION_LINE=$(grep -n -m 1 "^## $VERSION" "$CHANGELOG_FILE" | cut -d: -f1)
if [[ -z "$VERSION_LINE" ]]; then
  echo "Error: Could not find version $VERSION in $CHANGELOG_FILE" >&2
  exit 1
fi

NEXT_VERSION_LINE=$(grep -n "^## " "$CHANGELOG_FILE" | awk -F: -v version_line="$VERSION_LINE" '$1 > version_line { print $1; exit }')
if [[ -z "$NEXT_VERSION_LINE" ]]; then
  DESCRIPTION=$(sed -n "$((VERSION_LINE + 1)),\$p" "$CHANGELOG_FILE")
else
  DESCRIPTION=$(sed -n "$((VERSION_LINE + 1)),$((NEXT_VERSION_LINE - 1))p" "$CHANGELOG_FILE")
fi

DESCRIPTION=$(printf "%s\n" "$DESCRIPTION" | awk '
function ltrim(value) {
  sub(/^[[:space:]]+/, "", value)
  return value
}
function rtrim(value) {
  sub(/[[:space:]]+$/, "", value)
  return value
}
function flush_current() {
  if (current_line != "") {
    print current_line
    current_line = ""
  }
}
{
  line = rtrim($0)

  if (line ~ /^[[:space:]]*$/) {
    flush_current()
    if (!printed_blank_line) {
      print ""
      printed_blank_line = 1
    }
    next
  }

  printed_blank_line = 0

  if (line ~ /^[[:space:]]*-[[:space:]]+/) {
    flush_current()
    current_line = line
    next
  }

  if (current_line != "") {
    current_line = current_line " " ltrim(line)
    next
  }

  print line
}
END {
  flush_current()
}')

BODY_FILE=$(mktemp)
trap 'rm -f "$BODY_FILE"' EXIT

printf "%s\n" "$DESCRIPTION" | sed -e :a -e '/^$/N;/\n$/ba' > "$BODY_FILE"
echo "" >> "$BODY_FILE"

TAG="$VERSION"
if [[ "$PRERELEASE" == "true" ]]; then
  TAG="$VERSION-snapshot"
fi

echo "Docker image: [$DOCKER_IMAGE_NAME:$TAG](https://hub.docker.com/layers/$DOCKER_IMAGE_NAME/$TAG/images/$DOCKER_IMAGE_HASH)" >> "$BODY_FILE"

if gh release view "$TAG" >/dev/null 2>&1; then
  echo "Release $TAG already exists. Updating..."
  if [[ "$DRAFT" == "true" ]]; then
    gh release edit "$TAG" --notes-file "$BODY_FILE" --prerelease="$PRERELEASE" --draft=true
  else
    gh release edit "$TAG" --notes-file "$BODY_FILE" --prerelease="$PRERELEASE" --draft=false
  fi
else
  echo "Creating new release $TAG..."
  if [[ "$DRAFT" == "true" ]]; then
    gh release create "$TAG" --notes-file "$BODY_FILE" --prerelease="$PRERELEASE" --title "$VERSION" --draft
  else
    gh release create "$TAG" --notes-file "$BODY_FILE" --prerelease="$PRERELEASE" --title "$VERSION"
  fi
fi

if [[ "$PRERELEASE" == "false" ]]; then
  SNAPSHOT_TAG="$VERSION-snapshot"
  if gh release view "$SNAPSHOT_TAG" >/dev/null 2>&1; then
    echo "Promoting release. Deleting snapshot release $SNAPSHOT_TAG..."
    gh release delete "$SNAPSHOT_TAG" --yes --cleanup-tag
  fi
fi

echo "Release for $TAG created/updated successfully."
