#!/bin/bash

# Check if a commit message was provided
if [ -z "$1" ]; then
  echo "Error: No commit message provided."
  echo "Usage: $0 \"commit message\" [version]"
  exit 1
fi

commit_message=$1
manual_version=$2

# Validate manual version if provided
if [ ! -z "$manual_version" ]; then
  if [[ ! $manual_version =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Provided version is not a valid SemVer format (vMAJOR.MINOR.PATCH)."
    exit 1
  fi
fi

# Add all changes
git add .

# Commit with the provided message
git commit -m "$commit_message"

# Determine the new tag version
if [ -z "$manual_version" ]; then
  # Get the latest tag
  latest_tag=$(git tag -l "v*" --sort=-v:refname | head -n 1)

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
else
  new_tag="$manual_version"
fi

echo "Creating new tag: $new_tag"

# Tag the current commit
git tag $new_tag

# Push the new tag to the remote repository
git push origin $new_tag
