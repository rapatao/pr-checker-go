# pr-checker-go

Generate an output compatible with [XBar](https://xbarapp.com/) containing all pull requests pending.

It, currently, only support `GitHub`.

## Installation

1. Create a configuration file containing all repositories that must be monitored.

   The configuration file must be at `${HOME}/.pr-checker.yml`. Details on how to define it can be found in
   the [Configuration](#configuration) section.

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
   pr-checker-go --install=${HOME}/Library/Application Support/xbar/plugins/pr-checker.30m.sh
   ```

## Configuration

Currently, it supports filtering PRs by using 3 different parameters, `author`, `owner` and `repositories`.

The `repositories` filter allows adding a list of repositories, where each entry represents a `OR` statement, which
means that the result includes every PR found in all repositories.

Although, the `author` and `owner` work as a `AND` filtering, which means that, when added to the query, it only returns
PRs that match to all arguments.

So, in a filter that contains a list `Repositories` and an `Author` will return all open PRs in the listed repositories
that were created by the given `Author`

The same will happen with the `Owner` filter, but in this case, it will also restrict the `Repositories`. Which means
that, even if the configuration includes repositories from different owners, only the ones that belong to the configured
owner will be included in the response.

### Examples

* List all PRs created by one user

```yaml
  services:
    - name: GitHub
      provider: github
      token: <secret token>
      author: rapatao

    - name: ...
```

* List all PRs created in any repositories of a user

```yaml
  services:
    - name: GitHub
      provider: github
      token: <secret token>
      owner: rapatao

    - name: ...
```

* List PRs created in a list of repositories

```yaml
  services:
    - name: GitHub
      provider: github
      token: <secret token>
      repositories:
        - rapatao/pr-checker-go
        - org-name/repository-name

    - name: ...
```

## Customizing the menu

# Customizing the menu

It is possible to customize the generation menu by creating a custom template.

The application looks up for the following file `${HOME}/.pr-checker.tmpl`, when it is found, it is used instead of the
default one.

This template must be written using the [Go Template](https://pkg.go.dev/text/template) package.

* The default template can be checked [here](outgen/xbar_stdout.txt.tmpl)

The data provided to the template is a `map[string][]domain.PullRequest{}`, where the `key` is the repository's name and
the `value` is a list of all pull requested on it. The `domain` package can be checked [here](domain/pull_request.go)


