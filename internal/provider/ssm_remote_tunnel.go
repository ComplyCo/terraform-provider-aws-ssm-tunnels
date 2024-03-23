package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/ssmtunnels"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SSMRemoteTunnelDataSource{}

func NewSSMRemoteTunnelDataSource() datasource.DataSource {
	return &SSMRemoteTunnelDataSource{}
}

// SSMRemoteTunnelDataSource defines the data source implementation.
type SSMRemoteTunnelDataSource struct {
	svc *ssm.Client
}

// SSMRemoteTunnelDataSourceModel describes the data source data model.
type SSMRemoteTunnelDataSourceModel struct {
	Target     types.String `tfsdk:"target"`
	RemoteHost types.String `tfsdk:"remote_host"`
	RemotePort types.Int64  `tfsdk:"remote_port"`
	LocalPort  types.Int64  `tfsdk:"local_port"`
	Region     types.String `tfsdk:"region"`
	Id         types.String `tfsdk:"id"`
}

func (d *SSMRemoteTunnelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssm_remote_tunnel"
}

func (d *SSMRemoteTunnelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "AWSM SSM Remote Tunnel data source",

		Attributes: map[string]schema.Attribute{
			"target": schema.StringAttribute{
				MarkdownDescription: "The target to start the remote tunnel, such as an instance ID",
				Optional:            false,
			},
			"remote_host": schema.StringAttribute{
				MarkdownDescription: "The DNS name or IP address of the remote host",
				Optional:            false,
			},
			"remote_port": schema.Int64Attribute{
				MarkdownDescription: "The port number of the remote host",
				Optional:            false,
			},
			"local_port": schema.Int64Attribute{
				MarkdownDescription: "The local port number to use for the tunnel",
				Optional:            false,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The AWS region to use for the tunnel. This should match the region of the target",
				Optional:            false,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier", // TODO: Figure this out
				Computed:            true,
			},
		},
	}
}

func (d *SSMRemoteTunnelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	svc, ok := req.ProviderData.(*ssm.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.svc = svc
}

func (d *SSMRemoteTunnelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SSMRemoteTunnelDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Decide if this should be in the Configure step
	err := ssmtunnels.StartRemoteTunnel(context.Background(), ssmtunnels.RemoteTunnelConfig{
		Client:     d.svc,
		Target:     data.Target.ValueString(),
		Region:     data.Region.ValueString(),
		RemoteHost: data.RemoteHost.ValueString(),
		RemotePort: int(data.RemotePort.ValueInt64()),
		LocalPort:  int(data.LocalPort.ValueInt64()),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to start remote tunnel",
			fmt.Sprintf("Error: %s", err),
		)
	}
	// TODO: Figure out how to store a reference to the tunnel so it can be stopped later

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
