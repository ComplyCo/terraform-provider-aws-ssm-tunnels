package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/complyco/terraform-provider-aws-ssm-tunnels/internal/ports"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RemoteTunnelResource{}
var _ resource.ResourceWithImportState = &RemoteTunnelResource{}

func NewRemoteTunnelResource() resource.Resource {
	return &RemoteTunnelResource{}
}

// RemoteTunnelResource defines the resource implementation.
type RemoteTunnelResource struct {
	tracker *TunnelTracker
	region  string
	target  string
}

// SSMRemoteTunnelDataSourceModel describes the data source data model.
type SSMRemoteTunnelResourceModel struct {
	RefreshId  types.String `tfsdk:"refresh_id"`
	RemoteHost types.String `tfsdk:"remote_host"`
	RemotePort types.Int64  `tfsdk:"remote_port"`
	LocalPort  types.Int64  `tfsdk:"local_port"`
	LocalHost  types.String `tfsdk:"local_host"`
	Id         types.String `tfsdk:"id"`
}

func (d *RemoteTunnelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_tunnel"
}

func (d *RemoteTunnelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "AWSM SSM Remote Tunnel data source",

		Attributes: map[string]schema.Attribute{
			"refresh_id": schema.StringAttribute{
				MarkdownDescription: "Any value as this will trigger a refresh",
				Required:            true,
			},
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

func (d *RemoteTunnelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (d *RemoteTunnelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SSMRemoteTunnelResourceModel

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RemoteTunnelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SSMRemoteTunnelResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

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

	data.RefreshId = basetypes.NewStringValue(uuid.New().String()) // NOTE: We always change this in order to force an update
	data.LocalPort = basetypes.NewInt64Value(int64(tunnelInfo.LocalPort))
	data.LocalHost = basetypes.NewStringValue(tunnelInfo.LocalHost)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RemoteTunnelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SSMRemoteTunnelResourceModel

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *RemoteTunnelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SSMRemoteTunnelResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *RemoteTunnelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "|")
	// TODO: Decide if we need the local_host set. Also do we need the local_port?
	if len(parts) != 4 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID must be in the format `remote_host|remote_port|local_port|local_host`",
		)
		return
	}
	remoteHost := parts[0]
	remotePort := parts[1]
	localPort := parts[2]
	localHost := parts[3]

	localPortInt, _ := strconv.Atoi(localPort)
	remotePortInt, _ := strconv.Atoi(remotePort)

	resp.State.Set(ctx, &SSMRemoteTunnelResourceModel{
		// TODO: Figure out if we need to set the ID here
		Id:         basetypes.NewStringValue(uuid.New().String()),
		RemoteHost: basetypes.NewStringValue(remoteHost),
		RemotePort: basetypes.NewInt64Value(int64(remotePortInt)),
		LocalPort:  basetypes.NewInt64Value(int64(localPortInt)),
		LocalHost:  basetypes.NewStringValue(localHost),
	})
}
