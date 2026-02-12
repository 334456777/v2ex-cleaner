package models

// Member represents a v2ex user
// API: /api/members/show.json
type Member struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	URL          string `json:"url,omitempty"`
	Website      string `json:"website,omitempty"`
	Twitter      string `json:"twitter,omitempty"`
	Location     string `json:"location,omitempty"`
	Tagline      string `json:"tagline"`
	Bio          string `json:"bio,omitempty"`
	AvatarMini   string `json:"avatar_mini"`
	AvatarNormal string `json:"avatar_normal"`
	AvatarLarge  string `json:"avatar_large"`
	Created      int64  `json:"created"`
	Status       string `json:"status,omitempty"`
}
