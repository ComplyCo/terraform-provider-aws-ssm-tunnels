// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure AwsSSMTunnelsProvider satisfies various provider interfaces.
var _ provider.Provider = &AwsSSMTunnelsProvider{}

// AwsSSMTunnelsProvider defines the provider implementation.
type AwsSSMTunnelsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AwsSSMTunnelsProviderModel describes the provider data model.
type AwsSSMTunnelsProviderModel struct {
	Region            types.String   `tfsdk:"region"`
	AccessKey         types.String   `tfsdk:"access_key"`
	SecretKey         types.String   `tfsdk:"secret_key"`
	SessionToken      types.String   `tfsdk:"token"`
	SharedConfigFiles []types.String `tfsdk:"shared_config_files"`
}

func (p *AwsSSMTunnelsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "aws_ssm_tunnels"
	resp.Version = p.version
}

func (p *AwsSSMTunnelsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	// TODO: Figure out how to support more auth modes. Maybe import from the AWS provider
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Optional: true,
				Description: "The region where AWS operations will take place. Examples\n" +
					"are us-east-1, us-west-2, etc.",
			},
			"access_key": schema.StringAttribute{
				Optional: true,
				Description: "The access key for API operations. You can retrieve this\n" +
					"from the 'Security & Credentials' section of the AWS console.",
			},
			"secret_key": schema.StringAttribute{
				Optional: true,
				Description: "The secret key for API operations. You can retrieve this\n" +
					"from the 'Security & Credentials' section of the AWS console.",
			},
			"token": schema.StringAttribute{
				Optional: true,
				Description: "session token. A session token is only required if you are\n" +
					"using temporary security credentials.",
			},
			"shared_config_files": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of paths to shared config files. If not set, defaults to [~/.aws/config].",
			},
		},
	}
}

func (p *AwsSSMTunnelsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AwsSSMTunnelsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var awsCfg aws.Config
	var err error
	if len(data.SharedConfigFiles) > 0 {
		sharedConfigFilesAsString := []string{}
		for _, path := range data.SharedConfigFiles {
			sharedConfigFilesAsString = append(sharedConfigFilesAsString, path.ValueString())
		}

		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithSharedConfigFiles(sharedConfigFilesAsString),
		)

		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to load AWS configuration",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}
	} else {
		awsCfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(data.Region.String()),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					data.AccessKey.ValueString(),
					data.SecretKey.ValueString(),
					data.SessionToken.ValueString(), // NOTE: SessionToken can be an empty string
				),
			),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to load AWS configuration",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}
	}

	svc := ssm.NewFromConfig(awsCfg)
	// NOTE: We should make a "client" struct which hides the SSM client, and has a method to start a tunnel and it keeps track of the tunnel session
	// It should also handle the cancellation via context signalling

	resp.DataSourceData = svc
	resp.ResourceData = svc
}

func (p *AwsSSMTunnelsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *AwsSSMTunnelsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSSMRemoteTunnelDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AwsSSMTunnelsProvider{
			version: version,
		}
	}
}
