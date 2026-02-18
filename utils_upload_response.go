package gointrum

import (
	"encoding/json"
	"fmt"
)

type UtilsUploadResponse struct {
	*Response
	Data UploadData `json:"data"`
}

type UploadData struct {
	Names []string
}

func (u *UploadData) UnmarshalJSON(b []byte) error {
	// 1) {"name":"..."}
	var single struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(b, &single); err == nil && single.Name != "" {
		u.Names = []string{single.Name}
		return nil
	}

	// 2) {"name":[...]}
	var multi struct {
		Name []string `json:"name"`
	}
	if err := json.Unmarshal(b, &multi); err == nil && len(multi.Name) > 0 {
		u.Names = multi.Name
		return nil
	}

	// 3) ["..."] (на всякий)
	var arr []string
	if err := json.Unmarshal(b, &arr); err == nil && len(arr) > 0 {
		u.Names = arr
		return nil
	}

	// 4) fallback
	return fmt.Errorf("unknown upload data format: %s", string(b))
}
