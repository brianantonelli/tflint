package detector

import (
	"reflect"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/wata727/tflint/config"
	"github.com/wata727/tflint/issue"
)

func TestDetectAwsInstancePreviousType(t *testing.T) {
	cases := []struct {
		Name   string
		Src    string
		Issues []*issue.Issue
	}{
		{
			Name: "t1.micro is previous type",
			Src: `
resource "aws_instance" "web" {
    instance_type = "t1.micro"
}`,
			Issues: []*issue.Issue{
				{
					Type:    "WARNING",
					Message: "\"t1.micro\" is previous generation instance type.",
					Line:    3,
					File:    "test.tf",
				},
			},
		},
		{
			Name: "t2.micro is not previous type",
			Src: `
resource "aws_instance" "web" {
    instance_type = "t2.micro"
}`,
			Issues: []*issue.Issue{},
		},
	}

	for _, tc := range cases {
		var issues = []*issue.Issue{}
		err := TestDetectByCreatorName(
			"CreateAwsInstancePreviousTypeDetector",
			tc.Src,
			"",
			config.Init(),
			config.Init().NewAwsClient(),
			&issues,
		)
		if err != nil {
			t.Fatalf("\nERROR: %s", err)
		}

		if !reflect.DeepEqual(issues, tc.Issues) {
			t.Fatalf("\nBad: %s\nExpected: %s\n\ntestcase: %s", pp.Sprint(issues), pp.Sprint(tc.Issues), tc.Name)
		}
	}
}
