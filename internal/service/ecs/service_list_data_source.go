package ecs

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

func DataSourceServiceList() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceServiceListRead,

		Schema: map[string]*schema.Schema{
			//This is an attempt to fix the "no id found in attributes" error
			"cluster_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			// "filter": CustomFiltersSchema(),
			"service_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					ReadWithoutTimeout: dataSourceServiceRead,

					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// "arn": {
						// 	Type:     schema.TypeString,
						// 	Computed: true,
						// },
						// "cluster_arn": {
						// 	Type:     schema.TypeString,
						// 	Computed: true,
						// },
						// "desired_count": {
						// 	Type:     schema.TypeInt,
						// 	Computed: true,
						// },
						// "launch_type": {
						// 	Type:     schema.TypeString,
						// 	Computed: true,
						// },
						// "scheduling_strategy": {
						// 	Type:     schema.TypeString,
						// 	Computed: true,
						// },
						// "task_definition": {
						// 	Type:     schema.TypeString,
						// 	Computed: true,
						// },
						// "tags": tftags.TagsSchemaComputed(),
					},
				}},
		},
	}
}

func dataSourceServiceListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ECSConn()

	//Check this is the for loop when setting tags of the individual services within the list
	// ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	//Setting cluster arn from input (d)
	clusterArn := d.Get("cluster_arn").(string)

	d.SetId(aws.StringValue(&clusterArn))

	params := &ecs.ListServicesInput{
		Cluster: aws.String(clusterArn),
	}

	// //Applies filters to the sdk query
	// if filters, filtersOk := d.GetOk("filter"); filtersOk {
	// 	params.Filters = append(input.Filters,
	// 		BuildCustomFiltersDataSource(filters.(*schema.Set))...)
	// }

	//If filters param is set,but no filters are present, set to nil
	// if len(params.Filters) == 0 {
	// 	params.Filters = nil
	// }

	//This only returns a max of 10 services, and uses nextToken for pagination...this will require some looping logic as a single cluster can maintain up to 5000 services in a single region
	//This also means that we need to implement AS MUCH filtering as possible when making this initial call for performance purposes
	log.Printf("[DEBUG] Reading all services for cluster_arn: %s", clusterArn)
	desc, err := conn.ListServicesWithContext(ctx, params)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading services in cluster: %s", clusterArn)
	}

	if len(desc.ServiceArns) == 0 || desc == nil {
		return sdkdiag.AppendErrorf(diags, "Either no services matching filters found in cluster: %s, or cluster not found", clusterArn)
	}
	fmt.Printf("desc line 99: %s", *desc.ServiceArns[0])

	var serviceList []schema.ResourceData

	//Parse list of arns - pull service data using sdk - use that data to set attributes of each service
	for _, v := range desc.ServiceArns {
		newService := schema.ResourceData{}

		//I think the problem here is this should accept service name instaed of service arn?
		params := &ecs.DescribeServicesInput{
			Cluster:  aws.String(clusterArn),
			Services: []*string{aws.String(*v)},
		}

		describeServicesResponse, err := conn.DescribeServicesWithContext(ctx, params)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "Reading ECS Service %s", *v)
		}

		//Pull servcie from the response
		serviceData := describeServicesResponse.Services[0]
		fmt.Printf("serviceData.ServiceName ln 128 %s", *serviceData.ServiceName)

		//Maybe try reimplementing the map here? Might be able to get that to work now - not sure exactly what TF guys will have to say about it though...
		newService.Set("service_name", *serviceData.ServiceName)
		// newService.Set("arn", serviceData.ServiceArn)
		// newService.Set("cluster_arn", serviceData.ClusterArn)
		// newService.Set("desired_count", serviceData.DesiredCount)
		// newService.Set("launch_type", serviceData.LaunchType)
		// newService.Set("scheduling_strategy", serviceData.SchedulingStrategy)
		// newService.Set("task_definition", serviceData.TaskDefinition)
		// if err := newService.Set("tags", KeyValueTags(serviceData.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		// 	return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
		// }

		// Append each service to serviceList
		// This functionality will need tested...not sure that this will work
		serviceList = append(serviceList, newService)
		fmt.Printf("newService.service_name ln 144 %s", newService.Get("service_name"))
	}

	d.Set("service_list", serviceList)
	fmt.Print(d.Get("service_list"))

	return diags
}
