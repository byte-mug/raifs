/*
Copyright (c) 2018 Simon Schmidt

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/


package raifs

import "github.com/nu7hatch/gouuid"
import "encoding/asn1"

type objectStore interface{
	Store(id *uuid.UUID, value []byte)
	Load(id *uuid.UUID) []byte
}

type Header struct{
	Name string
	Mime string
	Offset int64
	Length int64 // Length of this shard.
	Shard   int
	Nshards int
	Pshards int
}

func Split(n,p int, data []byte, name,mime string) (shards [][]byte) {
	var part []byte
	offset := int64(0)
	t := len(data)/n
	shards = make([][]byte,n,n+p)
	lst := n-1
	for i := range shards {
		if lst==i {
			part,data = data,nil
		} else if len(data)>t {
			part,data = data[:t],data[t:]
		} else {
			part,data = data,nil
		}
		lp := int64(len(part))
		header,_ := asn1.Marshal(Header{name,mime,offset,lp,i,n,p})
		offset += lp
		shards[i] = append(header,part...)
	}
	return
}

func Fill(shards [][]byte) {
	min,max := -1,0
	for _,shard := range shards {
		l := len(shard)
		if min<0 || l<min { min = l }
		if max<l { max = l }
	}
	diff := make([]byte,max-min)
	for i := range shards {
		l := len(shards[i])
		shards[i] = append(shards[i],diff[:max-l]...)
	}
}

func Peek(shard []byte) (hdr Header,err error) {
	_,err = asn1.Unmarshal(shard,&hdr)
	return
}
func Pull(shard []byte) (rest []byte,err error) {
	var hdr Header
	rest,err = asn1.Unmarshal(shard,&hdr)
	return
}

/* =========== */

