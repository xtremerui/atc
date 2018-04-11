package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/concourse/atc"
	"github.com/concourse/atc/api/accessor/accessorfakes"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/db/dbfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func jsonEncode(object interface{}) *bytes.Buffer {
	reqPayload, err := json.Marshal(object)
	Expect(err).NotTo(HaveOccurred())

	return bytes.NewBuffer(reqPayload)
}

var _ = Describe("Teams API", func() {
	var (
		fakeTeam   *dbfakes.FakeTeam
		fakeaccess *accessorfakes.FakeAccess
	)

	BeforeEach(func() {
		fakeTeam = new(dbfakes.FakeTeam)
		fakeaccess = new(accessorfakes.FakeAccess)
	})

	JustBeforeEach(func() {
		fakeAccessor.CreateReturns(fakeaccess)
	})

	Describe("GET /api/v1/teams", func() {
		var response *http.Response

		JustBeforeEach(func() {
			path := fmt.Sprintf("%s/api/v1/teams", server.URL)

			request, err := http.NewRequest("GET", path, nil)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the database returns an error", func() {
			var disaster error

			BeforeEach(func() {
				disaster = errors.New("some error")
				dbTeamFactory.GetTeamsReturns(nil, disaster)
			})

			It("returns 500 Internal Server Error", func() {
				Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})

		Context("when the database returns teams", func() {
			var (
				fakeTeamOne   *dbfakes.FakeTeam
				fakeTeamTwo   *dbfakes.FakeTeam
				fakeTeamThree *dbfakes.FakeTeam
			)
			BeforeEach(func() {
				fakeTeamOne = new(dbfakes.FakeTeam)
				fakeTeamTwo = new(dbfakes.FakeTeam)
				fakeTeamThree = new(dbfakes.FakeTeam)

				fakeTeamOne.IDReturns(5)
				fakeTeamOne.NameReturns("avengers")

				fakeTeamTwo.IDReturns(9)
				fakeTeamTwo.NameReturns("aliens")
				fakeTeamTwo.AuthReturns(map[string][]string{
					"groups": []string{"github:org:team"},
				})

				fakeTeamThree.IDReturns(22)
				fakeTeamThree.NameReturns("predators")
				fakeTeamThree.AuthReturns(map[string][]string{
					"users": []string{"local:username"},
				})

				dbTeamFactory.GetTeamsReturns([]db.Team{fakeTeamOne, fakeTeamTwo, fakeTeamThree}, nil)
			})

			It("returns 200 OK", func() {
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})

			It("returns the teams", func() {
				body, err := ioutil.ReadAll(response.Body)
				Expect(err).NotTo(HaveOccurred())

				Expect(body).To(MatchJSON(`[
 					{
 						"id": 5,
 						"name": "avengers"
 					},
 					{
 						"id": 9,
 						"name": "aliens"
 					},
 					{
 						"id": 22,
 						"name": "predators"
 					}
 				]`))
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name", func() {
		var (
			response *http.Response
			atcTeam  atc.Team
		)

		BeforeEach(func() {
			fakeTeam.IDReturns(5)
			fakeTeam.NameReturns("some-team")

			atcTeam = atc.Team{}
		})

		JustBeforeEach(func() {
			path := fmt.Sprintf("%s/api/v1/teams/some-team", server.URL)

			var err error
			request, err := http.NewRequest("PUT", path, jsonEncode(atcTeam))
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		authorizedTeamTests := func() {
			Context("when the team exists", func() {
				BeforeEach(func() {
					atcTeam = atc.Team{
						Auth: map[string][]string{
							"users": []string{"local:username"},
						},
					}
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				})

				It("updates provider auth", func() {
					Expect(response.StatusCode).To(Equal(http.StatusOK))
					Expect(fakeTeam.UpdateProviderAuthCallCount()).To(Equal(1))

					updatedProviderAuth := fakeTeam.UpdateProviderAuthArgsForCall(0)
					Expect(updatedProviderAuth).To(Equal(atcTeam.Auth))
				})

				Context("when updating provider auth fails", func() {
					BeforeEach(func() {
						fakeTeam.UpdateProviderAuthReturns(errors.New("stop trying to make fetch happen"))
					})

					It("returns 500 Internal Server error", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})
			})
		}

		Context("when the requester team is authorized as an admin team", func() {
			BeforeEach(func() {
				fakeaccess.IsAuthenticatedReturns(true)
				fakeaccess.IsAuthenticatedReturns(true)
				fakeaccess.IsAdminReturns(true)
			})

			authorizedTeamTests()

			Context("when the team is not found", func() {
				BeforeEach(func() {
					dbTeamFactory.FindTeamReturns(nil, false, nil)
					dbTeamFactory.CreateTeamReturns(fakeTeam, nil)
				})

				It("creates the team", func() {
					Expect(response.StatusCode).To(Equal(http.StatusCreated))
					Expect(dbTeamFactory.CreateTeamCallCount()).To(Equal(1))

					createdTeam := dbTeamFactory.CreateTeamArgsForCall(0)
					Expect(createdTeam).To(Equal(atc.Team{
						Name: "some-team",
					}))
				})

				Context("when it fails to create team", func() {
					BeforeEach(func() {
						dbTeamFactory.CreateTeamReturns(nil, errors.New("it is never going to happen"))
					})

					It("returns a 500 Internal Server error", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})
			})
		})

		Context("when the requester team is authorized as the team being set", func() {
			BeforeEach(func() {
				fakeaccess.IsAuthenticatedReturns(true)
				fakeaccess.IsAuthorizedReturns(true)
			})

			authorizedTeamTests()

			Context("when the team is not found", func() {
				BeforeEach(func() {
					dbTeamFactory.FindTeamReturns(nil, false, nil)
					dbTeamFactory.CreateTeamReturns(fakeTeam, nil)
				})

				It("does not create the team", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
					Expect(dbTeamFactory.CreateTeamCallCount()).To(Equal(0))
				})
			})
		})
	})

	Describe("DELETE /api/v1/teams/:team_name", func() {
		var request *http.Request
		var response *http.Response

		var teamName string

		BeforeEach(func() {
			teamName = "team venture"

			fakeTeam.IDReturns(2)
			fakeTeam.NameReturns(teamName)
		})

		Context("when the requester is authenticated for some admin team", func() {
			JustBeforeEach(func() {
				path := fmt.Sprintf("%s/api/v1/teams/%s", server.URL, teamName)

				var err error
				request, err = http.NewRequest("DELETE", path, nil)
				Expect(err).NotTo(HaveOccurred())

				response, err = client.Do(request)
				Expect(err).NotTo(HaveOccurred())
			})

			BeforeEach(func() {
				fakeaccess.IsAuthenticatedReturns(true)
				fakeaccess.IsAdminReturns(true)
			})

			Context("when there's a problem finding teams", func() {
				BeforeEach(func() {
					dbTeamFactory.FindTeamReturns(nil, false, errors.New("a dingo ate my baby!"))
				})

				It("returns 500 Internal Server Error", func() {
					Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})

			Context("when team exists", func() {
				BeforeEach(func() {
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				})

				It("returns 204 No Content", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNoContent))
				})

				It("receives the correct team name", func() {
					Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
					Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal(teamName))
				})
				It("deletes the team from the DB", func() {
					Expect(fakeTeam.DeleteCallCount()).To(Equal(1))
					//TODO delete the build events via a table drop rather
				})

				Context("when trying to delete the admin team", func() {
					BeforeEach(func() {
						teamName = atc.DefaultTeamName
						fakeTeam.AdminReturns(true)
						dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
						dbTeamFactory.GetTeamsReturns([]db.Team{fakeTeam}, nil)
					})

					It("returns 403 Forbidden and backs off", func() {
						Expect(response.StatusCode).To(Equal(http.StatusForbidden))
						Expect(fakeTeam.DeleteCallCount()).To(Equal(0))
					})
				})

				Context("when there's a problem deleting the team", func() {
					BeforeEach(func() {
						fakeTeam.DeleteReturns(errors.New("disaster"))
					})

					It("returns 500 Internal Server Error", func() {
						Expect(response.StatusCode).To(Equal(http.StatusInternalServerError))
					})
				})
			})

			Context("when team does not exist", func() {
				BeforeEach(func() {
					dbTeamFactory.FindTeamReturns(nil, false, nil)
				})

				It("returns 404 Not Found", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNotFound))
				})
			})
		})

		Context("when the requester belongs to a non-admin team", func() {
			JustBeforeEach(func() {
				path := fmt.Sprintf("%s/api/v1/teams/%s", server.URL, "non-admin-team")

				var err error
				request, err = http.NewRequest("DELETE", path, nil)
				Expect(err).NotTo(HaveOccurred())

				response, err = client.Do(request)
				Expect(err).NotTo(HaveOccurred())

			})

			BeforeEach(func() {
				fakeaccess.IsAuthenticatedReturns(true)
				fakeaccess.IsAdminReturns(false)
			})

			It("returns 403 forbidden", func() {
				Expect(response.StatusCode).To(Equal(http.StatusForbidden))
			})
		})
	})

	Describe("PUT /api/v1/teams/:team_name/rename", func() {
		var response *http.Response
		var teamName string

		JustBeforeEach(func() {
			request, err := http.NewRequest(
				"PUT",
				server.URL+"/api/v1/teams/"+teamName+"/rename",
				bytes.NewBufferString(`{"name":"some-new-name"}`),
			)
			Expect(err).NotTo(HaveOccurred())

			response, err = client.Do(request)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			fakeTeam.IDReturns(2)
		})

		Context("when authenticated", func() {
			BeforeEach(func() {
				fakeaccess.IsAuthenticatedReturns(true)
			})
			Context("when requester belongs to an admin team", func() {
				BeforeEach(func() {
					teamName = "a-team"
					fakeTeam.NameReturns(teamName)
					fakeaccess.IsAdminReturns(true)
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				})

				It("constructs teamDB with provided team name", func() {
					Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
					Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("a-team"))
				})

				It("renames the team to the name provided", func() {
					Expect(fakeTeam.RenameCallCount()).To(Equal(1))
					Expect(fakeTeam.RenameArgsForCall(0)).To(Equal("some-new-name"))
				})

				It("returns 204 no content", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNoContent))
				})
			})

			Context("when requester belongs to the team", func() {
				BeforeEach(func() {
					teamName = "a-team"
					fakeTeam.NameReturns(teamName)
					fakeaccess.IsAuthorizedReturns(true)
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				})

				It("constructs teamDB with provided team name", func() {
					Expect(dbTeamFactory.FindTeamCallCount()).To(Equal(1))
					Expect(dbTeamFactory.FindTeamArgsForCall(0)).To(Equal("a-team"))
				})

				It("renames the team to the name provided", func() {
					Expect(fakeTeam.RenameCallCount()).To(Equal(1))
					Expect(fakeTeam.RenameArgsForCall(0)).To(Equal("some-new-name"))
				})

				It("returns 204 no content", func() {
					Expect(response.StatusCode).To(Equal(http.StatusNoContent))
				})
			})

			Context("when requester does not belong to the team", func() {
				BeforeEach(func() {
					teamName = "a-team"
					fakeTeam.NameReturns(teamName)
					fakeaccess.IsAuthorizedReturns(false)
					dbTeamFactory.FindTeamReturns(fakeTeam, true, nil)
				})

				It("returns 403 Forbidden", func() {
					Expect(response.StatusCode).To(Equal(http.StatusForbidden))
					Expect(fakeTeam.RenameCallCount()).To(Equal(0))
				})
			})
		})

		Context("when not authenticated", func() {
			BeforeEach(func() {
				fakeaccess.IsAuthenticatedReturns(false)
			})

			It("returns 401 Unauthorized", func() {
				Expect(response.StatusCode).To(Equal(http.StatusUnauthorized))
				Expect(fakeTeam.RenameCallCount()).To(Equal(0))
			})
		})
	})
})
