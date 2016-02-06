/* (C) Christof Fetzer, 2016

 A simple extreme values generator that can be used in combination
 with testing/quick.
 
 Example:
 
 func TestAbs2e(t *testing.T) {
	f := func(x int) bool {
		return Abs2(x) >= 0       // post-condition of Abs2
	}
	if err := quick.Check(f, &quick.Config{Values: extremeValues.ExtremeValues(f)}); err != nil {
		t.Error(err)
	}
}

*/


package extreme

import (
	"math"
	"math/rand"
	"testing/quick"
	"reflect"
    "log"
)

const maxSize = 50

type xValueType int
const (
    xMin xValueType = iota      // minimum value
    xMax                        // maximum value
    xZero                       // zero value
    xRnd                        // arbitrary value
)

type zValueType int
const (
    zMin zValueType = iota      // slightly above
    zMax                        // slightly below
    zZero                       // zero value
)

// rndValueType returns random value of type xValueType
func rndValueType(rnd *rand.Rand) xValueType {
    return xValueType(rnd.Int31n(int32(xRnd+1)))
}

// rndZValueType returns random value of type zValueType
func rndZValueType(rnd *rand.Rand) zValueType {
    return zValueType(rnd.Int31n(int32(zZero+1)))
}


// ExtremeValues returns a function that returns either a random value or an extreme value
func Values(f interface{}) func([]reflect.Value, *rand.Rand) {
    v := reflect.ValueOf(f)
    if v.Kind() != reflect.Func {
        return nil        
    }
    g := func(a []reflect.Value, r *rand.Rand)  {
        values(a, v.Type(), r)
    }
    return g
}

// values returns
func values(args []reflect.Value, f reflect.Type, rand *rand.Rand) {
    for j := 0; j < len(args); j++ {
        var ok bool
        switch rndValueType(rand) {
        case xMin:
           args[j], ok = minValue(f.In(j), rand)
        case xMax: 
           args[j], ok = maxValue(f.In(j), rand)
        case xZero:
           args[j], ok = zeroValue(f.In(j), rand)
        default:
           args[j], ok = quick.Value(f.In(j), rand)
        }
		if !ok {
            log.Printf("Error creating random value for type %#v\n", f.In(j))
		}
    }
    return
}


// minValue returns a minimal value of type t (or random value if t is non-scalar type)
func minValue(t reflect.Type, rand *rand.Rand) (value reflect.Value, ok bool) {
	v := reflect.New(t).Elem()
	switch concrete := t; concrete.Kind() {
	case reflect.Bool:
		v.SetBool(false)
	case reflect.Float32:
		v.SetFloat(-math.MaxFloat32)
	case reflect.Float64:
		v.SetFloat(-math.MaxFloat64)
	case reflect.Complex64:
		v.SetComplex(complex(-math.MaxFloat32, -math.MaxFloat32))
	case reflect.Complex128:
		v.SetComplex(complex(math.MaxFloat64, math.MaxFloat64))
	case reflect.Int16:
		v.SetInt(math.MinInt16)
	case reflect.Int32:
		v.SetInt(math.MinInt32)
	case reflect.Int64:
		v.SetInt(math.MinInt64)
	case reflect.Int8:
		v.SetInt(math.MinInt8)
	case reflect.Int:
		v.SetInt(math.MinInt64)
	case reflect.Uint16:
		v.SetUint(uint64(1))
	case reflect.Uint32:
		v.SetUint(uint64(1))
	case reflect.Uint64:
		v.SetUint(uint64(1))
	case reflect.Uint8:
		v.SetUint(uint64(1))
	case reflect.Uint:
		v.SetUint(uint64(1))
	case reflect.Uintptr:
		v.SetUint(uint64(1))
	default:
		return quick.Value(t, rand)
	}

	return v, true
}


