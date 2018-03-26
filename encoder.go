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

import "io/ioutil"
import "github.com/nu7hatch/gouuid"
import "github.com/klauspost/reedsolomon"

func Encode(s *Storage,u *uuid.UUID, data []byte,name,mime string) error {
	const maxSize = 64<<20
	
	var enc reedsolomon.Encoder
	var err error
	
	nds,nps := s.NumShards( (len(data)+maxSize)/maxSize )
	
	if nps>0 {
		enc,err = reedsolomon.New(nds,nps,reedsolomon.WithCauchyMatrix())
		if err!=nil { return err }
	}
	
	shards := Split(nds,nps,data,name,mime)
	Fill(shards)
	
	if nps>0 {
		for i := 0; i<nps;i++ {
			shards = append(shards,make([]byte,len(shards[0])))
		}
		err = enc.Encode(shards)
		if err!=nil { return err }
	}
	
	for i,shard := range shards {
		ioutil.WriteFile(s.Shard(u,i),shard,0600)
	}
	return nil
}



