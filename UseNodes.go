package main

import (
	"math/rand"
	//"fmt"
)

func isUsable(p float32) bool {
	n := rand.Float32()
	return n < p
}

func useNodes(ns []*NodeInfo) {
	used := 0
	for i := 0; i < len(ns); i++ {
		//fmt.Println(i, (*node).Mac)
		if isUsable(usable_rate){
			//fmt.Println("Use a node!")
			//fmt.Print((*node).RemainBW_real, " ")
			(*ns[i]).RemainBW_real -= real_usage
			(*ns[i]).linkHistory[0] ++
			//fmt.Println((*node).RemainBW_real)
			used ++
			if(used >= use_nodes_num){
				break
			}
		}
	}
}