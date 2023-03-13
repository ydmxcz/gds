package truple

type KV[K, V any] struct {
	Val V
	Key K
}
