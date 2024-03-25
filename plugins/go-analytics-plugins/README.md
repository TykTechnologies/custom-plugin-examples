# Analytics Plugin - Modify Tyk Analytic Records

This plugin allows you to modify or remove data within the [Tyk Analytics Record](https://github.com/TykTechnologies/tyk/blob/9cb830488690e983fc45c9f94ebafe4976ad7605/rpc/rpc_analytics_purger.go#L14) that is stored temporarily in Redis before it's pumped to a configured backend datastore. Here is the [Analytics Plugin Documentation](https://tyk.io/docs/plugins/plugin-types/analytics-plugins/).

## Configuration

The assumption here is that you've properly configured Redis and validated that the Tyk has established a connection. It's important to understand the values that are stored within the [Tyk Analytics Record Fields](https://tyk.io/docs/tyk-stack/tyk-pump/tyk-analytics-record-fields/) as well consider that Analytics Plugins can only be served via the filesystem and **NOT** using plugin bundles. Another important consideration is that the Analytics Plugin **DOES NOT** have access to the HTTP request or response objects.

## Usage

Suppose you need detailed logging enabled for all APIs but you only need to capture either the `RawRequest` or `RawResponse` depending on the situation -- you can utilize this plugin to delete either values to reduce the amount of unwanted data injested by the backend datastore.

Below is an example of an `API definition` body for the AnalyticsPlugin used to purge `RawRequest` data:
```
{
    "analytics_plugin": {
        "enable": true,
        "func_name": "AnalyticsDeleteRawResponse",
        "plugin_path": "<path>/analytics_plugin.so"
    }
}
```