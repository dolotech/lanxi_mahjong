package data

import (
	"testing"
	"time"
)

func Test_save(t *testing.T) {
	data := &Statistics{Name: "1234567", Total: 1234, Date: uint32(time.Now().Unix())}
	t.Log(data.Save())
}
