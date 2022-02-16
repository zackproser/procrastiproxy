#!/usr/bin/env bash
set -e
binary_name="procrastiproxy"
# Determine architecture
if [[ $(uname -s) == Darwin && $(uname -m) == amd64  ]]
then
	platform='Darwin_amd64'
elif [[ $(uname -s) == Darwin && $(uname -m) == arm64  ]]
then
	platform='darwin_arm64'
elif [[ $(uname -s) == Linux ]]
then
	platform='linux_amd64'
else
	echo "No supported architecture found"
	exit 1
fi

get_release_url () {
jq_cmd=".assets[] | select(.name | endswith(\"${platform}.tar.gz\")).browser_download_url"
# Find the latest binary release URL for this platform
url="$(curl -sL https://api.github.com/repos/zackproser/procrastiproxy/releases/latest | jq -r "${jq_cmd}")"
echo $url
}

target_url="$(get_release_url)"
curl -LO $target_url
#Rename and copy to your procrastiproxy folder
filename=$(basename $target_url)
tar xvzf ${filename}
filename="procrastiproxy"
chmod +x ${filename}

PROCRASTIPROXY_DIR=/usr/local/bin/procrastiproxy
sudo mv $filename ${PROCRASTIPROXY_DIR}
if [[ $? -eq 0 ]]; then
  echo "Successfully installed $binary_name at $PROCRASTIPROXY_DIR"
else
  echo "ERROR: could not install $binary_name at $PROCRASTIPROXY_DIR"
fi
