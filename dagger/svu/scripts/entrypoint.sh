#!/bin/sh

set -eu

readonly SVU_COMMAND="${SVU_COMMAND}"
readonly SVU_METADATA="${SVU_METADATA}"
readonly SVU_PATTERN="${SVU_PATTERN}"
readonly SVU_PRE_RELEASE="${SVU_PRE_RELEASE}"
readonly SVU_BUILD="${SVU_BUILD}"
readonly SVU_PREFIX="${SVU_PREFIX}"
readonly SVU_SUFFIX="${SVU_SUFFIX}"
readonly SVU_TAG_MODE="${SVU_TAG_MODE}"
readonly SVU_VERSION_INFO_FILE="${SVU_VERSION_INFO_FILE}"
readonly SVU_VERSION_WITHOUT_PREFIX_INFO_FILE="${SVU_VERSION_WITHOUT_PREFIX_INFO_FILE}"

run_cmd() {
  # Variables for optional parameters.
  metadata=""
  pattern=""
  prerelease=""
  build=""
  suffix=""
  prefix=""

  # Set optional parameters.
  if [ -n "${SVU_PATTERN}" ]; then
    pattern="--pattern=${SVU_PATTERN}"
  fi

  if [ -n "${SVU_PREFIX}" ]; then
    prefix="--prefix=${SVU_PREFIX}"
  fi

  if [ -n "${SVU_SUFFIX}" ]; then
    suffix="--suffix=${SVU_SUFFIX}"
  fi

  if [ -n "${SVU_TAG_MODE}" ]; then
    tag_mode="--tag-mode=${SVU_TAG_MODE}"
  fi

  if [ -n "${SVU_METADATA}" ]; then
    if [ "${SVU_METADATA}" = "true" ]; then
      metadata="--metadata"
    elif [ "${SVU_METADATA}" = "false" ]; then
      metadata="--no-metadata"
    else
      echo "::error Unknown metadata value: ${SVU_METADATA}"
      exit 1
    fi
  fi

  if [ -n "${SVU_PRE_RELEASE}" ]; then
    if [ "${SVU_PRE_RELEASE}" = "true" ]; then
      prerelease="--pre-release"
    elif [ "${SVU_PRE_RELEASE}" = "false" ]; then
      prerelease="--no-pre-release"
    else
      echo "::error Unknown pre-release value: ${SVU_PRE_RELEASE}"
      exit 1
    fi
  fi

  if [ -n "${SVU_BUILD}" ]; then
    if [ "${SVU_BUILD}" = "true" ]; then
      build="--build"
    elif [ "${SVU_BUILD}" = "false" ]; then
      build="--no-build"
    else
      echo "::error Unknown build value: ${SVU_BUILD}"
      exit 1
    fi
  fi

  # shellcheck disable=SC2086
  version=$(svu ${SVU_COMMAND} ${pattern} ${prefix} ${suffix} ${metadata} ${prerelease} ${tag_mode-mode} ${build})

  # shellcheck disable=SC2086
  version_without_prefix=$(svu ${SVU_COMMAND} ${pattern} ${prefix} ${suffix} ${metadata} ${prerelease} ${tag_mode-mode} ${build} --strip-prefix)

  # Write version to file.
  echo "${version}" > "${SVU_VERSION_INFO_FILE}"
  echo "${version_without_prefix}" > "${SVU_VERSION_WITHOUT_PREFIX_INFO_FILE}"
}

run_cmd "$@"