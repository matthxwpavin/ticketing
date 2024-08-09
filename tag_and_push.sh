#!/bin/bash

# Parse options
while getopts ":m:" opt; do
  case $opt in
    m) commit_message="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
        exit 1
    ;;
  esac
done

# Check if a commit message was provided
if [ -z "$commit_message" ]; then
  echo "Error: No commit message provided."
  echo "Usage: $0 -m \"commit message\""
  exit 1
fi

# Add all changes
git add .

# Commit with the provided message
git commit -m "$commit_message"

# Get the latest tag
latest_tag=$(git describe --tags `git rev-list --tags --max-count=1`)

# If there are no tags, start with v1.0.0
if [ -z "$latest_tag" ]; then
  new_tag="v1.0.0"
else
  # Use semver tool or logic to increment the tag
  # Here we're assuming the format vMAJOR.MINOR.PATCH
  IFS='.' read -r -a parts <<< "${latest_tag//v/}"
  major=${parts[0]}
  minor=${parts[1]}
  patch=${parts[2]}

  # Increment the patch version
  patch=$((patch + 1))

  # Form the new tag
  new_tag="v$major.$minor.$patch"
fi

echo "Creating new tag: $new_tag"

# Tag the current commit
git tag $new_tag

# Push the new tag to the remote repository
git push origin $new_tag
