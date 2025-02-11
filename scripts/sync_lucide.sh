#!/bin/bash

set -e

if [ -f .env ]; then
    source .env
fi

function error_exit {
  echo "Error: $1" >&2
  exit 1
}

if [ "$#" -gt 0 ]; then
  VERSION="$1"
else
  if [ -z "${LUCIDE_VERSION:-}" ]; then
    error_exit "Lucide version not specified. Provide a version via command line argument or set LUCIDE_VERSION in .env."
  fi
  VERSION="$LUCIDE_VERSION"
fi

echo "Using Lucide version: $VERSION"

# Create necessary directories.
mkdir -p tmp lucide

# Define the URL and output ZIP file path.
ZIP_URL="https://github.com/lucide-icons/lucide/archive/refs/tags/${VERSION}.zip"
ZIP_PATH="./tmp/lucide-${VERSION}.zip"

echo "Downloading ${ZIP_URL} to ${ZIP_PATH} ..."

# Download the ZIP using curl or wget.
if command -v curl &> /dev/null; then
  curl -L -o "$ZIP_PATH" "$ZIP_URL"
elif command -v wget &> /dev/null; then
  wget -O "$ZIP_PATH" "$ZIP_URL"
else
  error_exit "Neither curl nor wget is installed."
fi

echo "Download completed."

# Unzip the downloaded file into the tmp directory.
echo "Unzipping ${ZIP_PATH} ..."
rm -rf "./tmp/lucide-${VERSION}"
unzip -q "$ZIP_PATH" -d tmp

# Determine the name of the extracted directory.
# Typically, GitHub creates a folder named lucide-[VERSION] or similar.
EXTRACTED_DIR=$(find tmp -mindepth 1 -maxdepth 1 -type d -name "lucide*" | head -n 1)

if [ -z "$EXTRACTED_DIR" ]; then
  error_exit "Could not locate the extracted directory."
fi

echo "Extracted directory: $EXTRACTED_DIR"

# Check if the "icons" folder exists within the extracted directory.
if [ ! -d "${EXTRACTED_DIR}/icons" ]; then
  error_exit "'icons' folder not found in ${EXTRACTED_DIR}."
fi

# Copy the "icons" folder into the ./lucide folder.
echo "Copying icons to ./lucide ..."
rm -rf "lucide/icons" &&
    cp -r "${EXTRACTED_DIR}/icons" "./lucide/" && 
    rm ./lucide/icons/*.json
echo "${VERSION}" > "./lucide/version.txt"

echo "Lucide icons have been successfully copied to lucide/icons."
