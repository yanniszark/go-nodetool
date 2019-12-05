package nodetool

import (
	"fmt"
	"testing"

	"github.com/yanniszark/go-nodetool/client"
)

func TestVersion(t *testing.T) {

	tests := []struct {
		name            string
		response        string
		error           error
		expectedVersion string
	}{
		{
			name:            "ReleaseVersion is returned",
			response:        `{"ReleaseVersion": "3.9"}`,
			error:           nil,
			expectedVersion: "3.9",
		},
		{
			name:     "client returns error",
			response: "",
			error:    fmt.Errorf("simulated error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nodetool := New(
				&client.FakeClient{
					Response: []byte(test.response),
					Error:    test.error,
				},
			)
			v, err := nodetool.Version()

			if err != nil && test.error == nil {
				t.Fatalf("got unexpected error: %v)", err)
			}
			if err == nil && v.String() != test.expectedVersion {
				t.Fatalf("expected version %s, got %s", test.expectedVersion, v.String())
			}
		})
	}
}
