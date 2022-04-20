## Go plugin Template for v3.2.2+

Tyk Gateway version 3.2.2+ uses Go Modules for custom plugins.

Starting in v3.2.2+, The Tyk Plugin Compiler, used to compile go plugins for runtime execution uses the following steps.

This is because Tyk uses Go v1.15 starting in this version, which required you to use go mods or vendor paths to compile your plugin properly.

## Running
To run, you can simply compile the plugin for your Tyk Gateway Version and use it.

```
$ docker run --rm -v $(pwd):/plugin-source tykio/tyk-plugin-compiler:v3.2.2 my-post-plugin.so
```

This generates a file called `my-post-plugin.so` which you can mount onto your Gateway and then inject into the API request lifecycle.

## Development

If you want to start from scratch to replicate this repo, here are the steps taken to get here:

#### 0. go mod init
We first run `go mod init mycustomplugin` to create a go.mod file for the plugin development environment.

####  1. Identify your gateway version.
For our purposes, we will use `v3.2.2-rc7`, which is an unstable build.
The branch is available [here][1].

 
#### 2. Copy the "replace" line from go.mod
In the branch that corresponds to Gateway/Compiler version, we want to copy [the exact][2] `replace` line and paste it into our plugin's "go.mod" file we created in step #0  


#### 3. Find Tyk SHA1 commit ID

We will use the Tyk SHA1 commit ID to "go get" the correct version of Tyk.

on [this][3] page, that's easy to copy: `bda54b0f790c9bc11297c96fe8f2a5b370f39e05`
And then run `go get`
```bash
$ go get github.com/TykTechnologies/tyk@bda54b0f790c9bc11297c96fe8f2a5b370f39e05 
```

#### 4. Run "go mod tidy && go get" for the rest of your dependencies

This will fetch the rest of your dependencies by introspecting your source code:
```
$ go mod tidy && go get
```

Then, we are ready to compile the plugin:

```
$ docker run --rm -v $(pwd):/plugin-source tykio/tyk-plugin-compiler:v3.2.2-rc7 my-post-plugin-322rc7.so
```

Voilà!  We have a `plugin.so` we can mount onto our Gateway's file system.


## FAQ
1. What's with `replace` line in step #2?
This is a necessary step that resolves Tyk dependencies.

[1]: [https://github.com/TykTechnologies/tyk/tree/release-3.2.2]
[2]: [https://github.com/TykTechnologies/tyk/blob/release-3.2.2/go.mod#L101]
[3]: [https://github.com/TykTechnologies/tyk/commits/release-3.2.2]
