package server_test

import (
	"sacco/server"
	"sacco/wscli"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"server": server.Main,
		"wscli":  wscli.Main,
	})
}

func TestMembershipApplication(t *testing.T) {
	t.Skip()
	
	testscript.Run(t, testscript.Params{
		Dir: "testdata/membershipApplication",
	})
}
