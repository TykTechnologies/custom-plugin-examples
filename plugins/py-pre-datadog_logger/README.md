# tyk-plugin-datadog

## Overview

This plugin sends log data to a [Datadog](https://www.datadoghq.com/) agent.

It will increment a counter on every request, using this format: `tyk.requests.api-APIID` (where `APIID` is your API ID).

## Requirements

- [Datadog Python module](http://datadogpy.readthedocs.io/en/latest/): `pip3 install datadog`
- [Datadog Agent](http://docs.datadoghq.com/guides/basic_agent_usage/)

## License

This project is released under the MPL v2.0. See [full version of the license](LICENSE.md).
