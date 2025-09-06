
# 🔐 Secure Wipe Certificate Flow (Hackathon Final)

## 1. Boot & Setup
- User boots system from **custom ISO/USB**.
- Tool loads in a minimal Linux environment (no OS dependency).

---

## 2. Device Detection
- Detect connected drives (HDD/SSD/NVMe).
- Fetch details: **manufacturer, model, serial, capacity**.

---

## 3. Wiping Method (Our Prototype)
- **AES Overwrite**: Overwrite all user-accessible sectors with AES-random data.
- **Key Destruction**: Session AES key used for overwrite is securely deleted from memory.
- Result:
  - User data replaced with random ciphertext.
  - Deleted key ensures overwrite patterns cannot be reversed.

⚠️ Note: By NIST SP 800-88, this is **Clear** (not full Purge), but simulates cryptographic erase.  
Future: With device-native Secure Erase / hardware CE → full **Purge**.

---

## 4. Logging
- Record wiping session:
  - Device details (model, serial, capacity)
  - Method: *AES Overwrite + Key Destruction*
  - Start & End timestamps
  - Verification results (sampled entropy check)

---

## 5. Certificate Generation
- JSON Certificate:
  ```json
  {
    "device": "Samsung SSD 860",
    "serial": "S3Z9NB0K123456",
    "method": "AES-256 Overwrite + Key Destruction",
    "nist_level": "Clear (Prototype) → Roadmap to Purge",
    "start_time": "2025-09-07T10:00:00Z",
    "end_time": "2025-09-07T10:20:30Z",
    "status": "Success"
  }
  ```

- Convert to PDF for human readability:
  - Device info, wipe method, timestamps
  - QR code of JSON hash
  - Digital signature summary

---

## 6. Hashing & Digital Signature
- Compute **SHA-256 hash** of certificate JSON.  
- Sign hash using **private RSA/ECDSA key**.  
- Append hash + signature to certificate.

Final JSON:
```json
{
  "device": "Samsung SSD 860",
  "serial": "S3Z9NB0K123456",
  "method": "AES-256 Overwrite + Key Destruction",
  "nist_level": "Clear",
  "status": "Success",
  "hash": "4f3ac89f...",
  "signature": "MEUCIQC9..."
}
```

---

## 7. Verification (Third-Party)
- Third party loads certificate + public key.  
- Verifier recomputes hash and validates signature.  
- If match → ✅ VALID.  
- If altered → ❌ INVALID.

---

## 8. Roadmap for Full Purge
- Add ATA Secure Erase / NVMe Sanitize commands.  
- Add **true Cryptographic Erase** for self-encrypting drives.  
- Support **degaussing / destroy** for damaged legacy drives.  

---

## 🎯 Hackathon Pitch Angle
- **Now:** AES overwrite + key destruction = secure, user-friendly, NIST Clear.  
- **Future:** Add device-native + crypto erase → full NIST Purge compliance.  
- **Impact:** Tamper-proof certificates → trust for recyclers, enterprises, and regulators.  
