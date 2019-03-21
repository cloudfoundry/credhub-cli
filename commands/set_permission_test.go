package commands_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = FDescribe("Set Permission", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete")
	ItRequiresAnAPIToBeSet("set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete")

	testAutoLogIns := []TestAutoLogin{
		{
			method:              "GET",
			responseFixtureFile: "get_permission_response.json",
			responseStatus:      http.StatusNotFound,
			endpoint:            "/api/v2/permissions",
		},
		{
			method:              "POST",
			responseFixtureFile: "set_permission_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v2/permissions",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogIns, "set-permission", "-a", "some-actor", "-p", "'/some-path'", "-o", "read, write, delete")

	Describe("Help", func() {
		ItBehavesLikeHelp("set-permission", "", func(session *Session) {
			Expect(session.Err).To(Say("set-permission"))
			Expect(session.Err).To(Say("actor"))
			Expect(session.Err).To(Say("path"))
		})
	})
})
