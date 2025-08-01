package modelarmor_test

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// Helper function to expand a template
func expandTemplate(tmplStr string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("config").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestAccModelArmorTemplate_basic(t *testing.T) {
	t.Parallel()

	templateId := "modelarmor-test-basic-" + acctest.RandString(t, 10)

	basicContext := map[string]interface{}{
		"location":   "us-central1",
		"templateId": templateId,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckModelArmorTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: func() string {
					cfg, err := testAccModelArmorTemplate_basic_config(basicContext)
					if err != nil {
						t.Fatalf("Failed to expand basic config template: %v", err)
					}
					return cfg
				}(),
			},
			{
				ResourceName:      "google_model_armor_template.template-basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccModelArmorTemplate_basic_config(context map[string]interface{}) (string, error) {
	const basic_template = `
resource "google_model_armor_template" "template-basic" {
  location    = "{{.location}}"
  template_id = "{{.templateId}}"
  filter_config {
  
  }
  template_metadata {
  
  }
}`
	return expandTemplate(basic_template, context)
}

func TestAccModelArmorTemplate_update(t *testing.T) {
	t.Parallel()

	templateId := fmt.Sprintf("modelarmor-test-update-%s", acctest.RandString(t, 5))

	context := map[string]interface{}{
		"templateId": templateId,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckModelArmorTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccModelArmorTemplate_initial(context),
			},
			{
				ResourceName:            "google_model_armor_template.test-resource",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "template_id", "terraform_labels"},
			},
			{
				Config: testAccModelArmorTemplate_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_model_armor_template.test-resource", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_model_armor_template.test-resource",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "template_id", "terraform_labels"},
			},
		},
	})
}

func testAccModelArmorTemplate_initial(context map[string]interface{}) string {
	return acctest.Nprintf(`
      resource "google_model_armor_template" "test-resource" {
        location    = "us-central1"
        template_id = "%{templateId}"
        labels = {
            "test-label" = "env-testing-initial"
        }
        filter_config {
          rai_settings {
            rai_filters {
              filter_type      = "HATE_SPEECH"
              confidence_level = "MEDIUM_AND_ABOVE"
            }
          }
          sdp_settings {
            advanced_config {
              inspect_template     = "projects/llm-firewall-demo/locations/us-central1/inspectTemplates/t2"
              deidentify_template  = "projects/llm-firewall-demo/locations/us-central1/deidentifyTemplates/t3"
            }
          }
          pi_and_jailbreak_filter_settings {
            filter_enforcement = "ENABLED"
            confidence_level   = "HIGH"
          }
          malicious_uri_filter_settings {
            filter_enforcement = "ENABLED"
          }
        }
        template_metadata {
          custom_llm_response_safety_error_message = "This is a custom error message for LLM response"
          log_template_operations                  = true
          log_sanitize_operations                  = true
          multi_language_detection {
            enable_multi_language_detection        = true
          }
          ignore_partial_invocation_failures       = true
          custom_prompt_safety_error_code          = 400
          custom_prompt_safety_error_message       = "This is a custom error message for prompt"
          custom_llm_response_safety_error_code    = 401
          enforcement_type                         = "INSPECT_ONLY"
        }
      }
    `, context)
}

func testAccModelArmorTemplate_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
      resource "google_model_armor_template" "test-resource" {
        location    = "us-central1"
        template_id = "%{templateId}"
        labels = {
            "test-label" = "env-testing-updated"
        }
        filter_config {
          rai_settings {
            rai_filters {
              filter_type      = "DANGEROUS"
              confidence_level = "LOW_AND_ABOVE"
            }
          }
          sdp_settings {
            basic_config{
              filter_enforcement = "ENABLED"
            }
          }
          pi_and_jailbreak_filter_settings {
            filter_enforcement = "DISABLED"
            confidence_level   = "MEDIUM_AND_ABOVE"
          }
          malicious_uri_filter_settings {
            filter_enforcement = "DISABLED"
          }
        }
        template_metadata {
          custom_llm_response_safety_error_message = "Updated LLM error message"
          log_template_operations                  = false
          log_sanitize_operations                  = false
          multi_language_detection {
            enable_multi_language_detection        = false
          }
          ignore_partial_invocation_failures       = false
          custom_prompt_safety_error_code          = 404
          custom_prompt_safety_error_message       = "Updated prompt error message"
          custom_llm_response_safety_error_code    = 500
          enforcement_type                         = "INSPECT_AND_BLOCK"
        }
      }
    `, context)
}
