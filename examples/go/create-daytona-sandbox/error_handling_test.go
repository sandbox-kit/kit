package main

import (
	"bytes"
	"errors"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"strings"
	"testing"
)

func TestHandleErrorUsesCommonKind(t *testing.T) {
	for _, tc := range []struct {
		kind sandbox.ErrorKind
		want string
	}{
		{sandbox.ErrorKindInvalidArgument, "Check resources.cpu_cores"},
		{sandbox.ErrorKindUnsupported, "Check resources.cpu_cores"},
		{sandbox.ErrorKindAuthentication, "credentials"},
		{sandbox.ErrorKindPermissionDenied, "permission"},
		{sandbox.ErrorKindNotFound, "not found"},
		{sandbox.ErrorKindRateLimited, "limit"},
		{sandbox.ErrorKindResourceExhausted, "limit"},
		{sandbox.ErrorKindCanceled, "canceled"},
		{sandbox.ErrorKindTimeout, "timed out"},
		{sandbox.ErrorKindUnknown, "unexpected failure"},
	} {
		var output bytes.Buffer
		err := sandbox.NewError(sandbox.ErrorInfo{Kind: tc.kind, Field: sandbox.CreateFieldResourcesCPUCores, Message: "unexpected failure"}, nil)
		handleError(&output, err)
		if !strings.Contains(output.String(), tc.want) {
			t.Fatalf("%s: %s", tc.kind, output.String())
		}
	}
	var output bytes.Buffer
	handleError(&output, errors.New("dotenv failure"))
	if !strings.Contains(output.String(), "dotenv failure") {
		t.Fatal("ordinary error hidden")
	}
}
