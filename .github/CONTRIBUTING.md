# Contributing to Taskfleet

## Setting up your Development Environment

For running recurring tasks, this repository uses [Task](https://taskfile.dev) as a simpler and
more expressive alternative to `make`. Installation instructions can be found
[here](https://taskfile.dev/installation/). If you are using MacOS, you can simply install the task
runner as follows:

```
brew install go-task
```

To get a list of all available tasks, simply run `task` in your terminal from an arbitrary
directory within this repository. To get up and running on MacOS, you can then run

```
task dev:install-tools
```

which checks that you have installed required software and installs additional development software
via `brew` or `go install`.
