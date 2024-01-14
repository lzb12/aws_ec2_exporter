package main

import (
	awsec2 "aws_ec2_exporter/AwsEC2"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	StatusUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ec2_status",
		Help: "1 if the ec2_status is up, 0 if it's down.",
	},
		[]string{"region", "ip", "instanceid"},
	)
)

func main() {

	var region string
	flag.StringVar(&region, "r", "ap-northeast-1", "aws region,The default is ap-northeast-1")
	flag.Parse()

	parts := strings.Split(region, ",")

	prometheus.MustRegister(StatusUp)
	var wg sync.WaitGroup
	wg.Add(1)
	for _, part := range parts {
		go monitorProxyStatus(&wg, part)
	}

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9097", nil)
	wg.Wait()

	//awsec2.GetEc2Status(region, corpid, corpsecret, toUser, agentid)

}

func monitorProxyStatus(wg *sync.WaitGroup, region string) {
	defer wg.Done()
	for {
		data := awsec2.GetEc2StatusMonitor(region)

		//fmt.Println(data)
		for _, v := range data {
			StatusUp.WithLabelValues(v.Region, v.Ip, v.InstanceId).Set(v.Status)
		}
		time.Sleep(3 * time.Minute)

	}

}
