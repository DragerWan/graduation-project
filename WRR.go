package main

import (
	"sort"
	"fmt"
)

var (
	nodes_wrr NodeInfos_WRR
	weight_total = 0
)

type NodeInfos_WRR []NodeInfo

func (ns NodeInfos_WRR) Len() int {
	return len(ns)
}
func (ns NodeInfos_WRR) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}
func (ns NodeInfos_WRR) Less(i, j int) bool {
	return ns[i].RemainBW > ns[j].RemainBW
}

func test_WRR (times, number int){
	nodes_wrr = NodeInfos_WRR(nodes)
	fmt.Println("WRR")

	for _, node := range nodes_wrr {
		weight_total += node.Weight
	}
	//input := bufio.NewScanner(os.Stdin)
	for j := 0; j < times; j++ {
		//n, _ := strconv.Atoi(input.Text())
		for i := 0; i < 1; i++ {
			//getNodes_LoadBalance(30)
			getNodes_WRR(number)
			//fmt.Println(ni)
		}

		fmt.Println(loadBalancingIndex(nodes_wrr))
	}

}

func getNodes_WRR(n int) []*NodeInfo {
	sort.Sort(nodes_wrr)
	ret := make([]*NodeInfo, n)

	for i := 0; i < n; i++ {
		index := -1
		for j, node := range nodes_wrr {
			node.CurrentWeight += node.Weight
			if index == -1 || nodes_wrr[index].CurrentWeight < node.CurrentWeight{
				index = j
			}
		}
		nodes_wrr[index].RemainBW -= 0.2
		nodes_wrr[index].CurrentWeight -= weight_total
		ret[i] = &nodes_wrr[index]
	}

	return ret
}
