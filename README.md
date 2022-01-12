# bldr


## Config
- Support .bldrrc for setting repo defaults
  - Alternatively, these could just go in the pipeline-config.yaml ?
  - Main branch name
  - Container registry

## Plugins

- Some default initial plugins:
  - build:
    - Init artefact directory (and copy pipeline config)
    - DockerLogin
    - EnsureDeploymentBaseExists (?)
    - Manifest writer (after all pipeline plugins have run)
    - DeploymentContainer builder

  - deploy
    - DockerRun
    - K8sDeploy


- Could look to use the go plugin system: https://medium.com/learning-the-go-programming-language/writing-modular-go-programs-with-plugins-ec46381ee1a9
  - Initialise plugins by building them using `go build --buildmode=plugin -o /plugins/something.so github.com/something/plugin.go`

- Detecting diffs should be a plugin, so that we can be smarter than just relying on devs to populate `include` properly
