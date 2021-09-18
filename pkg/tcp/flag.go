package tcp

type ControlFlag uint8

func (f ControlFlag) String() string {
	var flags string

	flags += f.check(URG, "U")
	flags += f.check(ACK, "A")
	flags += f.check(PSH, "P")
	flags += f.check(RST, "R")
	flags += f.check(SYN, "S")
	flags += f.check(FIN, "F")

	return flags
}

func (f ControlFlag) check(bit ControlFlag, c string) string {
	if bit&f != 0 {
		return c
	} else {
		return "-"
	}
}

func (f ControlFlag) isSet(bit ControlFlag) bool {
	return bit&f != 0
}
