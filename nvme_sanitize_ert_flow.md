
# NVMe Sanitize + Evidence Report of Erasure (ERT) Flow

This document explains how the **Evidence Report of Erasure (ERT)** is generated when using **NVMe Sanize**.

---

## 1. Boot Environment
- User boots from the **Secure Wipe ISO/USB**.
- Minimal Linux + `nvme-cli` environment is available.

---

## 2. Device Metadata Collection
Before wiping, collect essential device details:

Command:
```bash
nvme list
nvme id-ctrl /dev/nvme0
```

Fields extracted:
- **Model number** → `mn`
- **Serial number** → `sn`
- **Firmware version** → `fr`
- **Capacity** → from `nvme list`

This ensures the certificate is linked to the specific device.

---

## 3. Sanitize Capabilities
Command:
```bash
nvme id-ctrl /dev/nvme0 --vendor-specific
```

Field:
- **SANICAP** → tells which sanitize methods are supported (crypto erase, block erase, overwrite).

This is logged in the JSON.

---

## 4. Start Time Logging
- Record system time before starting sanitize: `start_time = datetime.now()`.

---

## 5. Execute NVMe Sanitize
Command:
```bash
nvme sanitize /dev/nvme0 --sanitize=block-erase --ause
```
Options:
- `block-erase` → Erases NAND blocks (Purge-level).
- `crypto-erase` → Deletes encryption keys (if supported, very fast).
- `overwrite` → Writes a pattern to all blocks.

---

## 6. Monitor Progress
Poll status while sanitize is running:
```bash
nvme get-log /dev/nvme0 sanitize
```

Fields to check:
- `Sanitize Status` → In progress, complete, or failed.
- `Sanitize Progress %`.

---

## 7. End Time Logging
- Record system time after sanitize completes: `end_time = datetime.now()`.

---

## 8. Verification (Optional)
- Read random sectors with `nvme read`.
- Ensure sectors are cleared or show the sanitize pattern.
- Mark `"verification": "PASS"` if consistent.

---

## 9. Build JSON ERT
Example structure:
```json
{
  "device": {
    "model": "Samsung PM981a",
    "serial": "S4GNNE0R123456",
    "firmware": "EXA7301Q",
    "capacity": "1TB"
  },
  "sanitize": {
    "method": "NVMe Block Erase",
    "capability": "SANICAP: Crypto+Block",
    "status": "SUCCESS",
    "start_time": "2025-09-06T10:00:00Z",
    "end_time": "2025-09-06T10:15:00Z",
    "controller_status": "Sanitize Completed"
  },
  "verification": {
    "sample_check": "PASS",
    "progress": "100%"
  },
  "certificate": {
    "issued_by": "SecureWipe Boot ISO v1.0",
    "key_fingerprint": "a7f4...e11c"
  }
}
```

---

## 10. Digital Signing
- Sign JSON with app’s private key (RSA-2048 or ECDSA).
- Signature block is attached:

```json
{
  "ERT": { ... },
  "signature": "MEQCIG82...."
}
```

- Public key is distributed for third-party verification.

---

## 11. Generate Human-Readable PDF
- Convert JSON into a formatted PDF certificate.
- Include:
  - Device metadata
  - Wipe method used
  - Start/end time
  - Controller-reported success
  - Digital signature

---

## 12. Deliver Certificate
- Save `ERT.json` + `ERT.pdf` to external USB or cloud.
- User/recycler now has **tamper-proof proof of erasure**.

---

## 📌 Hackathon Pitch Point
- **Now:** Working prototype with NVMe sanitize + JSON/PDF signed report.
- **Future:** Extend to SATA Secure Erase, Cryptographic Erase, AES overwrite fallback for universality.
