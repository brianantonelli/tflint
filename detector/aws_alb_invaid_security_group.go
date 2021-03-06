package detector

import (
	"fmt"

	"github.com/hashicorp/hcl/hcl/token"
	"github.com/wata727/tflint/issue"
	"github.com/wata727/tflint/schema"
)

type AwsALBInvalidSecurityGroupDetector struct {
	*Detector
	IssueType      string
	Target         string
	DeepCheck      bool
	securityGroups map[string]bool
}

func (d *Detector) CreateAwsALBInvalidSecurityGroupDetector() *AwsALBInvalidSecurityGroupDetector {
	return &AwsALBInvalidSecurityGroupDetector{
		Detector:       d,
		IssueType:      issue.ERROR,
		Target:         "aws_alb",
		DeepCheck:      true,
		securityGroups: map[string]bool{},
	}
}

func (d *AwsALBInvalidSecurityGroupDetector) PreProcess() {
	resp, err := d.AwsClient.DescribeSecurityGroups()
	if err != nil {
		d.Logger.Error(err)
		d.Error = true
		return
	}

	for _, securityGroup := range resp.SecurityGroups {
		d.securityGroups[*securityGroup.GroupId] = true
	}
}

func (d *AwsALBInvalidSecurityGroupDetector) Detect(resource *schema.Resource, issues *[]*issue.Issue) {
	var varToken token.Token
	var securityGroupTokens []token.Token
	var ok bool
	if varToken, ok = resource.GetToken("security_groups"); ok {
		var err error
		securityGroupTokens, err = d.evalToStringTokens(varToken)
		if err != nil {
			d.Logger.Error(err)
			return
		}
	} else {
		securityGroupTokens, ok = resource.GetListToken("security_groups")
		if !ok {
			return
		}
	}

	for _, securityGroupToken := range securityGroupTokens {
		securityGroup, err := d.evalToString(securityGroupToken.Text)
		if err != nil {
			d.Logger.Error(err)
			continue
		}

		// If `security_groups` is interpolated by list variable, Filename is empty.
		if securityGroupToken.Pos.Filename == "" {
			securityGroupToken.Pos.Filename = varToken.Pos.Filename
		}
		if !d.securityGroups[securityGroup] {
			issue := &issue.Issue{
				Type:    d.IssueType,
				Message: fmt.Sprintf("\"%s\" is invalid security group.", securityGroup),
				Line:    securityGroupToken.Pos.Line,
				File:    securityGroupToken.Pos.Filename,
			}
			*issues = append(*issues, issue)
		}
	}
}
