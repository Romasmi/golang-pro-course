#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet

(
  cat <<EOF
Hello
From
NC
EOF
  sleep 5
) | nc -l localhost 4242 >/tmp/nc.out &
NC_PID=$!

sleep 1
(
  cat <<EOF
I
am
TELNET client
EOF
  sleep 5
) | ./go-telnet --timeout=5s localhost 4242 >/tmp/telnet.out &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  if [ "${fileData}" != "${2}" ]; then
    printf "unexpected output, %s:\n%s\n" "$1" "${fileData}"
    exit 1
  fi
}

expected_nc_out='I
am
TELNET client'
fileEquals /tmp/nc.out "${expected_nc_out}"

expected_telnet_out='Hello
From
NC'
fileEquals /tmp/telnet.out "${expected_telnet_out}"

rm -f go-telnet
echo "PASS"
