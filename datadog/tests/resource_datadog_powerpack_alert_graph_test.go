package test

import (
	"testing"
)

const datadogPowerpackAlertGraphTest = `
resource "datadog_powerpack" "alert_graph_powerpack" {
	name         = "{{uniq}}"
    tags = ["tag:foo1"]
	description   = "Created using the Datadog provider in Terraform"
	template_variables {
		defaults = ["defaults"]
		name     = "datacenter"
	}
    widget {
      alert_graph_definition {
        alert_id  = "895605"
        viz_type  = "timeseries"
        title     = "Widget Title"
        title_align = "center"
        title_size = "20"
      }
    }
}
`

var datadogPowerpackAlertGraphTestAsserts = []string{
	// Powerpack metadata
	"description = Created using the Datadog provider in Terraform",
	"widget.# = 1",
	"tags.# = 1",
	"tags.0 = tag:foo1",
	// Alert Graph widget
	"widget.0.alert_graph_definition.0.alert_id = 895605",
	"widget.0.alert_graph_definition.0.viz_type = timeseries",
	"widget.0.alert_graph_definition.0.title = Widget Title",
	"widget.0.alert_graph_definition.0.title_align = center",
	"widget.0.alert_graph_definition.0.title_size = 20",

	// Template Variables
	"template_variables.# = 1",
	"template_variables.0.name = datacenter",
	"template_variables.0.defaults.# = 1",
	"template_variables.0.defaults.0 = defaults",
}

func TestAccDatadogPowerpackAlertGraph(t *testing.T) {
	testAccDatadogPowerpackWidgetUtil(t, datadogPowerpackAlertGraphTest, "datadog_powerpack.alert_graph_powerpack", datadogPowerpackAlertGraphTestAsserts)
}
