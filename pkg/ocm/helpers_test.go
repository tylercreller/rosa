/**
Copyright (c) 2023 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ocm

import (
	"fmt"
	"net/url"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ocmCommonValidations "github.com/openshift-online/ocm-common/pkg/ocm/validations"
	commonUtils "github.com/openshift-online/ocm-common/pkg/utils"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

var _ = Describe("Http tokens", func() {
	Context("Http tokens variable validations", func() {
		It("OK: Validates successfully http tokens required", func() {
			err := ValidateHttpTokensValue(string(cmv1.Ec2MetadataHttpTokensRequired))
			Expect(err).NotTo(HaveOccurred())
		})
		It("OK: Validates successfully http tokens optional", func() {
			err := ValidateHttpTokensValue(string(cmv1.Ec2MetadataHttpTokensOptional))
			Expect(err).NotTo(HaveOccurred())
		})
		It("OK: Validates successfully http tokens empty string", func() {
			err := ValidateHttpTokensValue("")
			Expect(err).NotTo(HaveOccurred())
		})
		It("Error: Validates error for http tokens bad string", func() {
			err := ValidateHttpTokensValue("dummy")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(fmt.Sprintf("ec2-metadata-http-tokens value should be one of '%s', '%s'",
				cmv1.Ec2MetadataHttpTokensRequired, cmv1.Ec2MetadataHttpTokensOptional)))
		})
	})
})

var _ = Describe("Validate Issuer Url Matches Assume Policy Document", func() {
	const (
		fakeOperatorRoleArn = "arn:aws:iam::765374464689:role/fake-arn-openshift-cluster-csi-drivers-ebs-cloud-credentials"
	)
	It("OK: Matching", func() {
		//nolint
		fakeAssumePolicyDocument := `%7B%22Version%22%3A%222012-10-17%22%2C%22Statement%22%3A%5B%7B%22Effect%22%3A%22Allow%22%2C%22Principal%22%3A%7B%22Federated%22%3A%22arn%3Aaws%3Aiam%3A%3A765374464689%3Aoidc-provider%2Ffake-oidc.s3.us-east-1.amazonaws.com%22%7D%2C%22Action%22%3A%22sts%3AAssumeRoleWithWebIdentity%22%2C%22Condition%22%3A%7B%22StringEquals%22%3A%7B%22fake.s3.us-east-1.amazonaws.com%3Asub%22%3A%5B%22system%3Aserviceaccount%3Aopenshift-image-registry%3Acluster-image-registry-operator%22%2C%22system%3Aserviceaccount%3Aopenshift-image-registry%3Aregistry%22%5D%7D%7D%7D%5D%7D`
		parsedUrl, _ := url.Parse("https://fake-oidc.s3.us-east-1.amazonaws.com")
		err := ocmCommonValidations.ValidateIssuerUrlMatchesAssumePolicyDocument(
			fakeOperatorRoleArn, parsedUrl, fakeAssumePolicyDocument)
		Expect(err).NotTo(HaveOccurred())
	})
	It("OK: Matching with path", func() {
		//nolint
		fakeAssumePolicyDocument := `%7B%22Version%22%3A%222012-10-17%22%2C%22Statement%22%3A%5B%7B%22Effect%22%3A%22Allow%22%2C%22Principal%22%3A%7B%22Federated%22%3A%22arn%3Aaws%3Aiam%3A%3A765374464689%3Aoidc-provider%2Ffake-oidc.s3.us-east-1.amazonaws.com%2F23g84jr4cdfpej0ghlr4teqiog8747gt%22%7D%2C%22Action%22%3A%22sts%3AAssumeRoleWithWebIdentity%22%2C%22Condition%22%3A%7B%22StringEquals%22%3A%7B%22fake.s3.us-east-1.amazonaws.com%2F23g84jr4cdfpej0ghlr4teqiog8747gt%3Asub%22%3A%5B%22system%3Aserviceaccount%3Aopenshift-image-registry%3Acluster-image-registry-operator%22%2C%22system%3Aserviceaccount%3Aopenshift-image-registry%3Aregistry%22%5D%7D%7D%7D%5D%7D`
		parsedUrl, _ := url.Parse("https://fake-oidc.s3.us-east-1.amazonaws.com/23g84jr4cdfpej0ghlr4teqiog8747gt")
		err := ocmCommonValidations.ValidateIssuerUrlMatchesAssumePolicyDocument(
			fakeOperatorRoleArn, parsedUrl, fakeAssumePolicyDocument)
		Expect(err).NotTo(HaveOccurred())
	})
	It("KO: Not matching", func() {
		//nolint
		fakeAssumePolicyDocument := `%7B%22Version%22%3A%222012-10-17%22%2C%22Statement%22%3A%5B%7B%22Effect%22%3A%22Allow%22%2C%22Principal%22%3A%7B%22Federated%22%3A%22arn%3Aaws%3Aiam%3A%3A765374464689%3Aoidc-provider%2Ffake-oidc.s3.us-east-1.amazonaws.com%22%7D%2C%22Action%22%3A%22sts%3AAssumeRoleWithWebIdentity%22%2C%22Condition%22%3A%7B%22StringEquals%22%3A%7B%22fake.s3.us-east-1.amazonaws.com%3Asub%22%3A%5B%22system%3Aserviceaccount%3Aopenshift-image-registry%3Acluster-image-registry-operator%22%2C%22system%3Aserviceaccount%3Aopenshift-image-registry%3Aregistry%22%5D%7D%7D%7D%5D%7D`
		fakeIssuerUrl := "https://fake-oidc-2.s3.us-east-1.amazonaws.com"
		parsedUrl, _ := url.Parse(fakeIssuerUrl)
		err := ocmCommonValidations.ValidateIssuerUrlMatchesAssumePolicyDocument(
			fakeOperatorRoleArn, parsedUrl, fakeAssumePolicyDocument)
		Expect(err).To(HaveOccurred())
		//nolint
		Expect(
			fmt.Sprintf(
				"Operator role '%s' does not have trusted relationship to '%s' issuer URL",
				fakeOperatorRoleArn,
				parsedUrl.Host,
			),
		).To(Equal(err.Error()))
	})
	It("KO: Not matching with path", func() {
		//nolint
		fakeAssumePolicyDocument := `%7B%22Version%22%3A%222012-10-17%22%2C%22Statement%22%3A%5B%7B%22Effect%22%3A%22Allow%22%2C%22Principal%22%3A%7B%22Federated%22%3A%22arn%3Aaws%3Aiam%3A%3A765374464689%3Aoidc-provider%2Ffake-oidc.s3.us-east-1.amazonaws.com%2F23g84jr4cdfpej0ghlr4teqiog8747gt%22%7D%2C%22Action%22%3A%22sts%3AAssumeRoleWithWebIdentity%22%2C%22Condition%22%3A%7B%22StringEquals%22%3A%7B%22fake.s3.us-east-1.amazonaws.com%2F23g84jr4cdfpej0ghlr4teqiog8747gt%3Asub%22%3A%5B%22system%3Aserviceaccount%3Aopenshift-image-registry%3Acluster-image-registry-operator%22%2C%22system%3Aserviceaccount%3Aopenshift-image-registry%3Aregistry%22%5D%7D%7D%7D%5D%7D`
		fakeIssuerUrl := "https://fake-oidc-2.s3.us-east-1.amazonaws.com/23g84jr4cdfpej0ghlr4teqiog8747g"
		parsedUrl, _ := url.Parse(fakeIssuerUrl)
		err := ocmCommonValidations.ValidateIssuerUrlMatchesAssumePolicyDocument(
			fakeOperatorRoleArn, parsedUrl, fakeAssumePolicyDocument)
		Expect(err).To(HaveOccurred())
		//nolint
		Expect(
			fmt.Sprintf(
				"Operator role '%s' does not have trusted relationship to '%s' issuer URL",
				fakeOperatorRoleArn,
				parsedUrl.Host+parsedUrl.Path,
			),
		).To(Equal(err.Error()))
	})
})

var _ = Describe("ParseDiskSizeToGigibyte", func() {
	It("returns an error for invalid unit: 1foo", func() {
		size := "1foo"
		_, err := ParseDiskSizeToGigibyte(size)
		Expect(err).To(HaveOccurred())
	})

	It("returns 0 for valid unit: 0", func() {
		size := "0"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(0))
	})

	It("returns 0 for invalid unit no suffix: 1 but return 0", func() {
		size := "0"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(0))
	})

	It("returns an error for invalid unit: 1K", func() {
		size := "1K"
		_, err := ParseDiskSizeToGigibyte(size)
		Expect(err).To(HaveOccurred())
	})

	It("returns an error for invalid unit: 1KiB", func() {
		size := "1KiB"
		_, err := ParseDiskSizeToGigibyte(size)
		Expect(err).To(HaveOccurred())
	})

	It("returns an error for invalid unit: 1 MiB", func() {
		size := "1 MiB"
		_, err := ParseDiskSizeToGigibyte(size)
		Expect(err).To(HaveOccurred())
	})

	It("returns an error for invalid unit: 1 mib", func() {
		size := "1 mib"
		_, err := ParseDiskSizeToGigibyte(size)
		Expect(err).To(HaveOccurred())
	})

	It("returns 0 for invalid unit: 0 GiB", func() {
		size := "0 GiB"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(0))
	})

	It("returns the correct value for valid unit: 100 G", func() {
		size := "100 G"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(93))
	})

	It("returns the correct value for valid unit: 100GB", func() {
		size := "100GB"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(93))
	})

	It("returns the correct value for valid unit: 100Gb", func() {
		size := "100Gb"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(93))
	})

	It("returns the correct value for valid unit: 100g", func() {
		size := "100g"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(93))
	})

	It("returns the correct value for valid unit: 100GiB", func() {
		size := "100GiB"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(100))
	})

	//
	It("returns the correct value for valid unit: 100gib", func() {
		size := "100gib"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(100))
	})

	It("returns the correct value for valid unit: 100 gib", func() {
		size := "100 gib"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(100))
	})

	It("returns the correct value for valid unit: 100 TB", func() {
		size := "100 TB"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(93132))
	})

	It("returns the correct value for valid unit: 100 T ", func() {
		size := "100 T "
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(93132))
	})

	It("returns the correct value for valid unit: 1000 Ti", func() {
		size := "1000 Ti"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(1024000))
	})

	It("returns the correct value for valid unit: empty string", func() {
		size := ""
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).NotTo(HaveOccurred())
		Expect(got).To(Equal(0))
	})

	It("returns the correct value for valid unit: -1", func() {
		size := "-1"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).To(HaveOccurred())
		Expect(got).To(Equal(0))
	})

	It("returns the correct value for valid unit: 200000000000000 Ti", func() {
		// Hitting the max int64 value
		size := "200000000000000 Ti"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).To(HaveOccurred())
		Expect(got).To(Equal(0))
	})

	It("returns the correct value for valid unit: 200000000000000000000Ti", func() {
		// Hitting the max int64 value
		size := "200000000000000000000Ti"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).To(HaveOccurred())
		Expect(got).To(Equal(0))
	})

	It("returns the correct value for valid unit: -200000000000000000000Ti", func() {
		// Hitting the max int64 value
		size := "-200000000000000000000Ti"
		got, err := ParseDiskSizeToGigibyte(size)
		Expect(err).To(HaveOccurred())
		Expect(got).To(Equal(0))
	})

})

var _ = Describe("ValidateBalancingIgnoredLabels", func() {
	It("returns an error if didn't got a string", func() {
		var val interface{} = 1
		err := ValidateBalancingIgnoredLabels(val)
		Expect(err).To(HaveOccurred())
	})

	It("passes for an empty string", func() {
		var val interface{} = ""
		err := ValidateBalancingIgnoredLabels(val)
		Expect(err).ToNot(HaveOccurred())
	})

	It("passes for valid label keys", func() {
		var val interface{} = "eks.amazonaws.com/nodegroup,alpha.eksctl.io/nodegroup-name"
		err := ValidateBalancingIgnoredLabels(val)
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns an error for a label that doesn't start with an alphanumeric character", func() {
		var val interface{} = ".t"
		err := ValidateBalancingIgnoredLabels(val)
		Expect(err).To(HaveOccurred())
	})

	It("returns an error for a label that has illegal characters", func() {
		var val interface{} = "a%"
		err := ValidateBalancingIgnoredLabels(val)
		Expect(err).To(HaveOccurred())
	})

	It("returns an error for a label that exceeds 63 characters", func() {
		var val interface{} = strings.Repeat("a", commonUtils.MaxByteSize)
		err := ValidateBalancingIgnoredLabels(val)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("expectedSubnetsCount", func() {
	When("multiAZ and privateLink are true", func() {
		It("Should return privateLinkMultiAZSubnetsCount", func() {
			Expect(expectedSubnetsCount(true, true)).To(Equal(privateLinkMultiAZSubnetsCount))
		})
	})

	When("multiAZ is true and privateLink is false", func() {
		It("Should return privateLinkSingleAZSubnetsCount", func() {
			Expect(expectedSubnetsCount(true, false)).To(Equal(BYOVPCMultiAZSubnetsCount))
		})
	})

	When("multiAZ is false and privateLink is true", func() {
		It("Should return BYOVPCMultiAZSubnetsCount", func() {
			Expect(expectedSubnetsCount(false, true)).To(Equal(privateLinkSingleAZSubnetsCount))
		})
	})

	When("multiAZ and privateLink are false", func() {
		It("Should return BYOVPCSingleAZSubnetsCount", func() {
			Expect(expectedSubnetsCount(false, false)).To(Equal(BYOVPCSingleAZSubnetsCount))
		})
	})
})

var _ = Describe("ValidateSubnetsCount", func() {
	When("When privateLink is true", func() {
		When("multiAZ is true", func() {
			It("should return an error if subnetsInputCount is not equal to privateLinkMultiAZSubnetsCount", func() {
				err := ValidateSubnetsCount(true, true, privateLinkMultiAZSubnetsCount+1)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(fmt.Sprintf("The number of subnets for a 'multi-AZ' 'private link cluster' should be"+
					" '%d', instead received: '%d'", privateLinkMultiAZSubnetsCount, privateLinkMultiAZSubnetsCount+1)))
			})

			It("should not return an error if subnetsInputCount is equal to privateLinkMultiAZSubnetsCount", func() {
				err := ValidateSubnetsCount(true, true, privateLinkMultiAZSubnetsCount)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("multiAZ is false", func() {
			It("should return an error if subnetsInputCount is not equal to privateLinkSingleAZSubnetsCount", func() {
				err := ValidateSubnetsCount(false, true, privateLinkSingleAZSubnetsCount+1)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(fmt.Sprintf("The number of subnets for a 'single AZ' 'private link cluster' should be"+
					" '%d', instead received: '%d'", privateLinkSingleAZSubnetsCount, privateLinkSingleAZSubnetsCount+1)))
			})

			It("should not return an error if subnetsInputCount is equal to privateLinkSingleAZSubnetsCount", func() {
				err := ValidateSubnetsCount(false, true, privateLinkSingleAZSubnetsCount)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	When("privateLink is false", func() {
		When("multiAZ is true", func() {
			It("should return an error if subnetsInputCount is not equal to BYOVPCMultiAZSubnetsCount", func() {
				err := ValidateSubnetsCount(true, false, BYOVPCMultiAZSubnetsCount+1)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(fmt.Sprintf("The number of subnets for a 'multi-AZ' 'cluster' should be"+
					" '%d', instead received: '%d'", BYOVPCMultiAZSubnetsCount, BYOVPCMultiAZSubnetsCount+1)))
			})

			It("should not return an error if subnetsInputCount is equal to BYOVPCMultiAZSubnetsCount", func() {
				err := ValidateSubnetsCount(true, false, BYOVPCMultiAZSubnetsCount)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("multiAZ is false", func() {
			It("should return an error if subnetsInputCount is not equal to BYOVPCSingleAZSubnetsCount", func() {
				err := ValidateSubnetsCount(false, false, BYOVPCSingleAZSubnetsCount+1)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(fmt.Sprintf("The number of subnets for a 'single AZ' 'cluster' should"+
					" be '%d', instead received: '%d'", BYOVPCSingleAZSubnetsCount, BYOVPCSingleAZSubnetsCount+1)))
			})

			It("should not return an error if subnetsInputCount is equal to BYOVPCSingleAZSubnetsCount", func() {
				err := ValidateSubnetsCount(false, false, BYOVPCSingleAZSubnetsCount)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
