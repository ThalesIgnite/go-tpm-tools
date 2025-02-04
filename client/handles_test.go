package client_test

import (
	"reflect"
	"testing"

	"github.com/ThalesIgnite/go-tpm-tools/client"
	"github.com/ThalesIgnite/go-tpm-tools/notinternal/test"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

const (
	// Maximum number of handles to keys tests can create within a simulator.
	maxHandles = 3
)

func TestHandles(t *testing.T) {
	rwc := test.GetTPM(t)
	defer client.CheckedClose(t, rwc)

	expected := make([]tpmutil.Handle, 0)
	for i := 0; i < maxHandles; i++ {
		expected = append(expected, test.LoadRandomExternalKey(t, rwc))

		handles, err := client.Handles(rwc, tpm2.HandleTypeTransient)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(handles, expected) {
			t.Errorf("Handles mismatch got: %v; want: %v", handles, expected)
		}
	}

	// Don't leak our handles
	for _, handle := range expected {
		if err := tpm2.FlushContext(rwc, handle); err != nil {
			t.Error(err)
		}
	}
}
