package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "default certificate with minimum protocol version",
			Content: `
resource "aws_cloudfront_distribution" "this" {
  enabled = true

  origin {
    domain_name = "example-bucket.s3.amazonaws.com"
    origin_id   = "s3-origin"
  }

  default_cache_behavior {
    target_origin_id       = "s3-origin"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
    minimum_protocol_version       = "TLSv1.2_2021"
  }
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewAwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule(),
					Message: "`viewer_certificate.minimum_protocol_version` cannot be set when `viewer_certificate.cloudfront_default_certificate` is `true`.",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 33, Column: 38},
						End:      hcl.Pos{Line: 33, Column: 52},
					},
				},
			},
		},
		{
			Name: "custom certificate with minimum protocol version",
			Content: `
resource "aws_cloudfront_distribution" "this" {
  enabled = true

  origin {
    domain_name = "example-bucket.s3.amazonaws.com"
    origin_id   = "s3-origin"
  }

  default_cache_behavior {
    target_origin_id       = "s3-origin"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate.this.arn
    cloudfront_default_certificate = false
    minimum_protocol_version = "TLSv1.2_2021"
    ssl_support_method       = "sni-only"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "default certificate without minimum protocol version",
			Content: `
resource "aws_cloudfront_distribution" "this" {
  enabled = true

  origin {
    domain_name = "example-bucket.s3.amazonaws.com"
    origin_id   = "s3-origin"
  }

  default_cache_behavior {
    target_origin_id       = "s3-origin"
    viewer_protocol_policy = "allow-all"
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "unknown default certificate value",
			Content: `
variable "use_default_certificate" {
  type = bool
}

resource "aws_cloudfront_distribution" "this" {
  enabled = true

  origin {
    domain_name = "example-bucket.s3.amazonaws.com"
    origin_id   = "s3-origin"
  }

  default_cache_behavior {
    target_origin_id       = "s3-origin"
    viewer_protocol_policy = "allow-all"
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = var.use_default_certificate
    minimum_protocol_version       = "TLSv1.2_2021"
  }
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewAwsCloudFrontDistributionMinimumProtocolVersionDefaultCertificateRule()

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"resource.tf": tc.Content})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("Unexpected error occurred: %s", err)
		}

		helper.AssertIssues(t, tc.Expected, runner.Issues)
	}
}
