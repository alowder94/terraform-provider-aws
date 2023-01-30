package appmesh

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appmesh"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func DataSourceVirtualRouter() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceVirtualRouterRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mesh_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mesh_owner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//Create Spec: required here (for listener)

			"tags": tftags.TagsSchema(),
		},
	}
}

func dataSourceVirtualRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).AppMeshConn()
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	req := &appmesh.DescribeVirtualRouterInput{
		MeshName: aws.String(d.Get("mesh_name").(string)),
		VirtualRouterName: aws.String(d.Get("name").(string)),
	},

	if v, ok := d.GetOk("mesh_owner"); ok {
		req.MeshOwner = aws.String(v.(string))
	},

	resp, err := conn.DescribeVirtualRouterWithContext(ctx, req)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading App Mesh Virtual Router: %s", err)
	}

	arn := aws.StringValue(resp.VirtualRouter.Metadata.Arn)
}
