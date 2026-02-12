package models

// SiteInfo represents v2ex site information
// API: /api/site/info.json
type SiteInfo struct {
	Title       string `json:"title"`
	Slogan      string `json:"slogan"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
}

// SiteStats represents v2ex site statistics
// API: /api/site/stats.json
type SiteStats struct {
	TopicMax  int `json:"topic_max"`
	MemberMax int `json:"member_max"`
}
