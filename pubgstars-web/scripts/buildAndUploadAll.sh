#!/bin/bash

find ../cmd/* -type d -exec basename {} \; |
  while read -r name; do
    echo "### building: ${name}"
    ./buildAndUpload.sh "${name}"
    echo "### completed: ${name}"
  done
