package tcp

type OptionKind uint8

type Option interface {
	Kind() OptionKind
	Length() int // byte Kindのみのoptionは1
	Data() []byte
	Format() []byte // そのOptionのformatに沿って[]byteで生成
}

type Options []Option

// TODO: Options check from bytes

type EndOptionList struct{} // dataがないので空の構造体とする

func (eol EndOptionList) Kind() OptionKind {
	return EOL
}

func (eol EndOptionList) Length() int {
	return 1
}

func (eol EndOptionList) Data() []byte {
	return nil
}

func (eol EndOptionList) Format() []byte {
	return []byte{byte(EOL)}
}

type NoOperation struct{} // dataがないので空の構造体とする

func (nop NoOperation) Kind() OptionKind {
	return NOP
}

func (nop NoOperation) Length() int {
	return 1
}

func (nop NoOperation) Data() []byte {
	return nil
}

func (nop NoOperation) Format() []byte {
	return []byte{byte(NOP)}
}

type MaximumSegmentSize uint16 // 2byte

func (mss MaximumSegmentSize) Kind() OptionKind {
	return MSS
}

func (mss MaximumSegmentSize) Length() int {
	return 4
}

func (mss MaximumSegmentSize) Data() []byte {
	return []byte{byte(mss >> 8), byte(mss & 0xff)} // []byteなので8bitごとに区切る
}

func (mss MaximumSegmentSize) Format() []byte {
	tmp := []byte{byte(MSS), byte(mss.Length())}
	return append(tmp, mss.Data()...)
}

type WindowScale uint8

func (ws WindowScale) Kind() OptionKind {
	return WS
}

func (ws WindowScale) Length() int {
	return 3
}

func (ws WindowScale) Data() []byte {
	return []byte{byte(ws)}
}

func (ws WindowScale) Format() []byte {
	return []byte{byte(WS), byte(ws.Length()), byte(ws)}
}

type SACKPermitted struct{} // dataがないので空の構造体とする

func (sp SACKPermitted) Kind() OptionKind {
	return SP
}

func (sp SACKPermitted) Length() int {
	return 2
}

func (sp SACKPermitted) Data() []byte {
	return nil
}

func (sp SACKPermitted) Format() []byte {
	return []byte{byte(SP), byte(sp.Length())}
}

type SelectiveACK []byte // 8byte * N

func (sack SelectiveACK) Kind() OptionKind {
	return SACK
}

func (sack SelectiveACK) Length() int {
	return len(sack) + 2 // 2はkindとlength部分の分
}

func (sack SelectiveACK) Data() []byte {
	return sack
}

func (sack SelectiveACK) Format() []byte {
	return append([]byte{byte(SACK), byte(sack.Length())}, sack.Data()...)
}

type TimeStamp []byte // 8byte

func (ts TimeStamp) Kind() OptionKind {
	return TS
}

func (ts TimeStamp) Length() int {
	return 10
}

func (ts TimeStamp) Data() []byte {
	return ts
}

func (ts TimeStamp) Format() []byte {
	return append([]byte{byte(8), byte(10)}, ts.Data()...)
}
