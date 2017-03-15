package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"os"
	"time"
)

var ACCESS_KEY_ID = os.Getenv("ACCESS_KEY_ID")
var ACCESS_KEY_SECRET = os.Getenv("ACCESS_KEY_SECRET")

const region = common.Region("cn-shanghai")

func main() {
	client := ecs.NewClient(ACCESS_KEY_ID, ACCESS_KEY_SECRET)

	// 获取vpcId
	vpcs, _, _ := client.DescribeVpcs(&ecs.DescribeVpcsArgs{RegionId: region})
	vpcId := vpcs[0].VpcId

	// 获取 zoneId
	zones, _ := client.DescribeZones(region)
	zoneId := zones[0].ZoneId

	spew.Dump(vpcId)
	spew.Dump(zoneId)

	vswidthId, _ := client.CreateVSwitch(&ecs.CreateVSwitchArgs{ZoneId: zoneId, VpcId: vpcId, CidrBlock: "10.99.0.0/24"})

	spew.Dump(vswidthId)
	if vswidthId == "" {
		return
	}

	vswitches, _, _ := client.DescribeVSwitches(&ecs.DescribeVSwitchesArgs{VpcId: vpcId, VSwitchId: vswidthId})
	for vswitches[0].Status == "Pending" {
		spew.Dump("Sleep 5 second")
		time.Sleep(5 * time.Second)
		vsws, _, _ := client.DescribeVSwitches(&ecs.DescribeVSwitchesArgs{VpcId: vpcId, VSwitchId: vswidthId})
		vswitches = vsws
	}

	instanceId, e := client.CreateInstance(&ecs.CreateInstanceArgs{RegionId: region, ImageId: "centos_7_2_64_40G_base_20170222.vhd", InstanceType: "ecs.n1.tiny", HostName: "test-ecs", VSwitchId: vswidthId})

	spew.Dump("Done", instanceId, e)
}
