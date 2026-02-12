package api

import (
	"fmt"
	"strconv"

	"github.com/yusteven/v2ex-cleaner/internal/models"
)

// GetSiteInfo fetches site information
// API: /api/site/info.json
func (c *Client) GetSiteInfo() (*models.SiteInfo, error) {
	data, err := c.do("/site/info.json")
	if err != nil {
		return nil, err
	}

	var info models.SiteInfo
	if err := c.unmarshal(data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetSiteStats fetches site statistics
// API: /api/site/stats.json
func (c *Client) GetSiteStats() (*models.SiteStats, error) {
	data, err := c.do("/site/stats.json")
	if err != nil {
		return nil, err
	}

	var stats models.SiteStats
	if err := c.unmarshal(data, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetAllNodes fetches all nodes
// API: /api/nodes/all.json
func (c *Client) GetAllNodes() ([]models.Node, error) {
	data, err := c.do("/nodes/all.json")
	if err != nil {
		return nil, err
	}

	var nodes []models.Node
	if err := c.unmarshal(data, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

// GetNode fetches a specific node by id or name
// API: /api/nodes/show.json?id=X or /api/nodes/show.json?name=X
func (c *Client) GetNode(id int, name string) (*models.Node, error) {
	endpoint := "/nodes/show.json?"
	if id > 0 {
		endpoint += "id=" + strconv.Itoa(id)
	} else if name != "" {
		endpoint += "name=" + name
	} else {
		return nil, fmt.Errorf("either id or name must be provided")
	}

	data, err := c.do(endpoint)
	if err != nil {
		return nil, err
	}

	var node models.Node
	if err := c.unmarshal(data, &node); err != nil {
		return nil, err
	}
	return &node, nil
}

// GetTopicsLatest fetches latest topics
// API: /api/topics/latest.json
func (c *Client) GetTopicsLatest() ([]models.Topic, error) {
	data, err := c.do("/topics/latest.json")
	if err != nil {
		return nil, err
	}

	var topics []models.Topic
	if err := c.unmarshal(data, &topics); err != nil {
		return nil, err
	}
	return topics, nil
}

// GetTopicsHot fetches hot topics
// API: /api/topics/hot.json
func (c *Client) GetTopicsHot() ([]models.Topic, error) {
	data, err := c.do("/topics/hot.json")
	if err != nil {
		return nil, err
	}

	var topics []models.Topic
	if err := c.unmarshal(data, &topics); err != nil {
		return nil, err
	}
	return topics, nil
}

// GetTopic fetches a specific topic by id
// API: /api/topics/show.json?id=X
func (c *Client) GetTopic(id int) (*models.Topic, error) {
	endpoint := "/topics/show.json?id=" + strconv.Itoa(id)

	data, err := c.do(endpoint)
	if err != nil {
		return nil, err
	}

	var topics []models.Topic
	if err := c.unmarshal(data, &topics); err != nil {
		return nil, err
	}
	if len(topics) == 0 {
		return nil, fmt.Errorf("topic not found")
	}
	return &topics[0], nil
}

// GetTopicsByUser fetches topics by username
// API: /api/topics/show.json?username=X
func (c *Client) GetTopicsByUser(username string) ([]models.Topic, error) {
	endpoint := "/topics/show.json?username=" + username

	data, err := c.do(endpoint)
	if err != nil {
		return nil, err
	}

	var topics []models.Topic
	if err := c.unmarshal(data, &topics); err != nil {
		return nil, err
	}
	return topics, nil
}

// GetTopicsByNode fetches topics by node id or name
// API: /api/topics/show.json?node_id=X or /api/topics/show.json?node_name=X
func (c *Client) GetTopicsByNode(nodeID int, nodeName string) ([]models.Topic, error) {
	endpoint := "/topics/show.json?"
	if nodeID > 0 {
		endpoint += "node_id=" + strconv.Itoa(nodeID)
	} else if nodeName != "" {
		endpoint += "node_name=" + nodeName
	} else {
		return nil, fmt.Errorf("either node_id or node_name must be provided")
	}

	data, err := c.do(endpoint)
	if err != nil {
		return nil, err
	}

	var topics []models.Topic
	if err := c.unmarshal(data, &topics); err != nil {
		return nil, err
	}
	return topics, nil
}

// GetReplies fetches replies for a topic
// API: /api/replies/show.json?topic_id=X&page=X&page_size=X
func (c *Client) GetReplies(topicID int, page int, pageSize int) ([]models.Reply, error) {
	endpoint := fmt.Sprintf("/replies/show.json?topic_id=%d", topicID)
	if page > 0 {
		endpoint += "&page=" + strconv.Itoa(page)
	}
	if pageSize > 0 {
		endpoint += "&page_size=" + strconv.Itoa(pageSize)
	}

	data, err := c.do(endpoint)
	if err != nil {
		return nil, err
	}

	var replies []models.Reply
	if err := c.unmarshal(data, &replies); err != nil {
		return nil, err
	}

	// Set topic ID for each reply
	for i := range replies {
		replies[i].TopicID = topicID
	}

	return replies, nil
}

// GetMember fetches member information by username
// API: /api/members/show.json?username=X
func (c *Client) GetMember(username string) (*models.Member, error) {
	endpoint := "/members/show.json?username=" + username

	data, err := c.do(endpoint)
	if err != nil {
		return nil, err
	}

	var member models.Member
	if err := c.unmarshal(data, &member); err != nil {
		return nil, err
	}

	if member.Status == "not_found" || member.ID == 0 {
		return nil, fmt.Errorf("member not found: %s", username)
	}

	return &member, nil
}
