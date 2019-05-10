// +build js

// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webav

import (
	"errors"

	"syscall/js"
)

type ContextAttributes struct {
	// If Alpha is true, the drawing buffer has an alpha channel for
	// the purposes of performing OpenGL destination alpha operations
	// and compositing with the page.
	Alpha bool

	// If Depth is true, the drawing buffer has a depth buffer of at least 16 bits.
	Depth bool

	// If Stencil is true, the drawing buffer has a stencil buffer of at least 8 bits.
	Stencil bool

	// If Antialias is true and the implementation supports antialiasing
	// the drawing buffer will perform antialiasing using its choice of
	// technique (multisample/supersample) and quality.
	Antialias bool

	// If PremultipliedAlpha is true the page compositor will assume the
	// drawing buffer contains colors with premultiplied alpha.
	// This flag is ignored if the alpha flag is false.
	PremultipliedAlpha bool

	// If the value is true the buffers will not be cleared and will preserve
	// their values until cleared or overwritten by the author.
	PreserveDrawingBuffer bool
}

func mb2mi(x map[string]bool) map[string]interface{} {
	y := make(map[string]interface{})
	for k, v := range x {
		y[k] = v
	}
	return y
}

// Returns a copy of the default WebGL context attributes.
func DefaultAttributes() *ContextAttributes {
	return &ContextAttributes{true, true, false, true, true, false}
}

