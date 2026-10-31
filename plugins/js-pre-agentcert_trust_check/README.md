### Agent-Trust Pre-Auth Check (AgentCert + TrustScore)

Our goal:

-- verify a calling **agent's** AgentCert + live TrustScore on each request, and allow / deny / annotate on the verdict

-- keep all crypto out of the gateway: a JSVM `pre` hook just calls an external verifier sidecar, so there is no Go plugin to compile

-- ship it **off by default** and **shadow-mode first**, opt-in per API

As more API callers are autonomous agents, operators increasingly ask "should I trust this agent right now," not just "is this key/cert valid." This is the "middleware plugin" path a Tyk maintainer suggested in [TykTechnologies/tyk#8732](https://github.com/TykTechnologies/tyk/issues/8732).

The verifier is [`agentcert-tag`](https://github.com/pdxlab/agentcert-tag) (MIT, `ghcr.io/pdxlab/agentcert-tag/sidecar`). It verifies a presented AgentCert — signature chain + revocation + a live TrustScore — over its ext_authz-native `/ext_authz` endpoint and returns allow/deny plus `x-agent-trust-verdict` / `x-agent-trustscore` headers (`VERIFIED` / `REVOKED` / `UNVERIFIED`). The gateway stays lean; this recipe mirrors the [agentgateway](https://github.com/pdxlab/agentcert-tag/tree/main/integrations/agentgateway) one.

## Setup Instructions

1) Run the verifier sidecar (shadow-mode is the default — it allows everything and only annotates):

```
docker run -p 8080:8080 -e TAG_MODE=shadow \
  ghcr.io/pdxlab/agentcert-tag/sidecar:latest
```

2) Make sure your JSVM is enabled in `tyk.conf`: `"enable_jsvm": true`.

3) Add the [JavaScript plugin](pre-agent-trust-check.js) to the local filesystem and set the relative path (to the Gateway binary) in the API Definition, **or** [import the API `apidef.json`](apidef.json). The plugin runs in the `pre` phase, so it verifies the caller before authentication. Point the `SIDECAR` constant in the JS at your running sidecar.

https://tyk.io/docs/getting-started/key-concepts/middleware-execution-order/

## Roll out

1. **Shadow (default).** With `TAG_MODE=shadow` the sidecar returns `200` on every request, so the `pre` hook never denies — it just injects `x-agent-trust-verdict` / `x-agent-trustscore` for downstream. Zero request impact; use it to observe.
2. **Enforce.** Set `TAG_MODE=enforce` on the sidecar; now non-`VERIFIED` agents get a `403` from the hook.

#### Try it out

```
# shadow mode — allowed, annotated for downstream
$ curl -k https://localhost:8081/agent-trust-api/get \
    --header "X-Agent-Cert: <presented-cert>" --header "X-Agent-Id: did:agent:acme/billing"
{ "headers": { "X-Agent-Trust-Verdict": "VERIFIED", "X-Agent-Trustscore": "91", ... } }

# enforce mode — revoked / unverified agent
$ curl -k https://localhost:8081/agent-trust-api/get --header "X-Agent-Cert: <revoked-cert>"
{ "error": "agent-trust denied: REVOKED" }
```

> **Version note:** the JSVM API (`NewMiddleware` / `NewProcessRequest` / `ReturnData`, `request.SetHeaders`, `request.ReturnOverrides`) and the `TykMakeHttpRequest` response shape can vary by Tyk version — validate `pre-agent-trust-check.js` against your gateway. If you'd rather keep the verifier call out-of-process at the auth stage, the same `/ext_authz` call works from a gRPC coprocess plugin. The ext_authz header contract is documented in the [agentcert-tag](https://github.com/pdxlab/agentcert-tag) repo.
