package spine

type SlotData struct {
	name           string
	boneData       *BoneData
	r, g, b, a     float32
	attachmentName string
}

func NewSlotData(name string, boneData *BoneData) *SlotData {
	slotData := new(SlotData)
	slotData.name = name
	slotData.boneData = boneData
	slotData.r = 1
	slotData.g = 1
	slotData.b = 1
	slotData.a = 1
	return slotData
}

type Slot struct {
	data           *SlotData
	skeleton       *Skeleton
	Bone           *Bone
	R, G, B, A     float32
	attachmentTime float32
	Attachment     Attachment
}

func NewSlot(slotData *SlotData, skeleton *Skeleton, bone *Bone) *Slot {
	slot := new(Slot)
	slot.data = slotData
	slot.skeleton = skeleton
	slot.Bone = bone
	slot.R = 1
	slot.G = 1
	slot.B = 1
	slot.A = 1
	slot.SetToSetupPose()
	return slot
}

func (s *Slot) SetToSetupPose() {
	data := s.data
	s.R = data.r
	s.G = data.g
	s.B = data.b
	s.A = data.a

	for i, slotData := range s.skeleton.data.slots {
		if slotData == data {
			s.SetAttachment(s.skeleton.AttachmentBySlotIndex(i, data.attachmentName))
			return
		}
	}
}

func (s *Slot) SetAttachment(attachment Attachment) {
	s.Attachment = attachment
	s.attachmentTime = s.skeleton.time
}

func (s *Slot) SetAttachmentTime(time float32) {
	s.attachmentTime = s.skeleton.time - time
}

func (s *Slot) AttachmentTime() float32 {
	return s.skeleton.time - s.attachmentTime
}

func (s *Slot) Skeleton() *Skeleton {
	return s.skeleton
}
