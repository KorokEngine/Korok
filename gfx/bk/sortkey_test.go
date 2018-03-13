package bk

import (
	"testing"
)

func TestEncoder(t *testing.T) {
	sk := SortKey{
		Layer:1,
		Order:2,
		Shader:3,
		Blend:4,
		Texture:5,
	}

	key := sk.Encode()

	if key != 0x10087005 {
		t.Error("encoder err")
	}

	sk1 := &SortKey{}; sk1.Decode(key)


	if sk.Layer != sk1.Layer {
		t.Error("layer decoder error")
	}
	if sk.Order != sk1.Order {
		t.Error("order decoder error")
	}
	if sk.Shader != sk1.Shader {
		t.Error("shader decoder error")
	}
	if sk.Blend != sk1.Blend {
		t.Error("blend decoder error")
	}
	if sk.Texture != sk1.Texture {
		t.Error("texture decdoer error")
	}

}
