package main

import (
	"fmt"
	"math/rand"
	//"strings"

	"math"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang"
	"github.com/oschwald/geoip2-golang"
	//"sort"
	//"bufio"
	"time"
)

/*func passwd(d * redis.dialOptions){
    d.password = "123456"
}
61.144.172.130
*/

type NodeInfo struct {
	Mac       string  `json:"mac_addr"`
	PublicIp  string  `json:"public_ip"`
	LocalIp   string  `json:"Local_ip"`
	HttpPort  int     `json:"http_port"`
	HttpsPort int     `json:"https_port"`
	UploadBW  float32 `json:"upload_bw"`
	RemainBW  float32
	RemainBW_real float32
	NodeType  string `json:"type"`
	ISP       string `json:"isp"`
	ASN       int    `json:"asn"`
	Dis       uint   `json:"distance"`
	Weight    int
	CurrentWeight int
	linkHistory []int
}
type NodeInfo1 struct {
	Mac      string `json:"mac_addr"`
	PublicIp string `json:"public_ip"`
}

type NodeInfos []NodeInfo


var (
	nodes = NodeInfos{} //服务器预测节点信息
	return_nodes_num = 30 //每个请求返回节点数
	use_nodes_num = 5//实际使用节点数
	usable_rate = float32(0.8) //每个节点可用概率
	predicted_usage = float32(0.2) //预测每个节点被调用一次占用的带宽
	real_usage = float32(0.5) //实际每个节点被使用一次占用的带宽
	minutes_total = 20//总分钟数
	minutes_video = 3//视频播放时间
	minutes_update = 1//节点信息同步时间间隔
	request_per_minute = 500//每分钟请求数
)



func test(opt ...interface{}) {
	if len(opt) == 0 {
		fmt.Println("empty")
	} else {
		n1 := opt[0].(float32)
		n2 := opt[1].(float64)
		fmt.Println(n1, n2)
	}
}

func updateNodeInfo() {
	for i := 0; i < len(nodes); i++ {
		nodes[i].RemainBW = nodes[i].RemainBW_real
	}
}
//每分钟更新节点服务完成回收的带宽
func checkNodeInfo() {
	//fmt.Println("checkNodeInfo")
	//fmt.Print(loadBalancingIndexReal(nodes), "")
	for i := 0; i < len(nodes); i++ {
		if nodes[i].linkHistory[minutes_video - 1] != 0 {
			nodes[i].RemainBW_real += float32(nodes[i].linkHistory[minutes_video - 1]) * real_usage
		}
		for j := minutes_video - 1; j > 0; j-- {
			nodes[i].linkHistory[j] = nodes[i].linkHistory[j - 1]
		}
		nodes[i].linkHistory[0] = 0
	}
	//fmt.Println(loadBalancingIndexReal(nodes))
}

func main() {
	rand.Seed(time.Now().Unix())

	initNodes()


	for i := 0; i <= minutes_total; i++ {
		if i % minutes_update == 0 && i > 0 {
			//fmt.Println("updateNodeInfo")
			//fmt.Print(loadBalancingIndex(nodes), "")
			updateNodeInfo()
			//fmt.Println(loadBalancingIndex(nodes))
		}
		if i > 0 {
			checkNodeInfo()
		}

		test_Balance1(request_per_minute, 30)
	}
	//updateRedis()
	//makeIP()
	//testScores()

	/*c, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("123456"))
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()
	c.Do("SELECT", 1)
	res, err := redis.Strings(c.Do("ZRANGEBYLEX", "test", "(11", "+", "LIMIT", 0, 2))
	fmt.Println(res)
	fmt.Println(err)
	if err != nil {
		fmt.Println("err is not nil")
	}*/

	//test()

	/*initNodes()
	  fmt.Println("Random")
	  ni := getNodes_random(200)
	  fmt.Println(ni)
	  testRTT(ni, 20)
	  fmt.Println("IP")
	  ni = getNodes_ip("61.141.252.238", 200)
	  fmt.Println(ni)
	  testRTT(ni, 20)
	  fmt.Println("Mine")
	  ni = getNodes_mine("61.141.253.218", 200)
	  fmt.Println(ni)
	  testRTT(ni, 20)*/






	//clientIP:= "61.141.252.238"

	/*nodes = nodes[:20]
	  fmt.Println("Results:\n", nodes)
	  fmt.Printf("Len: %v\n", len(nodes))
	  for i := 0; i < len(nodes); i++{
	      ip := strings.Replace(nodes[i].Mac, ":", "", -1) + ".webrtc.win"
	      resTime := ping(ip, 4, 32, 1000, false)
	      fmt.Printf("ping [ %s ]: %dms\n", ip, resTime)
	  }*/

}









func Shuffle(vals []NodeInfo) []NodeInfo {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]NodeInfo, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

