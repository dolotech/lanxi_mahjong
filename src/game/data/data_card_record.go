package data

import "time"

type CarRecord struct {
	Id     string   `bson:"_id"`    //
	RoomId uint32   `bson:"RoomId"` //
	Record []uint64 `bson:"Record"` //
	CTime  uint32   `bson:"CTime"`  //
}

func (this *CarRecord) Add() error {
	this.CTime = uint32(time.Now().Unix())
	return C(_CARD_RECORD).Insert(this)
}

func (this *CarRecord) Get() error {
	return C(_CARD_RECORD).FindId(this.Id).One(this)
}
