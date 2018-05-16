package main

import (
	"math/rand"
	"fmt"
)

func test_Random(times, number int) {
	fmt.Println("Random")
	//input := bufio.NewScanner(os.Stdin)
	for j := 0; j < times; j++ {
		//n, _ := strconv.Atoi(input.Text())
		for i := 0; i < 1; i++ {
			//getNodes_LoadBalance(30)
			getNodes_Random(number)
			//fmt.Println(ni)
		}

		fmt.Println(loadBalancingIndex(nodes))
	}
}

func getNodes_Random(n int) []*NodeInfo {

	length := len(nodes)
	ret := make([]*NodeInfo, n)
	for i := 0; i < n; i++ {
		index := rand.Intn(length)
		nodes[index].RemainBW -= 0.2
		ret[i] = &nodes[index]
	}

	return ret
}