#!/bin/sh -e

TAG_PATTERN='[0-9].[0-9]'
DEFAULT_NEXT_VERSION="0.1"
RELEASE="${RELEASE:-false}"

# Get last tag that complies with TAG_PATTERN.
get_last_tag() {
    git describe --tags $(git rev-list --tags="${TAG_PATTERN}" --max-count=1)
}

# Get next version by using the last tag and incrementing the second digit.
# If there is no tag yet, then it selects DEFAULT_NEXT_VERSION.
get_next_version() {
    local last_tag=$(get_last_tag)
    [ "${last_tag}" == "" ] && echo "${DEFAULT_NEXT_VERSION}" \
                            || echo "${last_tag}" | awk -F. '{print $1"."$2+1}'
}

get_version() {
    [ "${RELEASE}" == "true" ] && get_next_version || get_last_tag
}

get_revision() {
    local sha=$(git rev-parse --short HEAD)
    [ "${RELEASE}" == "true" ] && echo "0.g${sha}" \
                               || echo "$(git rev-list --count $(get_last_tag)..HEAD).g${sha}"
}

# Get the release notes from a last tag to HEAD.
get_release_notes() {
    local range="$(get_last_tag)"
    [ "${range}" == "" ] && range="HEAD" || range="${range}...HEAD"
    git log "${range}" --pretty=format:' - [%h] %s' --reverse
}

$@
