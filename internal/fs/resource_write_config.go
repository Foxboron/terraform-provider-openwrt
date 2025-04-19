package fs

import (
	"bytes"
	"context"
	"fmt"
	"path"

	"github.com/foxboron/terraform-provider-openwrt/internal/api"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ConfigFileModel struct {
	Name    types.String `tfsdk:"name"`
	Content types.String `tfsdk:"content"`
	Commit  types.Bool   `tfsdk:"commit"`
}

// ConfigFileResource represent Incus project resource.
type ConfigFileResource struct {
	provider *api.Client
}

// NewProjectResource return new project resource.
func NewConfigFileResource() resource.Resource {
	return &ConfigFileResource{}
}

// Metadata for project resource.
func (c ConfigFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_configfile", req.ProviderTypeName)
}

func (c ConfigFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},

			"content": schema.StringAttribute{
				Required: true,
			},

			"commit": schema.BoolAttribute{
				Optional: true,
				Default:  booldefault.StaticBool(true),
				Computed: true,
			},
		},
	}
}

func (c *ConfigFileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	data := req.ProviderData
	if data == nil {
		return
	}
	provider, ok := data.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError("Failed to get api client", "")
		return
	}
	c.provider = provider
}

func (c ConfigFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConfigFileModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := path.Join("/etc/config", plan.Name.ValueString())

	if _, err := c.provider.Writefile(path, []byte(plan.Content.ValueString())); err != nil {
		resp.Diagnostics.AddError("Failed to write file", err.Error())
		return
	}

	if plan.Commit.ValueBool() {
		if _, d := c.provider.UCICommitAndRevert(plan.Name.ValueString()); d != nil {
			resp.Diagnostics.Append(d...)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (c ConfigFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ConfigFileModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := path.Join("/etc/config", state.Name.ValueString())
	b, err := c.provider.ReadFile(path)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read config %q", state.Name.ValueString()), err.Error())
		return
	}

	// Logic taken from the local_file provider
	if !bytes.Equal(b, []byte(state.Content.ValueString())) {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (c ConfigFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConfigFileModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := path.Join("/etc/config", plan.Name.ValueString())
	if _, err := c.provider.Writefile(path, []byte(plan.Content.ValueString())); err != nil {
		resp.Diagnostics.AddError("Failed to write file", err.Error())
	}

	if plan.Commit.ValueBool() {
		if _, d := c.provider.UCICommitAndRevert(plan.Name.ValueString()); d != nil {
			resp.Diagnostics.Append(d...)
			return
		}
	}

	plan.Commit = types.BoolValue(true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (c ConfigFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ConfigFileModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	path := path.Join("/etc/config", state.Name.ValueString())
	_, err := c.provider.RemoveFile(path)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to delete config file %q", state.Name.ValueString()), err.Error())
		return
	}

	if state.Commit.ValueBool() {
		if _, d := c.provider.UCICommitAndRevert(state.Name.ValueString()); d != nil {
			resp.Diagnostics.Append(d...)
			return
		}
	}
}

func (c *ConfigFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state ConfigFileModel

	path := path.Join("/etc/config", req.ID)
	b, err := c.provider.ReadFile(path)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to read config %q", state.Name.ValueString()), err.Error())
		return
	}

	state.Name = types.StringValue(req.ID)
	state.Content = types.StringValue(string(b))
	state.Commit = types.BoolValue(true)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
