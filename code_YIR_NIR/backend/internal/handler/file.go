package handler

import (
	"io"
	"mime/multipart"
	"net/http"
)

const maxUploadSize = 10 << 20 // 10 MiB

func readUploadFile(w http.ResponseWriter, r *http.Request) ([]byte, *multipart.FileHeader, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		return nil, nil, err
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	return data, header, nil
}
