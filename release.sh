#!/bin/bash

# Ensure the script exits on any error
set -e

# Commit message as an argument
COMMIT_MSG=$1
INCREMENT_MAJOR=false

# Check if the --major flag is provided
if [[ "$2" == "--major" ]]; then
  INCREMENT_MAJOR=true
fi

# Check if the commit message is provided
if [ -z "$COMMIT_MSG" ]; then
  echo "Usage: $0 <commit-message> [--major]"
  exit 1
fi

# Get the latest tag
LATEST_TAG=$(git describe --tags `git rev-list --tags --max-count=1`)

# Extract the major, minor, and patch numbers
IFS='.' read -r -a VERSION_PARTS <<< "${LATEST_TAG#v}"

MAJOR=${VERSION_PARTS[0]}
MINOR=${VERSION_PARTS[1]}
PATCH=${VERSION_PARTS[2]}

# Increment the version number
if [ "$INCREMENT_MAJOR" = true ]; then
  MAJOR=$((MAJOR + 1))
  MINOR=0
  PATCH=0
else
  MINOR=$((MINOR + 1))
  PATCH=0
fi

NEW_VERSION="v$MAJOR.$MINOR.$PATCH"

# Stage and commit changes
git add .
git commit -m "$COMMIT_MSG"

# Create a new tag
git tag -a "$NEW_VERSION" -m "Release version $NEW_VERSION"

# Push commit and tag to the remote repository
git push origin master
git push origin "$NEW_VERSION"

echo "Released version $NEW_VERSION"