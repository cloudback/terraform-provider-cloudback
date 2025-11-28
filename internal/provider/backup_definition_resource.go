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
	Platform    types.String                  `tfsdk:"platform"`
	Account     types.String                  `tfsdk:"account"`
	SubjectType types.String                  `tfsdk:"subject_type"`
	SubjectName types.String                  `tfsdk:"subject_name"`
	Repository  types.String                  `tfsdk:"repository"`
	Settings    BackupDefinitionSettingsModel `tfsdk:"settings"`
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
				MarkdownDescription: "Platform name (e.g., GitHub, GitLab, AzureDevOps)",
				Required:            true,
			},
			"account": schema.StringAttribute{
				MarkdownDescription: "Account name",
				Required:            true,
			},
			"subject_type": schema.StringAttribute{
				MarkdownDescription: "Subject type (e.g., Repository, Project)",
				Optional:            true,
			},
			"subject_name": schema.StringAttribute{
				MarkdownDescription: "Subject name (repository name, project name, etc.)",
				Optional:            true,
			},
			"repository": schema.StringAttribute{
				MarkdownDescription: "Repository name (deprecated: use subject_type and subject_name instead)",
				Optional:            true,
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

	// Determine subject_type and subject_name for API call (backward compatibility)
	var subjectType, subjectName string
	if !data.SubjectType.IsNull() && !data.SubjectName.IsNull() {
		// New format: use subject fields directly
		subjectType = data.SubjectType.ValueString()
		subjectName = data.SubjectName.ValueString()
	} else if !data.Repository.IsNull() {
		// Old format: derive from repository
		subjectType = "Repository"
		subjectName = data.Repository.ValueString()
	} else {
		resp.Diagnostics.AddError(
			"Missing Required Fields",
			"Either 'repository' or both 'subject_type' and 'subject_name' must be provided.",
		)
		return
	}

	err := r.client.UpdateBackupDefinition(
		data.Platform.ValueString(),
		data.Account.ValueString(),
		subjectType,
		subjectName,
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

	// Determine subject_type and subject_name for API call (backward compatibility)
	var subjectType, subjectName string
	if !data.SubjectType.IsNull() && !data.SubjectName.IsNull() {
		subjectType = data.SubjectType.ValueString()
		subjectName = data.SubjectName.ValueString()
	} else if !data.Repository.IsNull() {
		subjectType = "Repository"
		subjectName = data.Repository.ValueString()
	}

	backupDefinition, err := r.client.GetBackupDefinition(data.Platform.ValueString(), data.Account.ValueString(), subjectType, subjectName)
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

	// Determine subject_type and subject_name for API call (backward compatibility)
	var subjectType, subjectName string
	if !data.SubjectType.IsNull() && !data.SubjectName.IsNull() {
		// New format: use subject fields directly
		subjectType = data.SubjectType.ValueString()
		subjectName = data.SubjectName.ValueString()
	} else if !data.Repository.IsNull() {
		// Old format: derive from repository
		subjectType = "Repository"
		subjectName = data.Repository.ValueString()
	} else {
		resp.Diagnostics.AddError(
			"Missing Required Fields",
			"Either 'repository' or both 'subject_type' and 'subject_name' must be provided.",
		)
		return
	}

	err := r.client.UpdateBackupDefinition(
		data.Platform.ValueString(),
		data.Account.ValueString(),
		subjectType,
		subjectName,
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

	// Determine subject_type and subject_name for API call (backward compatibility)
	var subjectType, subjectName string
	if !data.SubjectType.IsNull() && !data.SubjectName.IsNull() {
		subjectType = data.SubjectType.ValueString()
		subjectName = data.SubjectName.ValueString()
	} else if !data.Repository.IsNull() {
		subjectType = "Repository"
		subjectName = data.Repository.ValueString()
	}

	err := r.client.UpdateBackupDefinition(
		data.Platform.ValueString(),
		data.Account.ValueString(),
		subjectType,
		subjectName,
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

	var data BackupDefinitionResourceModel

	// Support both old format (platform/account/repository) and new format (platform/account/subject_type/subject_name)
	if len(idParts) == 3 {
		// Old format: platform/account/repository - assume Repository subject type
		if idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier",
				fmt.Sprintf("Expected import identifier with format: platform/account/repository or platform/account/subject_type/subject_name. Got: %q", req.ID),
			)
			return
		}
		data.Platform = types.StringValue(idParts[0])
		data.Account = types.StringValue(idParts[1])
		data.Repository = types.StringValue(idParts[2])
	} else if len(idParts) == 4 {
		// New format: platform/account/subject_type/subject_name
		if idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier",
				fmt.Sprintf("Expected import identifier with format: platform/account/repository or platform/account/subject_type/subject_name. Got: %q", req.ID),
			)
			return
		}
		data.Platform = types.StringValue(idParts[0])
		data.Account = types.StringValue(idParts[1])
		data.SubjectType = types.StringValue(idParts[2])
		data.SubjectName = types.StringValue(idParts[3])
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: platform/account/repository or platform/account/subject_type/subject_name. Got: %q", req.ID),
		)
		return
	}

	// Determine subject_type and subject_name for API call
	var subjectType, subjectName string
	if !data.SubjectType.IsNull() && !data.SubjectName.IsNull() {
		subjectType = data.SubjectType.ValueString()
		subjectName = data.SubjectName.ValueString()
	} else if !data.Repository.IsNull() {
		subjectType = "Repository"
		subjectName = data.Repository.ValueString()
	}

	backupDefinition, err := r.client.GetBackupDefinition(
		data.Platform.ValueString(),
		data.Account.ValueString(),
		subjectType,
		subjectName)

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
	logData := map[string]interface{}{
		"platform":  data.Platform.ValueString(),
		"account":   data.Account.ValueString(),
		"enabled":   data.Settings.Enabled.ValueBool(),
		"schedule":  data.Settings.Schedule.ValueString(),
		"storage":   data.Settings.Storage.ValueString(),
		"retention": data.Settings.Retention.ValueString(),
	}
	if !data.SubjectType.IsNull() {
		logData["subject_type"] = data.SubjectType.ValueString()
	}
	if !data.SubjectName.IsNull() {
		logData["subject_name"] = data.SubjectName.ValueString()
	}
	if !data.Repository.IsNull() {
		logData["repository"] = data.Repository.ValueString()
	}
	tflog.Trace(ctx, "updated backup definition", logData)
}
