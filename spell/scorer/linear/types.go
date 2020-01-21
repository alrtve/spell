package linear

type LearningAlgorithm interface {
	Learn(vectorSystems []*VectorSystem) *Vector
}
