package main

import "fmt"

var index_RR = -1

func test_RR(times, number int) {
	index_RR = -1
	fmt.Println("RR")
	//input := bufio.NewScanner(os.Stdin)
	for j := 0; j < times; j++ {
		//n, _ := strconv.Atoi(input.Text())
		for i := 0; i < 1; i++ {
			//getNodes_LoadBalance(30)
			getNodes_RR(number)
			//fmt.Println(ni)
		}

		fmt.Println(loadBalancingIndex(nodes))
	}
}
func getNodes_RR(n int) []*NodeInfo {
	length := len(nodes)
	ret := make([]*NodeInfo, n)
	for i := 0; i < n; i++ {
		index_RR = (index_RR + 1) % length
		nodes[index_RR].RemainBW -= 0.2
		ret[i] = &nodes[index_RR]
	}

	return ret
}
