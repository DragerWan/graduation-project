package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"encoding/json"
	"strconv"
	"strings"
)

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
		node.RemainBW = node.UploadBW
		node.RemainBW_real = node.RemainBW
		node.Weight = int(node.UploadBW) + 1
		node.CurrentWeight = 0
		node.linkHistory = make([]int, minutes_video)
		/*if (strings.Compare(clientISP, node.Region.ISP) != 0) {
		    node.Dis += 128 // 2^7
		}*/
		nodes = append(nodes, *node)
	}
	nodes = Shuffle(nodes)
}
