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

import "fmt"
import "github.com/nu7hatch/gouuid"

type Storage struct {
	Pathes   []string
	Data,Par int
}
func (s *Storage) Len() int {
	return s.Data+s.Par
}
func (s *Storage) SetRedundancy(i int) {
	l := len(s.Pathes)-1
	if i>l { i=l }
	l++
	s.Data = len(s.Pathes)-i
	s.Par  = i
}
func (s *Storage) Shard(u *uuid.UUID,i int) string {
	return fmt.Sprintf("%s/%v.s%d",s.Pathes[i%len(s.Pathes)],u,i)
}
func (s *Storage) NumShards(min int) (data,par int) {
	l := s.Data+s.Par
	k := l-1
	r := (min+k)/l
	data = r*s.Data
	par  = r*s.Par
	return
}


