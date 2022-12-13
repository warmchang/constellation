#!/usr/bin/env bash
# Copyright (c) Edgeless Systems GmbH
#
# SPDX-License-Identifier: AGPL-3.0-only

set -euo pipefail
shopt -s inherit_errexit

cdn_url="https://cdn.confidential.cloud"

function usage() {
  cat << 'EOF'
Remove images from cloud provider, image and versions API.

Usage: rm-image.sh <command> [options]

Options:
  --help              Show this help

Commands:
  ref <ref>           Ref to delete, deletes all versions of ref
  version <version>   Single version to delete
EOF
}

POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
  case $1 in
  --help)
    usage
    exit 0
    ;;
  -*)
    echo "Unknown option $1"
    echo
    usage
    exit 1
    ;;
  *)
    POSITIONAL_ARGS+=("$1") # save positional arg
    shift                   # past argument
    ;;
  esac
done

set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters

function error() {
  echo "[Error] $1"
}

function warn() {
  echo "[Warn] $1"
}

function setup() {
  # Install image-gallery extension if not already installed
  if ! az extension list-available --query "[?name=='image-gallery'].name" | grep -q image-gallery; then
    az extension add --name image-gallery
    if [[ $? -ne 0 ]]; then
      error "failed to install Azure image-gallery extension"
      return 1
    fi
  fi
}

function delete() {
  local ref="$1"
  local stream="$2"
  local version="$3"

  image_info_path="constellation/v1/ref/${ref}/stream/${stream}/image/${version}/info.json"
  image_info_url="${cdn_url}/${image_info_path}"

  status=$(curl -s -o /dev/null -w "%{http_code}" "${image_info_url}")
  if [[ ${status} != "200" ]]; then
    error "image info not found at ${image_info_url}"
    return 1
  fi

  local image_info
  image_info=$(curl -sL "${image_info_url}")
  jq <<< "${image_info}"

  local trash_after_error=false

  #
  # Delete AWS images
  #

  local aws_regions
  aws_regions=$(jq -r '.aws | keys[]' <<< "${image_info}")

  for region in ${aws_regions}; do
    local aws_ami
    aws_ami=$(jq -r ".aws[\"${region}\"]" <<< "${image_info}")

    echo "Deleting AWS image ${aws_ami} in region ${region}"

    out=$(
      aws ec2 deregister-image \
        --image-id "${aws_ami}" \
        --region "${region}" \
        --output json \
        --dry-run 2>&1 ||
        :
    )
    res=$?
    if [[ -n ${res} && ${out} == *"InvalidAMIID.NotFound"* ]]; then
      warn "not found: AWS image ${aws_ami} in region ${region}"
    elif [[ -n ${res} ]]; then
      error "failed to delete AWS image ${aws_ami} in region ${region}: ${out}"
      trash_after_error=true
    fi
  done

  #
  # Delete Azure images
  #

  local az_image_types
  az_image_types=$(jq -r '.azure | keys[]' <<< "${image_info}")

  for az_image_type in ${az_image_types}; do
    local az_image_uri
    az_image_uri=$(jq -r ".azure[\"${az_image_type}\"]" <<< "${image_info}")

    echo "Deleting Azure image ${az_image} of type ${az_image_type}"

    case "${az_image_type}" in
    "cvm") ;;
      # local imageRegexp='^/CommunityGalleries/ConstellationCVM'

    "trustedlaunch")
      local image_regexp='^/subscriptions/([[:alnum:]-]+)/resourceGroups/([[:alnum:]-]+)/providers/Microsoft.Compute/galleries/([[:alnum:]-]+)/images/([[:alnum:]-]+)$'
      if [[ ${az_image_uri} =~ ${image_regexp} ]]; then
        local az_subscription="${BASH_REMATCH[1]}"
        local az_resource_group="${BASH_REMATCH[2]}"
        local az_gallery="${BASH_REMATCH[3]}"
        local az_image="${BASH_REMATCH[4]}"
      else
        error "invalid Azure image format, expected /subscriptions/<subscription>/resourceGroups/<resource_group>/providers/Microsoft.Compute/galleries/<gallery>/images/<image>"
        trash_after_error=true
        continue
      fi

      out=$(
        az sig image-definition delete \
          --gallery-image-definition "${az_image}" \
          --gallery-name "${az_gallery}" \
          --resource-group "${az_resource_group}" \
          --subscription "${az_subscription}" 2>&1 ||
          :
      )
      res=$?
      ;;
    *)
      error "unknown Azure image type ${az_image_type}"
      trash_after_error=true
      continue
      ;;
    esac

    if [[ -n ${res} && ${out} == *"ResourceNotFound"* ]]; then
      warn "not found: Azure image ${az_image} of type ${az_image_type}"
    elif [[ -n ${res} ]]; then
      error "failed to delete Azure image ${az_image} of type ${az_image_type}: ${out}"
      trash_after_error=true
    fi
  done

  #
  # Delete GCP images
  #

  local gcp_image_type
  gcp_image_type=$(jq -r '.gcp | keys[]' <<< "${image_info}")

  for gcp_image in ${gcp_image_type}; do
    local gcp_image_uri
    gcp_image_uri=$(jq -r ".gcp[\"${gcp_image_type}\"]" <<< "${image_info}")

    local image_regexp='^projects/([[:alnum:]-]+)/global/images/([[:alnum:]-]+)$'
    if [[ ${gcp_image_uri} =~ ${image_regexp} ]]; then
      local gcp_project="${BASH_REMATCH[1]}"
      local gcp_image="${BASH_REMATCH[2]}"
    else
      error "invalid GCP image format, expected projects/<project>/global/images/<image>"
      trash_after_error=true
      continue
    fi

    out=$(
      gcloud compute images delete "${gcp_image}" \
        --quiet \
        --project "${gcp_project}" 2>&1 ||
        :
    )
    res=$?
    if [[ -n ${res} && ${out} == *"notFound"* ]]; then
      warn "not found: GCP image ${gcp_image}"
    elif [[ -n ${res} ]]; then
      error "failed to delete GCP image ${gcp_image}: ${out}"
      trash_after_error=true
    fi
  done

  # TODO handle error, move to trash
  echo "${trash_after_error}"
}

function delete_single_version() {
  local version_str="$1"
  if [[ -z ${version_str} ]]; then
    error "missing version value"
  fi

  #shellcheck disable=SC2310
  if setup; then
    return 1
  fi

  local versionRegexp='^ref/([[:alnum:]-]+)/stream/([[:alnum:]-]+)/([[:alnum:].-]+)$'
  if [[ ${version_str} =~ ${versionRegexp} ]]; then
    local ref="${BASH_REMATCH[1]}"
    local stream="${BASH_REMATCH[2]}"
    local version="${BASH_REMATCH[3]}"
  else
    error "invalid version format, expected ref/<ref>/stream/<stream>/<version>"
  fi

  echo "Deleting version ${version} of ref ${ref} and stream ${stream}"

  delete "${ref}" "${stream}" "${version}"
  return $?
}

case $1 in
ref)
  # Canonicalize ref format (e.g. "feat/foo/bar" -> "feat-foo-bar")
  ref=$(echo -n "$2" | tr -c '[:alnum:]' '-')
  shift # past argument
  shift # past value

  echo "Not implemented yet"
  exit 1
  ;;
version)
  delete_single_version "${2-}"
  exit $?
  ;;
*)
  echo "Unknown command $1"
  echo
  usage
  exit 1
  ;;
esac
