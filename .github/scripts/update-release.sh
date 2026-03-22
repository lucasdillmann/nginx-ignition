#!/bin/bash
set -e

CHANGELOG_FILE="CHANGELOG.md"
DOCKER_IMAGE_NAME="dillmann/nginx-ignition"
DOCKER_IMAGE_HASH=$1
PRERELEASE=$2

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

VERSION=$(grep -m 1 "^## " "$CHANGELOG_FILE" | sed 's/## //')
if [[ -z "$VERSION" ]]; then
  echo "Error: Could not find version in $CHANGELOG_FILE" >&2
  exit 1
fi

VERSION_LINE=$(grep -n -m 1 "^## $VERSION" "$CHANGELOG_FILE" | cut -d: -f1)
NEXT_VERSION_LINE=$(grep -n "^## " "$CHANGELOG_FILE" | sed -n '2p' | cut -d: -f1)

if [[ -z "$NEXT_VERSION_LINE" ]]; then
  DESCRIPTION=$(sed -n "$((VERSION_LINE + 1)),\$p" "$CHANGELOG_FILE")
else
  DESCRIPTION=$(sed -n "$((VERSION_LINE + 1)),$((NEXT_VERSION_LINE - 1))p" "$CHANGELOG_FILE")
fi

BODY_FILE=$(mktemp)
echo "$DESCRIPTION" | sed -e :a -e '/^\n*$/{$d;N;ba' -e '}' > "$BODY_FILE"
echo "" >> "$BODY_FILE"

TAG="$VERSION"
if [[ "$PRERELEASE" == "true" ]]; then
  TAG="$VERSION-snapshot"
fi

echo "Docker image: [$DOCKER_IMAGE_NAME:$TAG](https://hub.docker.com/layers/$DOCKER_IMAGE_NAME/$TAG/images/$DOCKER_IMAGE_HASH)" >> "$BODY_FILE"

if gh release view "$VERSION" >/dev/null 2>&1; then
  echo "Release $VERSION already exists. Updating..."
  gh release edit "$VERSION" --notes-file "$BODY_FILE" --prerelease="$PRERELEASE"
else
  echo "Creating new release $VERSION..."
  gh release create "$VERSION" --notes-file "$BODY_FILE" --prerelease="$PRERELEASE" --title "$VERSION"
fi

rm "$BODY_FILE"
echo "Release $VERSION created/updated successfully."
