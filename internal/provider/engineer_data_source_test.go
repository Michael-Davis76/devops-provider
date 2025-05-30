package provider

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestEngineerDataSource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            // Read testing
            {
                Config: providerConfig + `data "devops_engineer" "test" {}`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    // Verify number of coffees returned
                    resource.TestCheckResourceAttr("data.devops_engineer.test", "name", "testie_mctestface"),
                    // Verify the first coffee to ensure all attributes are set
                    resource.TestCheckResourceAttr("data.devops_engineer.test", "email", "testie@liatriolife.com"),
                ),
            },
        },
    })
}
