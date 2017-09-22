// Copyright 2009 Bill Zissimopoulos. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json

import (
	"bytes"
	"encoding/base64"
	"reflect"
)

// NetjsonEncoder is the interface used to encode channels.
// An object that implements NetjsonEncoder will be called
// when a channel needs to be encoded.
type NetjsonEncoder interface {
	// NetjsonEncode returns a byte slice representing the encoding of the
	// passed channel.
	NetjsonEncode(reflect.Value) ([]byte, error)
}

// NetjsonDecoder is the interface used to decode channels.
// An object that implements NetjsonDecoder will be called
// when a channel needs to be decoded.
type NetjsonDecoder interface {
	// NetjsonDecode overwrites the passed channel, which must be a pointer,
	// with the value represented by the byte slice.
	NetjsonDecode(reflect.Value, []byte) error
}

func netjsonEncoder(e *encodeState, v reflect.Value, opts encOpts) {
	if nil == e.netjsonEnc {
		e.error(&UnsupportedTypeError{v.Type()})
		return
	}
	buf, err := e.netjsonEnc.NetjsonEncode(v)
	if err == nil {
		n := base64.RawURLEncoding.EncodedLen(len(buf)) +
			len(netjsonChanQPrefix) + len(netjsonChanQSuffix)
		s := make([]byte, n)
		copy(s, netjsonChanQPrefix)
		copy(s[n-len(netjsonChanQSuffix):], netjsonChanQSuffix)
		base64.RawURLEncoding.Encode(s[len(netjsonChanQPrefix):n-len(netjsonChanQSuffix)], buf)
		// copy JSON into buffer, checking validity.
		err = compact(&e.Buffer, s, opts.escapeHTML)
	}
	if err != nil {
		e.error(&MarshalerError{v.Type(), err})
	}
}

func netjsonDecoder(d *decodeState, s []byte, v reflect.Value) {
	if !netjsonIsEncoded(d, s) {
		d.saveError(&UnmarshalTypeError{Value: "string", Type: v.Type(), Offset: int64(d.off)})
		return
	}
	s = s[len(netjsonChanPrefix) : len(s)-len(netjsonChanSuffix)]
	buf := make([]byte, base64.RawURLEncoding.DecodedLen(len(s)))
	n, err := base64.RawURLEncoding.Decode(buf, s)
	if err == nil {
		err = d.netjsonDec.NetjsonDecode(v.Addr(), buf[:n])
	}
	if err != nil {
		d.saveError(err)
	}
}

func netjsonIsEncoded(d *decodeState, s []byte) bool {
	return nil != d.netjsonDec &&
		len(s) >= len(netjsonChanPrefix) && //len(s) >= len(netjsonChanSuffix) &&
		bytes.Equal(s[0:len(netjsonChanPrefix)], netjsonChanPrefix) &&
		bytes.Equal(s[len(s)-len(netjsonChanSuffix):], netjsonChanSuffix)
}

var (
	netjsonChanQPrefix = []byte(`"\//chan(`)
	netjsonChanQSuffix = []byte(`)\//"`)
	netjsonChanPrefix  = []byte(`//chan(`)
	netjsonChanSuffix  = []byte(`)//`)
)
