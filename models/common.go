package models

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusDeleted  Status = "deleted"
)

func (s Status) String() string {
	return string(s)
}

type SortDirection string

const (
	SortDirectionASC  SortDirection = "ASC"
	SortDirectionDESC SortDirection = "DESC"
)

func (s SortDirection) String() string {
	return string(s)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}
