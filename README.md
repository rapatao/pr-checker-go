# pr-checker-go

A macOS menu bar application that lists your open pull requests from GitHub and Gitea.

The menu bar shows the number of open PRs. The menu groups them by repository, and each PR shows its creation date, review decision and CI check status.

## Requirements

* macOS
* Go (version from `go.mod`) and the Xcode Command Line Tools, to build from source

## Installation

### Using Go Install

```shell
go install github.com/rapatao/pr-checker-go@latest
```

### From Source

```shell
git clone https://github.com/rapatao/pr-checker-go.git
cd pr-checker-go
make build
```

This produces the `pr-checker-go` binary in the repository root, built from the vendored dependencies.

## Configuration

Create `~/.pr-checker.yml`. Each entry under `services` is one account to query:

```yaml
services:
  - name: GitHub
    provider: github
    token: <your-personal-access-token>
    # Optional filters:
    # author: <username>
    # owner: <org-or-user>
    # repositories:
    #   - owner/repository-name
  - name: Gitea
    provider: gitea
    url: https://gitea.example.com
    token: <your-access-token>
    # Optional filters:
    # author: <username>
    # owner: <org-or-user>
    # repositories:
    #   - owner/repository-name
```

| Key            | Required    | Description                                                    |
|----------------|-------------|----------------------------------------------------------------|
| `name`         | no          | Label for the service.                                         |
| `provider`     | yes         | `github` or `gitea`. Any other value makes the app exit.       |
| `token`        | yes         | Access token for the provider.                                 |
| `url`          | Gitea only  | Base URL of the Gitea instance.                                |
| `author`       | no          | Only PRs opened by this user.                                  |
| `owner`        | no          | Only PRs in repositories owned by this user or organization.   |
| `repositories` | no          | Only PRs in these repositories, as `owner/name`.               |

The config file is read on every refresh, so changes apply on the next refresh without restarting the app.

### Provider notes

* **GitHub** uses the search API and returns at most 100 PRs per service. The token needs read access to the repositories you want to see (the `repo` scope for private repositories on a classic token).
* **Gitea** lists open PRs in repositories the token can access. The token needs the `read:repository` and `read:issue` scopes.

## Usage

Run the binary:

```shell
pr-checker-go
```

* PRs refresh at startup and then every 30 minutes.
* Click a repository to open it in the browser, or a PR to open the PR.
* **Refresh now** forces an update. Starting a refresh cancels any refresh still in progress.
* The footer shows the time of the last update.

### Run on Login (Launch Agent)

```shell
pr-checker-go --install
```

This writes `~/Library/LaunchAgents/com.rapatao.pr-checker.plist` pointing to the current binary path and loads it with `launchctl`. If you move the binary, run `--install` again.

The agent uses `KeepAlive`, so launchd restarts the app after **Quit**. To stop it, remove the agent:

```shell
pr-checker-go --uninstall
```

This unloads the agent from `launchctl` and deletes the plist file.

## Development

```shell
make deps    # go mod tidy && go mod vendor
make build   # build the binary
make lint    # golangci-lint run
go test ./...
```

Every push to `main` runs build, lint and tests, then creates a patch version tag.