type Context struct {
	js.Value
	ARRAY_BUFFER                                 int `js:"ARRAY_BUFFER"`
	ARRAY_BUFFER_BINDING                         int `js:"ARRAY_BUFFER_BINDING"`
	ATTACHED_SHADERS                             int `js:"ATTACHED_SHADERS"`
	BACK                                         int `js:"BACK"`
	BLEND                                        int `js:"BLEND"`
	BLEND_COLOR                                  int `js:"BLEND_COLOR"`
	BLEND_DST_ALPHA                              int `js:"BLEND_DST_ALPHA"`
	BLEND_DST_RGB                                int `js:"BLEND_DST_RGB"`
	BLEND_EQUATION                               int `js:"BLEND_EQUATION"`
	BLEND_EQUATION_ALPHA                         int `js:"BLEND_EQUATION_ALPHA"`
	BLEND_EQUATION_RGB                           int `js:"BLEND_EQUATION_RGB"`
	BLEND_SRC_ALPHA                              int `js:"BLEND_SRC_ALPHA"`
	BLEND_SRC_RGB                                int `js:"BLEND_SRC_RGB"`
	BLUE_BITS                                    int `js:"BLUE_BITS"`
	BOOL                                         int `js:"BOOL"`
	BOOL_VEC2                                    int `js:"BOOL_VEC2"`
	BOOL_VEC3                                    int `js:"BOOL_VEC3"`
	BOOL_VEC4                                    int `js:"BOOL_VEC4"`
	BROWSER_DEFAULT_WEBGL                        int `js:"BROWSER_DEFAULT_WEBGL"`
	BUFFER_SIZE                                  int `js:"BUFFER_SIZE"`
	BUFFER_USAGE                                 int `js:"BUFFER_USAGE"`
	BYTE                                         int `js:"BYTE"`
	CCW                                          int `js:"CCW"`
	CLAMP_TO_EDGE                                int `js:"CLAMP_TO_EDGE"`
	COLOR_ATTACHMENT0                            int `js:"COLOR_ATTACHMENT0"`
	COLOR_BUFFER_BIT                             int `js:"COLOR_BUFFER_BIT"`
	COLOR_CLEAR_VALUE                            int `js:"COLOR_CLEAR_VALUE"`
	COLOR_WRITEMASK                              int `js:"COLOR_WRITEMASK"`
	COMPILE_STATUS                               int `js:"COMPILE_STATUS"`
	COMPRESSED_TEXTURE_FORMATS                   int `js:"COMPRESSED_TEXTURE_FORMATS"`
	CONSTANT_ALPHA                               int `js:"CONSTANT_ALPHA"`
	CONSTANT_COLOR                               int `js:"CONSTANT_COLOR"`
	CONTEXT_LOST_WEBGL                           int `js:"CONTEXT_LOST_WEBGL"`
	CULL_FACE                                    int `js:"CULL_FACE"`
	CULL_FACE_MODE                               int `js:"CULL_FACE_MODE"`
	CURRENT_PROGRAM                              int `js:"CURRENT_PROGRAM"`
	CURRENT_VERTEX_ATTRIB                        int `js:"CURRENT_VERTEX_ATTRIB"`
	CW                                           int `js:"CW"`
	DECR                                         int `js:"DECR"`
	DECR_WRAP                                    int `js:"DECR_WRAP"`
	DELETE_STATUS                                int `js:"DELETE_STATUS"`
	DEPTH_ATTACHMENT                             int `js:"DEPTH_ATTACHMENT"`
	DEPTH_BITS                                   int `js:"DEPTH_BITS"`
	DEPTH_BUFFER_BIT                             int `js:"DEPTH_BUFFER_BIT"`
	DEPTH_CLEAR_VALUE                            int `js:"DEPTH_CLEAR_VALUE"`
	DEPTH_COMPONENT                              int `js:"DEPTH_COMPONENT"`
	DEPTH_COMPONENT16                            int `js:"DEPTH_COMPONENT16"`
	DEPTH_FUNC                                   int `js:"DEPTH_FUNC"`
	DEPTH_RANGE                                  int `js:"DEPTH_RANGE"`
	DEPTH_STENCIL                                int `js:"DEPTH_STENCIL"`
	DEPTH_STENCIL_ATTACHMENT                     int `js:"DEPTH_STENCIL_ATTACHMENT"`
	DEPTH_TEST                                   int `js:"DEPTH_TEST"`
	DEPTH_WRITEMASK                              int `js:"DEPTH_WRITEMASK"`
	DITHER                                       int `js:"DITHER"`
	DONT_CARE                                    int `js:"DONT_CARE"`
	DST_ALPHA                                    int `js:"DST_ALPHA"`
	DST_COLOR                                    int `js:"DST_COLOR"`
	DYNAMIC_DRAW                                 int `js:"DYNAMIC_DRAW"`
	ELEMENT_ARRAY_BUFFER                         int `js:"ELEMENT_ARRAY_BUFFER"`
	ELEMENT_ARRAY_BUFFER_BINDING                 int `js:"ELEMENT_ARRAY_BUFFER_BINDING"`
	EQUAL                                        int `js:"EQUAL"`
	FASTEST                                      int `js:"FASTEST"`
	FLOAT                                        int `js:"FLOAT"`
	FLOAT_MAT2                                   int `js:"FLOAT_MAT2"`
	FLOAT_MAT3                                   int `js:"FLOAT_MAT3"`
	FLOAT_MAT4                                   int `js:"FLOAT_MAT4"`
	FLOAT_VEC2                                   int `js:"FLOAT_VEC2"`
	FLOAT_VEC3                                   int `js:"FLOAT_VEC3"`
	FLOAT_VEC4                                   int `js:"FLOAT_VEC4"`
	FRAGMENT_SHADER                              int `js:"FRAGMENT_SHADER"`
	FRAMEBUFFER                                  int `js:"FRAMEBUFFER"`
	FRAMEBUFFER_ATTACHMENT_OBJECT_NAME           int `js:"FRAMEBUFFER_ATTACHMENT_OBJECT_NAME"`
	FRAMEBUFFER_ATTACHMENT_OBJECT_TYPE           int `js:"FRAMEBUFFER_ATTACHMENT_OBJECT_TYPE"`
	FRAMEBUFFER_ATTACHMENT_TEXTURE_CUBE_MAP_FACE int `js:"FRAMEBUFFER_ATTACHMENT_TEXTURE_CUBE_MAP_FACE"`
	FRAMEBUFFER_ATTACHMENT_TEXTURE_LEVEL         int `js:"FRAMEBUFFER_ATTACHMENT_TEXTURE_LEVEL"`
	FRAMEBUFFER_BINDING                          int `js:"FRAMEBUFFER_BINDING"`
	FRAMEBUFFER_COMPLETE                         int `js:"FRAMEBUFFER_COMPLETE"`
	FRAMEBUFFER_INCOMPLETE_ATTACHMENT            int `js:"FRAMEBUFFER_INCOMPLETE_ATTACHMENT"`
	FRAMEBUFFER_INCOMPLETE_DIMENSIONS            int `js:"FRAMEBUFFER_INCOMPLETE_DIMENSIONS"`
	FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT    int `js:"FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT"`
	FRAMEBUFFER_UNSUPPORTED                      int `js:"FRAMEBUFFER_UNSUPPORTED"`
	FRONT                                        int `js:"FRONT"`
	FRONT_AND_BACK                               int `js:"FRONT_AND_BACK"`
	FRONT_FACE                                   int `js:"FRONT_FACE"`
	FUNC_ADD                                     int `js:"FUNC_ADD"`
	FUNC_REVERSE_SUBTRACT                        int `js:"FUNC_REVERSE_SUBTRACT"`
	FUNC_SUBTRACT                                int `js:"FUNC_SUBTRACT"`
	GENERATE_MIPMAP_HINT                         int `js:"GENERATE_MIPMAP_HINT"`
	GEQUAL                                       int `js:"GEQUAL"`
	GREATER                                      int `js:"GREATER"`
	GREEN_BITS                                   int `js:"GREEN_BITS"`
	HIGH_FLOAT                                   int `js:"HIGH_FLOAT"`
	HIGH_INT                                     int `js:"HIGH_INT"`
	INCR                                         int `js:"INCR"`
	INCR_WRAP                                    int `js:"INCR_WRAP"`
	INFO_LOG_LENGTH                              int `js:"INFO_LOG_LENGTH"`
	INT                                          int `js:"INT"`
	INT_VEC2                                     int `js:"INT_VEC2"`
	INT_VEC3                                     int `js:"INT_VEC3"`
	INT_VEC4                                     int `js:"INT_VEC4"`
	INVALID_ENUM                                 int `js:"INVALID_ENUM"`
	INVALID_FRAMEBUFFER_OPERATION                int `js:"INVALID_FRAMEBUFFER_OPERATION"`
	INVALID_OPERATION                            int `js:"INVALID_OPERATION"`
	INVALID_VALUE                                int `js:"INVALID_VALUE"`
	INVERT                                       int `js:"INVERT"`
	KEEP                                         int `js:"KEEP"`
	LEQUAL                                       int `js:"LEQUAL"`
	LESS                                         int `js:"LESS"`
	LINEAR                                       int `js:"LINEAR"`
	LINEAR_MIPMAP_LINEAR                         int `js:"LINEAR_MIPMAP_LINEAR"`
	LINEAR_MIPMAP_NEAREST                        int `js:"LINEAR_MIPMAP_NEAREST"`
	LINES                                        int `js:"LINES"`
	LINE_LOOP                                    int `js:"LINE_LOOP"`
	LINE_STRIP                                   int `js:"LINE_STRIP"`
	LINE_WIDTH                                   int `js:"LINE_WIDTH"`
	LINK_STATUS                                  int `js:"LINK_STATUS"`
	LOW_FLOAT                                    int `js:"LOW_FLOAT"`
	LOW_INT                                      int `js:"LOW_INT"`
	LUMINANCE                                    int `js:"LUMINANCE"`
	LUMINANCE_ALPHA                              int `js:"LUMINANCE_ALPHA"`
	MAX_COMBINED_TEXTURE_IMAGE_UNITS             int `js:"MAX_COMBINED_TEXTURE_IMAGE_UNITS"`
	MAX_CUBE_MAP_TEXTURE_SIZE                    int `js:"MAX_CUBE_MAP_TEXTURE_SIZE"`
	MAX_FRAGMENT_UNIFORM_VECTORS                 int `js:"MAX_FRAGMENT_UNIFORM_VECTORS"`
	MAX_RENDERBUFFER_SIZE                        int `js:"MAX_RENDERBUFFER_SIZE"`
	MAX_TEXTURE_IMAGE_UNITS                      int `js:"MAX_TEXTURE_IMAGE_UNITS"`
	MAX_TEXTURE_SIZE                             int `js:"MAX_TEXTURE_SIZE"`
	MAX_VARYING_VECTORS                          int `js:"MAX_VARYING_VECTORS"`
	MAX_VERTEX_ATTRIBS                           int `js:"MAX_VERTEX_ATTRIBS"`
	MAX_VERTEX_TEXTURE_IMAGE_UNITS               int `js:"MAX_VERTEX_TEXTURE_IMAGE_UNITS"`
	MAX_VERTEX_UNIFORM_VECTORS                   int `js:"MAX_VERTEX_UNIFORM_VECTORS"`
	MAX_VIEWPORT_DIMS                            int `js:"MAX_VIEWPORT_DIMS"`
	MEDIUM_FLOAT                                 int `js:"MEDIUM_FLOAT"`
	MEDIUM_INT                                   int `js:"MEDIUM_INT"`
	MIRRORED_REPEAT                              int `js:"MIRRORED_REPEAT"`
	NEAREST                                      int `js:"NEAREST"`
	NEAREST_MIPMAP_LINEAR                        int `js:"NEAREST_MIPMAP_LINEAR"`
	NEAREST_MIPMAP_NEAREST                       int `js:"NEAREST_MIPMAP_NEAREST"`
	NEVER                                        int `js:"NEVER"`
	NICEST                                       int `js:"NICEST"`
	NONE                                         int `js:"NONE"`
	NOTEQUAL                                     int `js:"NOTEQUAL"`
	NO_ERROR                                     int `js:"NO_ERROR"`
	NUM_COMPRESSED_TEXTURE_FORMATS               int `js:"NUM_COMPRESSED_TEXTURE_FORMATS"`
	ONE                                          int `js:"ONE"`
	ONE_MINUS_CONSTANT_ALPHA                     int `js:"ONE_MINUS_CONSTANT_ALPHA"`
	ONE_MINUS_CONSTANT_COLOR                     int `js:"ONE_MINUS_CONSTANT_COLOR"`
	ONE_MINUS_DST_ALPHA                          int `js:"ONE_MINUS_DST_ALPHA"`
	ONE_MINUS_DST_COLOR                          int `js:"ONE_MINUS_DST_COLOR"`
	ONE_MINUS_SRC_ALPHA                          int `js:"ONE_MINUS_SRC_ALPHA"`
	ONE_MINUS_SRC_COLOR                          int `js:"ONE_MINUS_SRC_COLOR"`
	OUT_OF_MEMORY                                int `js:"OUT_OF_MEMORY"`
	PACK_ALIGNMENT                               int `js:"PACK_ALIGNMENT"`
	POINTS                                       int `js:"POINTS"`
	POLYGON_OFFSET_FACTOR                        int `js:"POLYGON_OFFSET_FACTOR"`
	POLYGON_OFFSET_FILL                          int `js:"POLYGON_OFFSET_FILL"`
	POLYGON_OFFSET_UNITS                         int `js:"POLYGON_OFFSET_UNITS"`
	RED_BITS                                     int `js:"RED_BITS"`
	RENDERBUFFER                                 int `js:"RENDERBUFFER"`
	RENDERBUFFER_ALPHA_SIZE                      int `js:"RENDERBUFFER_ALPHA_SIZE"`
	RENDERBUFFER_BINDING                         int `js:"RENDERBUFFER_BINDING"`
	RENDERBUFFER_BLUE_SIZE                       int `js:"RENDERBUFFER_BLUE_SIZE"`
	RENDERBUFFER_DEPTH_SIZE                      int `js:"RENDERBUFFER_DEPTH_SIZE"`
	RENDERBUFFER_GREEN_SIZE                      int `js:"RENDERBUFFER_GREEN_SIZE"`
	RENDERBUFFER_HEIGHT                          int `js:"RENDERBUFFER_HEIGHT"`
	RENDERBUFFER_INTERNAL_FORMAT                 int `js:"RENDERBUFFER_INTERNAL_FORMAT"`
	RENDERBUFFER_RED_SIZE                        int `js:"RENDERBUFFER_RED_SIZE"`
	RENDERBUFFER_STENCIL_SIZE                    int `js:"RENDERBUFFER_STENCIL_SIZE"`
	RENDERBUFFER_WIDTH                           int `js:"RENDERBUFFER_WIDTH"`
	RENDERER                                     int `js:"RENDERER"`
	REPEAT                                       int `js:"REPEAT"`
	REPLACE                                      int `js:"REPLACE"`
	RGB                                          int `js:"RGB"`
	RGB5_A1                                      int `js:"RGB5_A1"`
	RGB565                                       int `js:"RGB565"`
	RGBA                                         int `js:"RGBA"`
	RGBA4                                        int `js:"RGBA4"`
	SAMPLER_2D                                   int `js:"SAMPLER_2D"`
	SAMPLER_CUBE                                 int `js:"SAMPLER_CUBE"`
	SAMPLES                                      int `js:"SAMPLES"`
	SAMPLE_ALPHA_TO_COVERAGE                     int `js:"SAMPLE_ALPHA_TO_COVERAGE"`
	SAMPLE_BUFFERS                               int `js:"SAMPLE_BUFFERS"`
	SAMPLE_COVERAGE                              int `js:"SAMPLE_COVERAGE"`
	SAMPLE_COVERAGE_INVERT                       int `js:"SAMPLE_COVERAGE_INVERT"`
	SAMPLE_COVERAGE_VALUE                        int `js:"SAMPLE_COVERAGE_VALUE"`
	SCISSOR_BOX                                  int `js:"SCISSOR_BOX"`
	SCISSOR_TEST                                 int `js:"SCISSOR_TEST"`
	SHADER_COMPILER                              int `js:"SHADER_COMPILER"`
	SHADER_SOURCE_LENGTH                         int `js:"SHADER_SOURCE_LENGTH"`
	SHADER_TYPE                                  int `js:"SHADER_TYPE"`
	SHADING_LANGUAGE_VERSION                     int `js:"SHADING_LANGUAGE_VERSION"`
	SHORT                                        int `js:"SHORT"`
	SRC_ALPHA                                    int `js:"SRC_ALPHA"`
	SRC_ALPHA_SATURATE                           int `js:"SRC_ALPHA_SATURATE"`
	SRC_COLOR                                    int `js:"SRC_COLOR"`
	STATIC_DRAW                                  int `js:"STATIC_DRAW"`
	STENCIL_ATTACHMENT                           int `js:"STENCIL_ATTACHMENT"`
	STENCIL_BACK_FAIL                            int `js:"STENCIL_BACK_FAIL"`
	STENCIL_BACK_FUNC                            int `js:"STENCIL_BACK_FUNC"`
	STENCIL_BACK_PASS_DEPTH_FAIL                 int `js:"STENCIL_BACK_PASS_DEPTH_FAIL"`
	STENCIL_BACK_PASS_DEPTH_PASS                 int `js:"STENCIL_BACK_PASS_DEPTH_PASS"`
	STENCIL_BACK_REF                             int `js:"STENCIL_BACK_REF"`
	STENCIL_BACK_VALUE_MASK                      int `js:"STENCIL_BACK_VALUE_MASK"`
	STENCIL_BACK_WRITEMASK                       int `js:"STENCIL_BACK_WRITEMASK"`
	STENCIL_BITS                                 int `js:"STENCIL_BITS"`
	STENCIL_BUFFER_BIT                           int `js:"STENCIL_BUFFER_BIT"`
	STENCIL_CLEAR_VALUE                          int `js:"STENCIL_CLEAR_VALUE"`
	STENCIL_FAIL                                 int `js:"STENCIL_FAIL"`
	STENCIL_FUNC                                 int `js:"STENCIL_FUNC"`
	STENCIL_INDEX                                int `js:"STENCIL_INDEX"`
	STENCIL_INDEX8                               int `js:"STENCIL_INDEX8"`
	STENCIL_PASS_DEPTH_FAIL                      int `js:"STENCIL_PASS_DEPTH_FAIL"`
	STENCIL_PASS_DEPTH_PASS                      int `js:"STENCIL_PASS_DEPTH_PASS"`
	STENCIL_REF                                  int `js:"STENCIL_REF"`
	STENCIL_TEST                                 int `js:"STENCIL_TEST"`
	STENCIL_VALUE_MASK                           int `js:"STENCIL_VALUE_MASK"`
	STENCIL_WRITEMASK                            int `js:"STENCIL_WRITEMASK"`
	STREAM_DRAW                                  int `js:"STREAM_DRAW"`
	SUBPIXEL_BITS                                int `js:"SUBPIXEL_BITS"`
	TEXTURE                                      int `js:"TEXTURE"`
	TEXTURE0                                     int `js:"TEXTURE0"`
	TEXTURE1                                     int `js:"TEXTURE1"`
	TEXTURE2                                     int `js:"TEXTURE2"`
	TEXTURE3                                     int `js:"TEXTURE3"`
	TEXTURE4                                     int `js:"TEXTURE4"`
	TEXTURE5                                     int `js:"TEXTURE5"`
	TEXTURE6                                     int `js:"TEXTURE6"`
	TEXTURE7                                     int `js:"TEXTURE7"`
	TEXTURE8                                     int `js:"TEXTURE8"`
	TEXTURE9                                     int `js:"TEXTURE9"`
	TEXTURE10                                    int `js:"TEXTURE10"`
	TEXTURE11                                    int `js:"TEXTURE11"`
	TEXTURE12                                    int `js:"TEXTURE12"`
	TEXTURE13                                    int `js:"TEXTURE13"`
	TEXTURE14                                    int `js:"TEXTURE14"`
	TEXTURE15                                    int `js:"TEXTURE15"`
	TEXTURE16                                    int `js:"TEXTURE16"`
	TEXTURE17                                    int `js:"TEXTURE17"`
	TEXTURE18                                    int `js:"TEXTURE18"`
	TEXTURE19                                    int `js:"TEXTURE19"`
	TEXTURE20                                    int `js:"TEXTURE20"`
	TEXTURE21                                    int `js:"TEXTURE21"`
	TEXTURE22                                    int `js:"TEXTURE22"`
	TEXTURE23                                    int `js:"TEXTURE23"`
	TEXTURE24                                    int `js:"TEXTURE24"`
	TEXTURE25                                    int `js:"TEXTURE25"`
	TEXTURE26                                    int `js:"TEXTURE26"`
	TEXTURE27                                    int `js:"TEXTURE27"`
	TEXTURE28                                    int `js:"TEXTURE28"`
	TEXTURE29                                    int `js:"TEXTURE29"`
	TEXTURE30                                    int `js:"TEXTURE30"`
	TEXTURE31                                    int `js:"TEXTURE31"`
	TEXTURE_2D                                   int `js:"TEXTURE_2D"`
	TEXTURE_BINDING_2D                           int `js:"TEXTURE_BINDING_2D"`
	TEXTURE_BINDING_CUBE_MAP                     int `js:"TEXTURE_BINDING_CUBE_MAP"`
	TEXTURE_CUBE_MAP                             int `js:"TEXTURE_CUBE_MAP"`
	TEXTURE_CUBE_MAP_NEGATIVE_X                  int `js:"TEXTURE_CUBE_MAP_NEGATIVE_X"`
	TEXTURE_CUBE_MAP_NEGATIVE_Y                  int `js:"TEXTURE_CUBE_MAP_NEGATIVE_Y"`
	TEXTURE_CUBE_MAP_NEGATIVE_Z                  int `js:"TEXTURE_CUBE_MAP_NEGATIVE_Z"`
	TEXTURE_CUBE_MAP_POSITIVE_X                  int `js:"TEXTURE_CUBE_MAP_POSITIVE_X"`
	TEXTURE_CUBE_MAP_POSITIVE_Y                  int `js:"TEXTURE_CUBE_MAP_POSITIVE_Y"`
	TEXTURE_CUBE_MAP_POSITIVE_Z                  int `js:"TEXTURE_CUBE_MAP_POSITIVE_Z"`
	TEXTURE_MAG_FILTER                           int `js:"TEXTURE_MAG_FILTER"`
	TEXTURE_MIN_FILTER                           int `js:"TEXTURE_MIN_FILTER"`
	TEXTURE_WRAP_S                               int `js:"TEXTURE_WRAP_S"`
	TEXTURE_WRAP_T                               int `js:"TEXTURE_WRAP_T"`
	TRIANGLES                                    int `js:"TRIANGLES"`
	TRIANGLE_FAN                                 int `js:"TRIANGLE_FAN"`
	TRIANGLE_STRIP                               int `js:"TRIANGLE_STRIP"`
	UNPACK_ALIGNMENT                             int `js:"UNPACK_ALIGNMENT"`
	UNPACK_COLORSPACE_CONVERSION_WEBGL           int `js:"UNPACK_COLORSPACE_CONVERSION_WEBGL"`
	UNPACK_FLIP_Y_WEBGL                          int `js:"UNPACK_FLIP_Y_WEBGL"`
	UNPACK_PREMULTIPLY_ALPHA_WEBGL               int `js:"UNPACK_PREMULTIPLY_ALPHA_WEBGL"`
	UNSIGNED_BYTE                                int `js:"UNSIGNED_BYTE"`
	UNSIGNED_INT                                 int `js:"UNSIGNED_INT"`
	UNSIGNED_SHORT                               int `js:"UNSIGNED_SHORT"`
	UNSIGNED_SHORT_4_4_4_4                       int `js:"UNSIGNED_SHORT_4_4_4_4"`
	UNSIGNED_SHORT_5_5_5_1                       int `js:"UNSIGNED_SHORT_5_5_5_1"`
	UNSIGNED_SHORT_5_6_5                         int `js:"UNSIGNED_SHORT_5_6_5"`
	VALIDATE_STATUS                              int `js:"VALIDATE_STATUS"`
	VENDOR                                       int `js:"VENDOR"`
	VERSION                                      int `js:"VERSION"`
	VERTEX_ATTRIB_ARRAY_BUFFER_BINDING           int `js:"VERTEX_ATTRIB_ARRAY_BUFFER_BINDING"`
	VERTEX_ATTRIB_ARRAY_ENABLED                  int `js:"VERTEX_ATTRIB_ARRAY_ENABLED"`
	VERTEX_ATTRIB_ARRAY_NORMALIZED               int `js:"VERTEX_ATTRIB_ARRAY_NORMALIZED"`
	VERTEX_ATTRIB_ARRAY_POINTER                  int `js:"VERTEX_ATTRIB_ARRAY_POINTER"`
	VERTEX_ATTRIB_ARRAY_SIZE                     int `js:"VERTEX_ATTRIB_ARRAY_SIZE"`
	VERTEX_ATTRIB_ARRAY_STRIDE                   int `js:"VERTEX_ATTRIB_ARRAY_STRIDE"`
	VERTEX_ATTRIB_ARRAY_TYPE                     int `js:"VERTEX_ATTRIB_ARRAY_TYPE"`
	VERTEX_SHADER                                int `js:"VERTEX_SHADER"`
	VIEWPORT                                     int `js:"VIEWPORT"`
	ZERO                                         int `js:"ZERO"`
}

