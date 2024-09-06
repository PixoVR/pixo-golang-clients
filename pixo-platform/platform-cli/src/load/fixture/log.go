package fixture

func (t *Tester) LogSuccess(id int, msg string) {
	t.Printf(":check_mark_button: %d: %s\n", id, SuccessColor.Sprint(msg))
}

func (t *Tester) LogError(id int, msg string, err error) {
	t.Printf(":exclamation: %d: %s - %s\n", id, msg, ErrorColor.Sprint(err))
}
