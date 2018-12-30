package strutils

type StringList struct {
	RightMargin int
	Items       []*string
}

func (sl StringList) Clear() {
	sl.Items = make([]*string, 0)
}

func (sl StringList) Count() int {

	return len(sl.Items)
}

func (sl StringList) Add(s string) int {

	return 0
}

func (sl StringList) SetText(text string) error {

	return nil
}

func (sl StringList) Text() string {

	return ""
}

func (sl StringList) IndexOf(s string) int {

	return -1
}

func (sl StringList) Delete(index int) error {

	return nil
}