// NewContext takes an HTML5 canvas object and optional context attributes.
// If an error is returned it means you won't have access to WebGL
// functionality.
func NewContext(canvas js.Value, ca *ContextAttributes) (*Context, error) {
	if js.Global().Get("WebGLRenderingContext") == js.Undefined() {
		return nil, errors.New("Your browser doesn't appear to support webgl.")
	}

	if ca == nil {
		ca = DefaultAttributes()
	}

	attrs := map[string]bool{
		"alpha":                 ca.Alpha,
		"depth":                 ca.Depth,
		"stencil":               ca.Stencil,
		"antialias":             ca.Antialias,
		"premultipliedAlpha":    ca.PremultipliedAlpha,
		"preserveDrawingBuffer": ca.PreserveDrawingBuffer,
	}
	gl := canvas.Call("getContext", "webgl", mb2mi(attrs))
	if gl == js.Null() {
		gl = canvas.Call("getContext", "experimental-webgl", mb2mi(attrs))
		if gl == js.Null() {
			return nil, errors.New("Creating a webgl context has failed.")
		}
	}
	ctx := new(Context)
	ctx.Value = gl
	return ctx, nil
}

// Returns the context attributes active on the context. These values might
// be different than what was requested on context creation if the
// browser's implementation doesn't support a feature.
func (c *Context) GetContextAttributes() ContextAttributes {
	ca := c.Call("getContextAttributes")
	return ContextAttributes{
		ca.Get("alpha").Bool(),
		ca.Get("depth").Bool(),
		ca.Get("stencil").Bool(),
		ca.Get("antialias").Bool(),
		ca.Get("premultipliedAlpha").Bool(),
		ca.Get("preservedDrawingBuffer").Bool(),
	}
}

