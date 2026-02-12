package models

// Node represents a v2ex node (forum category)
// API: /api/nodes/all.json, /api/nodes/show.json
type Node struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	Title            string `json:"title"`
	TitleAlternative string `json:"title_alternative"`
	Topics           int    `json:"topics"`
	Header           string `json:"header"`
	Footer           string `json:"footer"`
	Created          int64  `json:"created"`
	AvatarMini       string `json:"avatar_mini,omitempty"`
	AvatarNormal     string `json:"avatar_normal,omitempty"`
	AvatarLarge      string `json:"avatar_large,omitempty"`
}
