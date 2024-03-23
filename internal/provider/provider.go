// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
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
	AwsRegion types.String `tfsdk:"aws_region"`
	AwsKey    types.String `tfsdk:"aws_key"`
	AwsSecret types.String `tfsdk:"aws_secret"`
}

func (p *AwsSSMTunnelsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "aws_ssm_tunnels"
	resp.Version = p.version
}

func (p *AwsSSMTunnelsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	// TODO: Figure out how to support more auth modes. Maybe import from the AWS provider
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"aws_region": schema.StringAttribute{
				MarkdownDescription: "The AWS region to use for the SSM tunnel",
				Optional:            false,
			},
			"aws_key": schema.StringAttribute{
				MarkdownDescription: "The AWS Access Key ID to use for the SSM tunnel",
				Sensitive:           true,
				Optional:            false,
			},
			"aws_secret": schema.StringAttribute{
				MarkdownDescription: "The AWS Secret Access Key to use for the SSM tunnel",
				Sensitive:           true,
				Optional:            false,
			},
			// TODO: Add session token
		},
	}
}

func (p *AwsSSMTunnelsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AwsSSMTunnelsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// config.Config{
	// 	Region: data.AwsRegion.String(),
	// 	Credentials: config.Credentials{

	// }

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(data.AwsRegion.String()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to load AWS configuration",
			fmt.Sprintf("Error: %s", err),
		)
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