// Specifies the active texture unit.
func (c *Context) ActiveTexture(texture int) {
	c.Call("activeTexture", texture)
}

// Attaches a WebGLShader object to a WebGLProgram object.
func (c *Context) AttachShader(program js.Value, shader js.Value) {
	c.Call("attachShader", program, shader)
}

// Binds a generic vertex index to a user-defined attribute variable.
func (c *Context) BindAttribLocation(program js.Value, index int, name string) {
	c.Call("bindAttribLocation", program, index, name)
}

// Associates a buffer with a buffer target.
func (c *Context) BindBuffer(target int, buffer js.Value) {
	c.Call("bindBuffer", target, buffer)
}

// Associates a WebGLFramebuffer object with the FRAMEBUFFER bind target.
func (c *Context) BindFramebuffer(target int, framebuffer js.Value) {
	c.Call("bindFramebuffer", target, framebuffer)
}

// Binds a WebGLRenderbuffer object to be used for rendering.
func (c *Context) BindRenderbuffer(target int, renderbuffer js.Value) {
	c.Call("bindRenderbuffer", target, renderbuffer)
}

// Binds a named texture object to a target.
func (c *Context) BindTexture(target int, texture js.Value) {
	c.Call("bindTexture", target, texture)
}

// The GL_BLEND_COLOR may be used to calculate the source and destination blending factors.
func (c *Context) BlendColor(r, g, b, a float64) {
	c.Call("blendColor", r, g, b, a)
}

