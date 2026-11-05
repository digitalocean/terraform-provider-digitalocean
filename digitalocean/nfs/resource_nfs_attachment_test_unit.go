package nfs

import (
	"testing"
)

func TestNfsShareHasVpcID(t *testing.T) {
	if !nfsShareHasVpcID([]string{"a", "b"}, "b") {
		t.Fatal("expected vpc b to be found")
	}
	if nfsShareHasVpcID([]string{"a"}, "b") {
		t.Fatal("expected vpc b to be missing")
	}
	if nfsShareHasVpcID(nil, "a") {
		t.Fatal("expected empty vpc list to return false")
	}
}
