package eventbus

type State int8

const (
	START State = iota
	RUNNING
	DONE
	OVERMAXTIME
)

func (ps State) IsStart() bool {
	return ps == START
}

func (ps State) IsStilRunning() bool {
	return ps == RUNNING
}

func (ps State) IsDone() bool {
	return ps == DONE
}

func (ps State) IsOverMaxTime() bool {
	return ps == OVERMAXTIME
}
