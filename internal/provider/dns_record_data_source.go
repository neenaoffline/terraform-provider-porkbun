package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DNSRecordDataSource{}

func NewDNSRecordDataSource() datasource.DataSource {
	return &DNSRecordDataSource{}
}

// DNSRecordDataSource defines the data source implementation.
type DNSRecordDataSource struct {
	client *Client
}

// DNSRecordDataSourceModel describes the data source data model.
type DNSRecordDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Domain  types.String `tfsdk:"domain"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
	Content types.String `tfsdk:"content"`
	TTL     types.String `tfsdk:"ttl"`
	Prio    types.String `tfsdk:"prio"`
	Notes   types.String `tfsdk:"notes"`
}

func (d *DNSRecordDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (d *DNSRecordDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a DNS record from Porkbun.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the DNS record.",
				Required:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The domain name for the DNS record (e.g., example.com).",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The subdomain for the record.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of DNS record.",
				Computed:    true,
			},
			"content": schema.StringAttribute{
				Description: "The answer content for the record.",
				Computed:    true,
			},
			"ttl": schema.StringAttribute{
				Description: "The time to live in seconds for the record.",
				Computed:    true,
			},
			"prio": schema.StringAttribute{
				Description: "The priority of the record.",
				Computed:    true,
			},
			"notes": schema.StringAttribute{
				Description: "Notes for the DNS record.",
				Computed:    true,
			},
		},
	}
}

func (d *DNSRecordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *DNSRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DNSRecordDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading DNS record", map[string]interface{}{
		"id":     data.ID.ValueString(),
		"domain": data.Domain.ValueString(),
	})

	record, err := d.client.GetDNSRecord(data.Domain.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read DNS record: %s", err))
		return
	}

	// Extract the subdomain from the full name
	name := record.Name
	domain := data.Domain.ValueString()
	if name == domain {
		name = ""
	} else if strings.HasSuffix(name, "."+domain) {
		name = strings.TrimSuffix(name, "."+domain)
	}

	data.Name = types.StringValue(name)
	data.Type = types.StringValue(record.Type)
	data.Content = types.StringValue(record.Content)
	data.TTL = types.StringValue(record.TTL)
	data.Prio = types.StringValue(record.Prio)
	data.Notes = types.StringValue(record.Notes)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
