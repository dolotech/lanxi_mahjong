package utils

import (
	"math/rand"
	"sync"
)

var o *rand.Rand = rand.New(rand.NewSource(TimestampNano()))
var random_mux_ sync.Mutex

func RandInt64() (r int64) {
	random_mux_.Lock()
	r = o.Int63()
	random_mux_.Unlock()
	return
}

func RandInt32() (r int32) {
	random_mux_.Lock()
	r = o.Int31()
	random_mux_.Unlock()
	return
}

func RandUint32() (r uint32) {
	random_mux_.Lock()
	r = o.Uint32()
	random_mux_.Unlock()
	return
}

func RandInt64N(n int64) (r int64) {
	random_mux_.Lock()
	r = o.Int63n(n)
	random_mux_.Unlock()
	return
}

func RandInt32N(n int32) (r int32) {
	random_mux_.Lock()
	r = o.Int31n(n)
	random_mux_.Unlock()
	return
}
