package gointrum

type WorkerFilterResponse struct {
	*Response
	Data map[string]*WorkerFilterData `json:"data"`
}

type WorkerFilterData struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	DivisionID  string `json:"division_id"`
	SubofficeID string `json:"suboffice_id"`
	Post        string `json:"post"`
	Boss        string `json:"boss"`
	Status      string `json:"status"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	// Secondname          string           `json:"secondname"`
	// Internalemail       []string         `json:"internalemail"`
	// Externalemail       []interface{}    `json:"externalemail"`
	// Internalphone       []interface{}    `json:"internalphone"`
	// Externalphone       []interface{}    `json:"externalphone"`
	// Mobilephone         []Mobilephone    `json:"mobilephone"`
	// Birthday            *string          `json:"birthday"`
	// Address             interface{}      `json:"address"`
	// About               string           `json:"about"`
	// Hobby               string           `json:"hobby"`
	// CreatedAt           *time.Time       `json:"created_at"`
	// Skype               string           `json:"skype"`
	// Facebook            *string          `json:"facebook"`
	// Vkontakte           string           `json:"vkontakte"`
	// Gender              string           `json:"gender"`
	// Fields              map[string]Field `json:"fields"`
	// GroupID             []string         `json:"group_id"`
	// Avatars             Avatars          `json:"avatars"`
	// AsteriskShortNumber []string         `json:"asterisk_short_number,omitempty"`
}
