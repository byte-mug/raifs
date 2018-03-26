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

import "io"
import "io/ioutil"
import "github.com/nu7hatch/gouuid"
import "encoding/asn1"
import "github.com/klauspost/reedsolomon"
import "fmt"

var ECOR = fmt.Errorf("raifs.ECOR")

type Decoder struct{
	St *Storage
	U  *uuid.UUID
	hdr Header
	ele []int64
	noff int64
	data []byte
	allData [][]byte
}
func (d *Decoder) Look() error {
	l := d.St.Len()
	max := 16
	
	if max<l { max = l*4 }
	
	//found := false
	//shards := make([][]byte,0,max)
	var hdr Header
	
	for i := 0 ; i<max ; i++ {
		fn := d.St.Shard(d.U,i)
		data,_ := ioutil.ReadFile(fn)
		//shards = append(shards,data)
		
		if len(data)==0 { continue } // Not found.
		//if found { continue }
		
		_,err := asn1.Unmarshal(data,&hdr)
		if err!=nil { return err } // Decode error. Too late.
		
		d.hdr = hdr
		d.hdr.Offset = 0
		d.hdr.Shard  = 0
		d.noff = 0
		return nil
	}
	return io.EOF
}
func (d *Decoder) fillN() error {
	shards := make([][]byte,d.hdr.Nshards+d.hdr.Pshards)
	for i := range shards {
		data,_ := ioutil.ReadFile(d.St.Shard(d.U,i))
		shards[i] = data
	}
	enc,err := reedsolomon.New(d.hdr.Nshards,d.hdr.Pshards,reedsolomon.WithCauchyMatrix())
	if err!=nil { return err }
	
	err = enc.ReconstructData(shards)
	if err!=nil { return err }
	
	d.allData = shards[:d.hdr.Nshards]
	
	return nil
}
func (d *Decoder) fill1() error {
	var hdr Header
	var rest []byte
	i := d.hdr.Shard
	if i>=d.hdr.Nshards { return io.EOF }
	
	var data []byte
	var err error
	
	if len(d.allData)>0 {
		data = d.allData[i]
	} else {
		data,err = ioutil.ReadFile(d.St.Shard(d.U,i))
		if err!=nil {
			err = d.fillN() // try fill-in
			if err==nil { data = d.allData[i] }
		}
	}
	
	if err==nil { rest,err = asn1.Unmarshal(data,&hdr) }
	
	if err==nil{
		if i!=hdr.Shard { return ECOR }
		if (d.noff)!=hdr.Offset { return ECOR } // error
		if int64(len(rest))<hdr.Length { return ECOR }
		d.data = rest[:hdr.Length]
		d.hdr = hdr
		d.hdr.Shard++
		d.noff = d.hdr.Length+d.hdr.Offset
		d.ele = append(d.ele,d.noff)
	}
	
	return err
}
func (d *Decoder) Read(buf []byte) (n int, err error) {
	m := len(buf)
restart:
	l := len(d.data)
	if m==0 { return }
	if l==0 {
		err = d.fill1()
		if err!=nil { return }
		l = len(d.data)
	}
	if l>0 {
		if l>m { l=m }
		copy(buf,d.data[:l])
		n+=l
		d.data = d.data[l:]
		if m==l { return }
		m-=l
		buf = buf[l:]
	}
	goto restart
}

//Seek(offset int64, whence int) (int64, error)


