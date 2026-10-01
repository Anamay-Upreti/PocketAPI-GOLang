package utils

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const MaxAvatarSize int64 = 2 * 1024 * 1024 

var AllowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func ValidateImage(filename string) error {

	extension := strings.ToLower(filepath.Ext(filename))

	if !AllowedImageExtensions[extension] {
		return fmt.Errorf("only jpg, jpeg and png images are allowed")
	}

	return nil
}

func GenerateAvatarFilename(filename string) string {

	extension := strings.ToLower(filepath.Ext(filename))

	return fmt.Sprintf(
		"avatar-%s%s",
		uuid.New().String(),
		extension,
	)
}