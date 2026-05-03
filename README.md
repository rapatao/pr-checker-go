# pr-checker-go

A standalone macOS menu bar application to monitor your GitHub pull requests.

## Installation

### Binary

1. Build the application:

   ```shell
   go build -o pr-checker-app main.go
   ```

2. (Optional) Install as a macOS Launch Agent to run in the background on login:

   ```shell
   ./pr-checker-app --install
   ```

   *Note: This will create a `launchd` plist file in `~/Library/LaunchAgents/` using the absolute path of the current binary.*

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
