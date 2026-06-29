package utils

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"mime/multipart"
	"os"
	"time"
)

var cld *cloudinary.Cloudinary

func InitCloudinary() {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return
	}

	var err error
	cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return
	}
}

func UploadToCloudinary(file *multipart.FileHeader) (string, error) {
	if cld == nil {
		return "", nil
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID: "events/" + time.Now().Format("20060102") + "/" + file.Filename,
		Folder:   "event-management",
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
