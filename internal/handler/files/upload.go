package files

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/eifzed/joona/internal/constant"
	"github.com/eifzed/joona/internal/entity/files"
	"github.com/eifzed/joona/lib/common"
	"github.com/eifzed/joona/lib/common/commonerr"
	"github.com/pkg/errors"
)

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	param, err := h.getAndValidateParams(r)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	resp, err := h.FileUC.UploadFile(ctx, param)
	if err != nil {
		commonwriterRespondError(ctx, w, err)
		return
	}
	commonwriterRespondOKWithData(ctx, w, resp)

}

func (h *FileHandler) getAndValidateParams(r *http.Request) (param files.UploadFileRequest, err error) {
	container := r.FormValue("container")
	if !h.checkIsValidContainer(container) {
		err = commonerr.ErrorBadRequest("upload_file", "invalid container")

	}

	maxFileSize := h.Config.File.MaxImageUploadSizeByte

	err = parseFileFromRequest(r, maxFileSize)
	if err != nil {
		err = commonerr.ErrorBadRequest("upload_file", "invalid file")
		return
	}

	file, fileInfo, err := getFileFromRequest(r, "file")
	if err != nil {
		err = commonerr.ErrorBadRequest("upload_file", "error parsing file")
		return
	}
	if fileInfo.Size > int64(maxFileSize) {
		err = commonerr.ErrorBadRequest("upload_file", fmt.Sprintf("file size exceeds maximum limit (%.1fMB)", float64(maxFileSize)/(1024*1024)))
		return
	}

	mimeType, err := getFileMimeType(file)
	if err != nil {
		return
	}

	isValid, err := h.checkIsValidMimeType(mimeType)
	if err != nil {
		err = commonerr.ErrorBadRequest("upload_file", "error parsing file mime type")
		return
	}
	if !isValid {
		err = commonerr.ErrorBadRequest("upload_file", "invalid mime type")
		return
	}
	fileByte, err := io.ReadAll(file)
	if err != nil {
		err = commonerr.ErrorBadRequest("upload_file", "error parsing file")
		return
	}
	param = files.UploadFileRequest{
		File:      fileByte,
		Container: container,
		Filename:  fmt.Sprintf("%s.%s", common.GenerateUUIDV7(), constant.MapMimeToExtension[mimeType]),
	}

	return
}

var parseFileFromRequest = func(r *http.Request, maxMemory int64) error {
	return r.ParseMultipartForm(maxMemory)
}

var getFileFromRequest = func(r *http.Request, fieldName string) (multipart.File, *multipart.FileHeader, error) {
	return r.FormFile(fieldName)
}

var getFileMimeType = func(file multipart.File) (mimeType string, err error) {
	// Determine the MIME type of the file
	fileHeader := make([]byte, 512)
	_, err = file.Read(fileHeader)
	if err != nil {
		err = errors.New("unable to read file header")
		return
	}

	// Reset the file read offset
	_, err = file.Seek(0, 0)
	if err != nil {
		err = errors.New("unable to reset file read offset")
		return
	}
	mimeType = http.DetectContentType(fileHeader)
	return
}

func (h *FileHandler) checkIsValidMimeType(mimeType string) (isValid bool, err error) {
	for _, validMimeType := range h.Config.File.MimeTypeWhitelist {
		if validMimeType == mimeType {
			isValid = true
			return
		}
	}
	isValid = false
	return
}

func (h *FileHandler) checkIsValidContainer(container string) (isValid bool) {
	for _, validContainer := range h.Config.File.ContainerWhitelist {
		if validContainer == container {
			isValid = true
			return
		}
	}
	isValid = false
	return
}