// Sets the equation used to blend RGB and Alpha values of an incoming source
// fragment with a destination values as stored in the fragment's frame buffer.
func (c *Context) BlendEquation(mode int) {
	c.Call("blendEquation", mode)
}

// Controls the blending of an incoming source fragment's R, G, B, and A values
// with a destination R, G, B, and A values as stored in the fragment's WebGLFramebuffer.
func (c *Context) BlendEquationSeparate(modeRGB, modeAlpha int) {
	c.Call("blendEquationSeparate", modeRGB, modeAlpha)
}

// Sets the blending factors used to combine source and destination pixels.
func (c *Context) BlendFunc(sfactor, dfactor int) {
	c.Call("blendFunc", sfactor, dfactor)
}

// Sets the weighting factors that are used by blendEquationSeparate.
func (c *Context) BlendFuncSeparate(srcRGB, dstRGB, srcAlpha, dstAlpha int) {
	c.Call("blendFuncSeparate", srcRGB, dstRGB, srcAlpha, dstAlpha)
}

// Creates a buffer in memory and initializes it with array data.
// If no array is provided, the contents of the buffer is initialized to 0.
func (c *Context) BufferData(target int, data interface{}, usage int) {
	c.Call("bufferData", target, data, usage)
}

// Used to modify or update some or all of a data store for a bound buffer object.
func (c *Context) BufferSubData(target int, offset int, data interface{}) {
	c.Call("bufferSubData", target, offset, data)
}

// Returns whether the currently bound WebGLFramebuffer is complete.
// If not complete, returns the reason why.
func (c *Context) CheckFramebufferStatus(target int) int {
	return c.Call("checkFramebufferStatus", target).Int()
}

// Sets all pixels in a specific buffer to the same value.
func (c *Context) Clear(flags int) {
	c.Call("clear", flags)
}

// Specifies color values to use by the clear method to clear the color buffer.
func (c *Context) ClearColor(r, g, b, a float32) {
	c.Call("clearColor", r, g, b, a)
}

// Clears the depth buffer to a specific value.
func (c *Context) ClearDepth(depth float64) {
	c.Call("clearDepth", depth)
}

func (c *Context) ClearStencil(s int) {
	c.Call("clearStencil", s)
}

// Lets you set whether individual colors can be written when
// drawing or rendering to a framebuffer.
func (c *Context) ColorMask(r, g, b, a bool) {
	c.Call("colorMask", r, g, b, a)
}

// Compiles the GLSL shader source into binary data used by the WebGLProgram object.
func (c *Context) CompileShader(shader js.Value) {
	c.Call("compileShader", shader)
}

// Copies a rectangle of pixels from the current WebGLFramebuffer into a texture image.
func (c *Context) CopyTexImage2D(target, level, internal, x, y, w, h, border int) {
	c.Call("copyTexImage2D", target, level, internal, x, y, w, h, border)
}

// Replaces a portion of an existing 2D texture image with data from the current framebuffer.
func (c *Context) CopyTexSubImage2D(target, level, xoffset, yoffset, x, y, w, h int) {
	c.Call("copyTexSubImage2D", target, level, xoffset, yoffset, x, y, w, h)
}

