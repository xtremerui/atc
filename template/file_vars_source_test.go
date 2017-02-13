package template_test

import (
	. "github.com/concourse/atc/template"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileVarsSource", func() {
	var source *FileVarsSource
	BeforeEach(func() {
		source = &FileVarsSource{
			ParamsContent: []byte(`
secret:
  concourse_repo:
    private_key: some-private-key
`,
			),
		}
	})

	Describe("Evaluate", func() {
		Context("when all of the variables are defined", func() {
			It("evaluates params", func() {
				evaluatedContent, err := source.Evaluate([]byte(`
resources:
- name: my-repo
  source:
    uri: git@github.com:concourse/concourse.git
    private_key: ((secret.concourse_repo.private_key))
`,
				))
				Expect(err).NotTo(HaveOccurred())
				Expect(evaluatedContent).To(MatchYAML([]byte(`
resources:
- name: my-repo
  source:
    uri: git@github.com:concourse/concourse.git
    private_key: some-private-key
`,
				)))
			})
		})

		Context("when not all of the variables are defined", func() {
			It("returns an error", func() {
				_, err := source.Evaluate([]byte(`
resources:
- name: my-repo
  source:
    uri: git@github.com:concourse/concourse.git
    private_key: ((secret.concourse_repo.private_key))

- name: env-state
  source:
    bucket: ((env))
    key: ((state))
`,
				))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Expected to find variables: env\nstate"))
			})
		})
	})
})
