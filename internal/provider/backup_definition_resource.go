package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &BackupDefinitionResource{}
var _ resource.ResourceWithImportState = &BackupDefinitionResource{}

func NewBackupDefinitionResource() resource.Resource {
	return &BackupDefinitionResource{}
}

// BackupDefinitionResource defines the resource implementation.
type BackupDefinitionResource struct {
	client *CloudbackClient
}

// BackupDefinitionResourceModel describes the resource data model.
type BackupDefinitionResourceModel struct {
	Platform   types.String                  `tfsdk:"platform"`
	Account    types.String                  `tfsdk:"account"`
	Repository types.String                  `tfsdk:"repository"`
	Settings   BackupDefinitionSettingsModel `tfsdk:"settings"`
}

type BackupDefinitionSettingsModel struct {
	Enabled   types.Bool   `tfsdk:"enabled"`
	Schedule  types.String `tfsdk:"schedule"`
	Storage   types.String `tfsdk:"storage"`
	Retention types.String `tfsdk:"retention"`
}

func (r *BackupDefinitionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_definition"
}

func (r *BackupDefinitionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cloudback backup definition resource",

		Attributes: map[string]schema.Attribute{
			"platform": schema.StringAttribute{
				MarkdownDescription: "Platform name (e.g., GitHub, GitLab)",
				Required:            true,
			},
			"account": schema.StringAttribute{
				MarkdownDescription: "Account name",
				Required:            true,
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "Repository name",
				Required:            true,
			},
			"settings": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether the backup is scheduled",
						Required:            true,
					},
					"schedule": schema.StringAttribute{
						MarkdownDescription: "Backup schedule name",
						Required:            true,
					},
					"storage": schema.StringAttribute{
						MarkdownDescription: "Storage name",
						Required:            true,
					},
					"retention": schema.StringAttribute{
						MarkdownDescription: "Retention policy name",
						Required:            true,
					},
				},
			},
		},
	}
}

func (r *BackupDefinitionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*CloudbackClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *CloudbackClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *BackupDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BackupDefinitionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateBackupDefinition(
		data.Platform.ValueString(),
		data.Account.ValueString(),
		data.Repository.ValueString(),
		BackupDefinitionSettings{
			Enabled:   data.Settings.Enabled.ValueBool(),
			Schedule:  data.Settings.Schedule.ValueString(),
			Storage:   data.Settings.Storage.ValueString(),
			Retention: data.Settings.Retention.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update backup definition, got error: %s", err))
		return
	}

	r.LogUpdatedBackupDefinition(ctx, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BackupDefinitionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	backupDefinition, err := r.client.GetBackupDefinition(data.Platform.ValueString(), data.Account.ValueString(), data.Repository.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read backup definition, got error: %s", err))
		return
	}

	data.Settings = BackupDefinitionSettingsModel{
		Enabled:   types.BoolValue(backupDefinition.Settings.Enabled),
		Schedule:  types.StringValue(backupDefinition.Settings.Schedule),
		Storage:   types.StringValue(backupDefinition.Settings.Storage),
		Retention: types.StringValue(backupDefinition.Settings.Retention),
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BackupDefinitionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateBackupDefinition(
		data.Platform.ValueString(),
		data.Account.ValueString(),
		data.Repository.ValueString(),
		BackupDefinitionSettings{
			Enabled:   data.Settings.Enabled.ValueBool(),
			Schedule:  data.Settings.Schedule.ValueString(),
			Storage:   data.Settings.Storage.ValueString(),
			Retention: data.Settings.Retention.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update backup definition, got error: %s", err))
		return
	}

	r.LogUpdatedBackupDefinition(ctx, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BackupDefinitionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateBackupDefinition(
		data.Platform.ValueString(),
		data.Account.ValueString(),
		data.Repository.ValueString(),
		BackupDefinitionSettings{
			Enabled: false,
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update backup definition, got error: %s", err))
		return
	}

	r.LogUpdatedBackupDefinition(ctx, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: platform/account/repository. Got: %q", req.ID),
		)
		return
	}

	var data BackupDefinitionResourceModel
	data.Platform = types.StringValue(idParts[0])
	data.Account = types.StringValue(idParts[1])
	data.Repository = types.StringValue(idParts[2])

	backupDefinition, err := r.client.GetBackupDefinition(
		data.Platform.ValueString(),
		data.Account.ValueString(),
		data.Repository.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read backup definition, got error: %s", err))
		return
	}

	data.Settings = BackupDefinitionSettingsModel{
		Enabled:   types.BoolValue(backupDefinition.Settings.Enabled),
		Schedule:  types.StringValue(backupDefinition.Settings.Schedule),
		Storage:   types.StringValue(backupDefinition.Settings.Storage),
		Retention: types.StringValue(backupDefinition.Settings.Retention),
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupDefinitionResource) LogUpdatedBackupDefinition(ctx context.Context, data BackupDefinitionResourceModel) {
	tflog.Trace(ctx, "updated backup definition", map[string]interface{}{
		"platform":   data.Platform.ValueString(),
		"account":    data.Account.ValueString(),
		"repository": data.Repository.ValueString(),
		"enabled":    data.Settings.Enabled.ValueBool(),
		"schedule":   data.Settings.Schedule.ValueString(),
		"storage":    data.Settings.Storage.ValueString(),
		"retention":  data.Settings.Retention.ValueString(),
	})
}
