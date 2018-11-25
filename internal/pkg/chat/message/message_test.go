package message

import "testing"

var TestMessages = []Message{
	{
		"",
		"",
		"hello world!2",
		0,
	},
	{
		"kek",
		"lol",
		"sosi hui1",
		0,
	},
	{
		"kek",
		"lol",
		"sosi hui2",
		0,
	},
}
func TestMessageInsert(t *testing.T) {
	err := TestMessages[0].Save()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMessageGetOne(t *testing.T) {
	messages, err := GetGlobalMessages(0, 1444)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(messages)
	if messages[0].Text != TestMessages[0].Text {
		t.Fatal("not Same Message")
	}
}

func TestDialogMessageInsert1(t *testing.T) {
	err := TestMessages[1].Save()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDialogMessageCheck(t *testing.T) {
	messages, err := GetDialogMessages(0, 1, TestMessages[1].ToLogin, TestMessages[1].FromLogin)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(messages)
	if messages[0].Text != TestMessages[1].Text {
		t.Fatal("not Same Message")
	}
}

func TestDialogMessageInsert2(t *testing.T) {
	err := TestMessages[2].Save()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDialogReverseMessageCheck(t *testing.T) {
	messages, err := GetDialogMessages(0, 2, TestMessages[2].FromLogin, TestMessages[2].ToLogin)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(messages)
	if messages[0].Text != TestMessages[2].Text {
		t.Fatal("not Same Message")
	}
}