package spine

type SkeletonData struct {
	bones       []*BoneData
	slots       []*SlotData
	skins       []*Skin
	animations  []*Animation
	defaultSkin *Skin
}

func NewSkeletonData() *SkeletonData {
	data := new(SkeletonData)
	data.bones = make([]*BoneData, 0)
	data.slots = make([]*SlotData, 0)
	data.skins = make([]*Skin, 0)
	data.animations = make([]*Animation, 0)
	return data
}

func (s *SkeletonData) findBone(name string) (int, *BoneData) {
	for i, bone := range s.bones {
		if bone.name == name {
			return i, bone
		}
	}
	return -1, nil
}

func (s *SkeletonData) findSlot(name string) (int, *SlotData) {
	for i, slot := range s.slots {
		if slot.name == name {
			return i, slot
		}
	}
	return -1, nil
}

func (s *SkeletonData) findSkin(name string) (int, *Skin) {
	for i, skin := range s.skins {
		if skin.name == name {
			return i, skin
		}
	}
	return -1, nil
}

func (s *SkeletonData) findAnimation(name string) (int, *Animation) {
	for i, animation := range s.animations {
		if animation.name == name {
			return i, animation
		}
	}
	return -1, nil
}

type Skeleton struct {
	data         *SkeletonData
	Bones        []*Bone
	Slots        []*Slot
	DrawOrder    []*Slot
	skin         *Skin
	X, Y         float32
	r, g, b, a   float32
	time         float32
	FlipX, FlipY bool
	DebugBones   bool
	DebugSlots   bool
}

func NewSkeleton(skeletonData *SkeletonData) *Skeleton {
	skeleton := new(Skeleton)
	skeleton.data = skeletonData
	skeleton.r = 1
	skeleton.g = 1
	skeleton.b = 1
	skeleton.a = 1

	skeleton.Bones = make([]*Bone, 0)
	for _, boneData := range skeletonData.bones {
		var parent *Bone
		if boneData.parent != nil {
			i, _ := skeletonData.findBone(boneData.parent.name)
			parent = skeleton.Bones[i]
		}
		skeleton.Bones = append(skeleton.Bones, NewBone(boneData, parent))
	}

	skeleton.Slots = make([]*Slot, 0)
	skeleton.DrawOrder = make([]*Slot, 0)
	for _, slotData := range skeletonData.slots {
		i, _ := skeletonData.findBone(slotData.boneData.name)
		bone := skeleton.Bones[i]
		slot := NewSlot(slotData, skeleton, bone)
		skeleton.Slots = append(skeleton.Slots, slot)
		skeleton.DrawOrder = append(skeleton.DrawOrder, slot)
	}

	return skeleton
}

func (s *Skeleton) UpdateWorldTransform() {
	for _, bone := range s.Bones {
		bone.UpdateWorldTransform(s.FlipX, s.FlipY)
	}
}

func (s *Skeleton) SetToSetupPose() {
	s.setBonesToSetupPose()
	s.setSlotsToSetupPose()
}

func (s *Skeleton) setBonesToSetupPose() {
	for _, bone := range s.Bones {
		bone.SetToSetupPose()
	}
}

func (s *Skeleton) setSlotsToSetupPose() {
	for _, slot := range s.Slots {
		slot.SetToSetupPose()
	}
}

func (s *Skeleton) RootBone() *Bone {
	if len(s.Bones) != 0 {
		return s.Bones[0]
	}
	return nil
}

func (s *Skeleton) FindBone(name string) (int, *Bone) {
	for i, bone := range s.Bones {
		if bone.name == name {
			return i, bone
		}
	}
	return -1, nil
}

func (s *Skeleton) FindSlot(name string) (int, *Slot) {
	for i, slot := range s.Slots {
		if slot.data.name == name {
			return i, slot
		}
	}
	return -1, nil
}

func (s *Skeleton) SetSkinByName(name string) {
	_, skin := s.data.findSkin(name)
	if skin == nil {
		panic("Skin not found: " + name)
	}
	s.SetSkin(skin)
}

func (s *Skeleton) SetSkin(skin *Skin) {
	if s.skin != nil && skin != nil {
		skin.attachAll(s, s.skin)
	}
	s.skin = skin
}

func (s *Skeleton) AttachmentBySlotName(slot string, attachment string) Attachment {
	i, _ := s.data.findSlot(slot)
	return s.AttachmentBySlotIndex(i, attachment)
}

func (s *Skeleton) AttachmentBySlotIndex(index int, name string) Attachment {
	if s.skin != nil {
		attachment := s.skin.Attachment(index, name)
		if attachment != nil {
			return attachment
		}
	}
	if s.data.defaultSkin != nil {
		return s.data.defaultSkin.Attachment(index, name)
	}
	return nil
}

func (s *Skeleton) SetAttachment(slotName, attachmentName string) {
	for i, slot := range s.Slots {
		if slot.data.name == slotName {
			var attachment Attachment
			if attachmentName != "" {
				attachment = s.AttachmentBySlotIndex(i, attachmentName)
				if attachment == nil {
					panic("Attachment not found: " + attachmentName + ", for slot: " + slotName)
				}
			}
			slot.SetAttachment(attachment)
			return
		}
	}
	panic("Slot not found: " + slotName)
}

func (s *Skeleton) FindAnimation(name string) *Animation {
	_, a := s.data.findAnimation(name)
	return a
}

func (s *Skeleton) Update(dt float32) {
	s.time += dt
}
