package spell

var OperationWeights = []OperationWeight{
	{
		1,
		0.8,
		[]int{0,-1},
	},
	{
		2,
		1.1,
		[]int{0},
	},
	{
		3,
		1.9,
		[]int{0},
	},
}

var CheckOperationWeight = append([]OperationWeight{{
	0,
	0.8,
	[]int{1},
}}, OperationWeights...)

var MinSpanningLen = 3.0

