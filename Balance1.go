package main

import (
	"sort"
	"fmt"
)

type NodeInfos_Balance1 []NodeInfo
var nodes_balance1 NodeInfos_Balance1

func (ns NodeInfos_Balance1) Len() int {
	return len(ns)
}
func (ns NodeInfos_Balance1) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}
func (ns NodeInfos_Balance1) Less(i, j int) bool {
	return ns[i].RemainBW/ns[i].UploadBW > ns[j].RemainBW/ns[j].UploadBW
	//return ns[i].RemainBW/ns[i].UploadBW > ns[j].RemainBW/ns[j].UploadBW || (ns[i].RemainBW/ns[i].UploadBW == ns[j].RemainBW/ns[j].UploadBW && ns[i].RemainBW > ns[j].RemainBW)
}

func test_Balance1 (times, number int){
	nodes_balance1 = NodeInfos_Balance1(nodes)
	//fmt.Println("Balance1")

	//input := bufio.NewScanner(os.Stdin)
	for j := 0; j < times; j++ {
		//n, _ := strconv.Atoi(input.Text())
		for i := 0; i < 1; i++ {
			//getNodes_LoadBalance(30)
			ns := getNodes_Balance1(number)
			useNodes(ns)
			//fmt.Println(ni)
		}

		//fmt.Println(loadBalancingIndex(nodes_balance1), loadBalancingIndexReal(nodes_balance1))
	}
	fmt.Println(loadBalancingIndex(nodes_balance1), loadBalancingIndexReal(nodes_balance1))

}

func getNodes_Balance1(n int) []*NodeInfo {
	sort.Sort(nodes_balance1)
	ret := make([]*NodeInfo, n)
	for i := 0; i < n; i++ {
		ret[i] = &nodes_balance1[i]
		nodes_balance1[i].RemainBW -= 0.2
		//nodes_balance1[i].RemainBW_real = 0
	}

	return ret
}
