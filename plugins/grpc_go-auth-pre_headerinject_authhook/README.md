# Go plugin

## Plugin overview

This repository provides a sample [Go](https://golang.org/) plugin for [Tyk](https://tyk.io).

The project implements a simple middleware for header injection (MyPreHook), using a **Pre** hook (see [Tyk custom middleware hooks](https://tyk.io/docs/tyk-api-gateway-v1-9/javascript-plugins/middleware-scripting/)). An authentication hook is also provided (MyAuthCheck), see [hooks.go](hooks.go).

## Requirements

Go compiler.

## Instructions

After checking the requirements, clone this repository:

```
$ git clone https://github.com/TykTechnologies/tyk-plugin-demo-golang.git
```

Enter the plugin directory:

```
$ cd tyk-plugin-demo-golang
```

## Building a bundle

Go plugins are delivered as plugin bundles. The manifest file (`manifest.json`) contains the custom middleware definition.

```
$ tyk-cli bundle build
```

You may check the [tyk-cli documentation](https://github.com/TykTechnologies/tyk-cli) for additional options.

