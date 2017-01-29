package erlang
//-*-Mode:Go;coding:utf-8;tab-width:4;c-basic-offset:4;indent-tabs-mode:()-*-
// ex: set ft=go fenc=utf-8 sts=4 ts=4 sw=4 et nomod:
//
// BSD LICENSE
// 
// Copyright (c) 2017, Michael Truog <mjtruog at gmail dot com>
// All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// 
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above copyright
//       notice, this list of conditions and the following disclaimer in
//       the documentation and/or other materials provided with the
//       distribution.
//     * All advertising materials mentioning features or use of this
//       software must display the following acknowledgment:
//         This product includes software developed by Michael Truog
//     * The name of the author may not be used to endorse or promote
//       products derived from this software without specific prior
//       written permission
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
    "bytes"
    "encoding/binary"
    "compress/zlib"
    "math"
    "math/big"
)

const (
    // tag values here http://www.erlang.org/doc/apps/erts/erl_ext_dist.html
    tagVersion = 131
    tagCompressedZlib = 80
    tagNewFloatExt = 70
    tagBitBinaryExt = 77
    tagAtomCacheRef = 78
    tagSmallIntegerExt = 97
    tagIntegerExt = 98
    tagFloatExt = 99
    tagAtomExt = 100
    tagReferenceExt = 101
    tagPortExt = 102
    tagPidExt = 103
    tagSmallTupleExt = 104
    tagLargeTupleExt = 105
    tagNilExt = 106
    tagStringExt = 107
    tagListExt = 108
    tagBinaryExt = 109
    tagSmallBigExt = 110
    tagLargeBigExt = 111
    tagNewFunExt = 112
    tagExportExt = 113
    tagNewReferenceExt = 114
    tagSmallAtomExt = 115
    tagMapExt = 116
    tagFunExt = 117
    tagAtomUtf8Ext = 118
    tagSmallAtomUtf8Ext = 119
)

// OtpErlangAtomCacheRef represents ATOM_CACHE_REF
type OtpErlangAtomCacheRef uint8

// OtpErlangAtom represents SMALL_ATOM_EXT or ATOM_EXT
type OtpErlangAtom string

// OtpErlangAtomUTF8 represents SMALL_ATOM_UTF8_EXT or ATOM_UTF8_EXT
type OtpErlangAtomUTF8 string

// OtpErlangList represents NIL_EXT or LIST_EXT
type OtpErlangList struct {
    Value []interface{}
    Improper bool
}

// OtpErlangTuple represents SMALL_TUPLE_EXT or LARGE_TUPLE_EXT
type OtpErlangTuple []interface{}

// OtpErlangBinary represents BIT_BINARY_EXT or BINARY_EXT
type OtpErlangBinary struct {
    Value []byte
    Bits uint8
}

// OtpErlangFunction represents FUN_EXT or NEW_FUN_EXT
type OtpErlangFunction struct {
    Tag uint8
    Value []byte
}

// OtpErlangReference represents REFERENCE_EXT or NEW_REFERENCE_EXT
type OtpErlangReference struct {
    Node OtpErlangAtom
    ID []byte
    Creation byte
}

// OtpErlangPort represents PORT_EXT
type OtpErlangPort struct {
    Node OtpErlangAtom
    ID []byte
    Creation byte
}

// OtpErlangPid represents PID_EXT
type OtpErlangPid struct {
    Node OtpErlangAtom
    ID []byte
    Serial []byte
    Creation byte
}

// ParsingError provides specific parsing failure information
type ParseError struct {
    message string
}
func parseErrorNew(message string) error {
    return &ParseError{message}
}
func (e *ParseError) Error() string {
    return e.message
}

// InputError describes problems with function input parameters
type InputError struct {
    message string
}
func inputErrorNew(message string) error {
    return &InputError{message}
}
func (e *InputError) Error() string {
    return e.message
}

// OutputError describes problems with creating function output data
type OutputError struct {
    message string
}
func outputErrorNew(message string) error {
    return &OutputError{message}
}
func (e *OutputError) Error() string {
    return e.message
}

// BinaryToTerm decodes the Erlang Binary Term Format into Go types
func BinaryToTerm(data []byte) (interface{}, error) {
    return nil, parseErrorNew("invalid")
}

