package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DNSRecordResource{}
var _ resource.ResourceWithImportState = &DNSRecordResource{}

func NewDNSRecordResource() resource.Resource {
	return &DNSRecordResource{}
}

// DNSRecordResource defines the resource implementation.
type DNSRecordResource struct {
	client *Client
}

// DNSRecordResourceModel describes the resource data model.
type DNSRecordResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Domain  types.String `tfsdk:"domain"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
	Content types.String `tfsdk:"content"`
	TTL     types.String `tfsdk:"ttl"`
	Prio    types.String `tfsdk:"prio"`
	Notes   types.String `tfsdk:"notes"`
}

func (r *DNSRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *DNSRecordResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DNS record in Porkbun.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the DNS record.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Description: "The domain name for the DNS record (e.g., example.com).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The subdomain for the record, not including the domain itself. Leave empty for root domain. Use * for wildcard.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				Description: "The type of DNS record. Valid types are: A, MX, CNAME, ALIAS, TXT, NS, AAAA, SRV, TLSA, CAA, HTTPS, SVCB.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("A", "MX", "CNAME", "ALIAS", "TXT", "NS", "AAAA", "SRV", "TLSA", "CAA", "HTTPS", "SVCB"),
				},
			},
			"content": schema.StringAttribute{
				Description: "The answer content for the record.",
				Required:    true,
			},
			"ttl": schema.StringAttribute{
				Description: "The time to live in seconds for the record. Minimum and default is 600.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("600"),
			},
			"prio": schema.StringAttribute{
				Description: "The priority of the record for those that support it (e.g., MX, SRV).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("0"),
			},
			"notes": schema.StringAttribute{
				Description: "Notes for the DNS record.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
		},
	}
}

func (r *DNSRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *DNSRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DNSRecordResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createReq := CreateDNSRecordRequest{
		Name:    data.Name.ValueString(),
		Type:    data.Type.ValueString(),
		Content: data.Content.ValueString(),
		TTL:     data.TTL.ValueString(),
		Prio:    data.Prio.ValueString(),
		Notes:   data.Notes.ValueString(),
	}

	tflog.Debug(ctx, "Creating DNS record", map[string]interface{}{
		"domain":  data.Domain.ValueString(),
		"name":    createReq.Name,
		"type":    createReq.Type,
		"content": createReq.Content,
	})

	id, err := r.client.CreateDNSRecord(data.Domain.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create DNS record: %s", err))
		return
	}

	data.ID = types.StringValue(id)

	tflog.Trace(ctx, "Created DNS record", map[string]interface{}{
		"id": id,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DNSRecordResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	record, err := r.client.GetDNSRecord(data.Domain.ValueString(), data.ID.ValueString())
	if err != nil {
		// Check if the record was deleted outside of Terraform
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DNSRecordResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	editReq := EditDNSRecordRequest{
		Name:    data.Name.ValueString(),
		Type:    data.Type.ValueString(),
		Content: data.Content.ValueString(),
		TTL:     data.TTL.ValueString(),
		Prio:    data.Prio.ValueString(),
		Notes:   data.Notes.ValueString(),
	}

	tflog.Debug(ctx, "Updating DNS record", map[string]interface{}{
		"id":      data.ID.ValueString(),
		"domain":  data.Domain.ValueString(),
		"name":    editReq.Name,
		"type":    editReq.Type,
		"content": editReq.Content,
	})

	err := r.client.EditDNSRecord(data.Domain.ValueString(), data.ID.ValueString(), editReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update DNS record: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DNSRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DNSRecordResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting DNS record", map[string]interface{}{
		"id":     data.ID.ValueString(),
		"domain": data.Domain.ValueString(),
	})

	err := r.client.DeleteDNSRecord(data.Domain.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete DNS record: %s", err))
		return
	}
}

func (r *DNSRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: domain/record_id
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'domain/record_id', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
