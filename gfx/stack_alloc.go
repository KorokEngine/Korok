package gfx

import "unsafe"

// A stack for per-frame memory allocation
type StackAllocator struct {
	data uintptr
	alloc int32

	// keep a reference
	ref unsafe.Pointer
}

func NewStackAllocator() *StackAllocator {
	stack := new(StackAllocator)
	mem := make([]byte, 100 * 1024)
	stack.ref = unsafe.Pointer(&mem[0])
	stack.data = uintptr(stack.ref)
	return stack
}

func (stack *StackAllocator) Alloc(size int32) unsafe.Pointer {
	ptr := unsafe.Pointer(stack.data + uintptr(stack.alloc))
	stack.alloc += size
	return ptr
}

func (stack *StackAllocator) Free(size int32) {
	stack.alloc -= size
	if stack.alloc < 0 {
		stack.alloc = 0
	}
	return
}

func (stack *StackAllocator) Reset() {
	stack.alloc = 0
}