// TermToBinary encodes Go types into the Erlang Binary Term Format
func TermToBinary(term interface{}, compressed int) ([]byte, error) {
    if compressed < -1 || compressed > 9 {
        return nil, inputErrorNew("compressed in [-1..9]")
    }
    dataUncompressed, err := termsToBinary(term, new(bytes.Buffer))
    if err != nil {
        return nil, err
    }
    if compressed == -1 {
        return append([]byte{tagVersion}, dataUncompressed.Bytes()...), nil
    }
    var dataCompressed *bytes.Buffer = new(bytes.Buffer)
    compress, err := zlib.NewWriterLevel(dataCompressed, compressed)
    if err != nil {
        return nil, err
    }
    _, err = compress.Write(dataUncompressed.Bytes())
    if err != nil {
        return nil, err
    }
    err = compress.Close()
    if err != nil {
        return nil, err
    }
    var result *bytes.Buffer = new(bytes.Buffer)
    _, err = result.Write([]byte{tagVersion, tagCompressedZlib})
    if err != nil {
        return nil, err
    }
    err = binary.Write(result, binary.BigEndian,
                       uint32(dataUncompressed.Len()))
    if err != nil {
        return nil, err
    }
    _, err = result.Write(dataCompressed.Bytes())
    if err != nil {
        return nil, err
    }
    return result.Bytes(), nil
}

func termsToBinary(term interface{},
                   buffer *bytes.Buffer) (*bytes.Buffer, error) {
    switch term.(type) {
        case uint8:
            _, err := buffer.Write([]byte{tagSmallIntegerExt, term.(uint8)})
            return buffer, err
        case uint16:
            return integerToBinary(int32(term.(uint16)), buffer)
        case uint32:
            return bignumToBinary(big.NewInt(int64(term.(uint32))), buffer)
        case uint64:
            var value *big.Int = new(big.Int)
            value.SetUint64(term.(uint64))
            return bignumToBinary(value, buffer)
        case int8:
            return integerToBinary(int32(term.(int8)), buffer)
        case int16:
            return integerToBinary(int32(term.(int16)), buffer)
        case int32:
            return integerToBinary(term.(int32), buffer)
        case int64:
            return bignumToBinary(big.NewInt(term.(int64)), buffer)
        case int:
            switch i := term.(int); {
                case i >= 0 && i <= math.MaxUint8:
                    return termsToBinary(uint8(i), buffer)
                case i >= math.MinInt32 && i <= math.MaxInt32:
                    return integerToBinary(int32(i), buffer)
                case i >= math.MinInt64 && i <= math.MaxInt64:
                    return termsToBinary(int64(i), buffer)
                default:
                    return buffer, outputErrorNew("invalid int")
            }
        case *big.Int:
            return bignumToBinary(term.(*big.Int), buffer)
        case float32:
            return floatToBinary(float64(term.(float32)), buffer)
        case float64:
            return floatToBinary(term.(float64), buffer)
        case bool:
            if (term.(bool)) {
                return atomToBinary("true", buffer)
            }
            return atomToBinary("false", buffer)
        case OtpErlangAtom:
            return atomToBinary(string(term.(OtpErlangAtom)), buffer)
        case OtpErlangAtomUTF8:
            return atomUtf8ToBinary(string(term.(OtpErlangAtomUTF8)), buffer)
        case OtpErlangAtomCacheRef:
            _, err := buffer.Write([]byte{tagAtomCacheRef,
                                          uint8(term.(OtpErlangAtomCacheRef))})
            return buffer, err
        case string:
            return stringToBinary(term.(string), buffer)
        case OtpErlangList:
            return listToBinary(term.(OtpErlangList), buffer)
        case []byte:
            return binaryToBinary(OtpErlangBinary{Value: term.([]byte),
                                                  Bits: 8}, buffer)
        case OtpErlangBinary:
            return binaryToBinary(term.(OtpErlangBinary), buffer)
        case OtpErlangTuple:
            return tupleToBinary(term.(OtpErlangTuple), buffer)
        case []interface{}:
            return tupleToBinary(term.([]interface{}), buffer)
        case OtpErlangFunction:
            return functionToBinary(term.(OtpErlangFunction), buffer)
        case OtpErlangReference:
            return referenceToBinary(term.(OtpErlangReference), buffer)
        case OtpErlangPort:
            return portToBinary(term.(OtpErlangPort), buffer)
        case OtpErlangPid:
            return pidToBinary(term.(OtpErlangPid), buffer)
        default:
            return buffer, outputErrorNew("unknown go type")
    }
}

func integerToBinary(term int32,
                     buffer *bytes.Buffer) (*bytes.Buffer, error) {
    err := buffer.WriteByte(tagIntegerExt)
    if err != nil {
        return buffer, err
    }
    err = binary.Write(buffer, binary.BigEndian, term)
    return buffer, err
}

