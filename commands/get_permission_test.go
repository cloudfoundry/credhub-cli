package commands_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = FDescribe("Get Permission", func() {
	BeforeEach(func() {
		login()
	})

	ItRequiresAuthentication("get-permission", "-a", "some-actor", "-p", "'/some-path'")
	ItRequiresAnAPIToBeSet("get-permission", "-a", "some-actor", "-p", "'/some-path'")

	testAutoLogIns := []TestAutoLogin{
		{
			method:              "GET",
			responseFixtureFile: "set_permission_response.json",
			responseStatus:      http.StatusOK,
			endpoint:            "/api/v2/permissions",
		},
	}
	ItAutomaticallyLogsIn(testAutoLogIns, "get-permission", "-a", "some-actor", "-p", "'/some-path'")

	Context("when help flag is used", func() {
		ItBehavesLikeHelp("get-permission", "", func(session *Session) {
			Expect(session.Err).To(Say("get-permission"))
			Expect(session.Err).To(Say("actor"))
			Expect(session.Err).To(Say("path"))
		})
	})
})
