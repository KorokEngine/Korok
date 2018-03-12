package f32

import (
	"golang.org/x/image/math/f32"
)

type Mat3 f32.Mat3
type Mat4 f32.Mat4

// Sets a Column within the Matrix, so it mutates the calling matrix.
func (m *Mat3) SetCol(col int, v Vec3) {
	m[col*3+0], m[col*3+1], m[col*3+2] = v[0], v[1], v[2]
}

// Sets a Row within the Matrix, so it mutates the calling matrix.
func (m *Mat3) SetRow(row int, v Vec3) {
	m[row+0], m[row+3], m[row+6] = v[0], v[1], v[2]
}

// Diag is a basic operation on a square matrix that simply
// returns main diagonal (meaning all elements such that row==col).
func (m Mat3) Diag() Vec3 {
	return Vec3{m[0], m[4], m[8]}
}

// Ident3 returns the 3x3 identity matrix.
// The identity matrix is a square matrix with the value 1 on its
// diagonals. The characteristic property of the identity matrix is that
// any matrix multiplied by it is itself. (MI = M; IN = N)
func Ident3() Mat3 {
	return Mat3{1, 0, 0, 0, 1, 0, 0, 0, 1}
}

// Sets a Column within the Matrix, so it mutates the calling matrix.
func (m *Mat4) SetCol(col int, v Vec4) {
	m[col*4+0], m[col*4+1], m[col*4+2], m[col*4+3] = v[0], v[1], v[2], v[3]
}

// Sets a Row within the Matrix, so it mutates the calling matrix.
func (m *Mat4) SetRow(row int, v Vec4) {
	m[row+0], m[row+4], m[row+8], m[row+12] = v[0], v[1], v[2], v[3]
}

// Diag is a basic operation on a square matrix that simply
// returns main diagonal (meaning all elements such that row==col).
func (m Mat4) Diag() Vec4 {
	return Vec4{m[0], m[5], m[10], m[15]}
}

// Ident4 returns the 4x4 identity matrix.
// The identity matrix is a square matrix with the value 1 on its
// diagonals. The characteristic property of the identity matrix is that
// any matrix multiplied by it is itself. (MI = M; IN = N)
func Ident4() Mat4 {
	return Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// At returns the matrix element at the given row and column.
// This is equivalent to mat[col * numRow + row] where numRow is constant
// (E.G. for a Mat3x2 it's equal to 3)
//
// This method is garbage-in garbage-out. For instance, on a Mat4 asking for
// At(5,0) will work just like At(1,1). Or it may panic if it's out of bounds.
func (m Mat4) At(row, col int) float32 {
	return m[col*4+row]
}

// Set sets the corresponding matrix element at the given row and column.
// This has a pointer receiver because it mutates the matrix.
//
// This method is garbage-in garbage-out. For instance, on a Mat4 asking for
// Set(5,0,val) will work just like Set(1,1,val). Or it may panic if it's out of bounds.
func (m *Mat4) Set(row, col int, value float32) {
	m[col*4+row] = value
}
