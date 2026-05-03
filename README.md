# pr-checker-go

A standalone macOS menu bar application to monitor your GitHub pull requests.

## Installation

### From Source
1. Clone the repository and build:
   ```shell
   go build -o pr-checker-app main.go
   ```

### Using Go Install
Alternatively, you can install the latest version directly:
```shell
go install github.com/rapatao/pr-checker-go@latest
```

### Background Execution (Launch Agent)
To run the app in the background on login:
```shell
pr-checker-app --install
```

To remove the background service:
```shell
pr-checker-app --uninstall
```
*Note: This unloads the service from `launchctl` and removes the associated plist file.*

### Configuration

Create a configuration file at `${HOME}/.pr-checker.yml`:

```yaml
services:
  - name: GitHub
    provider: github
    token: <your-personal-access-token>
    # Optional filters:
    # author: <username>
    # owner: <org-or-user>
    # repositories:
    #   - username/repository-name
```

## Features

* **Real-time Monitoring:** Automatically polls for PR updates every 30 minutes.
* **Grouped View:** PRs are organized by repository in the menu.
* **Interactive:** Click a repository to open its URL or a PR to open it directly in your browser.
* **Manual Refresh:** Force a check for new PRs at any time.
* **Status:** Easily view the total count of open PRs in your menu bar.