// Creates and initializes a WebGLBuffer.
func (c *Context) CreateBuffer() js.Value {
	return c.Call("createBuffer")
}

// Returns a WebGLFramebuffer object.
func (c *Context) CreateFramebuffer() js.Value {
	return c.Call("createFramebuffer")
}

// Creates an empty WebGLProgram object to which vector and fragment
// WebGLShader objects can be bound.
func (c *Context) CreateProgram() js.Value {
	return c.Call("createProgram")
}

// Creates and returns a WebGLRenderbuffer object.
func (c *Context) CreateRenderbuffer() js.Value {
	return c.Call("createRenderbuffer")
}

// Returns an empty vertex or fragment shader object based on the type specified.
func (c *Context) CreateShader(typ int) js.Value {
	return c.Call("createShader", typ)
}

// Used to generate a WebGLTexture object to which images can be bound.
func (c *Context) CreateTexture() js.Value {
	return c.Call("createTexture")
}

// Sets whether or not front, back, or both facing facets are able to be culled.
func (c *Context) CullFace(mode int) {
	c.Call("cullFace", mode)
}

// Delete a specific buffer.
func (c *Context) DeleteBuffer(buffer js.Value) {
	c.Call("deleteBuffer", buffer)
}

// Deletes a specific WebGLFramebuffer object. If you delete the
// currently bound framebuffer, the default framebuffer will be bound.
// Deleting a framebuffer detaches all of its attachments.
func (c *Context) DeleteFramebuffer(framebuffer js.Value) {
	c.Call("deleteFramebuffer", framebuffer)
}

// Flags a specific WebGLProgram object for deletion if currently active.
// It will be deleted when it is no longer being used.
// Any shader objects associated with the program will be detached.
// They will be deleted if they were already flagged for deletion.
func (c *Context) DeleteProgram(program js.Value) {
	c.Call("deleteProgram", program)
}

// Deletes the specified renderbuffer object. If the renderbuffer is
// currently bound, it will become unbound. If the renderbuffer is
// attached to the currently bound framebuffer, it is detached.
func (c *Context) DeleteRenderbuffer(renderbuffer js.Value) {
	c.Call("deleteRenderbuffer", renderbuffer)
}

// Deletes a specific shader object.
func (c *Context) DeleteShader(shader js.Value) {
	c.Call("deleteShader", shader)
}

// Deletes a specific texture object.
func (c *Context) DeleteTexture(texture js.Value) {
	c.Call("deleteTexture", texture)
}

// Sets a function to use to compare incoming pixel depth to the
// current depth buffer value.
func (c *Context) DepthFunc(fun int) {
	c.Call("depthFunc", fun)
}

// Sets whether or not you can write to the depth buffer.
func (c *Context) DepthMask(flag bool) {
	c.Call("depthMask", flag)
}

// Sets the depth range for normalized coordinates to canvas or viewport depth coordinates.
func (c *Context) DepthRange(zNear, zFar float64) {
	c.Call("depthRange", zNear, zFar)
}

// Detach a shader object from a program object.
func (c *Context) DetachShader(program, shader js.Value) {
	c.Call("detachShader", program, shader)
}

// Turns off specific WebGL capabilities for this context.
func (c *Context) Disable(cap int) {
	c.Call("disable", cap)
}

// Turns off a vertex attribute array at a specific index position.
func (c *Context) DisableVertexAttribArray(index int) {
	c.Call("disableVertexAttribArray", index)
}

// Render geometric primitives from bound and enabled vertex data.
func (c *Context) DrawArrays(mode, first, count int) {
	c.Call("drawArrays", mode, first, count)
}

// Renders geometric primitives indexed by element array data.
func (c *Context) DrawElements(mode, count, typ, offset int) {
	c.Call("drawElements", mode, count, typ, offset)
}

// Turns on specific WebGL capabilities for this context.
func (c *Context) Enable(cap int) {
	c.Call("enable", cap)
}

// Turns on a vertex attribute at a specific index position in
// a vertex attribute array.
func (c *Context) EnableVertexAttribArray(index int) {
	c.Call("enableVertexAttribArray", index)
}

func (c *Context) Finish() {
	c.Call("finish")
}

func (c *Context) Flush() {
	c.Call("flush")
}

// Attaches a WebGLRenderbuffer object as a logical buffer to the
// currently bound WebGLFramebuffer object.
func (c *Context) FrameBufferRenderBuffer(target, attachment, renderbufferTarget int, renderbuffer js.Value) {
	c.Call("framebufferRenderBuffer", target, attachment, renderbufferTarget, renderbuffer)
}

// Attaches a texture to a WebGLFramebuffer object.
func (c *Context) FramebufferTexture2D(target, attachment, textarget int, texture js.Value, level int) {
	c.Call("framebufferTexture2D", target, attachment, textarget, texture, level)
}

// Sets whether or not polygons are considered front-facing based
// on their winding direction.
func (c *Context) FrontFace(mode int) {
	c.Call("frontFace", mode)
}

// Creates a set of textures for a WebGLTexture object with image
// dimensions from the original size of the image down to a 1x1 image.
func (c *Context) GenerateMipmap(target int) {
	c.Call("generateMipmap", target)
}

// Returns an WebGLActiveInfo object containing the size, type, and name
// of a vertex attribute at a specific index position in a program object.
func (c *Context) GetActiveAttrib(program js.Value, index int) js.Value {
	return c.Call("getActiveAttrib", program, index)
}

// Returns an WebGLActiveInfo object containing the size, type, and name
// of a uniform attribute at a specific index position in a program object.
func (c *Context) GetActiveUniform(program js.Value, index int) js.Value {
	return c.Call("getActiveUniform", program, index)
}

// Returns a slice of WebGLShaders bound to a WebGLProgram.
func (c *Context) GetAttachedShaders(program js.Value) []js.Value {
	objs := c.Call("getAttachedShaders", program)
	shaders := make([]js.Value, objs.Length())
	for i := 0; i < objs.Length(); i++ {
		shaders[i] = objs.Index(i)
	}
	return shaders
}

// Returns an index to the location in a program of a named attribute variable.
func (c *Context) GetAttribLocation(program js.Value, name string) int {
	return c.Call("getAttribLocation", program, name).Int()
}

// TODO: Create type specific variations.
// Returns the type of a parameter for a given buffer.
func (c *Context) GetBufferParameter(target, pname int) js.Value {
	return c.Call("getBufferParameter", target, pname)
}

// TODO: Create type specific variations.
// Returns the natural type value for a constant parameter.
func (c *Context) GetParameter(pname int) js.Value {
	return c.Call("getParameter", pname)
}

