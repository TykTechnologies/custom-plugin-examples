# Send a custom error message in a Golang function.

This plugin allows the Gateway proxy to process a custom error message that is produced by this plugin, based on a specific condition being set. 

## Configuration
The function `sendErrorResponseFromMiddleware` contains a map of different error codes and their custom messages. You can add as many as you like, or read them from a different file into a map. The function takes a key as an argument so that we can dynamically pull the correct error message from the map. Then, we create a new map for the response containing the key and value, which is converted into JSON using the `json.Marshal` built in Go-function. We also write the header to explicitly use application/JSON, and the body to contain the new custom error message.

You can customize this to your liking, specifically by modifying the `ReturnError` function. In the case of this example, we are accepting a header that may contain a specific error code, then returning the custom error message from our map of errors. In your case, you might want to return a specific error code based on a redirect, path visited, query parameter located in the request, or other conditional. The sky is the limit.

## Usage
To be used in the post hook in your API definition. Call ReturnError as the function name, give the path for where the binary file is located on the Gateway file system, and you're off to the races.
