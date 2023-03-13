package constraints

type SortAble[T any] interface {
	Len() int

	Less(i, j int) bool

	Swap(i, j int)
}
