package ecs_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccECSServiceListDataSource_basic(t *testing.T) {
	resourceName1 := "aws_ecs_service.test"
	resourceName2 := "aws_ecs_service.test2"
	dataSourceName := "data.aws_ecs_service_list.test"
	//Changed for testing purposes...seeing if using name in describeServices call fixes my problem
	serviceName1 := "test" //sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	serviceName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{

		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, ecs.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccServicesDataSourceConfig_basic(serviceName1, serviceName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("aws_ecs_cluster.default", "arn", dataSourceName, "cluster_arn"),
					resource.TestCheckResourceAttrPair(resourceName1, "name", dataSourceName, "service_list[0].service_name"),
					resource.TestCheckResourceAttrPair(resourceName2, "name", dataSourceName, "service_list[1].service_name"),
				),
			},
		},
	})
}

// Split this up just like it is in ec2.vpc_subets_data_source_test.go...base set up and a secondary config for the test. Will require a change in the Steps above as well
func testAccServicesDataSourceConfig_basic(serviceName1 string, serviceName2 string) string {
	return fmt.Sprintf(`
resource "aws_ecs_cluster" "default" {
	name = "test"
}	

resource "aws_ecs_task_definition" "test" {
	family = "test_family"
  
	container_definitions = <<DEFINITION
  [
	{
	  "cpu": 128,
	  "essential": true,
	  "image": "mongo:latest",
	  "memory": 128,
	  "name": "mongodb"
	}
  ]
  DEFINITION
}

resource "aws_ecs_service" "test" {
	name = %[1]q
	cluster = aws_ecs_cluster.default.id
	task_definition = aws_ecs_task_definition.test.arn
	desired_count   = 1
}

resource "aws_ecs_service" "test2" {
	name = %[2]q
	cluster = aws_ecs_cluster.default.id
	task_definition = aws_ecs_task_definition.test.arn
	desired_count   = 1
}

data "aws_ecs_service_list" "test" {
	cluster_arn = aws_ecs_cluster.default.arn

	depends_on = [aws_ecs_service.test, aws_ecs_service.test2]
}
`, serviceName1, serviceName2)
}
