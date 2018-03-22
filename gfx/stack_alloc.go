package gfx

import "unsafe"

// // This is a stack allocator used for fast per step allocations.
type StackAllocator struct {
	// bottom of the stack
	data uintptr

	// current position and capacity
	alloc, cap int32

	// keep a reference
	ref unsafe.Pointer
}

func (stack *StackAllocator) initialize(max int) {
	mem := make([]byte, max)
	stack.ref = unsafe.Pointer(&mem[0])
	stack.data = uintptr(stack.ref)
}

// Alloc alloc memory on the stack.
func (stack *StackAllocator) Alloc(size int32) (ptr unsafe.Pointer) {
	if req := stack.alloc + size; req < stack.cap {
		ptr = unsafe.Pointer(stack.data + uintptr(stack.alloc))
		stack.alloc += size
	} else {
		tmp := make([]byte, size)
		ptr = unsafe.Pointer(&tmp[0])
	}
	return
}

// Free frees memory on the stack, it should pair with Alloc.
func (stack *StackAllocator) Free(size int32) {
	stack.alloc -= size
	if stack.alloc < 0 {
		stack.alloc = 0
	}
	return
}

// release will empty the stack each frame.
func (stack *StackAllocator) release() {
	stack.alloc = 0
}


