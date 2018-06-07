package Utilities

import (
	"math/rand"
	"reflect"
	"time"
)

// from: https://gist.github.com/rafkhan/6501567

type mapf func(interface{}) interface{}
type reducef func(interface{}, interface{}) interface{}
type filterf func(interface{}) bool

// Map(slice, func)
func Map(in interface{}, fn mapf) interface{} {
	val := reflect.ValueOf(in)
	out := make([]interface{}, val.Len())

	for i := 0; i < val.Len(); i++ {
		out[i] = fn(val.Index(i).Interface())
	}

	return out
}

// Reduce(slice, starting value, func)
func Reduce(in interface{}, memo interface{}, fn reducef) interface{} {
	val := reflect.ValueOf(in)

	for i := 0; i < val.Len(); i++ {
		memo = fn(val.Index(i).Interface(), memo)
	}

	return memo
}

// Filter(slice, predicate func)
func Filter(in interface{}, fn filterf) interface{} {
	val := reflect.ValueOf(in)
	out := make([]interface{}, 0, val.Len())

	for i := 0; i < val.Len(); i++ {
		current := val.Index(i).Interface()

		if fn(current) {
			out = append(out, current)
		}
	}

	return out
}

// RandomIn is for random number in unsign integer in 64 bit length
// ref: https://stackoverflow.com/a/47865090
func RandomIn(min, max uint64) uint64 {
	const maxInt64 uint64 = 1<<63 - 1
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	n := max - min
	if n < maxInt64 {
		return uint64(r.Int63n(int64(n + 1)))
	}
	x := r.Uint64()
	for x > n {
		x = r.Uint64()
	}
	return x
}
