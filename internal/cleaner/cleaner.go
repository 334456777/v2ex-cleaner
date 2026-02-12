package cleaner

import (
	"regexp"
	"strings"
	"time"

	"github.com/yusteven/v2ex-cleaner/internal/models"
)

// htmlTagRegex removes HTML tags
var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)

// multipleSpacesRegex normalizes multiple spaces
var multipleSpacesRegex = regexp.MustCompile(`\s+`)

// CleanNode cleans a node for output
func CleanNode(node models.Node) models.CleanedNode {
	return models.CleanedNode{
		ID:          node.ID,
		Name:        node.Name,
		Title:       node.Title,
		Topics:      node.Topics,
		Description: cleanText(node.Header),
		CreatedAt:   timestampToISO(node.Created),
	}
}

// CleanTopic cleans a topic for output
func CleanTopic(topic models.Topic) models.CleanedTopic {
	author := ""
	if topic.Member != nil {
		author = topic.Member.Username
	}

	nodeName := ""
	nodeTitle := ""
	if topic.Node != nil {
		nodeName = topic.Node.Name
		nodeTitle = topic.Node.Title
	}

	return models.CleanedTopic{
		ID:           topic.ID,
		Title:        cleanText(topic.Title),
		URL:          topic.URL,
		Content:      cleanHTML(topic.ContentRendered),
		Author:       author,
		NodeName:     nodeName,
		NodeTitle:    nodeTitle,
		ReplyCount:   topic.Replies,
		CreatedAt:    timestampToISO(topic.Created),
		LastModified: timestampToISO(topic.LastModified),
	}
}

// CleanReply cleans a reply for output
func CleanReply(reply models.Reply) models.CleanedReply {
	author := ""
	if reply.Member != nil {
		author = reply.Member.Username
	}

	return models.CleanedReply{
		ID:        reply.ID,
		TopicID:   reply.TopicID,
		Content:   cleanHTML(reply.ContentRendered),
		Author:    author,
		Thanks:    reply.Thanks,
		CreatedAt: timestampToISO(reply.Created),
	}
}

// CleanMember cleans a member for output
func CleanMember(member models.Member) models.CleanedMember {
	return models.CleanedMember{
		ID:        member.ID,
		Username:  member.Username,
		Tagline:   cleanText(member.Tagline),
		Bio:       cleanText(member.Bio),
		Website:   member.Website,
		Location:  member.Location,
		CreatedAt: timestampToISO(member.Created),
	}
}

// cleanText removes extra whitespace from plain text
func cleanText(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	// Replace multiple whitespace with single space
	text = multipleSpacesRegex.ReplaceAllString(text, " ")
	return text
}

// cleanHTML removes HTML tags and cleans the text
func cleanHTML(html string) string {
	if html == "" {
		return ""
	}
	// Remove HTML tags
	text := htmlTagRegex.ReplaceAllString(html, "")
	// Decode common HTML entities
	text = decodeHTMLEntities(text)
	// Clean whitespace
	return cleanText(text)
}

// decodeHTMLEntities decodes common HTML entities
func decodeHTMLEntities(text string) string {
	replacer := strings.NewReplacer(
		"&nbsp;", " ",
		"&lt;", "<",
		"&gt;", ">",
		"&amp;", "&",
		"&quot;", "\"",
		"&#39;", "'",
		"&ndash;", "–",
		"&mdash;", "—",
		"\r\n", "\n",
		"\r", "\n",
	)
	return replacer.Replace(text)
}

// timestampToISO converts Unix timestamp to ISO 8601 format
func timestampToISO(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).Format(time.RFC3339)
}

// CreateMetadata creates output metadata
func CreateMetadata() models.Metadata {
	return models.Metadata{
		FetchedAt: time.Now().Format(time.RFC3339),
		Source:    "v2ex.com",
		Version:   "1.0.0",
	}
}
