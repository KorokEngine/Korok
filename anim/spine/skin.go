package spine

import (
	"fmt"
)

type skinData struct {
	Index      int
	Name       string
	Attachment Attachment
}

type Skin struct {
	name        string
	attachments map[string]*skinData
}

func NewSkin(name string) *Skin {
	skin := new(Skin)
	skin.name = name
	skin.attachments = make(map[string]*skinData)
	return skin
}

func (s *Skin) AddAttachment(slotIndex int, name string, attachment Attachment) {
	data := &skinData{slotIndex, name, attachment}
	s.attachments[fmt.Sprintf("%v:%v", slotIndex, name)] = data
}

func (s *Skin) Attachment(slotIndex int, name string) Attachment {
	values, ok := s.attachments[fmt.Sprintf("%v:%v", slotIndex, name)]
	if !ok {
		return nil
	}
	return values.Attachment
}

func (s *Skin) attachAll(skeleton *Skeleton, oldSkin *Skin) {
	for _, val := range oldSkin.attachments {
		slot := skeleton.Slots[val.Index]
		attachment := s.Attachment(val.Index, val.Name)
		if attachment != nil {
			slot.SetAttachment(attachment)
		}
	}
}
