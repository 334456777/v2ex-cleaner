# V2EX 数据清洗工具

一个用 Go 语言编写的命令行工具，用于从 v2ex API 抓取数据、清洗数据并输出结构化的 JSON 格式，便于 AI 分析和人类阅读。

## 功能特性

- 🔍 完整 API 支持：站点信息、节点、主题、回复、用户信息
- 🧹 数据清洗：自动移除 HTML 标签、规范化空白字符
- 📅 时间格式化：Unix 时间戳转换为 ISO 8601 格式
- 📊 JSON 输出：格式化输出，便于阅读和解析
- 🎯 灵活过滤：支持按类型、节点、用户、主题 ID 过滤

## 安装

```bash
git clone https://github.com/334456777/v2ex-cleaner.git
cd v2ex-cleaner
go build -o v2ex-cleaner .
```

## 使用方法

### 获取最新主题

```bash
./v2ex-cleaner fetch --type topics --output ./output
```

### 获取热门主题

```bash
./v2ex-cleaner fetch --type hot --output ./output
```

### 获取所有节点

```bash
./v2ex-cleaner fetch --type nodes --output ./output
```

### 获取指定节点主题

```bash
./v2ex-cleaner fetch --node python --output ./output
```

### 获取用户信息和主题

```bash
./v2ex-cleaner fetch --user <username> --output ./output
```

### 获取指定主题及回复

```bash
./v2ex-cleaner fetch --topic-id <id> --output ./output
```

### 获取所有数据

```bash
./v2ex-cleaner fetch --output ./output
```

## 命令参数

| 参数 | 说明 |
|------|--------|
| `--base-url`| v2ex API 基础 URL |
| `--output`| 输出目录（默认: output） |
| `--type`| 数据类型: all, site, nodes, topics, hot, replies, members |
| `--node` | 获取指定节点的主题 |
| `--user` | 按用户名获取主题 |
| `--topic-id` | 获取指定主题及其回复 |
| `--clean` | 应用数据清洗（默认: true） |
| `--pretty`| 格式化 JSON 输出（默认: true） |

## 输出格式

输出数据保存在 `<输出目录>/v2ex_data.json`，格式如下：

```json
{
  "meta": {
    "fetched_at": "2025-02-12T16:00:00Z",
    "source": "v2ex.com",
    "version": "1.0.0"
  },
  "data": {
    "site_info": { ... },
    "site_stats": { ... },
    "nodes": [ ... ],
    "latest_topics": [ ... ],
    "hot_topics": [ ... ]
  }
}
```

## 数据清洗说明

- HTML 标签自动移除
- HTML 实体字符自动解码
- 多余空白字符规范化
- 时间戳转换为 ISO 8601 格式
- 保留原始内容供对比

## API 参考

完整的 v2ex API 文档请参考: [v2ex-api.md](./v2ex-api.md)

## 许可证

MIT License
