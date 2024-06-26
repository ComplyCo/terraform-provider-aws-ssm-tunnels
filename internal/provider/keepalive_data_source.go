package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &KeepaliveDataSource{}

func NewKeepaliveDataSource() datasource.DataSource {
	return &KeepaliveDataSource{}
}

// KeepaliveDataSource defines the data source implementation.
type KeepaliveDataSource struct {
}

// KeepaliveDataSourceModel describes the data source data model.
type KeepaliveDataSourceModel struct {
	Id types.String `tfsdk:"id"`
}

func (d *KeepaliveDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keepalive"
}

func (d *KeepaliveDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Data source used to keep the provider and tunnels alive",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier", // TODO: Figure this out
				Computed:            true,
			},
		},
	}
}

func (d *KeepaliveDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KeepaliveDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
