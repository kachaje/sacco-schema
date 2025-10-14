package tests

import (
	"os"
	"testing"

	"github.com/kachaje/sacco-schema/server"
	"github.com/kachaje/sacco-schema/wscli"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"server": server.Main,
		"wscli":  wscli.Main,
	})
}

func TestMembershipApplication(t *testing.T) {
	if os.Getenv("UITESTS") == "true" {
		testscript.Run(t, testscript.Params{
			Dir: "testdata/membershipApplication",
		})
	}
}
