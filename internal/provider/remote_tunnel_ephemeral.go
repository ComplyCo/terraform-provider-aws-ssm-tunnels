package provider

import (
	"context"
	"fmt"

	"github.com/complyco/terraform-provider-aws-ssm-tunnels/internal/ports"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ ephemeral.EphemeralResourceWithConfigure = &RemoteTunnelEphemeralResource{}
	_ ephemeral.EphemeralResource              = &RemoteTunnelEphemeralResource{}
)

func NewRemoteTunnelEphemeralResource() ephemeral.EphemeralResource {
	return &RemoteTunnelEphemeralResource{}
}

// RemoteTunnelEphemeralResource defines the resource implementation.
type RemoteTunnelEphemeralResource struct {
	tracker *TunnelTracker
	region  string
	target  string
}

// SSMRemoteTunnelDataSourceModel describes the data source data model.
type SSMRemoteTunnelEphemeralResourceModel struct {
	RemoteHost types.String `tfsdk:"remote_host"`
	RemotePort types.Int64  `tfsdk:"remote_port"`
	LocalPort  types.Int64  `tfsdk:"local_port"`
	LocalHost  types.String `tfsdk:"local_host"`
	Id         types.String `tfsdk:"id"`
}

func (d *RemoteTunnelEphemeralResource) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_tunnel"
}

func (d *RemoteTunnelEphemeralResource) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "AWSM SSM Remote Tunnel data source",

		Attributes: map[string]schema.Attribute{
			"remote_host": schema.StringAttribute{
				MarkdownDescription: "The DNS name or IP address of the remote host",
				Required:            true,
			},
			"remote_port": schema.Int64Attribute{
				MarkdownDescription: "The port number of the remote host",
				Required:            true,
			},
			"local_host": schema.StringAttribute{
				MarkdownDescription: "The DNS name or IP address of the local host",
				Computed:            true,
			},
			"local_port": schema.Int64Attribute{
				MarkdownDescription: "The local port number to use for the tunnel",
				Optional:            true,
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier", // TODO: Figure this out
				Computed:            true,
			},
		},
	}
}

func (d *RemoteTunnelEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	configData, ok := req.ProviderData.(*ProvidedConfigData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProvidedConfigData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.tracker = configData.Tracker
	d.region = configData.Region
	d.target = configData.Target
}

func (d *RemoteTunnelEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data SSMRemoteTunnelEphemeralResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var port int
	var err error
	port = int(data.LocalPort.ValueInt64())
	if port == 0 {
		port, err = ports.FindOpenPort(16000, 26000)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to find open port",
				fmt.Sprintf("Error: %s", err),
			)
			return
		}
	}

	tunnelInfo, err := d.tracker.StartTunnel(
		ctx,
		data.Id.ValueString(),
		d.target,
		data.RemoteHost.ValueString(),
		int(data.RemotePort.ValueInt64()),
		port,
		d.region,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to start remote tunnel",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	data.Id = basetypes.NewStringValue(uuid.New().String())
	data.LocalPort = basetypes.NewInt64Value(int64(tunnelInfo.LocalPort))
	data.LocalHost = basetypes.NewStringValue(tunnelInfo.LocalHost)

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
