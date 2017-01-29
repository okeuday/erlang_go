package erlang

//-*-Mode:Go;coding:utf-8;tab-width:4;c-basic-offset:4-*-
// ex: set ft=go fenc=utf-8 sts=4 ts=4 sw=4 noet nomod:
//
//
// BSD LICENSE
//
// Copyright (c) 2017, Michael Truog <mjtruog at gmail dot com>
// Copyright (c) 2009-2013, Dmitry Vasiliev <dima@hlabs.org>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//	 * Redistributions of source code must retain the above copyright
//	   notice, this list of conditions and the following disclaimer.
//	 * Redistributions in binary form must reproduce the above copyright
//	   notice, this list of conditions and the following disclaimer in
//	   the documentation and/or other materials provided with the
//	   distribution.
//	 * All advertising materials mentioning features or use of this
//	   software must display the following acknowledgment:
//		 This product includes software developed by Michael Truog
//	 * The name of the author may not be used to endorse or promote
//	   products derived from this software without specific prior
//	   written permission
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
// CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR
// CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH
// DAMAGE.
//

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"
	//"bytes"
)

func assertEqual(t *testing.T,
	expect interface{}, result interface{}, message string) {
	// Go doesn't believe in assertions (https://golang.org/doc/faq#assertions)
	// (Go isn't pursuing fail-fast development or fault-tolerance)
	if expect == result {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%#v != %#v", expect, result)
	}
	t.Fail()
	log.SetPrefix("\t")
	log.SetFlags(log.Lshortfile)
	log.Output(2, message)
}

func encode(t *testing.T, term interface{}, compressed int) string {
	b, err := TermToBinary(term, compressed)
	if err != nil {
		log.SetPrefix("\t")
		log.SetFlags(log.Lshortfile)
		log.Output(2, err.Error())
		t.FailNow()
		return ""
	}
	return string(b)
}

func TestAtom(t *testing.T) {
	atom1 := OtpErlangAtom("test")
	assertEqual(t, OtpErlangAtom("test"), atom1, "")
	assertEqual(t, strings.Repeat("X", 255),
		string(OtpErlangAtom(strings.Repeat("X", 255))), "")
	assertEqual(t, strings.Repeat("X", 256),
		string(OtpErlangAtom(strings.Repeat("X", 256))), "")
}

func TestList(t *testing.T) {
}

func TestImproperList(t *testing.T) {
}

func TestImproperListComparison(t *testing.T) {
}

func TestImproperListErrors(t *testing.T) {
}

func TestDecodeBinaryToTerm(t *testing.T) {
}

//...

