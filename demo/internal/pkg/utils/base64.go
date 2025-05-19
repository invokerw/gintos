package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
)

func DecodeImageDataURI(dataURI string) (data []byte, ext string, err error) {
	var mimeType string
	mimeType, data, err = DecodeDataURI(dataURI)
	if err != nil {
		return nil, "", err
	}
	switch mimeType {
	case "image/png":
		ext = "png"
	case "image/jpeg":
		ext = "jpg"
	case "image/gif":
		ext = "gif"
	case "image/webp":
		ext = "webp"
	case "image/bmp":
		ext = "bmp"
	case "image/svg+xml":
		ext = "svg"
	default:
		return nil, "", errors.New("unsupported image type")
	}
	return data, ext, nil
}

func DecodeDataURI(dataURI string) (mimeType string, data []byte, err error) {
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 {
		return "", nil, errors.New("invalid Data URI")
	}

	meta, dataStr := parts[0], parts[1]
	meta = strings.TrimPrefix(meta, "data:")
	isBase64 := strings.Contains(meta, ";base64")

	mimeType = strings.Split(meta, ";")[0]
	if mimeType == "" {
		mimeType = "text/plain"
	}

	if isBase64 {
		data, err = base64.StdEncoding.DecodeString(dataStr)
	} else {
		var decodedStr string
		decodedStr, err = url.QueryUnescape(dataStr)
		data = []byte(decodedStr)
	}
	if err != nil {
		return "", nil, err
	}

	// 检测实际类型
	if detected := detectMimeType(data); detected != "" {
		mimeType = detected
	}

	return mimeType, data, nil
}

func detectMimeType(data []byte) string {
	if len(data) < 8 {
		return ""
	}
	if bytes.HasPrefix(data, []byte("\x89PNG\x0D\x0A\x1A\x0A")) {
		return "image/png"
	}
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8}) {
		return "image/jpeg"
	}
	if bytes.HasPrefix(data, []byte("GIF8")) {
		return "image/gif"
	}
	return ""
}