func getISP(ipStr string) (string, error) {
	db := "ip2region.db"

	_, err := os.Stat(db)
	if os.IsNotExist(err) {
		panic("not found db " + db)
	}
	region, err := ip2region.New(db)
	defer region.Close()

	ip := ip2region.IpInfo{}

	ip, err = region.MemorySearch(ipStr)

	if err != nil {
		fmt.Println(fmt.Sprintf("\x1b[0;31m%s\x1b[0m", err.Error()))
	}
	return ip.ISP, err
}

func getASN(ipStr string) (asn int, err error) {
	db, err := geoip2.Open("GeoLite2-ASN.mmdb")
	if err != nil {
		fmt.Println("geoip2 open db failed: ", err)
		return 0, err
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(ipStr)
	record, err := db.ASN(ip)
	if err != nil {
		fmt.Println("geoip2 get record failed: ", err)
		return 0, err
	}
	asn = int(record.AutonomousSystemNumber)
	return
}

func DisIP(ip_a, ip_b string) (dis uint) {
	parts_a := strings.Split(ip_a, ".")
	parts_b := strings.Split(ip_b, ".")
	for i := 0; i < 4; i++ {
		a, _ := strconv.Atoi(parts_a[i])
		b, _ := strconv.Atoi(parts_b[i])
		c := a ^ b
		if c == 0 {
			continue
		}
		dis = uint(int(math.Sqrt(float64(a^b))) + 1 + 8*(3-i))
		return
	}
	return
}

/*func getNodes_mine(ipStr string, n int) []NodeInfo {

	clientISP, err := getISP(ipStr)
	if err != nil {
		fmt.Println("getISP failed: ", err)
		return nil
	}
	clientASN, err := getASN(ipStr)
	if err != nil {
		fmt.Println("getASN failed: ", err)
		return nil
	}

	for i := 0; i < len(nodes); i++ {
		nodes[i].Dis = 0
		if strings.Compare(clientISP, nodes[i].ISP) != 0 {
			nodes[i].Dis += 128
		}
		if clientASN != nodes[i].ASN {
			nodes[i].Dis += 64
		}
		nodes[i].Dis += DisIP(ipStr, nodes[i].PublicIp)
	}
	sort.Sort(nodes)
	return nodes[:n]
}*/



/*func getNodes_ip(ipStr string, n int) []NodeInfo {

	for i := 0; i < len(nodes); i++ {
		nodes[i].Dis = DisIP(ipStr, nodes[i].PublicIp)
	}
	sort.Sort(nodes)
	return nodes[:n]
}*/





func loadBalancingIndex( nodeInfos []NodeInfo) float32 {
	var a float32
	var b float32
	var c float32
	for _, node := range nodeInfos {
		c = (node.UploadBW - node.RemainBW) / node.UploadBW
		if c > 1 {
			c = 1
			node.RemainBW = 0
		}
		a += c
		b += c * c
	}
	return (a * a) / (float32(len(nodes)) * b)

}
func loadBalancingIndexReal( nodeInfos []NodeInfo) float32 {
	var a float32
	var b float32
	var c float32
	for _, node := range nodeInfos {
		c = (node.UploadBW - node.RemainBW_real) / node.UploadBW
		if c > 1 {
			c = 1
			node.RemainBW_real = 0
		}
		a += c
		b += c * c
	}
	if b == 0 {
		fmt.Println("b = 0!")
	}
	return (a * a) / (float32(len(nodes)) * b)

}

func testRTT(ni []NodeInfo, sgm int) {
	rtt := make(chan int)
	lossRate := make(chan float32)
	jump := make(chan int)
	num := 0
	sumRTT := 0
	sumLossRate := float32(0)
	sumJump := 0
	start := 0
	end := int(math.Min(float64(start+sgm), float64(len(ni))))
	for start < len(ni) {
		for i := start; i < end; i++ {
			ip := strings.Replace(ni[i].Mac, ":", "", -1) + ".webrtc.win"
			go ping(ip, 4, 32, 1000, false, rtt, lossRate, jump)
		}

		for i := 0; i < end-start; i++ {
			r := <-rtt
			l := <-lossRate
			j := <-jump
			if r != 0 {
				num++
				sumRTT += r
				sumLossRate += l
				sumJump += j
			}
			/*if(num == 20){
				break
			}*/
		}
		start = end
		end = int(math.Min(float64(start+sgm), float64(len(ni))))
	}

	fmt.Printf("%d of %d nodes are usable\n", num, len(ni))
	fmt.Printf("Average RTT: %d, jump number: %f, package loss rate: %f, connection rate: %f\n", sumRTT/num, float32(sumJump)/float32(num), sumLossRate/float32(num), float32(num)/float32(len(ni)))

}
