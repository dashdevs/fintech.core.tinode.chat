package main

import (
	"strings"
	"testing"
)

func TestPacketSummaryOmitsContent(t *testing.T) {
	cases := []struct {
		msg    *ClientComMessage
		want   string
		secret string
	}{
		{
			msg:    &ClientComMessage{Login: &MsgClientLogin{Id: "1", Scheme: "rest", Secret: []byte("access-token")}},
			want:   "{login:rest} id='1' topic=''",
			secret: "access-token",
		},
		{
			msg:    &ClientComMessage{Pub: &MsgClientPub{Id: "2", Topic: "grpABC", Content: "message text"}},
			want:   "{pub} id='2' topic='grpABC'",
			secret: "message text",
		},
		{
			msg:  &ClientComMessage{Note: &MsgClientNote{Topic: "grpABC", What: "kp"}},
			want: "{note} id='' topic='grpABC'",
		},
		{
			msg:  &ClientComMessage{},
			want: "{unknown} id='' topic=''",
		},
	}
	for _, tc := range cases {
		got := packetSummary(tc.msg)
		if got != tc.want {
			t.Errorf("packetSummary() = %q, want %q", got, tc.want)
		}
		if tc.secret != "" && strings.Contains(got, tc.secret) {
			t.Errorf("packetSummary() = %q leaks %q", got, tc.secret)
		}
	}
}
