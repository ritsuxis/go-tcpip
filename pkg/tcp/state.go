package tcp

// https://upload.wikimedia.org/wikipedia/commons/f/f6/Tcp_state_diagram_fixed_new.svg
// に合うようなステートマシンを作成する
const (
	Close       State = "CLOSE"
	SynSent     State = "SYNSENT"
	Established State = "ESTABLISHED"
)

type State string

func NewTcpState() State {
	return Close
}

func (s State) Transition(flag ControlFlag) State {
	switch s {
	case Close:
		if flag.isSet(SYN) { // これは自分がSYNを送るので必要ないかも
			return SynSent
		}
	case SynSent:
		if flag.isSet(SYN | ACK) {
			return Established
		}
	default:
		return Close
	}
	return Close
}
