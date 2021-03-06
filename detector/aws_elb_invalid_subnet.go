package detector

import (
	"fmt"

	"github.com/hashicorp/hcl/hcl/token"
	"github.com/wata727/tflint/issue"
	"github.com/wata727/tflint/schema"
)

type AwsELBInvalidSubnetDetector struct {
	*Detector
	IssueType string
	Target    string
	DeepCheck bool
	subnets   map[string]bool
}

func (d *Detector) CreateAwsELBInvalidSubnetDetector() *AwsELBInvalidSubnetDetector {
	return &AwsELBInvalidSubnetDetector{
		Detector:  d,
		IssueType: issue.ERROR,
		Target:    "aws_elb",
		DeepCheck: true,
		subnets:   map[string]bool{},
	}
}

func (d *AwsELBInvalidSubnetDetector) PreProcess() {
	resp, err := d.AwsClient.DescribeSubnets()
	if err != nil {
		d.Logger.Error(err)
		d.Error = true
		return
	}

	for _, subnet := range resp.Subnets {
		d.subnets[*subnet.SubnetId] = true
	}
}

func (d *AwsELBInvalidSubnetDetector) Detect(resource *schema.Resource, issues *[]*issue.Issue) {
	var varToken token.Token
	var subnetTokens []token.Token
	var ok bool
	if varToken, ok = resource.GetToken("subnets"); ok {
		var err error
		subnetTokens, err = d.evalToStringTokens(varToken)
		if err != nil {
			d.Logger.Error(err)
			return
		}
	} else {
		subnetTokens, ok = resource.GetListToken("subnets")
		if !ok {
			return
		}
	}

	for _, subnetToken := range subnetTokens {
		subnet, err := d.evalToString(subnetToken.Text)
		if err != nil {
			d.Logger.Error(err)
			continue
		}

		// If `subnets` is interpolated by list variable, Filename is empty
		if subnetToken.Pos.Filename == "" {
			subnetToken.Pos.Filename = varToken.Pos.Filename
		}
		if !d.subnets[subnet] {
			issue := &issue.Issue{
				Type:    d.IssueType,
				Message: fmt.Sprintf("\"%s\" is invalid subnet ID.", subnet),
				Line:    subnetToken.Pos.Line,
				File:    subnetToken.Pos.Filename,
			}
			*issues = append(*issues, issue)
		}
	}
}
