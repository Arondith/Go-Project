# RepoRadar

A portfolio-ready GitHub repository analyzer built with **Go** and a lightweight browser frontend.

RepoRadar accepts a public GitHub repository, retrieves live metadata from the GitHub REST API, and presents a compact engineering snapshot including stars, forks, open issues, license, update time, and language composition.

## Why this project exists

This repository was created to add **Go** to the portfolio while reinforcing existing web-development skills through REST APIs, asynchronous browser requests, testing, Docker, and CI.

## Tech stack

- Go 1.23
- GitHub REST API
- HTML, CSS, JavaScript
- Docker
- GitHub Actions

## Features

- Analyze any public GitHub repository
- Fetch repository metadata and language statistics concurrently
- Optional GitHub token support for higher API rate limits
- Responsive dashboard
- Health endpoint for deployments
- Unit-tested GitHub API client
- Docker-ready
- CI workflow for formatting, vetting, and tests

## Run locally

Requirements: Go 1.23+

```bash
git clone https://github.com/Arondith/Portfolio-project.git
cd Portfolio-project
go run ./cmd/server
```

Open http://localhost:8080 and try a repository such as `Arondith/KoroAR`.

### Optional GitHub token

Unauthenticated GitHub API calls are rate-limited. You can provide a token through an environment variable:

```bash
# PowerShell
$env:GITHUB_TOKEN="your_token_here"
go run ./cmd/server
```

Do not commit your token.

## API

`GET /api/analyze?owner=OWNER&repo=REPOSITORY`

`GET /api/health`

Example:

```bash
curl "http://localhost:8080/api/analyze?owner=Arondith&repo=KoroAR"
```

## Docker

```bash
docker build -t reporadar .
docker run --rm -p 8080:8080 -e GITHUB_TOKEN reporadar
```

## What this demonstrates

This project is intentionally small enough to understand quickly while still demonstrating backend engineering fundamentals: HTTP services, external API integration, JSON handling, concurrency, error handling, testability, environment configuration, containerization, and CI.

## Next improvements

- Repository comparison mode
- Commit activity trends
- Contributor insights
- Pull-request and issue analytics
- PostgreSQL history storage
- React/TypeScript dashboard
- Deployment to a cloud platform

## License

MIT
