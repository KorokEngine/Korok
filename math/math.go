package math

import (
	"unsafe"
	"math/rand"
)

/// This is A approximate yet fast inverse square-root.
func InvSqrt(x float32) float32 {
	xhalf := float32(0.5) * x
	i := *(*int32)(unsafe.Pointer(&x))
	i = int32(0x5f3759df) - int32(i>>1)
	x = *(*float32)(unsafe.Pointer(&i))
	x = x * (1.5 - (xhalf * x * x))
	return x
}

/// a faster way ?
func Random(low, high float32) float32 {
	return low + (high - low) * rand.Float32()
}

func UInt32_clamp(v, low, high uint32) uint32 {
	return 0
}

func UInt32_min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func UInt32_max(a, b uint32) uint32 {
	if a < b {
		return b
	}
	return a
}

func UInt16_min(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

func UInt16_max(a, b uint16) uint16 {
	if a < b {
		return b
	}
	return a
}

func UInt32_cnttz(v uint32) uint32{
	return 0
}