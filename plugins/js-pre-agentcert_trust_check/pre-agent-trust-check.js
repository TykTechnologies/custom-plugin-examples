// AgentCert / agent-trust pre-auth check (JSVM pre hook).
//
// Verifies a calling agent's AgentCert + live TrustScore on each request by
// forwarding the agent's cert headers to an external verifier sidecar
// (agentcert-tag, MIT, ghcr.io/pdxlab/agentcert-tag/sidecar) over its
// ext_authz-native endpoint, then allowing/denying on the verdict. All crypto
// (chain + revocation + score) stays in the sidecar, so there's no Go plugin to
// compile — the gateway only needs enable_jsvm.
//
// Shadow-mode by default: run the sidecar with TAG_MODE=shadow and /ext_authz
// returns 200 (allow) on every request, so this hook annotates but blocks
// nothing until you flip TAG_MODE=enforce on the sidecar. Opt-in per API.
//
// Context for this recipe: TykTechnologies/tyk#8732.

var agentCertTrustCheck = new TykJS.TykMiddleware.NewMiddleware({});

agentCertTrustCheck.NewProcessRequest(function(request, session, config) {
    // Point this at your running agentcert-tag sidecar.
    var SIDECAR = "http://tag-sidecar:8080";

    var verifyReq = {
        "Method": "GET",
        "Domain": SIDECAR,
        "Resource": "/ext_authz",
        "Headers": {
            // The agent presents its cert via x-agent-cert* (see the sidecar's
            // ext_authz header contract). Forward whatever the caller sent.
            "x-agent-cert": request.Headers["X-Agent-Cert"] || "",
            "x-agent-cert-proof": request.Headers["X-Agent-Cert-Proof"] || "",
            "x-agent-cert-carriage": request.Headers["X-Agent-Cert-Carriage"] || "header",
            "x-agent-id": request.Headers["X-Agent-Id"] || ""
        }
    };

    var resp;
    try {
        resp = JSON.parse(TykMakeHttpRequest(JSON.stringify(verifyReq)));
    } catch (e) {
        // Fail-open: if the sidecar is unreachable we don't want this example to
        // hard-block traffic. Flip to a 503 here if you prefer fail-closed.
        log("agentcert: verifier call failed, passing through: " + e);
        return agentCertTrustCheck.ReturnData(request, {});
    }

    var verdict = (resp.Headers && resp.Headers["X-Agent-Trust-Verdict"]) || "UNVERIFIED";
    var score = (resp.Headers && resp.Headers["X-Agent-Trustscore"]) || "";

    // ext_authz: non-200 == deny. In shadow mode the sidecar returns 200 on
    // everything, so this branch only fires under TAG_MODE=enforce.
    if (resp.Code !== 200) {
        request.ReturnOverrides.ResponseCode = 403;
        request.ReturnOverrides.ResponseError = "agent-trust denied: " + verdict;
        return agentCertTrustCheck.ReturnData(request, {});
    }

    // Allow — annotate the request so downstream services / auth can see the
    // verdict and score. In shadow mode this is the only thing that happens.
    request.SetHeaders["x-agent-trust-verdict"] = verdict;
    if (score) {
        request.SetHeaders["x-agent-trustscore"] = score;
    }

    // Return request + (empty) session metadata, per the Tyk JSVM convention.
    return agentCertTrustCheck.ReturnData(request, {});
});

// The last line is called with the name of the middleware here.
log("agentcert trust-check middleware initialised");
