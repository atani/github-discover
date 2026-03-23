# github-discover

Discover trending and interesting GitHub repositories from your terminal.

## Features

- **Trending**: View trending repositories by time range and language
- **Category Browsing**: Explore repositories organized by category (cli, web, ai, devops, security, data, mobile)
- **Random Discovery**: Get random repository recommendations
- **Enhanced Search**: Search repositories with popularity-sorted results
- **Similar Repos**: Find repositories similar to one you already know
- **Detailed Info**: View repository details with popularity stats
- **Bilingual**: Supports English and Japanese

## Installation

```bash
brew tap atani/tap
brew install github-discover
```

Or download the binary from [Releases](https://github.com/atani/github-discover/releases).

## Usage

### Trending Repositories

```bash
# Show weekly trending (default)
github-discover trending

# Daily trending in Go
github-discover trending --since daily -l go

# Monthly trending, top 10
github-discover trending --since monthly -n 10
```

### Browse by Category

```bash
# List all categories
github-discover browse

# Browse AI & Machine Learning repos
github-discover browse ai

# Available categories: cli, web, ai, devops, security, data, mobile
```

### Random Recommendations

```bash
# Get a random recommendation
github-discover random

# Get 5 random picks
github-discover random -n 5

# Random Rust repos
github-discover random -l rust
```

### Search

```bash
# Search repositories
github-discover search "web framework"

# Search Go repositories sorted by recent updates
github-discover search router -l go --sort updated
```

### Similar Repositories

```bash
# Find repos similar to a given repository
github-discover similar golang/go

# Limit results
github-discover similar denoland/deno -n 5
```

### Repository Info

```bash
# Show detailed info
github-discover info golang/go
```

### Language

```bash
# Use Japanese
github-discover trending --lang ja
```

## Authentication

github-discover works without authentication, but GitHub API rate limits are stricter for unauthenticated requests. Set a token for higher limits:

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
```

## Development

```bash
# Build
go build -o github-discover

# Run
./github-discover trending
```

## License

MIT
