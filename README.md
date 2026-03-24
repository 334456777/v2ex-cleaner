# V2EX Data Cleaner

A command-line tool written in Go for fetching data from v2ex API, cleaning it, and outputting structured JSON format for easy AI analysis and human reading.

## Features

- 🔍 **Complete API Support**: Site info, nodes, topics, replies, user information
- 🧹 **Data Cleaning**: Automatically removes HTML tags and normalizes whitespace
- 📅 **Time Formatting**: Converts Unix timestamps to ISO 8601 format
- 📊 **JSON Output**: Formatted output for easy reading and parsing
- 🎯 **Flexible Filtering**: Filter by type, node, user, or topic ID

## Installation

```bash
git clone https://github.com/334456777/v2ex-cleaner.git
cd v2ex-cleaner
go build -o v2ex-cleaner .
```

## Usage

### Fetch Latest Topics

```bash
./v2ex-cleaner fetch --type topics --output ./output
```

### Fetch Hot Topics

```bash
./v2ex-cleaner fetch --type hot --output ./output
```

### Fetch All Nodes

```bash
./v2ex-cleaner fetch --type nodes --output ./output
```

### Fetch Topics from Specific Node

```bash
./v2ex-cleaner fetch --node python --output ./output
```

### Fetch User Info and Topics

```bash
./v2ex-cleaner fetch --user <username> --output ./output
```

### Fetch Specific Topic and Replies

```bash
./v2ex-cleaner fetch --topic-id <id> --output ./output
```

### Fetch All Data

```bash
./v2ex-cleaner fetch --output ./output
```

## Command Parameters

| Parameter | Description |
|-----------|-------------|
| `--base-url` | v2ex API base URL |
| `--output` | Output directory (default: output) |
| `--type` | Data types: all, site, nodes, topics, hot, replies, members |
| `--node` | Fetch topics from specified node |
| `--user` | Fetch topics by username |
| `--topic-id` | Fetch specified topic and its replies |
| `--clean` | Apply data cleaning (default: true) |
| `--pretty` | Format JSON output (default: true) |

## Output Format

Output data is saved in `<output directory>/v2ex_data.json` with the following format:

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

## Data Cleaning

- Automatic HTML tag removal
- HTML entity character decoding
- Whitespace normalization
- Timestamp conversion to ISO 8601 format
- Original content preserved for comparison

## API Reference

For complete v2ex API documentation, see: [v2ex-api.md](./v2ex-api.md)

## License

MIT License
