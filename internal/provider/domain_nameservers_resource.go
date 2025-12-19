package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DomainNameServersResource{}
var _ resource.ResourceWithImportState = &DomainNameServersResource{}

func NewDomainNameServersResource() resource.Resource {
	return &DomainNameServersResource{}
}

// DomainNameServersResource defines the resource implementation.
type DomainNameServersResource struct {
	client *Client
}

// DomainNameServersResourceModel describes the resource data model.
type DomainNameServersResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	Domain      types.String   `tfsdk:"domain"`
	NameServers types.Set      `tfsdk:"nameservers"`
}

func (r *DomainNameServersResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_nameservers"
}

func (r *DomainNameServersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the name servers for a domain in Porkbun.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The domain name (used as identifier).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Description: "The domain name to configure name servers for (e.g., example.com).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nameservers": schema.SetAttribute{
				Description: "Set of name servers for the domain. If you want to use Porkbun's name servers, use: curitiba.ns.porkbun.com, fortaleza.ns.porkbun.com, maceio.ns.porkbun.com, salvador.ns.porkbun.com",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *DomainNameServersResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DomainNameServersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DomainNameServersResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert nameservers set to []string
	var nsElements []types.String
	resp.Diagnostics.Append(data.NameServers.ElementsAs(ctx, &nsElements, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nameservers := make([]string, len(nsElements))
	for i, ns := range nsElements {
		nameservers[i] = ns.ValueString()
	}
	sort.Strings(nameservers)

	tflog.Debug(ctx, "Updating domain name servers", map[string]interface{}{
		"domain":      data.Domain.ValueString(),
		"nameservers": nameservers,
	})

	err := r.client.UpdateNameServers(data.Domain.ValueString(), nameservers)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update name servers: %s", err))
		return
	}

	data.ID = data.Domain

	tflog.Trace(ctx, "Updated domain name servers", map[string]interface{}{
		"domain": data.Domain.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainNameServersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DomainNameServersResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nameservers, err := r.client.GetNameServers(data.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read name servers: %s", err))
		return
	}

	// Convert to types.Set
	nsValues := make([]types.String, len(nameservers))
	for i, ns := range nameservers {
		nsValues[i] = types.StringValue(ns)
	}
	nsSet, diags := types.SetValueFrom(ctx, types.StringType, nsValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.NameServers = nsSet

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainNameServersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DomainNameServersResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert nameservers set to []string
	var nsElements []types.String
	resp.Diagnostics.Append(data.NameServers.ElementsAs(ctx, &nsElements, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nameservers := make([]string, len(nsElements))
	for i, ns := range nsElements {
		nameservers[i] = ns.ValueString()
	}
	sort.Strings(nameservers)

	tflog.Debug(ctx, "Updating domain name servers", map[string]interface{}{
		"domain":      data.Domain.ValueString(),
		"nameservers": nameservers,
	})

	err := r.client.UpdateNameServers(data.Domain.ValueString(), nameservers)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update name servers: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainNameServersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DomainNameServersResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Resetting domain name servers to Porkbun defaults", map[string]interface{}{
		"domain": data.Domain.ValueString(),
	})

	// Reset to Porkbun's default name servers on delete
	defaultNS := []string{
		"curitiba.ns.porkbun.com",
		"fortaleza.ns.porkbun.com",
		"maceio.ns.porkbun.com",
		"salvador.ns.porkbun.com",
	}

	err := r.client.UpdateNameServers(data.Domain.ValueString(), defaultNS)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to reset name servers: %s", err))
		return
	}
}

func (r *DomainNameServersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: domain
	domain := req.ID

	// Fetch the current nameservers
	nameservers, err := r.client.GetNameServers(domain)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read name servers: %s", err))
		return
	}

	// Convert to types.Set
	nsValues := make([]types.String, len(nameservers))
	for i, ns := range nameservers {
		nsValues[i] = types.StringValue(ns)
	}
	nsSet, diags := types.SetValueFrom(ctx, types.StringType, nsValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data := DomainNameServersResourceModel{
		ID:          types.StringValue(domain),
		Domain:      types.StringValue(domain),
		NameServers: nsSet,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
