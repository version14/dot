package questions

type TextQuestion struct {
	Label       string
	Description string
	Placeholder string
	Value       string
	Next        *Next
}
