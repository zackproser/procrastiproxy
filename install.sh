#!/usr/bin/env bash
set -x

# Determine architecture
if [[ $(uname -s) == Darwin && $(uname -m) == amd64  ]]
then
	platform='Darwin_amd64'
elif [[ $(uname -s) == Darwin && $(uname -m) == arm64  ]]
then
	platform='Darwin_arm64'
elif [[ $(uname -s) == Linux ]]
then
	platform='Linux_amd64'
else
	echo "No supported architecture found"
	exit 1
fi

jq_cmd=".assets[] | select(.name | endswith(\"${platform}.tar.gz\")).browser_download_url"

# Find the latest binary release URL for this platform
url="$(curl -s https://api.github.com/repos/zackproser/procrastiproxy/releases/latest | jq -r "${jq_cmd}")"

# Download the tarball
curl -OL ${url}
# Rename and copy to your procrastiproxy folder
filename=$(basename $url)
tar xvzf ${filename}
filename="procrastiproxy"
chmod +x ${filename}

PROCRASTIPROXY_DIR=~/.procrastiproxy/$platform
mkdir -p $PROCRASTIPROXY_DIR
mv $filename ${PROCRASTIPROXY_DIR}/${filename%_${platform}}
echo ""
echo "Successfully installed procrastiproxy at: " ${PROCRASTIPROXY_DIR}
