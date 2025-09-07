#!/bin/bash

# Check if an argument was given
if [ $# -eq 0 ]; then
  echo "Usage: $0 <argument>"
  exit 1
fi

# Print the first argument
echo "You passed: $1"
