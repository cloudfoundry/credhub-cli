package commands_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Delete Permission", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("delete-permission", "-a", "some-actor", "-p", "'/some-path'")
	ItRequiresAnAPIToBeSet("delete-permission", "-a", "some-actor", "-p", "'/some-path'")

	testAutoLogIns := []TestAutoLogin{
		{
			method:              "GET",
			responseFixtureFile: "set_permission_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v2/permissions",
		},
		{
			method:              "DELETE",
			responseFixtureFile: "set_permission_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v2/permissions/" + UUID,
		},
	}
	ItAutomaticallyLogsIn(testAutoLogIns, "delete-permission", "-a", "some-actor", "-p", "'/some-path'")

	Describe("Help", func() {
		ItBehavesLikeHelp("delete-permission", "", func(session *Session) {
			Expect(session.Err).To(Say("delete-permission"))
			Expect(session.Err).To(Say("actor"))
			Expect(session.Err).To(Say("path"))
		})
	})

})
