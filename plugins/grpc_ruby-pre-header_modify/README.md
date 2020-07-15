#Tyk Plugin Demo Ruby

This plugin will allow users to write middleware for Tyk using [Ruby](https://www.ruby-lang.org).


##Build Requirements

* [Ruby 2.x](https://www.ruby-lang.org)
* [Go](https://golang.org/)
* [gRPC](https://www.grpc.io/) (Ruby gem) `gem install grpc:1.0.0 grpc-tools:1.0.0`
    * There are currently [issues with the latest and the precompiled versions](https://github.com/grpc/grpc/issues/7727) of the grpc gem so it is advised that you install version "1.0.0" for the time being.


##Usage

1. Clone and build [Tyk](https://github.com/TykTechnologies/tyk) locally with the `coprocess` grpc build tag:

        go build -tags 'coprocess grpc'

2. Clone this repo and `cd` into its location on your filesystem. This repo contains the Sample Server which you can run using the following command:

        ruby sample_server.rb

3. In another terminal, run your tyk instance. A simple way to do this would be to navigate into the tyk repository from step 1 and run the following:

        ./tyk

4. In a third terminal navigate back to your cloned version of [this repo](https://github.com/TykTechnologies/tyk-plugin-demo-ruby) and run the following cURL request:

        curl -v http://localhost:8080/grpc-api-test/headers

This will then make a request to the server using the middleware defined in the `grpc_app_sample.json` file.

## License

This project is released under the MPL v2.0. See [full version of the license](LICENSE.md).
