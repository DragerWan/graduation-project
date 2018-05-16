package main

import (
	"sort"
	"fmt"
)

type NodeInfos_Balance2 []NodeInfo
var nodes_balance2 NodeInfos_Balance2

func (ns NodeInfos_Balance2) Len() int {
	return len(ns)
}
func (ns NodeInfos_Balance2) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}
func (ns NodeInfos_Balance2) Less(i, j int) bool {
	//return ns[i].RemainBW/ns[i].UploadBW > ns[j].RemainBW/ns[j].UploadBW
	return ns[i].RemainBW/ns[i].UploadBW > ns[j].RemainBW/ns[j].UploadBW || (ns[i].RemainBW/ns[i].UploadBW == ns[j].RemainBW/ns[j].UploadBW && ns[i].RemainBW > ns[j].RemainBW)
}

func test_Balance2 (times, number int){
	nodes_balance2 = NodeInfos_Balance2(nodes)
	fmt.Println("Balance2")

	//input := bufio.NewScanner(os.Stdin)
	for j := 0; j < times; j++ {
		//n, _ := strconv.Atoi(input.Text())
		for i := 0; i < 1; i++ {
			//getNodes_LoadBalance(30)
			getNodes_Balance2(number)
			//fmt.Println(ni)
		}

		fmt.Println(loadBalancingIndex(nodes_balance2))
	}

}

func getNodes_Balance2(n int) []*NodeInfo {
	sort.Sort(nodes_balance2)
	ret := make([]*NodeInfo, n)
	for i, node := range nodes_balance2[:n] {
		ret[i] = &node
	}
	for i := 0; i < n; i++ {
		nodes_balance2[i].RemainBW -= 0.2
	}

	return ret
}
