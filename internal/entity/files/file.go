package files

type UploadFileRequest struct {
	File      []byte
	Filename  string
	Container string
}

type UploadFileResponse struct {
	URL string `json:"url"`
}
