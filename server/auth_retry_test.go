package main

import (
	"errors"
	"testing"
	"time"

	"github.com/tinode/chat/server/auth"
)

// flakyAuth fails RestrictedTags a given number of times, then succeeds.
type flakyAuth struct {
	auth.AuthHandler
	failures int
	calls    int
}

func (a *flakyAuth) RestrictedTags() ([]string, error) {
	a.calls++
	if a.calls <= a.failures {
		return nil, errors.New("connection refused")
	}
	return []string{"rest"}, nil
}

func withFastAuthRetry(t *testing.T, timeout time.Duration) {
	oldTimeout, oldInterval := authStartupRetryTimeout, authStartupRetryInterval
	authStartupRetryTimeout, authStartupRetryInterval = timeout, time.Millisecond
	t.Cleanup(func() { authStartupRetryTimeout, authStartupRetryInterval = oldTimeout, oldInterval })
}

func TestRestrictedTagsWithRetryRecovers(t *testing.T) {
	withFastAuthRetry(t, time.Second)
	a := &flakyAuth{failures: 3}
	tags, err := restrictedTagsWithRetry(a, "rest")
	if err != nil || len(tags) != 1 || tags[0] != "rest" || a.calls != 4 {
		t.Fatalf("got tags=%v err=%v calls=%d, want [rest] <nil> 4", tags, err, a.calls)
	}
}

func TestRestrictedTagsWithRetryGivesUp(t *testing.T) {
	withFastAuthRetry(t, 20*time.Millisecond)
	a := &flakyAuth{failures: 1 << 30}
	if _, err := restrictedTagsWithRetry(a, "rest"); err == nil {
		t.Fatal("expected an error after the retry timeout")
	}
}
