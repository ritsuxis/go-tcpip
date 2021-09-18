package tcp

// https://upload.wikimedia.org/wikipedia/commons/f/f6/Tcp_state_diagram_fixed_new.svg
// に合うようなステートマシンを作成する
const (
	Close       State = "CLOSE"
	SynSent     State = "SYNSENT"
	Established State = "ESTABLISHED"
	FirstSent State = "FIRSTSENT"
	Sent        State = "SENT"
)

type State string

func NewTcpState() State {
	return Close
}

// このままだとクライアント側の遷移しかできない
func (s State) TransitionRcv(flag ControlFlag) State {
	switch s {
	case Close:
		// server側も実装するならsynを受け取ったときはsynreceivedにする必要あり
		return Close
	case SynSent:
		// serverからsyn+ackが返ってきたとき
		if flag.isSet(SYN) && flag.isSet(ACK) {
			return Established
		}
		// 強制終了をもらったとき
		if flag.isSet(RST) {
			return Close
		}
		return SynSent
	case Established:
		{
			if flag.isSet(ACK) {
				return FirstSent
			}
			// 強制終了をもらったとき
			if flag.isSet(RST) {
				return Close
			}
			return Close
		}
	case FirstSent:
		return FirstSent
	case Sent:
		{
			if flag.isSet(ACK) {
				return Sent
			}
			if flag.isSet(RST) {
				return Close
			}
			return Close
		}
	default:
		return Close
	}
}

func (s State) TransitionSnd(flag ControlFlag) State {
	switch s {
	case Close:
		// synを送ったとき
		if flag.isSet(SYN) {
			return SynSent
		}
		return Close
	case SynSent:
		return SynSent
	case Established:
		return Established
	case FirstSent:
		return Sent
	case Sent:
		{
			if flag.isSet(FIN) {
				return Close
			}
			return Sent
		}
	default:
		return Close
	}
}
