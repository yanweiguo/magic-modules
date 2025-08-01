package dataplex_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataplexDatascanDataplexDatascanFullQuality_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexDatascanDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexDatascanDataplexDatascanFullQuality_full(context),
			},
			{
				ResourceName:            "google_dataplex_datascan.full_quality",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_scan_id", "labels", "location", "terraform_labels"},
			},
			{
				Config: testAccDataplexDatascanDataplexDatascanFullQuality_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataplex_datascan.full_quality", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_dataplex_datascan.full_quality",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"data_scan_id", "labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccDataplexDatascanDataplexDatascanFullQuality_full(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_bigquery_dataset" "tf_test_dataset" {
  dataset_id = "tf_test_dataset_id_%{random_suffix}"
  default_table_expiration_ms = 3600000
}

resource "google_bigquery_table" "tf_test_table" {
  dataset_id          = google_bigquery_dataset.tf_test_dataset.dataset_id
  table_id            = "tf_test_table_%{random_suffix}"
  deletion_protection = false
  schema              = <<EOF
    [
    {
      "name": "name",
      "type": "STRING",
      "mode": "NULLABLE"
    },
    {
      "name": "age",
      "type": "INTEGER",
      "mode": "NULLABLE",
      "description": "Age of the person"
    }
    ]
  EOF
}
  
resource "google_dataplex_datascan" "full_quality" {
  location = "us-central1"
  display_name = "Full Datascan Quality"
  data_scan_id = "tf-test-dataquality-full%{random_suffix}"
  description = "Example resource - Full Datascan Quality"
  labels = {
    author = "billing"
  }

  data {
    resource = "//bigquery.googleapis.com/projects/%{project_name}/datasets/${google_bigquery_dataset.tf_test_dataset.dataset_id}/tables/${google_bigquery_table.tf_test_table.table_id}"
  }

  execution_spec {
    trigger {
      schedule {
        cron = "TZ=America/New_York 1 1 * * *"
      }
    }
  }

  data_quality_spec {
    sampling_percent = 5
    row_filter = "age > 10"
    post_scan_actions {
      notification_report {
        recipients {
          emails = ["jane.doe@example.com"]
        }
        score_threshold_trigger {
          score_threshold = 86
        }
      }
    }
    
    rules {
      column = "name"
      dimension = "VALIDITY"
      threshold = 0.99
      non_null_expectation {}
    }

    rules {
      column = "age"
      dimension = "VALIDITY"
      ignore_null = true
      threshold = 0.9
      range_expectation {
        min_value = 1
        max_value = 100
        strict_min_enabled = true
        strict_max_enabled = false
      }
    }
  }

  project = "%{project_name}"
}
`, context)
}

func testAccDataplexDatascanDataplexDatascanFullQuality_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "tf_test_dataset" {
  dataset_id = "tf_test_dataset_id_%{random_suffix}"
  default_table_expiration_ms = 3600000
}

resource "google_bigquery_table" "tf_test_table" {
  dataset_id          = google_bigquery_dataset.tf_test_dataset.dataset_id
  table_id            = "tf_test_table_%{random_suffix}"
  deletion_protection = false
  schema              = <<EOF
    [
    {
      "name": "name",
      "type": "STRING",
      "mode": "NULLABLE"
    },
    {
      "name": "age",
      "type": "INTEGER",
      "mode": "NULLABLE",
      "description": "Age of the person"
    }
    ]
  EOF
}

resource "google_dataplex_datascan" "full_quality" {
  location = "us-central1"
  display_name = "Full Datascan Quality"
  data_scan_id = "tf-test-dataquality-full%{random_suffix}"
  description = "Example resource - Full Datascan Quality"
  labels = {
    author = "billing"
  }

  data {
    resource = "//bigquery.googleapis.com/projects/%{project_name}/datasets/${google_bigquery_dataset.tf_test_dataset.dataset_id}/tables/${google_bigquery_table.tf_test_table.table_id}"
  }

  execution_spec {
    trigger {
      schedule {
        cron = "TZ=America/New_York 1 1 * * *"
      }
    }
  }

  data_quality_spec {
    sampling_percent = 5
    row_filter = "age > 10"
    catalog_publishing_enabled = true
    post_scan_actions {
      notification_report {
        recipients {
          emails = ["jane.doe@example.com"]
        }
        score_threshold_trigger {
          score_threshold = 86
        }
      }
    }
    
    rules {
      column = "name"
      dimension = "VALIDITY"
      threshold = 0.99
      non_null_expectation {}
    }

    rules {
      column = "age"
      dimension = "VALIDITY"
      ignore_null = true
      threshold = 0.9
      range_expectation {
        min_value = 1
        max_value = 100
        strict_min_enabled = true
        strict_max_enabled = false
      }
      suspended = true
    }
  }

  project = "%{project_name}"
}
`, context)
}
