# tyk-plugin-queue

## Overview

This plugin sends a message to a queue server, it uses [kombu](https://github.com/celery/kombu) as the messaging library.

The following transports are supported: AMQP, QPID, Redis, MongoDB, SQS, ZooKeeper, SLMQ. You may check the [kombu documentation](https://kombu.readthedocs.io/) for additional options.

## Requirements

- [kombu](https://github.com/celery/kombu): `pip3 install kombu`

## License

This project is released under the MPL v2.0. See [full version of the license](LICENSE.md).
