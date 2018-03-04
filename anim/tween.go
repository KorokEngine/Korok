package anim

import "korok.io/korok/anim/tween"

type InterpolationType uint8

const (
	Linear InterpolationType = iota
)

func OfFloat(start, end float32) tween.Animator {
	return tween.Animator{}
}

func Default() *tween.Engine {
	return en
}

var en *tween.Engine

func init() {
	en = tween.NewEngine()
}


