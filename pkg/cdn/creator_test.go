package cdn

import (
	"context"
	"os"
	"testing"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/SlootSantos/janus-server/pkg/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	gomock "github.com/golang/mock/gomock"
)

func TestCDN_Create(t *testing.T) {
	t.Run("should create cdn", func(t *testing.T) {
		os.Setenv("DOMAIN_HOST", "test")
		ctrl := gomock.NewController(t)
		cdnMock := NewMockcdnandler(ctrl)
		dnsMock := NewMockdnshandler(ctrl)
		acmMock := NewMockcertificateHandler(ctrl)
		sqsMock := queue.NewMockqueueHandler(ctrl)
		qMock := queue.NewMockQ(sqsMock)

		c := CDN{
			cdn:   cdnMock,
			dns:   dnsMock,
			acm:   acmMock,
			queue: &qMock,
			config: &cdnConfig{
				domain:       "test.com",
				certARN:      "arn:cert:1234",
				hostedZoneID: "/hostedzone/12345",
			},
		}

		output := &jam.OutputParam{}
		input := &jam.CreationParam{
			CDN:    jam.StackCDN{Subdomain: "subdomain"},
			ID:     "ID_12345",
			Bucket: struct{ ID string }{"ID_12345"},
		}

		expectedCallParam := c.constructStandardDistroConfig(&constructDistroConfigInput{
			bucketID:       "ID_12345",
			subdomain:      "subdomain",
			stackID:        "ID_12345",
			originAccessID: "ABCDEFG",
		})
		returnCreateOrigin := &cloudfront.CreateCloudFrontOriginAccessIdentityOutput{
			CloudFrontOriginAccessIdentity: &cloudfront.OriginAccessIdentity{
				Id: aws.String("ABCDEFG"),
			},
		}
		returnCreateDistro := &cloudfront.CreateDistributionOutput{
			Distribution: &cloudfront.Distribution{
				DomainName: aws.String("xy.com"),
				Id:         aws.String("6778ghj"),
			},
		}

		expectedCertificateParam := &acm.RequestCertificateInput{
			DomainName:       aws.String("subdomain.test.com"),
			ValidationMethod: aws.String("DNS"),
			SubjectAlternativeNames: []*string{
				aws.String("*." + "subdomain.test.com"),
				aws.String("*.pr." + "subdomain.test.com"),
			},
		}
		returnCertificateParam := &acm.RequestCertificateOutput{
			CertificateArn: aws.String("arn:12345"),
		}
		returnCertifcateDescribeParam := &acm.DescribeCertificateOutput{
			Certificate: &acm.CertificateDetail{
				DomainValidationOptions: []*acm.DomainValidation{
					&acm.DomainValidation{
						ResourceRecord: &acm.ResourceRecord{
							Name:  aws.String("_name.aws.cert.com"),
							Type:  aws.String("CNAME"),
							Value: aws.String("_value.aws.cert.com"),
						},
					},
				},
			},
		}

		cdnMock.EXPECT().CreateCloudFrontOriginAccessIdentity(gomock.Any()).Times(1).Return(returnCreateOrigin, nil)
		cdnMock.EXPECT().CreateDistribution(expectedCallParam).Times(1).Return(returnCreateDistro, nil)
		dnsMock.EXPECT().ChangeResourceRecordSets(gomock.Any()).Times(1)
		acmMock.EXPECT().RequestCertificate(expectedCertificateParam).Times(1).Return(returnCertificateParam, nil)
		acmMock.EXPECT().DescribeCertificate(gomock.Any()).Return(returnCertifcateDescribeParam, nil)
		dnsMock.EXPECT().ChangeResourceRecordSets(gomock.Any()).Times(1)
		sqsMock.EXPECT().SendMessage(gomock.Any()).Times(2)

		ctx := context.WithValue(context.Background(), auth.ContextKeyIsThirdParty, false)
		ctx = context.WithValue(ctx, auth.ContextKeyUserName, "Tester")

		_, err := c.Create(ctx, input, output)
		if err != nil {
			t.Errorf("CDN.Create() error = %v, wantErr %v", err, nil)
			return
		}
	})

	t.Run("should mutate output param", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		cdnMock := NewMockcdnandler(ctrl)
		dnsMock := NewMockdnshandler(ctrl)
		acmMock := NewMockcertificateHandler(ctrl)
		sqsMock := queue.NewMockqueueHandler(ctrl)
		qMock := queue.NewMockQ(sqsMock)

		c := CDN{
			cdn:   cdnMock,
			dns:   dnsMock,
			acm:   acmMock,
			queue: &qMock,
			config: &cdnConfig{
				domain:       "test.com",
				certARN:      "arn:cert:1234",
				hostedZoneID: "/hostedzone/12345",
			},
		}

		output := &jam.OutputParam{}
		input := &jam.CreationParam{
			CDN:    jam.StackCDN{Subdomain: "subdomain"},
			ID:     "ID_12345",
			Bucket: struct{ ID string }{"ID_12345"},
		}

		expectedCallParam := c.constructStandardDistroConfig(&constructDistroConfigInput{
			bucketID:       "ID_12345",
			subdomain:      "subdomain",
			stackID:        "ID_12345",
			originAccessID: "ABCDEFG",
		})
		returnCreateOrigin := &cloudfront.CreateCloudFrontOriginAccessIdentityOutput{
			CloudFrontOriginAccessIdentity: &cloudfront.OriginAccessIdentity{
				Id: aws.String("ABCDEFG"),
			},
		}
		returnCreateDistro := &cloudfront.CreateDistributionOutput{
			Distribution: &cloudfront.Distribution{
				DomainName: aws.String("xy.com"),
				Id:         aws.String("6778ghj"),
			},
		}
		os.Setenv("DOMAIN_HOST", "test.com")
		expectedCertificateParam := &acm.RequestCertificateInput{
			DomainName:       aws.String("subdomain.test.com"),
			ValidationMethod: aws.String("DNS"),
			SubjectAlternativeNames: []*string{
				aws.String("*." + "subdomain.test.com"),
				aws.String("*.pr." + "subdomain.test.com"),
			},
		}
		returnCertificateParam := &acm.RequestCertificateOutput{
			CertificateArn: aws.String("arn:12345"),
		}
		returnCertifcateDescribeParam := &acm.DescribeCertificateOutput{
			Certificate: &acm.CertificateDetail{
				DomainValidationOptions: []*acm.DomainValidation{
					&acm.DomainValidation{
						ResourceRecord: &acm.ResourceRecord{
							Name:  aws.String("_name.aws.cert.com"),
							Type:  aws.String("CNAME"),
							Value: aws.String("_value.aws.cert.com"),
						},
					},
				},
			},
		}

		cdnMock.EXPECT().CreateCloudFrontOriginAccessIdentity(gomock.Any()).Times(1).Return(returnCreateOrigin, nil)
		cdnMock.EXPECT().CreateDistribution(expectedCallParam).Times(1).Return(returnCreateDistro, nil)
		dnsMock.EXPECT().ChangeResourceRecordSets(gomock.Any()).Times(1)
		acmMock.EXPECT().RequestCertificate(expectedCertificateParam).Times(1).Return(returnCertificateParam, nil)
		acmMock.EXPECT().DescribeCertificate(gomock.Any()).Return(returnCertifcateDescribeParam, nil)
		dnsMock.EXPECT().ChangeResourceRecordSets(gomock.Any()).Times(1)
		sqsMock.EXPECT().SendMessage(gomock.Any()).Times(2)

		ctx := context.WithValue(context.Background(), auth.ContextKeyIsThirdParty, false)
		ctx = context.WithValue(ctx, auth.ContextKeyUserName, "Tester")

		c.Create(ctx, input, output)
		if output.CDN.ID != "6778ghj" {
			t.Errorf("CDN.Create() error = %v, wantErr %v", output.CDN.ID, "6778ghj")
			return
		}
	})
}