func bignumToBinary(term *big.Int,
                    buffer *bytes.Buffer) (*bytes.Buffer, error) {
    var sign uint8
    if term.Sign() < 0 {
        sign = 1
    } else {
        sign = 0
    }
    value := term.Bytes()
    var length int
    var err error
    switch length = len(value); {
        case length <= math.MaxUint8:
            _, err = buffer.Write([]byte{tagSmallBigExt, uint8(length)})
            if err != nil {
                return buffer, err
            }
        case length <= math.MaxUint32:
            err = buffer.WriteByte(tagLargeBigExt)
            if err != nil {
                return buffer, err
            }
            err = binary.Write(buffer, binary.BigEndian, uint32(length))
            if err != nil {
                return buffer, err
            }
        default:
            return buffer, outputErrorNew("uint32 overflow")
    }
    err = buffer.WriteByte(sign)
    if err != nil {
        return buffer, err
    }
    // little-endian is required
    half := length >> 1
    iLast := length - 1
    for i := 0; i < half; i++ {
        j := iLast - i
        value[i], value[j] = value[j], value[i]
    }
    _, err = buffer.Write(value)
    return buffer, err
}

func floatToBinary(term float64,
                   buffer *bytes.Buffer) (*bytes.Buffer, error) {
    err := buffer.WriteByte(tagNewFloatExt)
    if err != nil {
        return buffer, err
    }
    err = binary.Write(buffer, binary.BigEndian, term)
    return buffer, err
}

func atomToBinary(term string,
                  buffer *bytes.Buffer) (*bytes.Buffer, error) {
    switch length := len(term); {
        case length <= math.MaxUint8:
            _, err := buffer.Write([]byte{tagSmallAtomExt, uint8(length)})
            if err != nil {
                return buffer, err
            }
            _, err = buffer.WriteString(term)
            return buffer, err
        case length <= math.MaxUint16:
            err := buffer.WriteByte(tagAtomExt)
            if err != nil {
                return buffer, err
            }
            err = binary.Write(buffer, binary.BigEndian, uint16(length))
            if err != nil {
                return buffer, err
            }
            _, err = buffer.WriteString(term)
            return buffer, err
        default:
            return buffer, outputErrorNew("uint16 overflow")
    }
}

func atomUtf8ToBinary(term string,
                      buffer *bytes.Buffer) (*bytes.Buffer, error) {
    switch length := len(term); {
        case length <= math.MaxUint8:
            _, err := buffer.Write([]byte{tagSmallAtomUtf8Ext, uint8(length)})
            if err != nil {
                return buffer, err
            }
            _, err = buffer.WriteString(term)
            return buffer, err
        case length <= math.MaxUint16:
            err := buffer.WriteByte(tagAtomUtf8Ext)
            if err != nil {
                return buffer, err
            }
            err = binary.Write(buffer, binary.BigEndian, uint16(length))
            if err != nil {
                return buffer, err
            }
            _, err = buffer.WriteString(term)
            return buffer, err
        default:
            return buffer, outputErrorNew("uint16 overflow")
    }
}

func stringToBinary(term string,
                    buffer *bytes.Buffer) (*bytes.Buffer, error) {
    switch length := len(term); {
        case length == 0:
            err := buffer.WriteByte(tagNilExt)
            return buffer, err
        case length <= math.MaxUint16:
            err := buffer.WriteByte(tagStringExt)
            if err != nil {
                return buffer, err
            }
            err = binary.Write(buffer, binary.BigEndian, uint16(length))
            if err != nil {
                return buffer, err
            }
            _, err = buffer.WriteString(term)
            return buffer, err
        case length <= math.MaxUint32:
            err := buffer.WriteByte(tagListExt)
            if err != nil {
                return buffer, err
            }
            err = binary.Write(buffer, binary.BigEndian, uint32(length))
            if err != nil {
                return buffer, err
            }
            for i := 0; i < length; i++ {
                _, err = buffer.Write([]byte{tagSmallIntegerExt, term[i]})
                if err != nil {
                    return buffer, err
                }
            }
            err = buffer.WriteByte(tagNilExt)
            return buffer, err
        default:
            return buffer, outputErrorNew("uint32 overflow")
    }
}

func listToBinary(term OtpErlangList,
                  buffer *bytes.Buffer) (*bytes.Buffer, error) {
    var length int
    var err error
    switch length = len(term.Value); {
        case length == 0:
            err = buffer.WriteByte(tagNilExt)
            return buffer, err
        case length <= math.MaxUint32:
            err = buffer.WriteByte(tagListExt)
            if err != nil {
                return buffer, err
            }
            if term.Improper {
                err = binary.Write(buffer, binary.BigEndian, uint32(length - 1))
                if err != nil {
                    return buffer, err
                }
            } else {
                err = binary.Write(buffer, binary.BigEndian, uint32(length))
                if err != nil {
                    return buffer, err
                }
            }
        default:
            return buffer, outputErrorNew("uint32 overflow")
    }
    for i := 0; i < length; i++ {
        buffer, err = termsToBinary(term.Value[i], buffer)
        if err != nil {
            return buffer, err
        }
    }
    if !term.Improper {
        err = buffer.WriteByte(tagNilExt)
    }
    return buffer, err
}

