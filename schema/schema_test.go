package schema

import (
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/hcl/hcl/token"
	"github.com/k0kubun/pp"
)

func TestMake(t *testing.T) {
	cases := []struct {
		Name   string
		Input  map[string]string
		Result []*Template
	}{
		{
			Name: "make templates",
			Input: map[string]string{
				"test.tf": `
resource "aws_instance" "web" {
  ami           = "ami-b73b63a0"
  instance_type = "t1.2xlarge" # invalid type!

  tags {
    Name = "HelloWorld"
  }
}

resource "aws_instance" "web2" {
  security_groups = ["sg-1", "sg-2"]
}

resource "aws_instance2" "web" {
  root_block_device = {
    volume_size = "24"
  }
}
`,
			},
			Result: []*Template{
				{
					File: "test.tf",
					Resources: []*Resource{
						{
							File: "test.tf",
							Type: "aws_instance2",
							Id:   "web",
							Pos: token.Pos{
								Filename: "test.tf",
								Offset:   258,
								Line:     15,
								Column:   32,
							},
							Attrs: map[string]*Attribute{
								"root_block_device": {
									Poses: []token.Pos{
										{
											Filename: "test.tf",
											Offset:   282,
											Line:     16,
											Column:   23,
										},
									},
									Vals: []interface{}{
										map[string]interface{}{
											"volume_size": token.Token{
												Type: 9,
												Pos: token.Pos{
													Filename: "test.tf",
													Offset:   302,
													Line:     17,
													Column:   19,
												},
												Text: "\"24\"",
												JSON: false,
											},
										},
									},
								},
							},
						},
						{
							File: "test.tf",
							Type: "aws_instance",
							Id:   "web",
							Pos: token.Pos{
								Filename: "test.tf",
								Offset:   31,
								Line:     2,
								Column:   31,
							},
							Attrs: map[string]*Attribute{
								"ami": {
									Poses: []token.Pos{
										{
											Filename: "test.tf",
											Offset:   51,
											Line:     3,
											Column:   19,
										},
									},
									Vals: []interface{}{
										token.Token{
											Type: 9,
											Pos: token.Pos{
												Filename: "test.tf",
												Offset:   51,
												Line:     3,
												Column:   19,
											},
											Text: "\"ami-b73b63a0\"",
											JSON: false,
										},
									},
								},
								"instance_type": {
									Poses: []token.Pos{
										{
											Filename: "test.tf",
											Offset:   84,
											Line:     4,
											Column:   19,
										},
									},
									Vals: []interface{}{
										token.Token{
											Type: 9,
											Pos: token.Pos{
												Filename: "test.tf",
												Offset:   84,
												Line:     4,
												Column:   19,
											},
											Text: "\"t1.2xlarge\"",
											JSON: false,
										},
									},
								},
								"tags": {
									Poses: []token.Pos{
										{
											Filename: "test.tf",
											Offset:   121,
											Line:     6,
											Column:   8,
										},
									},
									Vals: []interface{}{
										map[string]interface{}{
											"Name": token.Token{
												Type: 9,
												Pos: token.Pos{
													Filename: "test.tf",
													Offset:   134,
													Line:     7,
													Column:   12,
												},
												Text: "\"HelloWorld\"",
												JSON: false,
											},
										},
									},
								},
							},
						},
						{
							File: "test.tf",
							Type: "aws_instance",
							Id:   "web2",
							Pos: token.Pos{
								Filename: "test.tf",
								Offset:   185,
								Line:     11,
								Column:   32,
							},
							Attrs: map[string]*Attribute{
								"security_groups": {
									Poses: []token.Pos{
										{
											Filename: "test.tf",
											Offset:   207,
											Line:     12,
											Column:   21,
										},
									},
									Vals: []interface{}{
										[]interface{}{
											token.Token{
												Type: 9,
												Pos: token.Pos{
													Filename: "test.tf",
													Offset:   208,
													Line:     12,
													Column:   22,
												},
												Text: "\"sg-1\"",
												JSON: false,
											},
											token.Token{
												Type: 9,
												Pos: token.Pos{
													Filename: "test.tf",
													Offset:   216,
													Line:     12,
													Column:   30,
												},
												Text: "\"sg-2\"",
												JSON: false,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "override template",
			Input: map[string]string{
				"test.tf": `
resource "aws_instance" "web" {
  ami           = "ami-b73b63a0"
  instance_type = "t1.2xlarge" # invalid type!

  tags {
    Name = "HelloWorld"
  }
}
`,
				"test_override.tf": `
resource "aws_instance" "web" {
  ami           = "ami-override"
  instance_type = "t2.nano"

  tags {
    Version = "0.1"
  }
}
`,
				"override.tf": `
resource "aws_instance" "web" {
  instance_type = "t2.micro"
}
`,
			},
			Result: []*Template{
				{
					File: "test.tf",
					Resources: []*Resource{
						{
							File: "test.tf",
							Type: "aws_instance",
							Id:   "web",
							Pos: token.Pos{
								Filename: "test.tf",
								Offset:   31,
								Line:     2,
								Column:   31,
							},
							Attrs: map[string]*Attribute{
								"ami": {
									Poses: []token.Pos{
										{
											Filename: "test_override.tf",
											Offset:   51,
											Line:     3,
											Column:   19,
										},
									},
									Vals: []interface{}{
										token.Token{
											Type: 9,
											Pos: token.Pos{
												Filename: "test_override.tf",
												Offset:   51,
												Line:     3,
												Column:   19,
											},
											Text: "\"ami-override\"",
											JSON: false,
										},
									},
								},
								"instance_type": {
									Poses: []token.Pos{
										{
											Filename: "test_override.tf",
											Offset:   84,
											Line:     4,
											Column:   19,
										},
									},
									Vals: []interface{}{
										token.Token{
											Type: 9,
											Pos: token.Pos{
												Filename: "test_override.tf",
												Offset:   84,
												Line:     4,
												Column:   19,
											},
											Text: "\"t2.nano\"",
											JSON: false,
										},
									},
								},
								"tags": {
									Poses: []token.Pos{
										{
											Filename: "test_override.tf",
											Offset:   102,
											Line:     6,
											Column:   8,
										},
									},
									Vals: []interface{}{
										map[string]interface{}{
											"Version": token.Token{
												Type: 9,
												Pos: token.Pos{
													Filename: "test_override.tf",
													Offset:   118,
													Line:     7,
													Column:   15,
												},
												Text: "\"0.1\"",
												JSON: false,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		files := map[string][]byte{}
		for filename, body := range tc.Input {
			files[filename] = []byte(body)
		}
		templates, err := Make(files)
		if err != nil {
			t.Fatal(err)
		}

		for _, template := range templates {
			sort.Slice(template.Resources, func(i, j int) bool {
				return template.Resources[i].Type+template.Resources[i].Id < template.Resources[j].Type+template.Resources[j].Id
			})
		}

		if !reflect.DeepEqual(templates, tc.Result) {
			t.Fatalf("\nBad: %s\nExpected: %s\n\ntestcase: %s", pp.Sprint(templates), pp.Sprint(tc.Result), tc.Name)
		}
	}
}
