package jsonDiff

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"sort"
	"strconv"
)

type Difference int

const (
	FullMatch Difference = iota
	SupersetMatch
	NoMatch
	FirstArgIsInvalidJson
	SecondArgIsInvalidJson
	BothArgsAreInvalidJson
)

func (d Difference) String() string {
	switch d {
	case FullMatch:
		return "FullMatch"
	case SupersetMatch:
		return "SupersetMatch"
	case NoMatch:
		return "NoMatch"
	case FirstArgIsInvalidJson:
		return "FirstArgIsInvalidJson"
	case SecondArgIsInvalidJson:
		return "SecondArgIsInvalidJson"
	case BothArgsAreInvalidJson:
		return "BothArgsAreInvalidJson"
	}
	return "Invalid"
}

type context struct {
	buf     bytes.Buffer
	level   int
	diff    Difference
	Indent  string
}

func (ctx *context) newline(s string) {
	ctx.buf.WriteString(s)
}

func (ctx *context) key(k string) {
	ctx.buf.WriteString(strconv.Quote(k))
	ctx.buf.WriteString(": ")
}

func (ctx *context) writeValue(v interface{}, full bool, opCode string) {
	switch vv := v.(type) {
	case bool:
		ctx.buf.WriteString(strconv.Quote(string(opCode+strconv.FormatBool(vv))))
	case json.Number:
		ctx.buf.WriteString(strconv.Quote(string(opCode+string(vv))))
	case string:
		ctx.buf.WriteString(strconv.Quote(string(opCode+vv)))
	case []interface{}:
		if full {
			if len(vv) == 0 {
				ctx.buf.WriteString("[")
			} else {
				ctx.level++
				ctx.newline("[")
			}
			for i, v := range vv {
				ctx.writeValue(v, true, opCode)
				if i != len(vv)-1 {
					ctx.newline(",")
				} else {
					ctx.level--
					ctx.newline("")
				}
			}
			ctx.buf.WriteString("]")
		} else {
			ctx.buf.WriteString("[]")
		}
	case map[string]interface{}:
		if full {
			if len(vv) == 0 {
				ctx.buf.WriteString("{")
			} else {
				ctx.level++
				ctx.newline("{")
			}
			i := 0
			for k, v := range vv {
				ctx.key(k)
				ctx.writeValue(v, true, opCode)
				if i != len(vv)-1 {
					ctx.newline(",")
				} else {
					ctx.level--
					ctx.newline("")
				}
				i++
			}
			ctx.buf.WriteString("}")
		} else {
			ctx.buf.WriteString("{}")
		}
	default:
		ctx.buf.WriteString("null")
	}

}

func (ctx *context) writeMismatch(a, b interface{}) {
	ctx.buf.WriteString(strconv.Quote(a.(string) + " => "+ b.(string)))
}

func (ctx *context) result(d Difference) {
	if d == NoMatch {
		ctx.diff = NoMatch
	} else if d == SupersetMatch && ctx.diff != NoMatch {
		ctx.diff = SupersetMatch
	} else if ctx.diff != NoMatch && ctx.diff != SupersetMatch {
		ctx.diff = FullMatch
	}
}

func (ctx *context) printMismatch(a, b interface{}) {
	ctx.writeMismatch(a, b)
}

func (ctx *context) calculateDiff(a, b interface{}) {
	if a == nil || b == nil {
		if a == nil && b == nil {
			ctx.writeValue(a, false, "")
			ctx.result(FullMatch)
		} else {
			ctx.printMismatch(a, b)
			ctx.result(NoMatch)
		}
		return
	}

	ka := reflect.TypeOf(a).Kind()
	kb := reflect.TypeOf(b).Kind()
	if ka != kb {
		ctx.printMismatch(a, b)
		ctx.result(NoMatch)
		return
	}
	switch ka {
	case reflect.Bool:
		if a.(bool) != b.(bool) {
			ctx.printMismatch(a, b)
			ctx.result(NoMatch)
			return
		}
	case reflect.String:
		switch aa := a.(type) {
		case json.Number:
			bb, ok := b.(json.Number)
			if !ok || aa != bb {
				ctx.printMismatch(a, b)
				ctx.result(NoMatch)
				return
			}
		case string:
			bb, ok := b.(string)
			if !ok || aa != bb {
				ctx.printMismatch(a, b)
				ctx.result(NoMatch)
				return
			}
		}
	case reflect.Slice:
		sa, sb := a.([]interface{}), b.([]interface{})
		salen, sblen := len(sa), len(sb)
		max := salen
		if sblen > max {
			max = sblen
		}
		if max == 0 {
			ctx.buf.WriteString("[")
		} else {
			ctx.level++
			ctx.newline("[")
		}
		for i := 0; i < max; i++ {
			if i < salen && i < sblen {
				ctx.calculateDiff(sa[i], sb[i])
			} else if i < salen {
				ctx.writeValue(sa[i], true, "-")
				ctx.result(SupersetMatch)
			} else if i < sblen {
				ctx.writeValue(sb[i], true, "+")
				ctx.result(NoMatch)
			}
			if i != max-1 {
				ctx.newline(",")
			} else {
				ctx.level--
				ctx.newline("")
			}
		}
		ctx.buf.WriteString("]")
		return
	case reflect.Map:
		ma, mb := a.(map[string]interface{}), b.(map[string]interface{})
		keysMap := make(map[string]bool)
		for k := range ma {
			keysMap[k] = true
		}
		for k := range mb {
			keysMap[k] = true
		}
		keys := make([]string, 0, len(keysMap))
		for k := range keysMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		if len(keys) == 0 {
			ctx.buf.WriteString("{")
		} else {
			ctx.level++
			ctx.newline("{")
		}
		for i, k := range keys {
			va, aok := ma[k]
			vb, bok := mb[k]
			if aok && bok {
				ctx.key(k)
				ctx.calculateDiff(va, vb)
			} else if aok {
				ctx.key(k)
				ctx.writeValue(va, true, "-")
				ctx.result(SupersetMatch)
			} else if bok {
				ctx.key(k)
				ctx.writeValue(vb, true, "+")
				ctx.result(NoMatch)
			}
			if i != len(keys)-1 {
				ctx.newline(",")
			} else {
				ctx.level--
				ctx.newline("")
			}
		}
		ctx.buf.WriteString("}")
		return
	}
	ctx.writeValue(a, true, "")
	ctx.result(FullMatch)
}
func Compare(a, b []byte) ([]byte, error) {
	var av, bv interface{}
	da := json.NewDecoder(bytes.NewReader(a))
	da.UseNumber()
	db := json.NewDecoder(bytes.NewReader(b))
	db.UseNumber()
	errA := da.Decode(&av)
	errB := db.Decode(&bv)
	if errA != nil && errB != nil {
		return nil, errors.New("BothArgsAreInvalidJson both arguments are invalid json")
	}
	if errA != nil {
		return nil, errors.New("FirstArgIsInvalidJson first argument is invalid json")
	}
	if errB != nil {
		return nil, errors.New("SecondArgIsInvalidJson second argument is invalid json")
	}

	ctx := context{
		buf: bytes.Buffer{},
	}
	ctx.calculateDiff(av, bv)

	diff := ctx.buf.Bytes()
	diff1 := map[string]interface{}{}
	err := json.Unmarshal(diff, &diff1)
	if err != nil {
		log.Fatalln("结果json.unmarshal 失败：", err)
	}else{
		//fmt.Println("解析之后的", diff1)
	}
	return diff, nil
}