func TestEncodeTermToBinaryTuple(t *testing.T) {
	assertEqual(t, "\x83h\x00",
		encode(t, []interface{}{}, -1), "")
	assertEqual(t, "\x83h\x00",
		encode(t, OtpErlangTuple{}, -1), "")
	assertEqual(t, "\x83h\x02h\x00h\x00",
		encode(t, []interface{}{[]interface{}{},
			[]interface{}{}}, -1), "")
	tuple1 := make(OtpErlangTuple, 255)
	for i := 0; i < 255; i++ {
		tuple1[i] = OtpErlangTuple{}
	}
	assertEqual(t, "\x83h\xff"+strings.Repeat("h\x00", 255),
		encode(t, tuple1, -1), "")
	tuple2 := make(OtpErlangTuple, 256)
	for i := 0; i < 256; i++ {
		tuple2[i] = OtpErlangTuple{}
	}
	assertEqual(t, "\x83i\x00\x00\x01\x00"+strings.Repeat("h\x00", 256),
		encode(t, tuple2, -1), "")
}
func TestEncodeTermToBinaryEmptyList(t *testing.T) {
	assertEqual(t, "\x83j",
		encode(t, OtpErlangList{}, -1), "")
	assertEqual(t, "\x83j",
		encode(t, OtpErlangList{Value: []interface{}{}}, -1), "")
	assertEqual(t, "\x83j",
		encode(t, OtpErlangList{Value: []interface{}{},
			Improper: false}, -1), "")
}
func TestEncodeTermToBinaryStringList(t *testing.T) {
	assertEqual(t, "\x83j",
		encode(t, "", -1), "")
	assertEqual(t, "\x83k\x00\x01\x00",
		encode(t, "\x00", -1), "")
	s := "\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r" +
		"\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a" +
		"\x1b\x1c\x1d\x1e\x1f !\"#$%&'()*+,-./0123456789:;<=>" +
		"?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopq" +
		"rstuvwxyz{|}~\x7f\x80\x81\x82\x83\x84\x85\x86\x87\x88" +
		"\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95" +
		"\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2" +
		"\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf" +
		"\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc" +
		"\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9" +
		"\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6" +
		"\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3" +
		"\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0" +
		"\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff"
	assertEqual(t, "\x83k\x01\x00"+s,
		encode(t, s, -1), "")
}
func TestEncodeTermToBinaryListBasic(t *testing.T) {
	assertEqual(t, "\x83\x6A",
		encode(t, OtpErlangList{}, -1), "")
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x6A\x6A",
		encode(t, OtpErlangList{Value: []interface{}{""}}, -1), "")
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x61\x01\x6A",
		encode(t, OtpErlangList{Value: []interface{}{1}}, -1), "")
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x61\xFF\x6A",
		encode(t, OtpErlangList{Value: []interface{}{255}}, -1), "")
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x62\x00\x00\x01\x00\x6A",
		encode(t, OtpErlangList{Value: []interface{}{256}}, -1), "")
	i1 := 2147483647
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x62\x7F\xFF\xFF\xFF\x6A",
		encode(t, OtpErlangList{Value: []interface{}{i1}}, -1), "")
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x6E\x04\x00\x00\x00\x00\x80\x6A",
		encode(t, OtpErlangList{Value: []interface{}{i1 + 1}}, -1), "")
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x61\x00\x6A",
		encode(t, OtpErlangList{Value: []interface{}{0}}, -1), "")
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x62\xFF\xFF\xFF\xFF\x6A",
		encode(t, OtpErlangList{Value: []interface{}{-1}}, -1), "")
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x62\xFF\xFF\xFF\x00\x6A",
		encode(t, OtpErlangList{Value: []interface{}{-256}}, -1), "")
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x62\xFF\xFF\xFE\xFF\x6A",
		encode(t, OtpErlangList{Value: []interface{}{-257}}, -1), "")
	i2 := -2147483648
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x62\x80\x00\x00\x00\x6A",
		encode(t, OtpErlangList{Value: []interface{}{i2}}, -1), "")
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x6E\x04\x01\x01\x00\x00\x80\x6A",
		encode(t, OtpErlangList{Value: []interface{}{i2 - 1}}, -1), "")
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x6B\x00\x04\x74\x65\x73\x74\x6A",
		encode(t, OtpErlangList{Value: []interface{}{"test"}}, -1), "")
	assertEqual(t,
		"\x83\x6C\x00\x00\x00\x02\x62\x00\x00\x01\x75\x62"+
			"\x00\x00\x01\xC7\x6A",
		encode(t, OtpErlangList{Value: []interface{}{373, 455}}, -1), "")
	list1 := OtpErlangList{}
	list2 := OtpErlangList{Value: []interface{}{list1}}
	assertEqual(t, "\x83\x6C\x00\x00\x00\x01\x6A\x6A",
		encode(t, list2, -1), "")
	list3 := OtpErlangList{Value: []interface{}{list1, list1}}
	assertEqual(t, "\x83\x6C\x00\x00\x00\x02\x6A\x6A\x6A",
		encode(t, list3, -1), "")
	list4 := OtpErlangList{Value: []interface{}{
		OtpErlangList{Value: []interface{}{"this", "is"}},
		OtpErlangList{Value: []interface{}{
			OtpErlangList{Value: []interface{}{"a"}}}},
		"test"}}
	assertEqual(t,
		"\x83\x6C\x00\x00\x00\x03\x6C\x00\x00\x00\x02\x6B"+
			"\x00\x04\x74\x68\x69\x73\x6B\x00\x02\x69\x73\x6A"+
			"\x6C\x00\x00\x00\x01\x6C\x00\x00\x00\x01\x6B\x00"+
			"\x01\x61\x6A\x6A\x6B\x00\x04\x74\x65\x73\x74\x6A",
		encode(t, list4, -1), "")
}
func TestEncodeTermToBinaryList(t *testing.T) {
	list1 := OtpErlangList{}
	list2 := OtpErlangList{Value: []interface{}{list1}}
	assertEqual(t, "\x83l\x00\x00\x00\x01jj",
		encode(t, list2, -1), "")
	list3 := OtpErlangList{Value: []interface{}{list1, list1, list1,
		list1, list1}}
	assertEqual(t, "\x83l\x00\x00\x00\x05jjjjjj",
		encode(t, list3, -1), "")
}
func TestEncodeTermToBinaryImproperList(t *testing.T) {
	list1 := OtpErlangList{Value: []interface{}{OtpErlangTuple{},
		[]interface{}{}},
		Improper: true}
	assertEqual(t, "\x83l\x00\x00\x00\x01h\x00h\x00",
		encode(t, list1, -1), "")
	list2 := OtpErlangList{Value: []interface{}{0, 1},
		Improper: true}
	assertEqual(t, "\x83l\x00\x00\x00\x01a\x00a\x01",
		encode(t, list2, -1), "")
}
func TestEncodeTermToBinaryUnicode(t *testing.T) {
	assertEqual(t, "\x83j",
		encode(t, "", -1), "")
	assertEqual(t, "\x83k\x00\x04test",
		encode(t, "test", -1), "")
	assertEqual(t, "\x83k\x00\x03\x00\xc3\xbf",
		encode(t, "\x00\xc3\xbf", -1), "")
	assertEqual(t, "\x83k\x00\x02\xc4\x80",
		encode(t, "\xc4\x80", -1), "")
	assertEqual(t, "\x83k\x00\x08\xd1\x82\xd0\xb5\xd1\x81\xd1\x82",
		encode(t, "\xd1\x82\xd0\xb5\xd1\x81\xd1\x82", -1), "")
	// becomes a list of small integers
	assertEqual(t,
		"\x83l\x00\x02\x00\x00"+
			strings.Repeat("a\xd0a\x90", 65536)+"j",
		encode(t, strings.Repeat("\xd0\x90", 65536), -1), "")
}
func TestEncodeTermToBinaryAtom(t *testing.T) {
	assertEqual(t, "\x83s\x00",
		encode(t, OtpErlangAtom(""), -1), "")
	assertEqual(t, "\x83s\x04test",
		encode(t, OtpErlangAtom("test"), -1), "")
}
func TestEncodeTermToBinaryStringBasic(t *testing.T) {
	assertEqual(t, "\x83\x6A",
		encode(t, "", -1), "")
	assertEqual(t, "\x83\x6B\x00\x04\x74\x65\x73\x74",
		encode(t, "test", -1), "")
	assertEqual(t, "\x83\x6B\x00\x09\x74\x77\x6F\x20\x77\x6F\x72\x64\x73",
		encode(t, "two words", -1), "")
	assertEqual(t, "\x83\x6B\x00\x16\x74\x65\x73\x74\x69\x6E\x67\x20\x6D"+
		"\x75\x6C\x74\x69\x70\x6C\x65\x20\x77\x6F\x72\x64\x73",
		encode(t, "testing multiple words", -1), "")
	assertEqual(t, "\x83\x6B\x00\x01\x20",
		encode(t, " ", -1), "")
	assertEqual(t, "\x83\x6B\x00\x02\x20\x20",
		encode(t, "  ", -1), "")
	assertEqual(t, "\x83\x6B\x00\x01\x31",
		encode(t, "1", -1), "")
	assertEqual(t, "\x83\x6B\x00\x02\x33\x37",
		encode(t, "37", -1), "")
	assertEqual(t, "\x83\x6B\x00\x07\x6F\x6E\x65\x20\x3D\x20\x31",
		encode(t, "one = 1", -1), "")
	assertEqual(t, "\x83\x6B\x00\x20\x21\x40\x23\x24\x25\x5E\x26\x2A\x28"+
		"\x29\x5F\x2B\x2D\x3D\x5B\x5D\x7B\x7D\x5C\x7C\x3B\x27\x3A"+
		"\x22\x2C\x2E\x2F\x3C\x3E\x3F\x7E\x60",
		encode(t, "!@#$%^&*()_+-=[]{}\\|;':\",./<>?~`", -1), "")
	assertEqual(t, "\x83\x6B\x00\x09\x22\x08\x0C\x0A\x0D\x09\x0B\x53\x12",
		encode(t, "\"\b\f\n\r\t\v\123\x12", -1), "")
}
func TestEncodeTermToBinaryString(t *testing.T) {
	assertEqual(t, "\x83j",
		encode(t, "", -1), "")
	assertEqual(t, "\x83k\x00\x04test",
		encode(t, "test", -1), "")
}
func TestEncodeTermToBinaryBoolean(t *testing.T) {
	assertEqual(t, "\x83s\x04true",
		encode(t, true, -1), "")
	assertEqual(t, "\x83s\x05false",
		encode(t, false, -1), "")
}
func TestEncodeTermToBinaryShortInteger(t *testing.T) {
	assertEqual(t, "\x83a\x00",
		encode(t, 0, -1), "")
	assertEqual(t, "\x83a\xff",
		encode(t, 255, -1), "")
}
func TestEncodeTermToBinaryInteger(t *testing.T) {
	assertEqual(t, "\x83b\xff\xff\xff\xff",
		encode(t, -1, -1), "")
	assertEqual(t, "\x83b\x80\x00\x00\x00",
		encode(t, -2147483648, -1), "")
	assertEqual(t, "\x83b\x00\x00\x01\x00",
		encode(t, 256, -1), "")
	assertEqual(t, "\x83b\x7f\xff\xff\xff",
		encode(t, 2147483647, -1), "")
}
func TestEncodeTermToBinaryLongInteger(t *testing.T) {
	assertEqual(t, "\x83n\x04\x00\x00\x00\x00\x80",
		encode(t, 2147483648, -1), "")
	assertEqual(t, "\x83n\x04\x01\x01\x00\x00\x80",
		encode(t, -2147483649, -1), "")
	i1 := big.NewInt(0)
	i1.Exp(big.NewInt(2), big.NewInt(2040), nil)
	assertEqual(t, "\x83o\x00\x00\x01\x00\x00"+
		strings.Repeat("\x00", 255)+"\x01",
		encode(t, i1, -1), "")
	i2 := big.NewInt(0).Neg(i1)
	assertEqual(t, "\x83o\x00\x00\x01\x00\x01"+
		strings.Repeat("\x00", 255)+"\x01",
		encode(t, i2, -1), "")
}
func TestEncodeTermToBinaryFloat(t *testing.T) {
	assertEqual(t, "\x83F\x00\x00\x00\x00\x00\x00\x00\x00",
		encode(t, 0.0, -1), "")
	assertEqual(t, "\x83F?\xe0\x00\x00\x00\x00\x00\x00",
		encode(t, 0.5, -1), "")
	assertEqual(t, "\x83F\xbf\xe0\x00\x00\x00\x00\x00\x00",
		encode(t, -0.5, -1), "")
	assertEqual(t, "\x83F@\t!\xfbM\x12\xd8J",
		encode(t, 3.1415926, -1), "")
	assertEqual(t, "\x83F\xc0\t!\xfbM\x12\xd8J",
		encode(t, -3.1415926, -1), "")
}
func TestEncodeTermToBinaryCompressedTerm(t *testing.T) {
	list1 := OtpErlangList{}
	list2 := OtpErlangList{Value: []interface{}{list1, list1, list1,
		list1, list1, list1,
		list1, list1, list1,
		list1, list1, list1,
		list1, list1, list1}}
	assertEqual(t,
		"\x83P\x00\x00\x00\x15x\x9c\xcaa``\xe0\xcfB\x03\x80"+
			"\x00\x00\x00\xff\xffB@\a\x1c",
		encode(t, list2, 6), "")
	assertEqual(t,
		"\x83P\x00\x00\x00\x15x\xda\xcaa``\xe0\xcfB\x03\x80"+
			"\x00\x00\x00\xff\xffB@\a\x1c",
		encode(t, list2, 9), "")
	assertEqual(t,
		"\x83P\x00\x00\x00\x15x\x01\x00\x15\x00\xea\xffl\x00"+
			"\x00\x00\x0fjjjjjjjjjjjjjjjj\x01\x00\x00\xff\xffB@\a\x1c",
		encode(t, list2, 0), "")
	assertEqual(t,
		"\x83P\x00\x00\x00\x17x\xda\xcaf\x10I\xc1\x02\x00\x01"+
			"\x00\x00\xff\xff]`\bP",
		encode(t, strings.Repeat("d", 20), 9), "")
}
