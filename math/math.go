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