func TestCDN_Delete(t *testing.T) {
	t.Run("should delete cdn", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		cdnMock := NewMockcdnandler(ctrl)
		dnsMock := NewMockdnshandler(ctrl)
		acmMock := NewMockcertificateHandler(ctrl)
		sqsMock := queue.NewMockqueueHandler(ctrl)
		qMock := queue.NewMockQ(sqsMock)

		c := CDN{
			cdn:   cdnMock,
			dns:   dnsMock,
			acm:   acmMock,
			queue: &qMock,
			config: &cdnConfig{
				domain:       "test.com",
				certARN:      "arn:cert:1234",
				hostedZoneID: "/hostedzone/12345",
			},
		}

		input := &jam.DeletionParam{
			CDN: &jam.StackCDN{
				ID: "ADKQEAS",
			},
			Repo: &storage.RepoModel{
				Owner: "Tester",
			},
		}

		expectedGetDistroCall := &cloudfront.GetDistributionInput{
			Id: aws.String("ADKQEAS"),
		}
		returnGetDistro := &cloudfront.GetDistributionOutput{
			Distribution: &cloudfront.Distribution{
				Id: aws.String("ADKQEAS"),
				DistributionConfig: &cloudfront.DistributionConfig{
					Enabled: aws.Bool(false),
				},
			},
			ETag: aws.String("Ifmatch"),
		}

		expectedUpdateDistroCall := &cloudfront.UpdateDistributionInput{
			Id: aws.String("ADKQEAS"),
			DistributionConfig: &cloudfront.DistributionConfig{
				Enabled: aws.Bool(false),
			},
			IfMatch: aws.String("Ifmatch"),
		}
		returnUpdateDistro := &cloudfront.UpdateDistributionOutput{
			ETag: aws.String("Ifmatch"),
		}

		cdnMock.EXPECT().GetDistribution(expectedGetDistroCall).Times(1).Return(returnGetDistro, nil)
		cdnMock.EXPECT().UpdateDistribution(expectedUpdateDistroCall).Times(1).Return(returnUpdateDistro, nil)
		dnsMock.EXPECT().ChangeResourceRecordSets(gomock.Any()).Times(1)
		sqsMock.EXPECT().SendMessage(gomock.Any()).Times(1)

		ctx := context.WithValue(context.Background(), auth.ContextKeyIsThirdParty, false)
		ctx = context.WithValue(ctx, auth.ContextKeyUserName, "Tester")
		err := c.Destroy(ctx, input)
		if err != nil {
			t.Errorf("CDN.Create() error = %v, wantErr %v", err, nil)
			return
		}
	})
}
