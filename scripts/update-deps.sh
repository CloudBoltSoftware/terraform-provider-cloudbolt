#!/bin/sh
###
# Usage:
# Run this when updating the contents of `go.mod`.
# Specifically when updating the version of a cloudbolt-go-sdk library.
###

# We set GOPRIVATE to the cloudboltsoftware org until the Terraform Provider is fully open source
export GOPRIVATE="github.com/cloudboltsoftware"

echo "🔂  updating contents of 'go.sum'"

echo "🧹  tidying modules"

# Add missing and remove unused modules
# The `|| { ... }` runs `...` if the output of `go mod tidy` is non-zero
# In this case we print an informative error and exit with a non-zero status
go mod tidy || { echo "❌  Failed to tidy modules" ; exit 1; }

echo "🔎  verifying modules"

# Verify the contents of dependencies
# This command prints to the screen, but we want emojis in our messages so we supress that with a pipe to /dev/null
go mod verify > /dev/null || { echo "❌  Failed to verify modules" ; exit 1; }

echo "🎉  done"
