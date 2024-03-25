# Connect to Redis Cache Database

This plugin allows you to establish a connection to the Redis database configured for your Tyk Gateway using librairies provided by Tyk. It provides basic examples on retrieiving and modifying specific Redis data.

## Configuration

The assumption here is that you've properly configured a Redis database and validated that the Tyk Gateway is connected. You can configure this plugin at any level of the middleware execution chain (pre, auth, post, post-auth, response etc). The important thing to consider here is the variable `const pluginDefaultKeyPrefix = "Plugin-data:"` -- this value determines the prefix of the Redis key-value lookup. For example, you set the value to `apikey-` to retrieve API keys since API keys in Tyk are stored as such `apikey-{key_hash}`.

## Usage
As you can imagine there are so many permutations and use-cases where you can take advantage of the Redis connection plugin. Typically the Redis connection plugin is used to manipulate data or to add an additional layer of logic on top of the authentication layer etc.