# Tyk gRPC Plugin Demo (Python)
This plugin will allow users to write gRPC middleware for Tyk using [Python](https://www.python.org/).

## Requirements:
Please find the pip3 requirements located in the `requirements.txt` file.
* Python 3.*
* Tyk Gateway 4.0.*

## Deploying the gRPC Python Server
Assuming you are in the `custom-plugins/plugins/grpc_python-pre-insert_header/` directory, install the requirements, navigate into the `proto/` directory and launch the Python 3 server.
```bash
pip3 install -r requirements.txt

cd proto/ 

python3 sample_server.py
```
Please note that you may also launch the sample_server.py process by attaching a debugger to it. This allows you to intercept calls that the Tyk-Gateway makes to your gRPC server and is useful for debugging purposes.

## Enabling gRPC on the Tyk-Gateway
The server now will be running on tcp://localhost:5555. Please find the documentation for enabling gRPC plugins [here](https://tyk.io/docs/plugins/supported-languages/rich-plugins/grpc/write-grpc-plugin/).  

At minimum, the Tyk-Gateway needs to know that coprocess is enabled, as well as where the gRPC server lives. These options can be set in the configuration file with the following definitions:  
```json
"coprocess_options" : {
  "enable_coprocess":   true,
  "coprocess_grpc_server": "tcp://localhost:5555"
}
```

However, if you prefer using [environment variables](https://tyk.io/docs/tyk-oss-gateway/configuration/), ensure that you have the following environment variables configured for the Tyk-Gateway. Environment variables supercede the configuration files when both are set. 

```bash
TYK_GW_COPROCESSOPTIONS_ENABLECOPROCESS=true
TYK_GW_COPROCESSOPTIONS_COPROCESSGRPCSERVER=tcp://localhost:5555
```

## Adding gRPC to the API Definition
You'll notice that in `sample_server.py` there are a few example functions such as `MyPreMiddleware`, `MyPostMiddleware` and `MyAuthCheck`. These functions are where you can add your custom logic. Supposing you wish to inject a header at the `Custom pre-middlewares` section of the [request middleware chain](https://tyk.io/docs/concepts/middleware-execution-order/), you can add the following to your Tyk API definition:
```json
    "custom_middleware": {
      "pre": [
        {
          "name": "MyPreMiddleware",
          "path": "",
          "require_session": false,
          "raw_body_only": false
        }
      ],
      "driver": "grpc"
    },
```
Essentially this specifies that we are invoking some gRPC function `MyPreMiddlware` at the `pre` stage of the request. When making calls to your Tyk managed API, you will find that the header `MyPreMiddleware` gets inserted in the API request to your upstream service. 

## Updating protobuf definitions
Should you need to update the protobuf definitions, you simply need to generate them from the Tyk-Gateway repository [located here](https://github.com/TykTechnologies/tyk). The script to update the protobuf bindings is located [here](https://github.com/TykTechnologies/tyk/blob/master/coprocess/proto/update_bindings.sh), and the resultant protobuf definitions can be located under the directory `tyk/coprocess/bindings/python/*`.
