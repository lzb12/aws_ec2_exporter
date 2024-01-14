package awsec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var svc *ec2.EC2

type AwsEC2 struct {
	Region     string
	InstanceId string
	Ip         string
	Status     float64
}

func GetInstances(Region string) (*ec2.DescribeInstanceStatusOutput, error) {
	sess, _ := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("", "", ""),
		Region:      aws.String(Region),
	})
	svc = ec2.New(sess)

	result, err := svc.DescribeInstanceStatus(nil)
	// result, err := svc.DescribeInstances(nil)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetEc2Ipaddr(instanceids string) (publicipaddr, privateipaddr string) {

	result, err := svc.DescribeInstances(nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, r := range result.Reservations {
		for _, i := range r.Instances {
			if *i.InstanceId == instanceids {
				// fmt.Println(*i.PublicIpAddress)
				// fmt.Println(*i.PrivateIpAddress)
				return *i.PublicIpAddress, *i.PrivateIpAddress
			}
		}

	}
	return "", ""
}

func GetEc2StatusMonitor(Region string) []AwsEC2 {
	result, err := GetInstances(Region)
	if err != nil {
		fmt.Println("Got an error retrieving information about your Amazon EC2 instances:")
		fmt.Println(err.Error())
	}
	var data []AwsEC2
	for _, r := range result.InstanceStatuses {
		//fmt.Println(r) 打印所有返回信息
		if *r.InstanceState.Name == "running" && Region == "ap-southeast-2" {

			if *r.InstanceStatus.Status != "ok" || *r.SystemStatus.Status != "ok" {
				// fmt.Println(*r.InstanceId)
				_, privateipaddr := GetEc2Ipaddr(*r.InstanceId)
				data2 := AwsEC2{
					Region:     Region,
					InstanceId: *r.InstanceId,
					Ip:         privateipaddr,
					Status:     0,
				}
				data = append(data, data2)
				//return Region, *r.InstanceId, privateipaddr

			} else if *r.InstanceStatus.Status == "ok" && *r.SystemStatus.Status == "ok" {
				_, privateipaddr := GetEc2Ipaddr(*r.InstanceId)
				data2 := AwsEC2{
					Region:     Region,
					InstanceId: *r.InstanceId,
					Ip:         privateipaddr,
					Status:     1,
				}
				data = append(data, data2)
			}

		} else if *r.InstanceState.Name == "running" {
			if *r.InstanceStatus.Status != "ok" || *r.SystemStatus.Status != "ok" {
				// fmt.Println(*r.InstanceId)
				publicipaddr, _ := GetEc2Ipaddr(*r.InstanceId)
				data2 := AwsEC2{
					Region:     Region,
					InstanceId: *r.InstanceId,
					Ip:         publicipaddr,
					Status:     0,
				}
				data = append(data, data2)
				//return Region, *r.InstanceId, publicipaddr
			} else if *r.InstanceStatus.Status == "ok" && *r.SystemStatus.Status == "ok" {
				publicipaddr, _ := GetEc2Ipaddr(*r.InstanceId)
				data2 := AwsEC2{
					Region:     Region,
					InstanceId: *r.InstanceId,
					Ip:         publicipaddr,
					Status:     1,
				}
				data = append(data, data2)
			}
		} else {
			fmt.Println(fmt.Sprintf("AWS/EC2停止\n主机:%s", *r.InstanceId))
			//return Region, *r.InstanceId, ""
		}
	}
	return data

}
