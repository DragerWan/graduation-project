package main

import (
	"fmt"
	"math/rand"
	//"strings"
	"encoding/json"

	"math"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/lionsoul2014/ip2region/binding/golang"
	"github.com/oschwald/geoip2-golang"
	//"sort"
	"sort"
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
	NodeType  string  `json:"type"`
	ISP       string  `json:"isp"`
	ASN       int     `json:"asn"`
	Dis       uint    `json:"distance"`
}
type NodeInfo1 struct {
	Mac      string `json:"mac_addr"`
	PublicIp string `json:"public_ip"`
}

type NodeInfos []NodeInfo

var (
	nodes = NodeInfos{} //�ڵ���Ϣ
)

func (ns NodeInfos) Len() int {
	return len(ns)
}
func (ns NodeInfos) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}
func (ns NodeInfos) Less(i, j int) bool {
	return ns[i].Dis < ns[j].Dis
}

func test(opt ...interface{}) {
	if len(opt) == 0 {
		fmt.Println("empty")
	} else {
		n1 := opt[0].(float32)
		n2 := opt[1].(float64)
		fmt.Println(n1, n2)
	}
}
func main() {

	//updateRedis()
	//makeIP()
	//testScores()

	c, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("123456"))
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
	}

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

func makeIP() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("123456"))
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()
	c.Do("SELECT", 1)

	//ɾ��������
	values, _ := redis.Values(c.Do("KEYS", "fv:local:ip_nodes*"))
	//for i := 0; i < len(values); i++ {
	for i := 0; i < len(values); i++ {

		_, err := c.Do("DEL", string((values[i]).([]byte)))
		if err != nil {
			fmt.Println("DEL failed: ", err)
			return
		}

	}
	fmt.Println("DEL completed!")
	//��������
	values, _ = redis.Values(c.Do("KEYS", "fv:report:node_infos:*"))
	for i := 0; i < len(values); i++ {

		//val := strings.Replace(string([]byte(string(values[i].([]uint8)))[14:]), "-", ":", -1)
		node := &NodeInfo1{}
		nodeData, err := redis.String(c.Do("GET", string((values[i]).([]byte))))
		if err != nil {
			fmt.Println("Get nodeData failed: ", err)
			continue
		}
		err = json.Unmarshal([]byte(nodeData), node)
		if err != nil {
			fmt.Println("json.Unmarshal failed: ", err)
			continue
		}

		_, err = c.Do("SADD", "fv:local:ip_nodes:"+node.PublicIp, node.Mac)
		if err != nil {
			fmt.Println("Redis SADD failed: ", err)
			return
		}

	}
	fmt.Println("Update completed!")
}

func testScores() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("123456"))
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()
	c.Do("SELECT", 1)

	//ɾ��������
	values, _ := redis.Strings(c.Do("ZRANGE", "fv:cdn:file_nodes:qq.webrtc.win/tv/Pear-Demo-Corporate-Video.mp4", 0, -1, "WITHSCORES"))
	//for i := 0; i < len(values); i++ {
	/*for i := 0; i < len(values); i++ {
		if i/2 == 0 {
			fmt.Println(string((values[i]).([]byte)))
		}else {
		fmt.Println((values[i]).(int))
		}
	}*/
	fmt.Println(len(values))
	for i := 0; i < len(values); i++ {

		if i%2 == 0 {
			fmt.Print("even ")
			fmt.Println(values[i])
		} else {
			fmt.Print("odd ")
			fmt.Println(strconv.Atoi(values[i]))
		}

	}
	fmt.Println("End")

}

func updateRedis() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("123456"))
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()
	c.Do("SELECT", 1)

	//ɾ��������
	values, _ := redis.Values(c.Do("KEYS", "fv:local:small:node_infos*"))
	//for i := 0; i < len(values); i++ {
	for i := 0; i < len(values); i++ {

		_, err := c.Do("DEL", string((values[i]).([]byte)))
		if err != nil {
			fmt.Println("DEL failed: ", err)
			return
		}

	}
	fmt.Println("DEL completed!")
	//��������
	values, _ = redis.Values(c.Do("KEYS", "fv:report:node_infos:*"))
	for i := 0; i < len(values); i += 100 {

		//val := strings.Replace(string([]byte(string(values[i].([]uint8)))[14:]), "-", ":", -1)
		node := &NodeInfo{}
		nodeData, err := redis.String(c.Do("GET", string((values[i]).([]byte))))
		if err != nil {
			fmt.Println("Get nodeData failed: ", err)
			continue
		}
		err = json.Unmarshal([]byte(nodeData), node)
		if err != nil {
			fmt.Println("json.Unmarshal failed: ", err)
			continue
		}
		node.ISP, _ = getISP(node.PublicIp)
		node.ASN, _ = getASN(node.PublicIp)

		if node.UploadBW > 0.01 && (node.ISP == "����" || node.ISP == "�ƶ�" || node.ISP == "��ͨ") { //ֻ������upload_bw��������Ӫ�̵Ľڵ�
			nodeJS, err := json.Marshal(node)
			if err != nil {
				fmt.Println("json.Marshal failed: ", err)
				continue
			}
			_, err = c.Do("SET", "fv:local:small:node_infos:"+strings.Replace(node.Mac, ":", "-", -1), nodeJS)
			if err != nil {
				fmt.Println("Redis SET failed: ", err)
				return
			}
		}

	}
	fmt.Println("Update completed!")
}

