package email

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Attachment struct {
	Filename    string
	ContentType string
	Size        int
	S3Url       string
	Skipped     bool
}

func SaveAttachment(attachment *bytes.Buffer, filename string) (string, error) {
	// Create a unique filename
	hash := sha256.New()
	if _, err := io.Copy(hash, attachment); err != nil {
		return "", err
	}

	reader := bytes.NewReader(attachment.Bytes())
	if _, err := reader.Seek(0, 0); err != nil {
		return "", err
	}


	newFilename := fmt.Sprintf("%s-%s", hex.EncodeToString(hash.Sum(nil)), filename)

	// Create the uploads directory if it doesn't exist
	uploadsDir := "uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		if err := os.Mkdir(uploadsDir, 0755); err != nil {
			return "", err
		}
	}

	// Save the file
	filePath := filepath.Join(uploadsDir, newFilename)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", err
	}

	return filePath, nil
}
