package models

// Reply represents a reply to a topic
// API: /api/replies/show.json
type Reply struct {
	ID              int     `json:"id"`
	Thanks          int     `json:"thanks"`
	Content         string  `json:"content"`
	ContentRendered string  `json:"content_rendered"`
	Member          *Member `json:"member"`
	Created         int64   `json:"created"`
	LastModified    int64   `json:"last_modified"`
	TopicID         int     `json:"topic_id,omitempty"` // Added for linking
}
