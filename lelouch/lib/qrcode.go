package lib

import (
	"fmt"
	"github.com/skip2/go-qrcode"
)

// GenerateQRCodeASCII generates a QR code for the given text and returns it as ASCII art
func GenerateQRCodeASCII(text string) (string, error) {
	// Generate QR code
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("failed to create QR code: %v", err)
	}
	
	// Convert to ASCII string representation
	asciiQR := qr.ToSmallString(false)
	return asciiQR, nil
}

// GenerateQRCode generates a QR code for the given text and saves it as a PNG file
func GenerateQRCode(text, filename string) error {
	// Generate QR code with medium error correction level
	err := qrcode.WriteFile(text, qrcode.Medium, 256, filename)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %v", err)
	}
	return nil
}

// GenerateQRCodeBytes generates a QR code for the given text and returns the PNG bytes
func GenerateQRCodeBytes(text string) ([]byte, error) {
	// Generate QR code with medium error correction level
	png, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %v", err)
	}
	return png, nil
}
