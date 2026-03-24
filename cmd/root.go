package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yusteven/v2ex-cleaner/internal/api"
	"github.com/yusteven/v2ex-cleaner/internal/cleaner"
	"github.com/yusteven/v2ex-cleaner/internal/models"
	"github.com/yusteven/v2ex-cleaner/internal/output"
)

func init() {
	// Set custom help function
	cobra.EnableCommandSorting = false
}

var (
	baseURL     string
	outputDir   string
	pretty      bool
	dataType    string
	nodeName    string
	username    string
	topicID     int
	enableClean bool
)

// rootCmd represents base command
var rootCmd = &cobra.Command{
	Use:   "v2ex-cleaner",
	Short: "Tool to fetch and clean data from v2ex API",
	Long: `v2ex-cleaner is a command-line tool for fetching data from v2ex API,
cleaning it for AI and human analysis, and outputting structured JSON format.`,
	// Disable default help command generation (we use custom one)
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// fetchCmd represents fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch data from v2ex API",
	Long:  `Fetch data from v2ex API and save as cleaned JSON files.`,
	RunE:  runFetch,
}

// versionCmd represents version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v2ex-cleaner version 1.0.0")
	},
}

func init() {
	// Fetch command flags
	fetchCmd.Flags().StringVarP(&baseURL, "base-url", "u", "", "v2ex API base URL (default: https://www.v2ex.com/api)")
	fetchCmd.Flags().StringVarP(&outputDir, "output", "o", "output", "Output directory")
	fetchCmd.Flags().BoolVarP(&pretty, "pretty", "p", true, "Format JSON output")
	fetchCmd.Flags().StringVarP(&dataType, "type", "t", "all", "Data types: all, site, nodes, topics, hot, replies, members")
	fetchCmd.Flags().StringVar(&nodeName, "node", "", "Fetch topics from specified node")
	fetchCmd.Flags().StringVar(&username, "user", "", "Fetch topics by username")
	fetchCmd.Flags().IntVar(&topicID, "topic-id", 0, "Fetch specified topic and its replies")
	fetchCmd.Flags().BoolVar(&enableClean, "clean", true, "Apply data cleaning")

	// Add subcommands
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute runs root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runFetch(cmd *cobra.Command, args []string) error {
	client := api.NewClient(baseURL)
	writer := output.NewWriter(outputDir, pretty)

	fmt.Printf("Fetching v2ex data...\n")
	fmt.Printf("Output directory: %s\n\n", outputDir)

	data := make(map[string]interface{})
	data["fetched_at"] = time.Now().Format(time.RFC3339)

	// Parse data types
	types := parseDataTypes(dataType)
	hasType := func(t string) bool {
		for _, dt := range types {
			if dt == t {
				return true
			}
		}
		return false
	}

	// Fetch based on data type
	if hasType("all") || hasType("site") {
		fmt.Println("Fetching site info...")
		if info, err := client.GetSiteInfo(); err == nil {
			if enableClean {
				data["site_info"] = info
			} else {
				if err := writer.WriteSiteInfo(info); err != nil {
					fmt.Printf("Warning: %v\n", err)
				}
			}
		} else {
			fmt.Printf("Error fetching site info: %v\n", err)
		}

		fmt.Println("Fetching site stats...")
		if stats, err := client.GetSiteStats(); err == nil {
			if enableClean {
				data["site_stats"] = stats
			} else {
				if err := writer.WriteSiteStats(stats); err != nil {
					fmt.Printf("Warning: %v\n", err)
				}
			}
		} else {
			fmt.Printf("Error fetching site stats: %v\n", err)
		}
	}

	if hasType("all") || hasType("nodes") {
		fmt.Println("Fetching all nodes...")
		if nodes, err := client.GetAllNodes(); err == nil {
			cleanedNodes := make([]models.CleanedNode, len(nodes))
			for i, n := range nodes {
				cleanedNodes[i] = cleaner.CleanNode(n)
			}
			if enableClean {
				data["nodes"] = cleanedNodes
			} else {
				if err := writer.WriteNodes(cleanedNodes); err != nil {
					fmt.Printf("Warning: %v\n", err)
				}
			}
			fmt.Printf("  Fetched %d nodes\n", len(nodes))
		} else {
			fmt.Printf("Error fetching nodes: %v\n", err)
		}
	}

	if hasType("all") || hasType("topics") {
		fmt.Println("Fetching latest topics...")
		fetchTopics(client, writer, data, false)
	}

	if hasType("all") || hasType("hot") {
		fmt.Println("Fetching hot topics...")
		fetchTopics(client, writer, data, true)
	}

	if nodeName != "" {
		fmt.Printf("Fetching topics from node: %s\n", nodeName)
		if topics, err := client.GetTopicsByNode(0, nodeName); err == nil {
			cleanedTopics := make([]interface{}, len(topics))
			for i, t := range topics {
				cleanedTopics[i] = cleaner.CleanTopic(t)
			}
			data["node_"+nodeName+"_topics"] = cleanedTopics
			fmt.Printf("  Fetched %d topics\n", len(topics))
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}

	if username != "" {
		fmt.Printf("Fetching topics for user: %s\n", username)
		if topics, err := client.GetTopicsByUser(username); err == nil {
			cleanedTopics := make([]interface{}, len(topics))
			for i, t := range topics {
				cleanedTopics[i] = cleaner.CleanTopic(t)
			}
			data["user_"+username+"_topics"] = cleanedTopics
			fmt.Printf("  Fetched %d topics\n", len(topics))
		} else {
			fmt.Printf("Error: %v\n", err)
		}

		// Also fetch member info
		fmt.Printf("Fetching member info: %s\n", username)
		if member, err := client.GetMember(username); err == nil {
			data["member_"+username] = cleaner.CleanMember(*member)
		} else {
			fmt.Printf("Error fetching member: %v\n", err)
		}
	}

	if topicID > 0 {
		fmt.Printf("Fetching topic: %d\n", topicID)
		if topic, err := client.GetTopic(topicID); err == nil {
			cleanedTopic := cleaner.CleanTopic(*topic)
			data["topic"] = cleanedTopic

			// Fetch replies
			fmt.Printf("Fetching replies for topic %d...\n", topicID)
			if replies, err := client.GetReplies(topicID, 1, 100); err == nil {
				cleanedReplies := make([]interface{}, len(replies))
				for i, r := range replies {
					cleanedReplies[i] = cleaner.CleanReply(r)
				}
				data["replies"] = cleanedReplies
				fmt.Printf("  Fetched %d replies\n", len(replies))
			}
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}

	// Write combined output if cleaning is enabled
	if enableClean {
		fmt.Println("\nWriting combined output...")
		if err := writer.WriteCombined(data); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		fmt.Printf("Successfully wrote data to %s/v2ex_data.json\n", outputDir)
	}

	fmt.Println("\nDone!")
	return nil
}

func fetchTopics(client *api.Client, writer *output.Writer, data map[string]interface{}, hot bool) {
	var topics []models.Topic
	var err error

	if hot {
		topics, err = client.GetTopicsHot()
	} else {
		topics, err = client.GetTopicsLatest()
	}

	if err != nil {
		fmt.Printf("Error fetching topics: %v\n", err)
		return
	}

	cleanedTopics := make([]models.CleanedTopic, len(topics))
	for i, t := range topics {
		cleanedTopics[i] = cleaner.CleanTopic(t)
	}

	key := "latest_topics"
	if hot {
		key = "hot_topics"
	}
	data[key] = cleanedTopics
	fmt.Printf("  已获取 %d 条主题\n", len(topics))
}

// parseDataTypes parses comma-separated data types
func parseDataTypes(input string) []string {
	if input == "all" {
		return []string{"all"}
	}
	var types []string
	for _, t := range strings.Split(input, ",") {
		types = append(types, strings.TrimSpace(t))
	}
	return types
}
