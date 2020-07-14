# OAuth2 Scope evaluation with OPA Policy

The current status of this plugin is **proof-of-concept**.

This plugin allows the Tyk Gateway to evaluate OAuth2 Scopes to determine
whether the user (resource owner) has provided consent to perform an action.

It uses [Open Policy Agent](https://www.openpolicyagent.org/) (OPA) as the underlying policy enforcement mechanism. OPA was chosen as the policy
agent, because it is Open Source, and facilitates the concepts of Policies as Code.

The plugin assumes that the request has already been authenticated.

## Configuration

> Caveat: I am not a REGO or OPA expert - so there are probably much better ways, of achieving the same thing.
> Please familiarise yourself with the OPA docs to ensure that the policy works for your use-case.

The runtime configuration of OPA is configured in the Tyk API Definition object, `db`, `policy` 
& `module_name` to load. Just insert these as fields in the `config_data` json object of the API Definition
with the key `opa`.

```.json
"config_data": {
    "opa": {
        "db": {
            "scopes": [
                {
                    "operations": [
                        "GET"
                    ],
                    "resources": [
                        "/todos/todos"
                    ],
                    "scope": "todos:read"
                }
            ]
        },
        "module_name": "todos",
        "policy": "cGFja2FnZSB0b2RvcwoKaW1wb3J0IGRhdGEuc2NvcGVzCgpkZWZhdWx0IGFsbG93ID0gZmFsc2UKCiMgYWRtaW4gc2NvcGUgYmUgYWxsb3dlZCBmdWxsIGFjY2VzcwphbGxvdyB7CglzY29wZSA6PSBpbnB1dC5zY29wZXNbX10KCXNjb3BlID09ICJhZG1pbiIKfQoKYWxsb3cgewoJcmVxdWVzdGVkU2NvcGUgOj0gaW5wdXQuc2NvcGVzW19dCglzY29wZSA6PSBzY29wZXNbX10KCQkKCXJlcXVlc3RlZFNjb3BlID0gc2NvcGUuc2NvcGUKCWlucHV0LnBhdGggPSBzY29wZS5yZXNvdXJjZXNbX10KCWlucHV0Lm1ldGhvZCA9IHNjb3BlLm9wZXJhdGlvbnNbX10KfQ=="
    }
}
```

You will see in the above example, the policy has been base64 encoded. This is because the policy is actually
written in the rego language. Decoding the policy, you will see it looks like this:

```text
echo -n 'cGFja2FnZSB0b2RvcwoKaW1wb3J0IGRhdGEuc2NvcGVzCgpkZWZhdWx0IGFsbG93ID0gZmFsc2UKCiMgYWRtaW4gc2NvcGUgYmUgYWxsb3dlZCBmdWxsIGFjY2VzcwphbGxvdyB7CglzY29wZSA6PSBpbnB1dC5zY29wZXNbX10KCXNjb3BlID09ICJhZG1pbiIKfQoKYWxsb3cgewoJcmVxdWVzdGVkU2NvcGUgOj0gaW5wdXQuc2NvcGVzW19dCglzY29wZSA6PSBzY29wZXNbX10KCQkKCXJlcXVlc3RlZFNjb3BlID0gc2NvcGUuc2NvcGUKCWlucHV0LnBhdGggPSBzY29wZS5yZXNvdXJjZXNbX10KCWlucHV0Lm1ldGhvZCA9IHNjb3BlLm9wZXJhdGlvbnNbX10KfQ==' | base64 --decode
package todos

import data.scopes

default allow = false

# admin scope be allowed full access
allow {
	scope := input.scopes[_]
	scope == "admin"
}

allow {
	requestedScope := input.scopes[_]
	scope := scopes[_]

	requestedScope = scope.scope
	input.path = scope.resources[_]
	input.method = scope.operations[_]
}
```

## Usage

OAuth2: Because an access_token may be opaque, we should check for the scopes from an introspection response.
JWT: Can contain scopes claims of the token.
Tyk Bearer Token: Scopes may be present in the session metadata.

Because the scopes could be anywhere - it is assumed that a plugin, or middleware has performed some pre-processing
and injected the scope into an internal header `X-Tyk-Plugin-OAuth2Introspect-Scope`.

## Maybe TODO / Current Known Limitations

- [ ] Make the scope location generic / configurable so that it can extract from specified location.
- [ ] Standardise the required content field names of the `opa` object in the apidefinition config data.
- [ ] Currently a new instance of OPA is instantiated for every request. This needs to be optimised.
