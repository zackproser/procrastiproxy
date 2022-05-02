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

get_release_url () {
jq_cmd=".assets[] | select(.name | endswith(\"${platform}-${arch}.tar.gz\")).browser_download_url"
# Find the latest binary release URL for this platform
url="$(curl -sL https://api.github.com/repos/zackproser/procrastiproxy/releases/latest | jq -r "${jq_cmd}")"
echo $url
}

target_url="$(get_release_url)"
curl -LO $target_url
#Rename and copy to your procrastiproxy folder
filename=$(basename $target_url)
tar xvzf "${filename}"
binaryname="procrastiproxy"
unzippeddir="${filename/.tar.gz/""}"
fullpath="${unzippeddir}/${binaryname}"
chmod +x "${fullpath}"

PROCRASTIPROXY_DIR=/usr/local/bin/procrastiproxy
sudo mv "${fullpath}" "${PROCRASTIPROXY_DIR}"
if [[ $? -eq 0 ]]; then
  echo "Successfully installed $binary_name at $PROCRASTIPROXY_DIR"
else
  echo "ERROR: could not install $binary_name at $PROCRASTIPROXY_DIR"
fi