// Returns a value for the WebGL error flag and clears the flag.
func (c *Context) GetError() int {
	return c.Call("getError").Int()
}

// TODO: Create type specific variations.
// Enables a passed extension, otherwise returns null.
func (c *Context) GetExtension(name string) js.Value {
	return c.Call("getExtension", name)
}

// TODO: Create type specific variations.
// Gets a parameter value for a given target and attachment.
func (c *Context) GetFramebufferAttachmentParameter(target, attachment, pname int) js.Value {
	return c.Call("getFramebufferAttachmentParameter", target, attachment, pname)
}

// Returns the value of the program parameter that corresponds to a supplied pname
// which is interpreted as an int.
func (c *Context) GetProgramParameteri(program js.Value, pname int) int {
	return c.Call("getProgramParameter", program, pname).Int()
}

// Returns the value of the program parameter that corresponds to a supplied pname
// which is interpreted as a bool.
func (c *Context) GetProgramParameterb(program js.Value, pname int) bool {
	return c.Call("getProgramParameter", program, pname).Bool()
}

// Returns information about the last error that occurred during
// the failed linking or validation of a WebGL program object.
func (c *Context) GetProgramInfoLog(program js.Value) string {
	return c.Call("getProgramInfoLog", program).String()
}

// TODO: Create type specific variations.
// Returns a renderbuffer parameter from the currently bound WebGLRenderbuffer object.
func (c *Context) GetRenderbufferParameter(target, pname int) js.Value {
	return c.Call("getRenderbufferParameter", target, pname)
}

// TODO: Create type specific variations.
// Returns the value of the parameter associated with pname for a shader object.
func (c *Context) GetShaderParameter(shader js.Value, pname int) js.Value {
	return c.Call("getShaderParameter", shader, pname)
}

// Returns the value of the parameter associated with pname for a shader object.
func (c *Context) GetShaderParameterb(shader js.Value, pname int) bool {
	return c.Call("getShaderParameter", shader, pname).Bool()
}

// Returns errors which occur when compiling a shader.
func (c *Context) GetShaderInfoLog(shader js.Value) string {
	return c.Call("getShaderInfoLog", shader).String()
}

// Returns source code string associated with a shader object.
func (c *Context) GetShaderSource(shader js.Value) string {
	return c.Call("getShaderSource", shader).String()
}

// Returns a slice of supported extension strings.
func (c *Context) GetSupportedExtensions() []string {
	ext := c.Call("getSupportedExtensions")
	extensions := make([]string, ext.Length())
	for i := 0; i < ext.Length(); i++ {
		extensions[i] = ext.Index(i).String()
	}
	return extensions
}

// TODO: Create type specific variations.
// Returns the value for a parameter on an active texture unit.
func (c *Context) GetTexParameter(target, pname int) js.Value {
	return c.Call("getTexParameter", target, pname)
}

// TODO: Create type specific variations.
// Gets the uniform value for a specific location in a program.
func (c *Context) GetUniform(program, location js.Value) js.Value {
	return c.Call("getUniform", program, location)
}

// Returns a WebGLUniformLocation object for the location
// of a uniform variable within a WebGLProgram object.
func (c *Context) GetUniformLocation(program js.Value, name string) js.Value {
	return c.Call("getUniformLocation", program, name)
}

// TODO: Create type specific variations.
// Returns data for a particular characteristic of a vertex
// attribute at an index in a vertex attribute array.
func (c *Context) GetVertexAttrib(index, pname int) js.Value {
	return c.Call("getVertexAttrib", index, pname)
}

// Returns the address of a specified vertex attribute.
func (c *Context) GetVertexAttribOffset(index, pname int) int {
	return c.Call("getVertexAttribOffset", index, pname).Int()
}

// public function hint(target:GLenum, mode:GLenum) : Void;

// Returns true if buffer is valid, false otherwise.
func (c *Context) IsBuffer(buffer js.Value) bool {
	return c.Call("isBuffer", buffer).Bool()
}

// Returns whether the WebGL context has been lost.
func (c *Context) IsContextLost() bool {
	return c.Call("isContextLost").Bool()
}

// Returns true if buffer is valid, false otherwise.
func (c *Context) IsFramebuffer(framebuffer js.Value) bool {
	return c.Call("isFramebuffer", framebuffer).Bool()
}

// Returns true if program object is valid, false otherwise.
func (c *Context) IsProgram(program js.Value) bool {
	return c.Call("isProgram", program).Bool()
}

// Returns true if buffer is valid, false otherwise.
func (c *Context) IsRenderbuffer(renderbuffer js.Value) bool {
	return c.Call("isRenderbuffer", renderbuffer).Bool()
}

// Returns true if shader is valid, false otherwise.
func (c *Context) IsShader(shader js.Value) bool {
	return c.Call("isShader", shader).Bool()
}

// Returns true if texture is valid, false otherwise.
func (c *Context) IsTexture(texture js.Value) bool {
	return c.Call("isTexture", texture).Bool()
}

// Returns whether or not a WebGL capability is enabled for this context.
func (c *Context) IsEnabled(capability int) bool {
	return c.Call("isEnabled", capability).Bool()
}

// Sets the width of lines in WebGL.
func (c *Context) LineWidth(width float64) {
	c.Call("lineWidth", width)
}

// Links an attached vertex shader and an attached fragment shader
// to a program so it can be used by the graphics processing unit (GPU).
func (c *Context) LinkProgram(program js.Value) {
	c.Call("linkProgram", program)
}

// Sets pixel storage modes for readPixels and unpacking of textures
// with texImage2D and texSubImage2D.
func (c *Context) PixelStorei(pname, param int) {
	c.Call("pixelStorei", pname, param)
}

// Sets the implementation-specific units and scale factor
// used to calculate fragment depth values.
func (c *Context) PolygonOffset(factor, units float64) {
	c.Call("polygonOffset", factor, units)
}

// TODO: Figure out if pixels should be a slice.
// Reads pixel data into an ArrayBufferView object from a
// rectangular area in the color buffer of the active frame buffer.
func (c *Context) ReadPixels(x, y, width, height, format, typ int, pixels js.Value) {
	c.Call("readPixels", x, y, width, height, format, typ, pixels)
}

// Creates or replaces the data store for the currently bound WebGLRenderbuffer object.
func (c *Context) RenderbufferStorage(target, internalFormat, width, height int) {
	c.Call("renderbufferStorage", target, internalFormat, width, height)
}

//func (c *Context) SampleCoverage(value float64, invert bool) {
//	c.Call("sampleCoverage", value, invert)
//}

// Sets the dimensions of the scissor box.
func (c *Context) Scissor(x, y, width, height int) {
	c.Call("scissor", x, y, width, height)
}

