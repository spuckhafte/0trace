# Secure Data Wiping with Digital Certificate (NIST SP 800-88 Purge)

## 🚀 End-to-End Flow

### 1. Boot Environment
- User plugs in **bootable ISO/USB**.
- Device boots into a **custom Linux-based environment** (external to OS).

### 2. Device Detection
- System scans and lists connected drives (HDDs, SSDs, phones).
- Fetch details: model, serial, capacity, interface.

### 3. Data Wiping (Purge Method)
- Execute **ATA Secure Erase** (for HDDs/SSDs).
- Execute **NVMe Format/Secure Erase** (for NVMe SSDs).
- If unsupported, fall back to **AES-random overwrite**.
- Target: **All user + hidden areas (HPA/DCO, remapped sectors)**.

Result: Data is **cryptographically purged** → recovery not feasible.

### 4. Logging
- System records:
  - Device details (model, serial, capacity)
  - Wiping method used (e.g., *NIST SP 800-88 Purge*)
  - Start time, end time
  - Wipe status (success/failure)
  - Logs of commands executed

### 5. Certificate Generation
- Build a JSON structure:
  ```json
  {
    "device": "Samsung SSD 860",
    "serial": "S4YBNX0M12345",
    "method": "NIST SP 800-88 Purge (Secure Erase + AES overwrite)",
    "start_time": "2025-09-06T10:00:00Z",
    "end_time": "2025-09-06T10:15:30Z",
    "status": "Success"
  }
  ```

### 6. Hashing & Signing
- Apply **SHA-256 hash** on the JSON certificate.
- Sign the hash with the system’s **private RSA/ECDSA key**.
- Attach hash + signature to certificate.

Final JSON:
```json
{
  "device": "Samsung SSD 860",
  "serial": "S4YBNX0M12345",
  "method": "NIST SP 800-88 Purge (Secure Erase + AES overwrite)",
  "start_time": "2025-09-06T10:00:00Z",
  "end_time": "2025-09-06T10:15:30Z",
  "status": "Success",
  "hash": "4f3ac89f...",
  "signature": "b0a2e98d..."
}
```

### 7. Export Certificate
- Save certificate in:
  - **JSON** (for machine verification)
  - **PDF** (for human readability)

### 8. Verification (Third Party)
- Third party obtains:
  - Certificate JSON
  - Public key of issuing system
- Runs verification:
  - If signature is valid → Certificate authentic & untampered
  - If invalid → Possible tampering detected

---

## ✅ Benefits
- **Tamper-proof**: Cryptographically signed certificate.
- **Standards-compliant**: NIST SP 800-88 Purge method.
- **Cross-platform**: Works on HDDs, SSDs, NVMe.
- **Trustworthy**: Independent verification possible.

---

## 🔑 Short Pitch
> Our tool boots externally, securely purges drives using ATA/NVMe Secure Erase and AES overwrite, then generates a digitally signed wipe certificate. The certificate is tamper-proof and verifiable, ensuring trust and compliance in IT asset recycling.
