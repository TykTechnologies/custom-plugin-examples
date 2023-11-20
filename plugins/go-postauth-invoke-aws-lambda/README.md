.# Invoke AWS Lambda from the Tyk Gateway

This plugin allows the Tyk Gateway to execute an AWS Lambda function specified from a header value.
It requires a set of IAM credentials from AWS with execute permissions on your lambda of choice.

## Configuration
Configuration is simply done through static reference of the AWS credentials in the plugin.
Simply replace the `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_SESSION_TOKEN` from within the plugin code.

## Usage
Supposing you have the Tyk-Gateway deployed locally on port 8081 with an API configured to execute this 
plugin from the "post" hook, you can invoke a named lambda function `jia_test_custom_go_plugin` as follows:

```shell
curl --location 'http://localhost:8081/lambda/' \
--header 'function: jia_test_custom_go_plugin'
```
