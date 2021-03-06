package softlayer

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/softlayer/softlayer-go/services"
	"regexp"
)

func TestAccSoftLayerGlobalIp_Basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckSoftLayerGlobalIpConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerGlobalIpExists("softlayer_global_ip.test-global-ip"),
					resource.TestMatchResourceAttr("softlayer_global_ip.test-global-ip", "ip_address",
						regexp.MustCompile(`^(([01]?[0-9]?[0-9]|2([0-4][0-9]|5[0-5]))\.){3}([01]?[0-9]?[0-9]|2([0-4][0-9]|5[0-5]))$`)),
					testAccCheckSoftLayerResources("softlayer_global_ip.test-global-ip", "routes_to",
						"softlayer_virtual_guest.vm1", "ipv4_address"),
				),
			},

			resource.TestStep{
				Config: testAccCheckSoftLayerGlobalIpConfig_updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerResources("softlayer_global_ip.test-global-ip", "routes_to",
						"softlayer_virtual_guest.vm2", "ipv4_address"),
				),
			},

			resource.TestStep{
				Config: testAccCheckSoftLayerGlobalIpConfig_Ipv6Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerGlobalIpExists("softlayer_global_ip.test-global-ip-3"),
					resource.TestMatchResourceAttr("softlayer_global_ip.test-global-ip-3", "ip_address",
						regexp.MustCompile(`^(([[:xdigit:]]{4}:){7})([[:xdigit:]]{4})$`)),
					testAccCheckSoftLayerResources("softlayer_global_ip.test-global-ip-3", "routes_to",
						"softlayer_virtual_guest.vm3", "ipv6_address"),
				),
			},

			resource.TestStep{
				Config: testAccCheckSoftLayerGlobalIpConfig_Ipv6Updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSoftLayerResources("softlayer_global_ip.test-global-ip-3", "routes_to",
						"softlayer_virtual_guest.vm4", "ipv6_address"),
				),
			},
		},
	})
}

func testAccCheckSoftLayerGlobalIpExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		globalIpId, _ := strconv.Atoi(rs.Primary.ID)

		service := services.GetNetworkSubnetIpAddressGlobalService(testAccProvider.Meta().(ProviderConfig).SoftLayerSession())
		foundGlobalIp, err := service.Id(globalIpId).GetObject()

		if err != nil {
			return err
		}

		if strconv.Itoa(int(*foundGlobalIp.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		return nil
	}
}

func testAccCheckSoftLayerResources(srcResource, srcKey, tgtResource, tgtKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		sourceResource, ok := s.RootModule().Resources[srcResource]
		if !ok {
			return fmt.Errorf("Not found: %s", srcResource)
		}

		targetResource, ok := s.RootModule().Resources[tgtResource]
		if !ok {
			return fmt.Errorf("Not found: %s", tgtResource)
		}

		if sourceResource.Primary.Attributes[srcKey] != targetResource.Primary.Attributes[tgtKey] {
			return fmt.Errorf("Different values : Source : %s %s %s , Target : %s %s %s",
				srcResource, srcKey, sourceResource.Primary.Attributes[srcKey],
				tgtResource, tgtKey, targetResource.Primary.Attributes[tgtKey])
		}

		return nil
	}
}

const testAccCheckSoftLayerGlobalIpConfig_basic = `
resource "softlayer_virtual_guest" "vm1" {
    hostname = "vm1"
    domain = "example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "dal06"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
}

resource "softlayer_virtual_guest" "vm2" {
    hostname = "vm2"
    domain = "example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "tor01"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
}

resource "softlayer_global_ip" "test-global-ip" {
    routes_to = "${softlayer_virtual_guest.vm1.ipv4_address}"
}`

const testAccCheckSoftLayerGlobalIpConfig_updated = `
resource "softlayer_virtual_guest" "vm1" {
    hostname = "vm1"
    domain = "example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "dal06"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
}

resource "softlayer_virtual_guest" "vm2" {
    hostname = "vm2"
    domain = "example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "tor01"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
}

resource "softlayer_global_ip" "test-global-ip" {
    routes_to = "${softlayer_virtual_guest.vm2.ipv4_address}"
}`

const testAccCheckSoftLayerGlobalIpConfig_Ipv6Basic = `
resource "softlayer_virtual_guest" "vm3" {
    hostname = "vm3"
    domain = "example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "che01"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
    ipv6_enabled = true
}

resource "softlayer_virtual_guest" "vm4" {
    hostname = "vm4"
    domain = "example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "che01"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
    ipv6_enabled = true
}

resource "softlayer_global_ip" "test-global-ip-3" {
    routes_to = "${softlayer_virtual_guest.vm3.ipv6_address}"
}`

const testAccCheckSoftLayerGlobalIpConfig_Ipv6Updated = `
resource "softlayer_virtual_guest" "vm3" {
    hostname = "vm3"
    domain = "example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "che01"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
    ipv6_enabled = true
}

resource "softlayer_virtual_guest" "vm4" {
    hostname = "vm4"
    domain = "example.com"
    os_reference_code = "DEBIAN_7_64"
    datacenter = "che01"
    network_speed = 100
    hourly_billing = true
    private_network_only = false
    cores = 1
    memory = 1024
    disks = [25]
    local_disk = false
    ipv6_enabled = true
}

resource "softlayer_global_ip" "test-global-ip-3" {
    routes_to = "${softlayer_virtual_guest.vm4.ipv6_address}"
}`
