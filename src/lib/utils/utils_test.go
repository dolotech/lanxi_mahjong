/**
 * Created by Michael on 2015/8/4.
 */
package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_copy(t *testing.T) {
	a := AA{A: 999}
	b := AA{}
	//err :=Clone(b, a)
	t.Log(a, b)
}
func Test_AES(t *testing.T) {
	aesEnc := AesEncrypt{}
	aesEnc.SetKey([]byte("aalk;lkasjd;lkfj;alk"))
	doc := []byte("abcde号。")
	arrEncrypt, err := aesEnc.Encrypt(doc)
	fmt.Println(string(arrEncrypt))
	if err != nil {
		fmt.Println(string(arrEncrypt))
		return
	}
	strMsg, err := aesEnc.Decrypt(arrEncrypt)
	if err != nil {
		fmt.Println(string(arrEncrypt))
		return
	}
	fmt.Println(string(strMsg))
}
func Test_XXTEA(t *testing.T) {
	str := "Hello World! 你好，中国！"
	key := "1234567890"
	encrypt_data := Encrypt([]byte(str), []byte(key))
	//fmt.Println(base64.StdEncoding.EncodeToString(encrypt_data))
	decrypt_data := Decrypt(encrypt_data, []byte(key))
	t.Log(len(encrypt_data), len(decrypt_data))
	t.Log(string(encrypt_data))
}
func TestPWD(t *testing.T) {
	t.Log(AalidataPwd("dolo0425"))

}

func TestPhone(t *testing.T) {
	t.Log(PhoneRegexp("8601593533372"))
}

type AA struct {
	CC
	A int `json:"a"`
}

type BB interface {
	Decode(b *[]byte) error
	Encode() (*[]byte, error)
}

type CC struct{}

func (this *CC) Decode(b *[]byte) error {
	return json.Unmarshal(*b, this)
}

func (this *CC) Encode() (*[]byte, error) {
	data, err := json.Marshal(this)
	return &data, err
}

type At struct {
	Name string
	Nickname string
	Age  int
}

func TestRand(t *testing.T) {
	t.Log(Base62encode(uint64(TimestampNano())))
	t.Log(Base62encode(uint64(RandUint32())))
	a := &At{
		Name: "wang",
		Age: 29,
	}
	t.Log(Struct2Map(a))
}
