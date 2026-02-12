package models

// CleanedData represents the cleaned and structured output format
type CleanedData struct {
	Meta Metadata               `json:"meta"`
	Data map[string]interface{} `json:"data"`
}

// Metadata provides information about the fetched data
type Metadata struct {
	FetchedAt string `json:"fetched_at"`
	Source    string `json:"source"`
	Version   string `json:"version"`
}

// CleanedNode represents a cleaned node with readable fields
type CleanedNode struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Topics      int    `json:"topics"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// CleanedTopic represents a cleaned topic
type CleanedTopic struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Content      string `json:"content"`
	Author       string `json:"author"`
	NodeName     string `json:"node_name"`
	NodeTitle    string `json:"node_title"`
	ReplyCount   int    `json:"reply_count"`
	CreatedAt    string `json:"created_at"`
	LastModified string `json:"last_modified"`
}

// CleanedReply represents a cleaned reply
type CleanedReply struct {
	ID        int    `json:"id"`
	TopicID   int    `json:"topic_id"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	Thanks    int    `json:"thanks"`
	CreatedAt string `json:"created_at"`
}

// CleanedMember represents a cleaned member profile
type CleanedMember struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Tagline   string `json:"tagline,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Website   string `json:"website,omitempty"`
	Location  string `json:"location,omitempty"`
	CreatedAt string `json:"created_at"`
}
