package models

// Topic represents a v2ex topic/post
// API: /api/topics/latest.json, /api/topics/hot.json, /api/topics/show.json
type Topic struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	URL             string    `json:"url"`
	Content         string    `json:"content"`
	ContentRendered string    `json:"content_rendered"`
	Replies         int       `json:"replies"`
	Member          *Member   `json:"member"`
	Node            *NodeInfo `json:"node"`
	Created         int64     `json:"created"`
	LastModified    int64     `json:"last_modified"`
	LastTouched     int64     `json:"last_touched"`
}

// NodeInfo is a simplified node representation within a topic
type NodeInfo struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	Title            string `json:"title"`
	TitleAlternative string `json:"title_alternative"`
	Topics           int    `json:"topics"`
	AvatarMini       string `json:"avatar_mini,omitempty"`
	AvatarNormal     string `json:"avatar_normal,omitempty"`
	AvatarLarge      string `json:"avatar_large,omitempty"`
}
