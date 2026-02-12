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
	Short: "从 v2ex API 获取并清洗数据的工具",
	Long: `v2ex-cleaner 是一个命令行工具，用于从 v2ex API 获取数据，
清洗后供 AI 和人类分析，并输出结构化的 JSON 格式。`,
	// Disable default help command generation (we use custom one)
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// fetchCmd represents fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "从 v2ex API 获取数据",
	Long:  `从 v2ex API 获取数据并保存为清洗后的 JSON 文件。`,
	RunE:  runFetch,
}

// versionCmd represents version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "打印版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v2ex-cleaner 版本 1.0.0")
	},
}

func init() {
	// Fetch command flags
	fetchCmd.Flags().StringVarP(&baseURL, "base-url", "u", "", "v2ex API 基础 URL (默认: https://www.v2ex.com/api)")
	fetchCmd.Flags().StringVarP(&outputDir, "output", "o", "output", "输出目录")
	fetchCmd.Flags().BoolVarP(&pretty, "pretty", "p", true, "格式化 JSON 输出")
	fetchCmd.Flags().StringVarP(&dataType, "type", "t", "all", "数据类型: all, site, nodes, topics, hot, replies, members")
	fetchCmd.Flags().StringVar(&nodeName, "node", "", "获取指定节点的主题")
	fetchCmd.Flags().StringVar(&username, "user", "", "按用户名获取主题")
	fetchCmd.Flags().IntVar(&topicID, "topic-id", 0, "获取指定主题及其回复")
	fetchCmd.Flags().BoolVar(&enableClean, "clean", true, "应用数据清洗")

	// Add subcommands
	rootCmd.AddCommand(fetchCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute runs root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func runFetch(cmd *cobra.Command, args []string) error {
	client := api.NewClient(baseURL)
	writer := output.NewWriter(outputDir, pretty)

	fmt.Printf("正在获取 v2ex 数据...\n")
	fmt.Printf("输出目录: %s\n\n", outputDir)

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
		fmt.Println("正在获取站点信息...")
		if info, err := client.GetSiteInfo(); err == nil {
			if enableClean {
				data["site_info"] = info
			} else {
				if err := writer.WriteSiteInfo(info); err != nil {
					fmt.Printf("警告: %v\n", err)
				}
			}
		} else {
			fmt.Printf("获取站点信息错误: %v\n", err)
		}

		fmt.Println("正在获取站点统计...")
		if stats, err := client.GetSiteStats(); err == nil {
			if enableClean {
				data["site_stats"] = stats
			} else {
				if err := writer.WriteSiteStats(stats); err != nil {
					fmt.Printf("警告: %v\n", err)
				}
			}
		} else {
			fmt.Printf("获取站点统计错误: %v\n", err)
		}
	}

	if hasType("all") || hasType("nodes") {
		fmt.Println("正在获取所有节点...")
		if nodes, err := client.GetAllNodes(); err == nil {
			cleanedNodes := make([]models.CleanedNode, len(nodes))
			for i, n := range nodes {
				cleanedNodes[i] = cleaner.CleanNode(n)
			}
			if enableClean {
				data["nodes"] = cleanedNodes
			} else {
				if err := writer.WriteNodes(cleanedNodes); err != nil {
					fmt.Printf("警告: %v\n", err)
				}
			}
			fmt.Printf("  已获取 %d 个节点\n", len(nodes))
		} else {
			fmt.Printf("获取节点错误: %v\n", err)
		}
	}

	if hasType("all") || hasType("topics") {
		fmt.Println("正在获取最新主题...")
		fetchTopics(client, writer, data, false)
	}

	if hasType("all") || hasType("hot") {
		fmt.Println("正在获取热门主题...")
		fetchTopics(client, writer, data, true)
	}

	if nodeName != "" {
		fmt.Printf("正在从节点获取主题: %s\n", nodeName)
		if topics, err := client.GetTopicsByNode(0, nodeName); err == nil {
			cleanedTopics := make([]interface{}, len(topics))
			for i, t := range topics {
				cleanedTopics[i] = cleaner.CleanTopic(t)
			}
			data["node_"+nodeName+"_topics"] = cleanedTopics
			fmt.Printf("  已获取 %d 条主题\n", len(topics))
		} else {
			fmt.Printf("错误: %v\n", err)
		}
	}

	if username != "" {
		fmt.Printf("正在获取用户主题: %s\n", username)
		if topics, err := client.GetTopicsByUser(username); err == nil {
			cleanedTopics := make([]interface{}, len(topics))
			for i, t := range topics {
				cleanedTopics[i] = cleaner.CleanTopic(t)
			}
			data["user_"+username+"_topics"] = cleanedTopics
			fmt.Printf("  已获取 %d 条主题\n", len(topics))
		} else {
			fmt.Printf("错误: %v\n", err)
		}

		// Also fetch member info
		fmt.Printf("正在获取用户信息: %s\n", username)
		if member, err := client.GetMember(username); err == nil {
			data["member_"+username] = cleaner.CleanMember(*member)
		} else {
			fmt.Printf("获取用户错误: %v\n", err)
		}
	}

	if topicID > 0 {
		fmt.Printf("正在获取主题: %d\n", topicID)
		if topic, err := client.GetTopic(topicID); err == nil {
			cleanedTopic := cleaner.CleanTopic(*topic)
			data["topic"] = cleanedTopic

			// Fetch replies
			fmt.Printf("正在获取主题 %d 的回复...\n", topicID)
			if replies, err := client.GetReplies(topicID, 1, 100); err == nil {
				cleanedReplies := make([]interface{}, len(replies))
				for i, r := range replies {
					cleanedReplies[i] = cleaner.CleanReply(r)
				}
				data["replies"] = cleanedReplies
				fmt.Printf("  已获取 %d 条回复\n", len(replies))
			}
		} else {
			fmt.Printf("错误: %v\n", err)
		}
	}

	// Write combined output if cleaning is enabled
	if enableClean {
		fmt.Println("\n正在写入合并输出...")
		if err := writer.WriteCombined(data); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		fmt.Printf("成功写入数据到 %s/v2ex_data.json\n", outputDir)
	}

	fmt.Println("\n完成!")
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
		fmt.Printf("获取主题错误: %v\n", err)
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
