// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f32

import "math"

// Rotate3DY returns a 3x3 (non-homogeneous) Matrix that rotates by angle about the Y-axis
//
// Where c is cos(angle) and s is sin(angle)
//    [c 0 s]
//    [0 1 0]
//    [s 0 c ]
func Rotate3DY(angle float32) Mat3 {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := float32(math.Sin(float64(angle))), float32(math.Cos(float64(angle)))

	return Mat3{cos, 0, -sin, 0, 1, 0, sin, 0, cos}
}

// Rotate3DZ returns a 3x3 (non-homogeneous) Matrix that rotates by angle about the Z-axis
//
// Where c is cos(angle) and s is sin(angle)
//    [c -s 0]
//    [s c 0]
//    [0 0 1 ]
func Rotate3DZ(angle float32) Mat3 {
	//angle = (angle * math.Pi) / 180.0
	sin, cos := float32(math.Sin(float64(angle))), float32(math.Cos(float64(angle)))

	return Mat3{cos, sin, 0, -sin, cos, 0, 0, 0, 1}
}

// Translate2D returns a homogeneous (3x3 for 2D-space) Translation matrix that moves a point by Tx units in the x-direction and Ty units in the y-direction
//
//    [[1, 0, Tx]]
//    [[0, 1, Ty]]
//    [[0, 0, 1 ]]
func Translate2D(Tx, Ty float32) Mat3 {
	return Mat3{1, 0, 0, 0, 1, 0, float32(Tx), float32(Ty), 1}
}

// Translate3D returns a homogeneous (4x4 for 3D-space) Translation matrix that moves a point by Tx units in the x-direction, Ty units in the y-direction,
// and Tz units in the z-direction
//
//    [[1, 0, 0, Tx]]
//    [[0, 1, 0, Ty]]
//    [[0, 0, 1, Tz]]
//    [[0, 0, 0, 1 ]]
func Translate3D(Tx, Ty, Tz float32) Mat4 {
	return Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, float32(Tx), float32(Ty), float32(Tz), 1}
}
