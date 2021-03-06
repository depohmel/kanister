package function

import (
	. "gopkg.in/check.v1"

	crv1alpha1 "github.com/kanisterio/kanister/pkg/apis/cr/v1alpha1"
	"github.com/kanisterio/kanister/pkg/param"
)

type BackupDataSuite struct {
}

var _ = Suite(&BackupDataSuite{})

func newValidProfile() *param.Profile {
	return &param.Profile{
		Location: crv1alpha1.Location{
			Type:     crv1alpha1.LocationTypeS3Compliant,
			Bucket:   "test-bucket",
			Endpoint: "",
			Prefix:   "",
			Region:   "us-west-1",
		},
		Credential: param.Credential{
			Type: param.CredentialTypeKeyPair,
			KeyPair: &param.KeyPair{
				ID:     "test-id",
				Secret: "test-secret",
			},
		},
		SkipSSLVerify: false,
	}
}

func newInvalidProfile() *param.Profile {
	return &param.Profile{
		Location: crv1alpha1.Location{
			Type:     "foo-type",
			Bucket:   "test-bucket",
			Endpoint: "",
			Prefix:   "",
			Region:   "us-west-1",
		},
		Credential: param.Credential{
			Type: param.CredentialTypeKeyPair,
			KeyPair: &param.KeyPair{
				ID:     "test-id",
				Secret: "test-secret",
			},
		},
		SkipSSLVerify: false,
	}
}

func (s *BackupDataSuite) TestValidateProfile(c *C) {
	testCases := []struct {
		name       string
		profile    *param.Profile
		errChecker Checker
	}{
		{"Valid Profile", newValidProfile(), IsNil},
		{"Invalid Profile", newInvalidProfile(), NotNil},
		{"Nil Profile", nil, NotNil},
	}
	for _, tc := range testCases {
		err := validateProfile(tc.profile)
		c.Check(err, tc.errChecker, Commentf("Test %s Failed", tc.name))
	}
}
