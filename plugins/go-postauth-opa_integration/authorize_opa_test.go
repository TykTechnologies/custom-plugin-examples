package main

import (
	"encoding/json"
	"testing"
)

var configJSON = []byte(`{
  "opa": {
    "policy": "cGFja2FnZSB0b2RvcwoKCQlpbXBvcnQgZGF0YS5zY29wZXMKCgkJZGVmYXVsdCBhbGxvdyA9IGZhbHNlCgoJCSMgYWRtaW4gc2NvcGUgYmUgYWxsb3dlZCBmdWxsIGFjY2VzcwoJCWFsbG93IHsKCQkJc2NvcGUgOj0gaW5wdXQuc2NvcGVzW19dCgkJCXNjb3BlID09ICJhZG1pbiIKCQl9CgoJCWFsbG93IHsKCQkJcmVxdWVzdGVkU2NvcGUgOj0gaW5wdXQuc2NvcGVzW19dCgkJCXNjb3BlIDo9IHNjb3Blc1tfXQoJCQkKCQkJcmVxdWVzdGVkU2NvcGUgPSBzY29wZS5zY29wZQoJCQlpbnB1dC5wYXRoID0gc2NvcGUucmVzb3VyY2VzW19dCgkJCWlucHV0Lm1ldGhvZCA9IHNjb3BlLm9wZXJhdGlvbnNbX10KCQl9",
    "rules": {
      "scopes": [
        {
          "operations": [
            "GET"
          ],
          "resources": [
            "/todos",
            "/todos/*"
          ],
          "scope": "todos:read"
        },
        {
          "operations": [
            "PUT",
            "DELETE"
          ],
          "resources": [
            "/todos/*"
          ],
          "scope": "todos:manage"
        },
        {
          "operations": [
            "*"
          ],
          "resources": [
            "*"
          ],
          "scope": "admin"
        }
      ]
    }
  }
}`)

func Test_getOPAFromConfigData(t *testing.T) {
	configData := make(map[string]interface{}, 0)

	if err := json.Unmarshal(configJSON, &configData); err != nil {
		t.Logf("error unmarshalling JSON: %s", err.Error())
		t.FailNow()
	}

	testOpa, err := getOPAFromConfigData(configData)
	if err != nil {
		t.Logf("error from getOPAFromConfigData: %s", err.Error())
		t.FailNow()
	}
	if testOpa == nil {
		t.Log("should get back config from getOPAFromConfigData")
		t.FailNow()
	}
}