func binaryToBinary(term OtpErlangBinary,
                    buffer *bytes.Buffer) (*bytes.Buffer, error) {
    var err error
    switch length := len(term.Value); {
        case term.Bits < 1 || term.Bits > 8:
            return buffer, outputErrorNew("invalid OtpErlangBinary.Bits")
        case length <= math.MaxUint32:
            if term.Bits != 8 {
                err = buffer.WriteByte(tagBitBinaryExt)
                if err != nil {
                    return buffer, err
                }
                err = binary.Write(buffer, binary.BigEndian, uint32(length))
                if err != nil {
                    return buffer, err
                }
                err = buffer.WriteByte(term.Bits)
                if err != nil {
                    return buffer, err
                }
            } else {
                err = buffer.WriteByte(tagBinaryExt)
                if err != nil {
                    return buffer, err
                }
                err = binary.Write(buffer, binary.BigEndian, uint32(length))
                if err != nil {
                    return buffer, err
                }
            }
        default:
            return buffer, outputErrorNew("uint32 overflow")
    }
    _, err = buffer.Write(term.Value)
    return buffer, err
}

func tupleToBinary(term []interface{},
                   buffer *bytes.Buffer) (*bytes.Buffer, error) {
    var length int
    var err error
    switch length = len(term); {
        case length <= math.MaxUint8:
            _, err = buffer.Write([]byte{tagSmallTupleExt, byte(length)})
            if err != nil {
                return buffer, err
            }
        case length <= math.MaxUint32:
            err = buffer.WriteByte(tagLargeTupleExt)
            if err != nil {
                return buffer, err
            }
            err = binary.Write(buffer, binary.BigEndian, uint32(length))
            if err != nil {
                return buffer, err
            }
        default:
            return buffer, outputErrorNew("uint32 overflow")
    }
    for i := 0; i < length; i++ {
        buffer, err = termsToBinary(term[i], buffer)
        if err != nil {
            return buffer, err
        }
    }
    return buffer, nil
}

func functionToBinary(term OtpErlangFunction,
                      buffer *bytes.Buffer) (*bytes.Buffer, error) {
    err := buffer.WriteByte(term.Tag)
    if err != nil {
        return buffer, err
    }
    _, err = buffer.Write(term.Value)
    return buffer, err
}

func referenceToBinary(term OtpErlangReference,
                       buffer *bytes.Buffer) (*bytes.Buffer, error) {
    switch length := len(term.ID) / 4; {
        case length == 0:
            err := buffer.WriteByte(tagReferenceExt)
            if err != nil {
                return buffer, err
            }
            buffer, err = termsToBinary(term.Node, buffer)
            if err != nil {
                return buffer, err
            }
            _, err = buffer.Write(term.ID)
            if err != nil {
                return buffer, err
            }
            err = buffer.WriteByte(term.Creation)
            return buffer, err
        case length <= math.MaxUint16:
            err := buffer.WriteByte(tagNewReferenceExt)
            if err != nil {
                return buffer, err
            }
            err = binary.Write(buffer, binary.BigEndian, uint16(length))
            if err != nil {
                return buffer, err
            }
            buffer, err = termsToBinary(term.Node, buffer)
            if err != nil {
                return buffer, err
            }
            err = buffer.WriteByte(term.Creation)
            if err != nil {
                return buffer, err
            }
            _, err = buffer.Write(term.ID)
            return buffer, err
        default:
            return buffer, outputErrorNew("uint16 overflow")
    }
}

func portToBinary(term OtpErlangPort,
                  buffer *bytes.Buffer) (*bytes.Buffer, error) {
    err := buffer.WriteByte(tagPortExt)
    if err != nil {
        return buffer, err
    }
    buffer, err = termsToBinary(term.Node, buffer)
    if err != nil {
        return buffer, err
    }
    _, err = buffer.Write(term.ID)
    if err != nil {
        return buffer, err
    }
    err = buffer.WriteByte(term.Creation)
    return buffer, err
}

func pidToBinary(term OtpErlangPid,
                 buffer *bytes.Buffer) (*bytes.Buffer, error) {
    err := buffer.WriteByte(tagPidExt)
    if err != nil {
        return buffer, err
    }
    buffer, err = termsToBinary(term.Node, buffer)
    if err != nil {
        return buffer, err
    }
    _, err = buffer.Write(term.ID)
    if err != nil {
        return buffer, err
    }
    _, err = buffer.Write(term.Serial)
    if err != nil {
        return buffer, err
    }
    err = buffer.WriteByte(term.Creation)
    return buffer, err
}

