package seqio

import (
	"testing"
	"time"

	"github.com/go-gts/gts/internal/testutils"
)

func TestDate(t *testing.T) {
	now := time.Now()
	in := FromTime(now)
	out := FromTime(in.ToTime())
	testutils.Equals(t, in, out)
}
