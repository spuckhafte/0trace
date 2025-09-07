#!/bin/bash
# Usage: sudo ./wipe_disk.sh /dev/sdX

DEVICE="$1"

if [ -z "$DEVICE" ]; then
  echo "Usage: $0 /dev/sdX"
  exit 1
fi

START_TIME=$(date -Iseconds)
STATUS="success"
METHOD=""

# Detect if NVMe or SATA/ATA
if [[ "$DEVICE" == *"nvme"* ]]; then
  # NVMe device
  METHOD="nvme sanitize"
  nvme sanitize "$DEVICE" -a 1 >/dev/null 2>&1
  if [ $? -ne 0 ]; then
    METHOD="nvme format"
    nvme format "$DEVICE" >/dev/null 2>&1 || STATUS="failure"
  fi
else
  # SATA / HDD
  # First try ATA Secure Erase
  METHOD="hdparm secure erase"
  hdparm --user-master u --security-set-pass NULL "$DEVICE" >/dev/null 2>&1
  hdparm --user-master u --security-erase NULL "$DEVICE" >/dev/null 2>&1 || STATUS="failure"
fi

# Fallback to shred if secure erase fails
if [ "$STATUS" = "failure" ]; then
  METHOD="shred overwrite (fallback)"
  shred -n 1 -vz "$DEVICE" >/dev/null 2>&1 || STATUS="failure"
fi

END_TIME=$(date -Iseconds)

# Generate JSON log
LOGFILE="wipe_log.json"
cat >"$LOGFILE" <<EOF
{
  "device": "$DEVICE",
  "method": "$METHOD",
  "start_time": "$START_TIME",
  "end_time": "$END_TIME",
  "status": "$STATUS"
}
EOF

echo "Wipe complete. Log written to $LOGFILE"
