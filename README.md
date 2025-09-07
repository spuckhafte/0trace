# Secure Data Wiping Bootable ISO - Flow Documentation

## Overview
This project implements a **bootable ISO solution** for secure data wiping, compliant with **NIST SP 800-88 Purge standards**. It ensures user trust in IT asset recycling by performing tamper-proof, verifiable erasure and issuing certificates.

---

## Process Flow

1. **Boot**  
   - User inserts USB with Nix-built ISO and boots the target machine.  
   - ISO loads a minimal Linux environment with Go + Bash utilities preconfigured.

2. **Disk Detection & Selection**  
   - Go-based TUI scans connected drives (HDD, SSD, NVMe).  
   - User selects the target disk to erase.

3. **Wipe Execution**  
   - If supported: NVMe Sanitize (Crypto/Purge) or ATA Secure Erase (hdparm).  
   - If not supported: Bash script generates random AES-256 key (stored only in RAM), overwrites disk, and discards key after wipe.

4. **Completion Logging**  
   - System collects metadata:
     - Device details (model, serial, capacity).
     - Wipe method used (sanitize, overwrite, shred fallback).
     - Start + end timestamps.
     - Random session ID.

5. **Certificate Generation**  
   - Go application creates tamper-proof certificate in **JSON + PDF**.  
   - Contents include all metadata + SHA-256 hash.  
   - Digitally signed with project’s private key.  
   - QR code embedded for easy third-party verification.

6. **Verification**  
   - Independent parties use public key to verify authenticity of certificate.  
   - JSON/PDF hashes prove that logs were not tampered with.

---

## Technologies Used
- **Go** → TUI + certificate generation/verification.  
- **Bash** → Disk operations, AES key management.  
- **Nix** → Reproducible ISO image build.  
- **NVMe sanitize / hdparm / shred** → Device-native wipe operations.  
- **qrcode** → Embedded QR in certificate.  

---

## Deliverables
- Bootable ISO with full wiping + certification process.  
- JSON + PDF wipe certificates.  
- Go-based verification tool for auditors, recyclers, and users.

---

## Limitations (NIST Purge Scope)
- Wipe cannot guarantee results on physically damaged or non-responsive drives.  
- Degaussing / destruction not included (requires hardware tools).  
- Relies on drive compliance with sanitize/secure erase commands; otherwise AES overwrite fallback is used.

---

✅ End-to-end **flow-based documentation** for hackathon presentation.