func makeNewSet() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("123456"))
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()
	c.Do("SELECT", 1)

	//ɾ��������
	_, err = c.Do("ZREMRANGEBYRANK", "fv:local:file_nodes:qq.webrtc.win/tv/Pear-Demo-Yosemite_National_Park.mp4", 0, -1)
	if err != nil {
		fmt.Println("ZREMRANGEBYRANK failed: ", err.Error())
	}

	fmt.Println("DEL completed!")
	//��������
	values, _ := redis.Values(c.Do("KEYS", "fv:report:node_infos:*"))
	for i := 0; i < len(values); i += 100 {

		//val := strings.Replace(string([]byte(string(values[i].([]uint8)))[14:]), "-", ":", -1)
		node := &NodeInfo{}
		nodeData, err := redis.String(c.Do("GET", string((values[i]).([]byte))))
		if err != nil {
			fmt.Println("Get nodeData failed: ", err)
			continue
		}
		err = json.Unmarshal([]byte(nodeData), node)
		if err != nil {
			fmt.Println("json.Unmarshal failed: ", err)
			continue
		}
		node.ISP, _ = getISP(node.PublicIp)
		node.ASN, _ = getASN(node.PublicIp)

		ispNum, err := redis.Int(c.Do("HGET", "fv:local:isp_num", node.ISP))
		if err != nil {
			ispMax, err := redis.Int(c.Do("HGET", "fv:local:isp_num", "max"))
			if err != nil {
				ispMax = 0
				ispNum = 0
			}
			_, err = c.Do("HSET", "fv:local:isp_num", node.ISP, ispMax)
			_, err = c.Do("HSET", "fv:local:isp_num", "max", ispMax+1)
		}
		score := ispNum << 40
		score += node.ASN << 32 //TODO ����IP��ת��

		if node.UploadBW > 0.01 && (node.ISP == "����" || node.ISP == "�ƶ�" || node.ISP == "��ͨ") { //ֻ������upload_bw��������Ӫ�̵Ľڵ�
			nodeJS, err := json.Marshal(node)
			if err != nil {
				fmt.Println("json.Marshal failed: ", err)
				continue
			}
			_, err = c.Do("SET", "fv:local:small:node_infos:"+strings.Replace(node.Mac, ":", "-", -1), nodeJS)
			if err != nil {
				fmt.Println("Redis SET failed: ", err)
				return
			}
		}

	}
	fmt.Println("Update completed!")
}

func initNodes() {
	c, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("123456"))
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return
	}
	defer c.Close()
	c.Do("SELECT", 1)

	values, _ := redis.Values(c.Do("KEYS", "fv:local:small:node_infos:*"))
	if len(values) == 0 {
		updateRedis()
		values, _ = redis.Values(c.Do("KEYS", "fv:local:small:node_infos:*"))
	}
	for i := 0; i < len(values); i++ {

		//val := strings.Replace(string([]byte(string(values[i].([]uint8)))[14:]), "-", ":", -1)
		node := &NodeInfo{}
		nodeData, err := redis.String(c.Do("GET", string((values[i]).([]byte))))
		if err != nil {
			fmt.Println("Get nodeData failed: ", err)
			continue
		}
		err = json.Unmarshal([]byte(nodeData), node)
		if err != nil {
			fmt.Println("json.Unmarshal failed: ", err)
			continue
		}
		/*if (strings.Compare(clientISP, node.Region.ISP) != 0) {
		    node.Dis += 128 // 2^7
		}*/
		nodes = append(nodes, *node)
	}
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

func getNodes_mine(ipStr string, n int) []NodeInfo {

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
}

func getNodes_random(n int) []NodeInfo {
	dataLen := len(nodes)
	shuffled := make([]NodeInfo, dataLen)
	for i, j := range rand.Perm(dataLen) {
		shuffled[j] = nodes[i]
	}
	return shuffled[:n]
}

func getNodes_ip(ipStr string, n int) []NodeInfo {

	for i := 0; i < len(nodes); i++ {
		nodes[i].Dis = DisIP(ipStr, nodes[i].PublicIp)
	}
	sort.Sort(nodes)
	return nodes[:n]
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
