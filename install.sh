#!/usr/bin/env bash
set -e
binary_name="procrastiproxy"
# Determine architecture
platform=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m | tr '[:upper:]' '[:lower:]')

# We direct goreleaser to convert "amd64" to "x86_64", so perform the same rename 
# here if necessary
if [[ $arch == "amd64" ]]
then 
  arch="x86_64"
fi

# Convenience method for deleting the downloaded tarball and expanded folder after installation
cleanupInstallationDir () {
  echo "Cleaning up installation downloads..."
  rm -rf $1 $2 
}

# Fetch the latest release of procrastiproxy that is compatible with the machine this script is being executed on 
get_release_url () {
jq_cmd=".assets[] | select(.name | endswith(\"${platform}-${arch}.tar.gz\")).browser_download_url"
# Find the latest binary release URL for this platform
url="$(curl -sL https://api.github.com/repos/zackproser/procrastiproxy/releases/latest | jq -r "${jq_cmd}")"
echo $url
}

target_url="$(get_release_url)"
curl -LO $target_url
# Rename and copy to your procrastiproxy folder
filename=$(basename $target_url)
# Unpack the tarball to your pwd
tar xvzf "${filename}"
binaryname="procrastiproxy"
unzippeddir="${filename/.tar.gz/""}"
fullpath="${unzippeddir}/${binaryname}"
# Make binary executable
chmod +x "${fullpath}"

# Move the binary into the target directory
PROCRASTIPROXY_DIR=/usr/local/bin/procrastiproxy
sudo mv "${fullpath}" "${PROCRASTIPROXY_DIR}"
if [[ $? -eq 0 ]]; then
  echo "Successfully installed $binary_name at $PROCRASTIPROXY_DIR"
  # Delete the downloaded tarball and the expanded folder now that installation is complete
  cleanupInstallationDir "$unzippeddir" "$filename"
else
  echo "ERROR: could not install $binary_name at $PROCRASTIPROXY_DIR"
fi
