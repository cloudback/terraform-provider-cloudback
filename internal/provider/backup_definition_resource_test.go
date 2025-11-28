package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBackupDefinitionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "cloudback_backup_definition" "test" {
  platform = "GitHub"
  account = "testland"
  repository = "docs"
  settings = {
    enabled = true
    schedule = "Daily at 9 pm"
    storage = "Cloudback EU"
    retention = "Last 30 days"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudback_backup_definition.test", "platform", "GitHub"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test", "account", "testland"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test", "repository", "docs"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test", "settings.enabled", "true"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test", "settings.schedule", "Daily at 9 pm"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test", "settings.storage", "Cloudback EU"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test", "settings.retention", "Last 30 days"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "cloudback_backup_definition.test",
				ImportState:                          true,
				ImportStateId:                        "GitHub/testland/docs",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "repository",
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "cloudback_backup_definition" "test" {
  platform = "GitHub"
  account = "testland"
  repository = "docs"
  settings = {
    enabled = false
    schedule = "Daily at 9 pm"
    storage = "Cloudback EU"
    retention = "Last 30 days"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudback_backup_definition.test", "settings.enabled", "false"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test", "settings.schedule", "Daily at 9 pm"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test", "settings.storage", "Cloudback EU"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test", "settings.retention", "Last 30 days"),
				),
			},
			// Delete is implicitly tested at the end of the TestCase.
		},
	})
}

func TestAccBackupDefinitionResourceWithSubjectFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with subject_type and subject_name
			{
				Config: providerConfig + `
resource "cloudback_backup_definition" "test_subject" {
  platform = "GitHub"
  account = "testland"
  subject_type = "Repository"
  subject_name = "docs"
  settings = {
    enabled = true
    schedule = "Daily at 9 pm"
    storage = "Cloudback EU"
    retention = "Last 30 days"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudback_backup_definition.test_subject", "platform", "GitHub"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test_subject", "account", "testland"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test_subject", "subject_type", "Repository"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test_subject", "subject_name", "docs"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test_subject", "settings.enabled", "true"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test_subject", "settings.schedule", "Daily at 9 pm"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test_subject", "settings.storage", "Cloudback EU"),
					resource.TestCheckResourceAttr("cloudback_backup_definition.test_subject", "settings.retention", "Last 30 days"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "cloudback_backup_definition" "test_subject" {
  platform = "GitHub"
  account = "testland"
  subject_type = "Repository"
  subject_name = "docs"
  settings = {
    enabled = false
    schedule = "Daily at 9 pm"
    storage = "Cloudback EU"
    retention = "Last 30 days"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudback_backup_definition.test_subject", "settings.enabled", "false"),
				),
			},
		},
	})
}
