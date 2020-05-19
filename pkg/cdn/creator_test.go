package cdn

import (
	"context"
	"os"
	"testing"

	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	gomock "github.com/golang/mock/gomock"
)

func TestCDN_Create(t *testing.T) {
	t.Run("should create bucket", func(t *testing.T) {
		os.Setenv("DOMAIN_HOST", "test")
		ctrl := gomock.NewController(t)
		cdnMock := NewMockcdnandler(ctrl)
		dnsMock := NewMockdnshandler(ctrl)
		sqsMock := queue.NewMockqueueHandler(ctrl)
		qMock := queue.NewMockQ(sqsMock)

		c := CDN{
			cdn:   cdnMock,
			dns:   dnsMock,
			queue: &qMock,
		}

		output := &jam.OutputParam{}
		input := &jam.CreationParam{
			ID:     "ID_12345",
			Bucket: struct{ ID string }{"ID_12345"},
		}

		expectedCallParam := c.constructStandardDistroConfig("ID_12345", "ABCDEFG", "ID_12345")
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

		cdnMock.EXPECT().CreateCloudFrontOriginAccessIdentity(gomock.Any()).Times(1).Return(returnCreateOrigin, nil)
		cdnMock.EXPECT().CreateDistribution(expectedCallParam).Times(1).Return(returnCreateDistro, nil)
		dnsMock.EXPECT().ChangeResourceRecordSets(gomock.Any()).Times(1)
		sqsMock.EXPECT().SendMessage(gomock.Any()).Times(1)

		_, err := c.Create(context.Background(), input, output)
		if err != nil {
			t.Errorf("CDN.Create() error = %v, wantErr %v", err, nil)
			return
		}
	})

	t.Run("should mutate output param", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		cdnMock := NewMockcdnandler(ctrl)
		dnsMock := NewMockdnshandler(ctrl)
		sqsMock := queue.NewMockqueueHandler(ctrl)
		qMock := queue.NewMockQ(sqsMock)

		c := CDN{
			cdn:   cdnMock,
			dns:   dnsMock,
			queue: &qMock,
		}

		output := &jam.OutputParam{}
		input := &jam.CreationParam{
			ID:     "ID_12345",
			Bucket: struct{ ID string }{"ID_12345"},
		}

		expectedCallParam := c.constructStandardDistroConfig("ID_12345", "ABCDEFG", "ID_12345")
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

		cdnMock.EXPECT().CreateCloudFrontOriginAccessIdentity(gomock.Any()).Times(1).Return(returnCreateOrigin, nil)
		cdnMock.EXPECT().CreateDistribution(expectedCallParam).Times(1).Return(returnCreateDistro, nil)
		dnsMock.EXPECT().ChangeResourceRecordSets(gomock.Any()).Times(1)
		sqsMock.EXPECT().SendMessage(gomock.Any()).Times(1)

		c.Create(context.Background(), input, output)
		if output.CDN.ID != "6778ghj" {
			t.Errorf("CDN.Create() error = %v, wantErr %v", output.CDN.ID, "6778ghj")
			return
		}
	})
}

func TestCDN_Delete(t *testing.T) {
	t.Run("should delete bucket", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		cdnMock := NewMockcdnandler(ctrl)
		dnsMock := NewMockdnshandler(ctrl)
		sqsMock := queue.NewMockqueueHandler(ctrl)
		qMock := queue.NewMockQ(sqsMock)

		c := CDN{
			cdn:   cdnMock,
			dns:   dnsMock,
			queue: &qMock,
		}

		input := &jam.DeletionParam{
			CDN: &jam.StackCDN{
				ID: "ADKQEAS",
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
		// dnsMock.EXPECT().ChangeResourceRecordSets(gomock.Any()).Times(1)
		sqsMock.EXPECT().SendMessage(gomock.Any()).Times(1)

		err := c.Destroy(context.Background(), input)
		if err != nil {
			t.Errorf("CDN.Create() error = %v, wantErr %v", err, nil)
			return
		}
	})
}
