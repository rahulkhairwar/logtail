[![Workflow for Go Standard Action](https://github.com/codecov/go-standard/actions/workflows/go-standard.yml/badge.svg)](https://github.com/codecov/go-standard/actions/workflows/go-standard.yml)
[![codecov](https://codecov.io/gh/codecov/example-go/branch/master/graph/badge.svg)](https://codecov.io/gh/codecov/rahulkhairwar/logtail)

# logtail
Tail logs from a file

### TODOs:
#### Frontend:
<ol>
    <li>Store user's config in cookies.</li>
    - Use a buffer size of logs, so that only a limited number of logs are in memory.
    - Provide an API to load logs on scrolling - how to store past logs? DB vs file copy? Not ideal to create and write to a file being a simple UI tool. DB is a bigger overhead - user needs to have the DB (mongo useful here) running.
</ol>

#### Backend:
<ol>
    <li>Tests for server.go</li>
    <li>Github workflow for "Go Lint"</li>
    <li>Github workflow for "Go Fmt"</li>
    <li>Github workflow for "Go Vet"</li>
    <li>Github workflow for "Go Test"</li>
    <li>Github workflow for "Go Race Detector"</li>
    <li>Graceful shutdown for the HTTP server</li>
</ol>
