// Copyright 2009 Bill Zissimopoulos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"bytes"
	"reflect"
	"testing"
)

type testNetjsonCoder struct {
}

func (self *testNetjsonCoder) NetjsonEncode(i interface{}) ([]byte, error) {
	return nil, nil
}

func (self *testNetjsonCoder) NetjsonDecode(i interface{}, buf []byte) error {
	v := reflect.ValueOf(i).Elem()
	v.Set(reflect.MakeChan(v.Type(), 0))
	return nil
}

type testNetjson struct {
	I   int
	IP  *int
	IPP **int
	C   chan string
	CP  *chan string
	CPP **chan string
}

func TestNetjson(t *testing.T) {
	coder := &testNetjsonCoder{}

	c0 := make(chan int)
	s0 := testNetjson{}
	s0.I = 42
	s0.IP = &s0.I
	s0.IPP = &s0.IP
	s0.C = make(chan string)
	s0.CP = &s0.C
	s0.CPP = &s0.CP

	c1stg0 := make(chan int)
	c1stg1 := c1stg0
	s1stg0 := testNetjson{}
	s1stg1 := s1stg0
	c1 := &c1stg1
	s1 := &s1stg1

	b0, err := MarshalWithEncoder(c0, coder)
	if err != nil {
		t.Errorf("MarshalWithEncoder(%#v) error: %v", c0, err)
	}

	b1, err := MarshalWithEncoder(s0, coder)
	if err != nil {
		t.Errorf("MarshalWithEncoder(%#v) error: %v", s0, err)
	}

	err = UnmarshalWithDecoder(b0, c1, coder)
	if err != nil {
		t.Errorf("UnmarshalWithDecoder(%#v) error: %v", c1, err)
	}

	err = UnmarshalWithDecoder(b1, s1, coder)
	if err != nil {
		t.Errorf("Decode(%#v) error: %v", s1, err)
	}

	if c1stg0 == c1stg1 || s1stg0 == s1stg1 {
		t.Error("Decode check0: failed")
	}

	if &s1stg1.I == s1stg1.IP || &s1stg1.IP == s1stg1.IPP {
		t.Error("Decode check1: failed")
	}

	if &s1stg1.C == s1stg1.CP || &s1stg1.CP == s1stg1.CPP {
		t.Error("Decode check2: failed")
	}

	if 42 != s1stg1.I || 42 != *s1stg1.IP || 42 != **s1stg1.IPP {
		t.Error("Decode check3: failed")
	}

	if nil == s1stg1.C || nil == *s1stg1.CP || nil == **s1stg1.CPP {
		t.Error("Decode check3: failed")
	}
}

func TestNetjsonEncoderDecoder(t *testing.T) {
	coder := &testNetjsonCoder{}

	c0 := make(chan int)
	s0 := testNetjson{}
	s0.I = 42
	s0.IP = &s0.I
	s0.IPP = &s0.IP
	s0.C = make(chan string)
	s0.CP = &s0.C
	s0.CPP = &s0.CP

	c1stg0 := make(chan int)
	c1stg1 := c1stg0
	s1stg0 := testNetjson{}
	s1stg1 := s1stg0
	c1 := &c1stg1
	s1 := &s1stg1

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.SetNetjsonEncoder(coder)

	err := enc.Encode(c0)
	if err != nil {
		t.Errorf("Encode(%#v) error: %v", c0, err)
	}

	err = enc.Encode(s0)
	if err != nil {
		t.Errorf("Encode(%#v) error: %v", s0, err)
	}

	dec := NewDecoder(buf)
	dec.SetNetjsonDecoder(coder)

	err = dec.Decode(c1)
	if err != nil {
		t.Errorf("Decode(%#v) error: %v", c1, err)
	}

	err = dec.Decode(s1)
	if err != nil {
		t.Errorf("Decode(%#v) error: %v", s1, err)
	}

	if c1stg0 == c1stg1 || s1stg0 == s1stg1 {
		t.Error("Decode check0: failed")
	}

	if &s1stg1.I == s1stg1.IP || &s1stg1.IP == s1stg1.IPP {
		t.Error("Decode check1: failed")
	}

	if &s1stg1.C == s1stg1.CP || &s1stg1.CP == s1stg1.CPP {
		t.Error("Decode check2: failed")
	}

	if 42 != s1stg1.I || 42 != *s1stg1.IP || 42 != **s1stg1.IPP {
		t.Error("Decode check3: failed")
	}

	if nil == s1stg1.C || nil == *s1stg1.CP || nil == **s1stg1.CPP {
		t.Error("Decode check3: failed")
	}
}