// maxValue returns a maximum value of type t (or random value if t is non-scalar type)
func maxValue(t reflect.Type, rand *rand.Rand) (value reflect.Value, ok bool) {
	v := reflect.New(t).Elem()
	switch concrete := t; concrete.Kind() {
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Float32:
		v.SetFloat(math.MaxFloat32)
	case reflect.Float64:
		v.SetFloat(math.MaxFloat64)
	case reflect.Complex64:
		v.SetComplex(complex(float64(math.MaxFloat32), float64(math.MaxFloat32)))
	case reflect.Complex128:
		v.SetComplex(complex(math.MaxFloat64, math.MaxFloat64))
	case reflect.Int16:
		v.SetInt(math.MaxInt16)
	case reflect.Int32:
		v.SetInt(math.MaxInt32)
	case reflect.Int64:
		v.SetInt(math.MaxInt64)
	case reflect.Int8:
		v.SetInt(math.MaxInt8)
	case reflect.Int:
		v.SetInt(math.MaxInt64)
	case reflect.Uint16:
		v.SetUint(math.MaxUint16)
	case reflect.Uint32:
		v.SetUint(math.MaxUint32)
	case reflect.Uint64:
		v.SetUint(math.MaxUint64)
	case reflect.Uint8:
		v.SetUint(math.MaxUint8)
	case reflect.Uint:
		v.SetUint(math.MaxUint64)
	case reflect.Uintptr:
		v.SetUint(math.MaxUint64)
	default:
		return quick.Value(t, rand)
	}

	return v, true
}


// zeroValue returns a zero value or a value close to zero type t (or random value if t is non-scalar type)
func zeroValue(t reflect.Type, rand *rand.Rand) (value reflect.Value, ok bool) {
    zt := rndZValueType(rand)
	v := reflect.New(t).Elem()
	switch concrete := t; concrete.Kind() {
	case reflect.Bool:
		v.SetBool(false)
	case reflect.Float32:
        switch zt {
        case zMin:
		  v.SetFloat(-math.SmallestNonzeroFloat32)
        case zMax:
		  v.SetFloat(math.SmallestNonzeroFloat32)
        default:
		  v.SetFloat(0.0)
        }
	case reflect.Float64:
        switch zt {
        case zMin:
		  v.SetFloat(-math.SmallestNonzeroFloat64)
        case zMax:
		  v.SetFloat(math.SmallestNonzeroFloat64)
        default:
		  v.SetFloat(0.0)
        }
	case reflect.Complex64:
        switch zt {
        case zMin:
		  v.SetComplex(complex(-math.SmallestNonzeroFloat32, -math.SmallestNonzeroFloat32))
        case zMax:
		  v.SetComplex(complex(+math.SmallestNonzeroFloat32, +math.SmallestNonzeroFloat32))
        default:
		  v.SetComplex(complex(0.0, 0.0))
        }
	case reflect.Complex128:
        switch zt {
        case zMin:
		  v.SetComplex(complex(-math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64))
        case zMax:
		  v.SetComplex(complex(+math.SmallestNonzeroFloat64, +math.SmallestNonzeroFloat64))
        default:
		  v.SetComplex(complex(0.0, 0.0))
        }
	case reflect.Int16:
		v.SetInt(0)
	case reflect.Int32:
		v.SetInt(0)
	case reflect.Int64:
		v.SetInt(0)
	case reflect.Int8:
		v.SetInt(0)
	case reflect.Int:
		v.SetInt(0)
	case reflect.Uint16:
		v.SetUint(uint64(0))
	case reflect.Uint32:
		v.SetUint(uint64(0))
	case reflect.Uint64:
		v.SetUint(uint64(0))
	case reflect.Uint8:
		v.SetUint(uint64(0))
	case reflect.Uint:
		v.SetUint(uint64(0))
	case reflect.Uintptr:
		v.SetUint(uint64(0))
	default:
		return quick.Value(t, rand)
	}

	return v, true
}
