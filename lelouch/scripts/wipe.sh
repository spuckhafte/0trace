#!/bin/bash
# Usage: sudo ./wipe.sh /dev/sdX

DEVICE="$1"

if [ -z "$DEVICE" ]; then
    echo "Usage: $0 /dev/sdX"
    exit 1
fi

START_TIME=$(date -Iseconds)
STATUS="failure"

# Method flags
nvme_sanitize=0
nvme_format=0
hdparm_secure_erase=0
hdparm_dco_restore=0
blkdiscard_m=0
shred_m=0
dd_m=0

# Detect if NVMe or SATA/ATA
if [[ "$DEVICE" == *"nvme"* ]]; then
    # NVMe device
    if nvme sanitize "$DEVICE" -a 1 >/dev/null 2>&1; then
        nvme_sanitize=1
        STATUS="success"
    elif nvme format "$DEVICE" >/dev/null 2>&1; then
        nvme_format=1
        STATUS="success"
    fi
else
    # SATA / HDD or USB-SATA
    HDPARM_OUT=$(hdparm -N "$DEVICE" 2>/dev/null)
    MAX=$(echo "$HDPARM_OUT" | awk -F'[ =,/]+' '/max sectors/ {print $4}')
    NATIVE=$(echo "$HDPARM_OUT" | awk -F'[ =,/]+' '/max sectors/ {print $5}')

    # sanity check: must be real numbers > 100000
    if [[ "$MAX" =~ ^[0-9]+$ && "$NATIVE" =~ ^[0-9]+$ && "$MAX" -gt 100000 && "$NATIVE" -gt 100000 ]]; then
        # Looks like a real SATA disk, try hdparm secure erase
        if hdparm --user-master u --security-set-pass NULL "$DEVICE" >/dev/null 2>&1 &&
           hdparm --user-master u --security-erase NULL "$DEVICE" >/dev/null 2>&1; then
            hdparm_secure_erase=1
            STATUS="success"

            # Restore HPA if possible
            if hdparm -N p"$NATIVE" "$DEVICE" >/dev/null 2>&1; then
                hdparm_dco_restore=1
            fi
        fi
    fi

    # If still failure, try blkdiscard
    if [ "$STATUS" = "failure" ]; then
        if blkdiscard "$DEVICE" >/dev/null 2>&1; then
            blkdiscard_m=1
            STATUS="success"
        fi
    fi

    # If blkdiscard failed, try shred
    if [ "$STATUS" = "failure" ]; then
        if sudo shred -n 1 -vz "$DEVICE"; then
            shred_m=1
            STATUS="success"
        fi
    fi

    # If shred failed, try dd
    if [ "$STATUS" = "failure" ]; then
        if dd if=/dev/zero of="$DEVICE" bs=1M status=progress conv=fdatasync >/dev/null 2>&1; then
            dd_m=1
            STATUS="success"
        fi
    fi
fi

END_TIME=$(date -Iseconds)

# Generate JSON log
LOGFILE="wipe_log.json"
cat > "$LOGFILE" <<EOF
{
  "device": "$DEVICE",
  "start_time": "$START_TIME",
  "end_time": "$END_TIME",
  "status": "$STATUS",
  "method": {
    "nvme_sanitize": $nvme_sanitize,
    "nvme_format": $nvme_format,
    "hdparm_secure_erase": $hdparm_secure_erase,
    "hdparm_dco_restore": $hdparm_dco_restore,
    "blkdiscard": $blkdiscard_m,
    "shred": $shred_m,
    "dd": $dd_m
  }
}
EOF

echo "Wipe complete. Log written to $LOGFILE"

