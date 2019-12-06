package nodetool

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/pborman/uuid"
	"github.com/yanniszark/go-nodetool/client"
)

func TestNodeToolStatus(t *testing.T) {
	host1 := "192.0.2.101"
	uuid1 := uuid.Parse("11111111-1111-1111-1111-111111111111")

	tests := []struct {
		title                 string
		handler               func(t *testing.T, w http.ResponseWriter, r *http.Request)
		expectedResponse      NodeMap
		expectedError         bool
		closeClientConnection bool
	}{
		{
			title: "captured response",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				f, err := os.Open("../testdata/StorageService.json")
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				_, err = io.Copy(w, f)
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResponse: NodeMap{
				"192.168.1.14": &Node{
					Host:   "192.168.1.14",
					ID:     uuid.Parse("4f81ad39-8bc6-4c1f-9dee-778affbb5b90"),
					State:  NodeStateNormal,
					Status: NodeStatusUp,
				},
				"192.168.2.28": &Node{
					Host:   "192.168.2.28",
					ID:     uuid.Parse("f7010dd1-4fe6-4e0e-862b-f479fb734307"),
					State:  NodeStateNormal,
					Status: NodeStatusUp,
				},
			},
		},
		{
			title: "internal server error",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedError: true,
		},
		{
			title:                 "no server response",
			expectedError:         true,
			closeClientConnection: true,
		},
		{
			title: "unparsable json",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
			},
			expectedError: true,
		},
		{
			title: "missing Jolokia value",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write(
					[]byte(`{}`),
				)
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedError: true,
		},
		{
			title: "Empty Jolokia value",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write(
					[]byte(`{"value": {}}`),
				)
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResponse: NodeMap{},
		},
		{
			title: "Hosts unknown by default",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(
					fmt.Sprintf(
						`{"value": {"HostIdMap": {"%s": "%s"}}}`,
						host1, uuid1,
					),
				))
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResponse: NodeMap{
				host1: &Node{
					Host:   host1,
					ID:     uuid1,
					State:  NodeStateNormal,
					Status: NodeStatusUnknown,
				},
			},
		},
		{
			title: "Up node",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(
					fmt.Sprintf(
						`{"value": {"HostIdMap": {"%s": "%s"}, "LiveNodes": ["%s"]}}`,
						host1, uuid1, host1,
					),
				))
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResponse: NodeMap{
				host1: &Node{
					Host:   host1,
					ID:     uuid1,
					State:  NodeStateNormal,
					Status: NodeStatusUp,
				},
			},
		},
		{
			title: "Down node",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(
					fmt.Sprintf(
						`{"value": {"HostIdMap": {"%s": "%s"}, "UnreachableNodes": ["%s"]}}`,
						host1, uuid1, host1,
					),
				))
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResponse: NodeMap{
				host1: &Node{
					Host:   host1,
					ID:     uuid1,
					State:  NodeStateNormal,
					Status: NodeStatusDown,
				},
			},
		},
		{
			title: "Leaving node",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(
					fmt.Sprintf(
						`{"value": {"HostIdMap": {"%s": "%s"}, "LeavingNodes": ["%s"]}}`,
						host1, uuid1, host1,
					),
				))
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResponse: NodeMap{
				host1: &Node{
					Host:   host1,
					ID:     uuid1,
					State:  NodeStateLeaving,
					Status: NodeStatusUnknown,
				},
			},
		},
		{
			title: "Joining node",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(
					fmt.Sprintf(
						`{"value": {"HostIdMap": {"%s": "%s"}, "JoiningNodes": ["%s"]}}`,
						host1, uuid1, host1,
					),
				))
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResponse: NodeMap{
				host1: &Node{
					Host:   host1,
					ID:     uuid1,
					State:  NodeStateJoining,
					Status: NodeStatusUnknown,
				},
			},
		},
		{
			title: "Moving node",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(
					fmt.Sprintf(
						`{"value": {"HostIdMap": {"%s": "%s"}, "MovingNodes": ["%s"]}}`,
						host1, uuid1, host1,
					),
				))
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedResponse: NodeMap{
				host1: &Node{
					Host:   host1,
					ID:     uuid1,
					State:  NodeStateMoving,
					Status: NodeStatusUnknown,
				},
			},
		},
		{
			title: "Live intersects with unreachable",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(
					fmt.Sprintf(
						`{"value": {"HostIdMap": {"%s": "%s"}, "LiveNodes": ["%s"], "UnreachableNodes": ["%s"]}}`,
						host1, uuid1, host1, host1,
					),
				))
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedError: true,
		},
		{
			title: "Leaving intersects with joining",
			handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
				_, err := w.Write([]byte(
					fmt.Sprintf(
						`{"value": {"HostIdMap": {"%s": "%s"}, "LeavingNodes": ["%s"], "JoiningNodes": ["%s"]}}`,
						host1, uuid1, host1, host1,
					),
				))
				if err != nil {
					t.Fatal(err)
				}
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(
			test.title,
			func(t *testing.T) {
				ts := httptest.NewTLSServer(
					http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							test.handler(t, w, r)
						},
					),
				)
				defer ts.Close()
				if test.closeClientConnection {
					ts.Config.Handler = http.HandlerFunc(
						func(w http.ResponseWriter, r *http.Request) {
							ts.CloseClientConnections()
						},
					)
				}
				u, err := url.Parse(ts.URL)
				if err != nil {
					t.Fatal(err)
				}
				client := client.New(u, ts.Client())
				nt := New(client)
				ntResponse, ntErr := nt.Status()
				if ntErr == nil {
					if test.expectedError {
						t.Errorf("The expected error did not occur.")
					}
				} else {
					t.Logf("logged error: '%s'", ntErr)
					if !test.expectedError {
						t.Errorf("Unexpected error")
					}
				}
				if !reflect.DeepEqual(ntResponse, test.expectedResponse) {
					t.Errorf(
						"Unexpected response. Expected: %+v. Got %+v.",
						test.expectedResponse,
						ntResponse,
					)
				}
			},
		)
	}
}
