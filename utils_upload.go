package gointrum

import (
	"context"
	"fmt"
)

// Ссылка на метод: 	http://domainname.intrumnet.com:81/sharedapi/utils/upload
type UtilsUploadParams struct {
	ObjectType  string       // один из возможных вариантов: stock (объекты, продукты), applications (заявки), purchaser (контакты), sales (продажи)
	UploadFiles []UploadFile // поле и массив путей к загружаемым объектам
}

// Ссылка на метод: 	http://domainname.intrumnet.com:81/sharedapi/utils/upload
func UtilsUpload(ctx context.Context, subdomain, apiKey string, inputParams *UtilsUploadParams) (*UtilsUploadResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/utils/upload", subdomain)

	if inputParams == nil {
		return nil, fmt.Errorf("input params is nil")
	}
	if inputParams.ObjectType == "" {
		return nil, fmt.Errorf("params[object] is required")
	}
	if len(inputParams.UploadFiles) == 0 {
		return nil, fmt.Errorf("at least one file is required")
	}

	params := map[string]string{
		"params[object]": inputParams.ObjectType,
	}

	files := make([]UploadFile, 0, len(inputParams.UploadFiles))
	for i, f := range inputParams.UploadFiles {
		if f.FileName == "" {
			return nil, fmt.Errorf("upload file[%d]: FileName is required", i)
		}
		if f.R == nil {
			return nil, fmt.Errorf("upload file[%d]: R is nil", i)
		}
		if f.FieldName == "" {
			f.FieldName = fmt.Sprintf("upload[%d]", i)
		}
		files = append(files, f)
	}

	var resp UtilsUploadResponse
	if err := requestUploadFile(ctx, apiKey, methodURL, params, files, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