// Sets and replaces shader source code in a shader object.
func (c *Context) ShaderSource(shader js.Value, source string) {
	c.Call("shaderSource", shader, source)
}

// public function stencilFunc(func:GLenum, ref:GLint, mask:GLuint) : Void;
// public function stencilFuncSeparate(face:GLenum, func:GLenum, ref:GLint, mask:GLuint) : Void;
// public function stencilMask(mask:GLuint) : Void;
// public function stencilMaskSeparate(face:GLenum, mask:GLuint) : Void;
// public function stencilOp(fail:GLenum, zfail:GLenum, zpass:GLenum) : Void;
// public function stencilOpSeparate(face:GLenum, fail:GLenum, zfail:GLenum, zpass:GLenum) : Void;

// Loads the supplied pixel data into a texture.
func (c *Context) TexImage2D(target, level, internalFormat, width, height, border, format, kind int, image interface{}) {
	c.Call("texImage2D", target, level, internalFormat, width, height, border, format, kind, image)
}

// Sets texture parameters for the current texture unit.
func (c *Context) TexParameteri(target int, pname int, param int) {
	c.Call("texParameteri", target, pname, param)
}

// Replaces a portion of an existing 2D texture image with all of another image.
func (c *Context) TexSubImage2D(target, level, xoffset, yoffset, format, typ int, image interface{}) {
	c.Call("texSubImage2D", target, level, xoffset, yoffset, format, typ, image)
}

// Assigns a floating point value to a uniform variable for the current program object.
func (c *Context) Uniform1f(location js.Value, x float32) {
	c.Call("uniform1f", location, x)
}

// Assigns a integer value to a uniform variable for the current program object.
func (c *Context) Uniform1i(location js.Value, x int) {
	c.Call("uniform1i", location, x)
}

// Assigns 2 floating point values to a uniform variable for the current program object.
func (c *Context) Uniform2f(location js.Value, x, y float32) {
	c.Call("uniform2f", location, x, y)
}

// Assigns 2 integer values to a uniform variable for the current program object.
func (c *Context) Uniform2i(location js.Value, x, y int) {
	c.Call("uniform2i", location, x, y)
}

// Assigns 3 floating point values to a uniform variable for the current program object.
func (c *Context) Uniform3f(location js.Value, x, y, z float32) {
	c.Call("uniform3f", location, x, y, z)
}

// Assigns 3 integer values to a uniform variable for the current program object.
func (c *Context) Uniform3i(location js.Value, x, y, z int) {
	c.Call("uniform3i", location, x, y, z)
}

// Assigns 4 floating point values to a uniform variable for the current program object.
func (c *Context) Uniform4f(location js.Value, x, y, z, w float32) {
	c.Call("uniform4f", location, x, y, z, w)
}

// Assigns 4 integer values to a uniform variable for the current program object.
func (c *Context) Uniform4i(location js.Value, x, y, z, w int) {
	c.Call("uniform4i", location, x, y, z, w)
}

// Assigns a floating point value to a uniform variable for the current program object.
func (c *Context) Uniform1fv(location js.Value, src js.TypedArray) {
	c.Call("uniform1fv", location, src)
}

// Assigns a integer value to a uniform variable for the current program object.
func (c *Context) Uniform1iv(location js.Value, src js.TypedArray) {
	c.Call("uniform1iv", location, src)
}

// Assigns 2 floating point values to a uniform variable for the current program object.
func (c *Context) Uniform2fv(location js.Value, src js.TypedArray) {
	c.Call("uniform2fv", location, src)
}

// Assigns 2 integer values to a uniform variable for the current program object.
func (c *Context) Uniform2iv(location js.Value, src js.TypedArray) {
	c.Call("uniform2iv", location, src)
}

// Assigns 3 floating point values to a uniform variable for the current program object.
func (c *Context) Uniform3fv(location js.Value, src js.TypedArray) {
	c.Call("uniform3fv", location, src)
}

// Assigns 3 integer values to a uniform variable for the current program object.
func (c *Context) Uniform3iv(location js.Value, src js.TypedArray) {
	c.Call("uniform3iv", location, src)
}

// Assigns 4 floating point values to a uniform variable for the current program object.
func (c *Context) Uniform4fv(location js.Value, src js.TypedArray) {
	c.Call("uniform4fv", location, src)
}

// Assigns 4 integer values to a uniform variable for the current program object.
func (c *Context) Uniform4iv(location js.Value, src js.TypedArray) {
	c.Call("uniform4iv", location, src)
}

// Sets values for a 2x2 floating point vector matrix into a
// uniform location as a matrix or a matrix array.
func (c *Context) UniformMatrix2fv(location js.Value, transpose bool, value js.TypedArray) {
	c.Call("uniformMatrix2fv", location, transpose, value)
}

// Sets values for a 3x3 floating point vector matrix into a
// uniform location as a matrix or a matrix array.
func (c *Context) UniformMatrix3fv(location js.Value, transpose bool, value js.TypedArray) {
	c.Call("uniformMatrix3fv", location, transpose, value)
}

// Sets values for a 4x4 floating point vector matrix into a
// uniform location as a matrix or a matrix array.
func (c *Context) UniformMatrix4fv(location js.Value, transpose bool, value js.TypedArray) {
	c.Call("uniformMatrix4fv", location, transpose, value)
}

// Set the program object to use for rendering.
func (c *Context) UseProgram(program js.Value) {
	c.Call("useProgram", program)
}

// Returns whether a given program can run in the current WebGL state.
func (c *Context) ValidateProgram(program js.Value) {
	c.Call("validateProgram", program)
}

func (c *Context) VertexAttribPointer(index, size, typ int, normal bool, stride int, offset int) {
	c.Call("vertexAttribPointer", index, size, typ, normal, stride, offset)
}

// public function vertexAttrib1f(indx:GLuint, x:GLfloat) : Void;
// public function vertexAttrib2f(indx:GLuint, x:GLfloat, y:GLfloat) : Void;
// public function vertexAttrib3f(indx:GLuint, x:GLfloat, y:GLfloat, z:GLfloat) : Void;
// public function vertexAttrib4f(indx:GLuint, x:GLfloat, y:GLfloat, z:GLfloat, w:GLfloat) : Void;
// public function vertexAttrib1fv(indx:GLuint, values:ArrayAccess<Float>) : Void;
// public function vertexAttrib2fv(indx:GLuint, values:ArrayAccess<Float>) : Void;
// public function vertexAttrib3fv(indx:GLuint, values:ArrayAccess<Float>) : Void;
// public function vertexAttrib4fv(indx:GLuint, values:ArrayAccess<Float>) : Void;

// Represents a rectangular viewable area that contains
// the rendering results of the drawing buffer.
func (c *Context) Viewport(x, y, width, height int) {
	c.Call("viewport", x, y, width, height)
}
