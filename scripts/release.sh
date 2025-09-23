#!/bin/bash
# release.sh: Bump version, create tag, add release notes, and create draft GitHub release

set -e

ORIGIN=${ORIGIN:-origin}
COMMIT=${COMMIT:-HEAD}

if [[ $(git status --porcelain) != "" ]]; then
  echo "Error: repo is dirty. Run git status, clean repo and try again."
  exit 1
elif [[ $(git status --porcelain -b | grep -e "ahead" -e "behind") != "" ]]; then
  echo "Error: repo has unpushed commits. Push commits to remote and try again."
  exit 1
fi 


# Check if user has push access
if ! git ls-remote --exit-code "$ORIGIN" >/dev/null 2>&1; then
  echo "Error: Cannot access remote repository. Ensure you have push permissions."
  exit 1
fi

BUMP=${1:-patch}
latest_tag=$(git describe --tags --abbrev=0)
base_version=$(echo "$latest_tag" | sed 's/^v//')
IFS=. read major minor patch <<<"$base_version"
case "$BUMP" in
  major|breaking) new_version="v$((major+1)).0.0" ;;
  minor|feature) new_version="v$major.$((minor+1)).0" ;;
  patch|bugfix) new_version="v$major.$minor.$((patch+1))" ;;
  *) echo "Unknown bump type: $BUMP" >&2; exit 1 ;;
esac

echo "Creating tag $new_version"
git tag -a "$new_version" "$COMMIT" -m "$new_version"
git push "$ORIGIN" "$new_version"


