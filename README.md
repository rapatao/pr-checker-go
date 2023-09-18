# pr-checker-go

Generate an output compatible with [XBar](https://xbarapp.com/) containing all pull requests pending.

It, currently, only support `GitHub`.

## Installation

1. Create a configuration file containing all repositories that must be monitored. 
   
   The file must be placed at `${HOME}/.pr-checker.yml` and must have the following pattern:

   ```yaml
   services:
     - name: GitHub (personal token)
       provider: github
       token: <secret token>
       repositories:
         - username/repository-name
         - org-name/repository-name
     - name: GitHub (other token)
       provider: github
       token: <other secret token>
       repositories:
         - ...
   ```
   
2. Install `pr-checker-go` 
   
   ```shell
   go install github.com/rapatao/pr-checker-go@latest
   ```

3. Install/create the `XBar` plugin

   On Mac, the plugin is usually placed at: `${HOME}/Library/Application Support/xbar/plugins/`

   ```shell
   curl https://raw.githubusercontent.com/rapatao/pr-checker-go/main/examples/pr-checker.30m.sh -o ${HOME}/Library/Application Support/xbar/plugins/pr-checker.30m.sh
   ```
