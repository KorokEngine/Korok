package ap

import (
	"golang.org/x/mobile/exp/audio/al"
	"log"
)

const MAX_STREAM_BUFFER = 24

type Channel struct {
	State SourceState
	sound *Sound
	loop bool

	al.Source
}

func (ch *Channel) Create(xType SourceType) error{
	array := al.GenSources(1)
	ch.Source = array[0]
	return nil
}

func (ch *Channel) Play(sound *Sound) {
	if sound.Type == Static {
		log.Println("play sound:", *sound)

		d := sound.Data.(*StaticData)
		ch.Source.QueueBuffers(d.Static.Buffer)
	} else if sound.Type == Stream {
		log.Println("Play stream")

		var usedBuffers = 0
		d := sound.Data.(*StreamData)
		for i := 0; i < MAX_STREAM_BUFFER; i++ {
			if num := ch.fill(d.Stream.Buffer[i], d.Decoder); num == 0 {
				break
			}
			usedBuffers ++
		}
		if usedBuffers > 0 {
			ch.Source.QueueBuffers(d.Stream.Buffer[:usedBuffers]...)
		}
	}
	ch.sound = sound
	al.PlaySources(ch.Source)
}

var count = 0
func (ch *Channel) UpdateState() {
	//queued := ch.Source.BuffersQueued()
	//processed := ch.Source.BuffersProcessed()
	//st := ch.Source.State()
	//
	//// print state!!
	//fmt.Println("queued:", queued, " processed:", processed, " stat:", st)

	// NOT USED!
	if ch.sound == nil {
		return
	}

	// UPDATE STATIC
	if ch.sound.Type == Static {
		st := ch.Source.State()
		switch st {
		case al.Stopped:
			ch.State = STOP
		case al.Playing:
			ch.State = PLAYING
		case al.Paused:
			ch.State = PAUSED
		}
	} else {
	// UPDATE STREAM
		data := ch.sound.Data.(*StreamData)
		d    := data.Decoder
		src  := ch.Source
		buff := [1]al.Buffer{}

		if bp := src.BuffersProcessed(); bp > 0 {
			count ++
			log.Printf("refill buffer cout: %d bp: %d", count, bp)
		}

		for bp := src.BuffersProcessed(); bp > 0; bp--{
			offsetSample := src.OffsetSample()
			offsetSecond := offsetSample/d.SampleRate()

			// get a used buffer
			src.UnqueueBuffers(buff[:]...)

			// why ? from love2d
			newOffsetSample := src.OffsetSample()
			newOffsetSecond := newOffsetSample/d.SampleRate()

			offsetSample += offsetSample - newOffsetSample
			offsetSecond += offsetSecond - newOffsetSecond

			//
			if ch.fill(buff[0], d) > 0 {
				src.QueueBuffers(buff[:]...)
			}
		}
	}
}

func (ch *Channel) Halt() {
	al.StopSources(ch.Source)
}

func (ch *Channel) Destroy() {
	al.DeleteSources(ch.Source)
}

func (ch *Channel) fill(buffer al.Buffer, d Decoder) (num int) {
	decoded := d.Decode()

	if decoded > 0 {
		fmt := getFormat(d.NumOfChan(), d.BitDepth())
		if fmt != FORMAT_END {
			buffer.BufferData(formatCodes[fmt], d.Buffer(), d.SampleRate())
		} else {
			decoded = 0
		}
	}

	if d.ReachEnd() && ch.loop {
		// var queued, processed int
		// try to remind decoder!!!
	}
	return decoded
}
