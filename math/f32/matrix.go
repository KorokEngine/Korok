package f32

import (
	"golang.org/x/image/math/f32"
	"math"
)

type Mat3 f32.Mat3
type Mat4 f32.Mat4

// SetCol sets a Column within the Matrix, so it mutates the calling matrix.
func (m *Mat3) SetCol(col int, v Vec3) {
	m[col*3+0], m[col*3+1], m[col*3+2] = v[0], v[1], v[2]
}

// SetRow sets a Row within the Matrix, so it mutates the calling matrix.
func (m *Mat3) SetRow(row int, v Vec3) {
	m[row+0], m[row+3], m[row+6] = v[0], v[1], v[2]
}

// Diag returns main diagonal (meaning all elements such that row==col).
func (m Mat3) Diag() Vec3 {
	return Vec3{m[0], m[4], m[8]}
}

// Transform transforms (x, y) to (x1, y).
//
// matrix multiplication carried out on paper:
// |1    x| |c -s  | |sx     | | 1 ky  | |1   -ox|
// |  1  y| |s  c  | |   sy  | |kx  1  | |  1 -oy|
// |     1| |     1| |      1| |      1| |     1 |
//   move    rotate    scale     skew      origin
func (m Mat3) Transform(x, y float32) (x1, y1 float32) {
	x1 = m[0]*x + m[3]*y + m[6]
	y1 = m[1]*x + m[4]*y + m[7]
	return
}

// Initialize defines a 3-D Matrix.
//            | x |
//            | y |
//            | 1 |
// | e0 e3 e6 |
// | e1 e4 e7 |
// | e2 e5 e8 |
func (m *Mat3) Initialize(x, y, angle, sx, sy, ox, oy, kx, ky float32) {
	c, s := cos(angle), sin(angle)

	m[0] = c * sx - ky * s * sy // = a
	m[1] = s * sx + ky * c * sy // = b
	m[3] = kx * c * sx - s * sy // = c
	m[4] = kx * s * sx + c * sy // = d
	m[6] = x - ox * m[0] - oy * m[3]
	m[7] = y - ox * m[1] - oy * m[4]

	m[2], m[5] = 0, 0
	m[8] = 1.0
}

func (m *Mat3) InitializeScale1(x, y, angle, ox, oy float32) {
	c, s := cos(angle), sin(angle)

	m[0] = c // = a
	m[1] = s // = b
	m[3] = - s // = c
	m[4] = + c // = d
	m[6] = x - ox * m[0] - oy * m[3]
	m[7] = y - ox * m[1] - oy * m[4]

	m[2], m[5] = 0, 0
	m[8] = 1.0
}

func sin(r float32) float32 {
	return float32(math.Sin(float64(r)))
}

func cos(r float32) float32 {
	return float32(math.Cos(float64(r)))
}

// Ident3 returns the 3x3 identity matrix.
func Ident3() Mat3 {
	return Mat3{1, 0, 0, 0, 1, 0, 0, 0, 1}
}

// SetCol sets a Column within the Matrix.
func (m *Mat4) SetCol(col int, v Vec4) {
	m[col*4+0], m[col*4+1], m[col*4+2], m[col*4+3] = v[0], v[1], v[2], v[3]
}

// SetRow sets a Row within the Matrix.
func (m *Mat4) SetRow(row int, v Vec4) {
	m[row+0], m[row+4], m[row+8], m[row+12] = v[0], v[1], v[2], v[3]
}

// Diag returns main diagonal (meaning all elements such that row==col).
func (m Mat4) Diag() Vec4 {
	return Vec4{m[0], m[5], m[10], m[15]}
}

// Ident4 returns the 4x4 identity matrix.
func Ident4() Mat4 {
	return Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// At returns the matrix element at the given row and column.
func (m Mat4) At(row, col int) float32 {
	return m[col*4+row]
}

// Set sets the corresponding matrix element at the given row and column.
func (m *Mat4) Set(row, col int, value float32) {
	m[col*4+row] = value
}
