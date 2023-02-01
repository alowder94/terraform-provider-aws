package appmesh_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/appmesh"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccAppMeshVirtualRouterDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	virtualRouterName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_appmesh_virtual_router.test"
	dataSourceName := "data.aws_appmesh_virtual_router.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(appmesh.EndpointsID, t) },
		ErrorCheck:               acctest.ErrorCheck(t, appmesh.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy: 			  testAccCheckVirtualRouterDestroy(ctx),
	})
}
