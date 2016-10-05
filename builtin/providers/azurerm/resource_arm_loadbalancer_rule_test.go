package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMLoadbalancerRule_basic(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	lbRuleName := fmt.Sprintf("LbRule-%d", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadbalancerRule_basic(ri, lbRuleName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadbalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadbalancerRuleExists(lbRuleName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadbalancerRule_removal(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	lbRuleName := fmt.Sprintf("LbRule-%d", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadbalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadbalancerRule_basic(ri, lbRuleName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadbalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadbalancerRuleExists(lbRuleName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadbalancerRule_removal(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadbalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadbalancerRuleNotExists(lbRuleName, &lb),
				),
			},
		},
	})
}

func testCheckAzureRMLoadbalancerRuleExists(lbRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerRuleByName(lb, lbRuleName)
		if !exists {
			return fmt.Errorf("A LoadBalancer Rule with name %q cannot be found.", lbRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadbalancerRuleNotExists(lbRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerRuleByName(lb, lbRuleName)
		if exists {
			return fmt.Errorf("A LoadBalancer Rule with name %q has been found.", lbRuleName)
		}

		return nil
	}
}

func testAccAzureRMLoadbalancerRule_basic(rInt int, lbRuleName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestrg-%d"
    location = "West US"
}

resource "azurerm_public_ip" "test" {
    name = "test-ip-%d"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
    name = "arm-test-loadbalancer-%d"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"

    frontend_ip_configuration {
      name = "one-%d"
      public_ip_address_id = "${azurerm_public_ip.test.id}"
    }
}

resource "azurerm_lb_rule" "test" {
  location = "West US"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id = "${azurerm_lb.test.id}"
  name = "%s"
  protocol = "Tcp"
  frontend_port = 3389
  backend_port = 3389
  frontend_ip_configuration_name = "one-%d"
}

`, rInt, rInt, rInt, rInt, lbRuleName, rInt)
}

func testAccAzureRMLoadbalancerRule_removal(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestrg-%d"
    location = "West US"
}

resource "azurerm_public_ip" "test" {
    name = "test-ip-%d"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
    name = "arm-test-loadbalancer-%d"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"

    frontend_ip_configuration {
      name = "one-%d"
      public_ip_address_id = "${azurerm_public_ip.test.id}"
    }
}
`, rInt, rInt, rInt, rInt)
}
