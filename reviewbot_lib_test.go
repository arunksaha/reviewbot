package main

import (
	"strings"
	"testing"
)

func TestUserSaveReview1(t *testing.T) {
	user := UserInfo{}

	mesgs := []string{"I bought a book", "It was long", "I like it"}
	for _, mesg := range mesgs {
		user.HandleMessage(mesg)
	}

	user.SaveReview()

	nmessages := len(user.messages)
	if nmessages != 0 {
		t.Errorf("Incorrect: expected = %d, observed = %d\n", 0, nmessages)
	}

	if user.ReviewInProgress() {
		t.Errorf("Incorrect")
	}
}

func TestUserSaveReview2(t *testing.T) {
	user := UserInfo{}

	user.StartReview()

	if !user.ReviewInProgress() {
		t.Errorf("Incorrect")
	}

	mesgs := []string{"I bought a book", "It was long", "I like it"}
	for _, mesg := range mesgs {
		user.HandleMessage(mesg)
	}

	nmessages := len(user.messages)
	if nmessages != len(mesgs) {
		t.Errorf("Incorrect: expected = %d, observed = %d\n", len(mesgs), nmessages)
	}

	for idx := 0; idx < nmessages; idx++ {
		if user.messages[idx] != mesgs[idx] {
			t.Errorf("Incorrect: expected = %s, observed = %s\n", mesgs[idx], user.messages[idx])
		}
	}

	if len(user.reviews) != 0 {
		t.Errorf("Incorrect: expected = %d, observed = %d\n", 0, len(user.reviews))
	}

	user.SaveReview()

	if user.ReviewInProgress() {
		t.Errorf("Incorrect")
	}

	if len(user.reviews) != 1 {
		t.Errorf("Incorrect: expected = %d, observed = %d\n", 1, len(user.reviews))
	}

	allReviews := user.GetAllReviews()
	for _, mesg := range mesgs {
		contained := strings.Contains(allReviews, mesg)
		if !contained {
			t.Errorf("Incorrect: expected = %s, observed = %s\n", mesg, allReviews)
		}
	}
}

func TestUserHandleCommandStart(t *testing.T) {
	revbot := NewReviewBot()
	menu := revbot.Menu()

	userName := "Alan"
	user := UserInfo{}
	user.PutFirstName(userName)

	resp := user.HandleCommand(*revbot, "/start")
	if !strings.Contains(resp, userName) {
		t.Errorf("Incorrect: username %s missing in resp = %s\n", userName, resp)
	}
	if !strings.Contains(resp, menu) {
		t.Errorf("Incorrect: menu missing in resp = %s\n", resp)
	}
}

func TestUserHandleCommandMenu(t *testing.T) {
	revbot := NewReviewBot()
	menu := revbot.Menu()

	userName := "Alan"
	user := UserInfo{}
	user.PutFirstName(userName)

	resp := user.HandleCommand(*revbot, "/menu")
	if !strings.Contains(resp, menu) {
		t.Errorf("Incorrect: menu missing in resp = %s\n", resp)
	}
}

func TestUserHandleCommandPostReview(t *testing.T) {
	revbot := NewReviewBot()

	user := UserInfo{}

	user.HandleCommand(*revbot, "/postreview")
	if !user.ReviewInProgress() {
		t.Errorf("Incorrect")
	}

	mesgs := []string{"I bought a book", "It was long", "I like it"}
	for idx, mesg := range mesgs {
		user.HandleMessage(mesg)
		if len(user.messages) != (idx + 1) {
			t.Errorf("Incorrect: expected = %d, observed = %d\n", idx+1, len(user.messages))
		}
	}
}

func TestUserHandleCommandMyReviews(t *testing.T) {
	revbot := NewReviewBot()

	user := UserInfo{}

	user.HandleCommand(*revbot, "/postreview")

	mesgs := []string{"I bought a book", "It was long", "I like it"}
	for _, mesg := range mesgs {
		user.HandleMessage(mesg)
	}

	resp := user.HandleCommand(*revbot, "/myreviews")
	if user.ReviewInProgress() {
		t.Errorf("Incorrect")
	}

	for _, mesg := range mesgs {
		if !strings.Contains(resp, mesg) {
			t.Errorf("Incorrect: mesg %s expected in response %s\n", mesg, resp)
		}
	}
}

func TestRbUserFindOrAdd(t *testing.T) {
	revbot := NewReviewBot()
	userId := UserId(123)

	_, ok1 := revbot.userMap[userId]
	if ok1 {
		t.Errorf("Incorrect")
	}

	revbot.UserFindOrAdd(userId)

	_, ok2 := revbot.userMap[userId]
	if !ok2 {
		t.Errorf("Incorrect")
	}
}

func TestRbIsCommand(t *testing.T) {
	revbot := NewReviewBot()

	var data = []struct {
		text           string
		cmd            string
		startWithSlash bool
		validCmd       bool
	}{
		{"/start extra", "/start", true, true},
		{"/menu extra ", "/menu", true, true},
		{"/postreview abra ca dabra ", "/postreview", true, true},
		{"/myreviews abra ca dabra ", "/myreviews", true, true},
		{"/foo", "/foo", true, false},
		{"hello", "hello", false, false},
	}

	for _, datum := range data {
		cmd, sws, valCmd := (*revbot).IsCommand(datum.text)
		if cmd != datum.cmd || sws != datum.startWithSlash || valCmd != datum.validCmd {
			t.Errorf("Incorrect: Failed for input = %v\n", datum)
		}
	}
}

func TestRbHandleText(t *testing.T) {
	revbot := NewReviewBot()

	var data = []struct {
		text      string
		userId    int64
		firstName string
		resp      string
	}{
		{"/start extra", 123, "Alan", "Welcome Alan"},
		{"/menu extra ", 123, "Alan", "You can control me by sending these commands:"},
		{"/postreview abra ca dabra ", 123, "Alan", ""},
		{"/myreviews abra ca dabra ", 123, "Alan", ""},
		{"/foo", 123, "Alan", "is not a valid command"},
		{"hello", 123, "Alan", ""},
	}

	for _, datum := range data {
		resp := revbot.HandleText(datum.text, UserId(datum.userId), datum.firstName)
		if !strings.Contains(resp, datum.resp) {
			t.Errorf("Incorrect: expected = %s, observed = %s\n", datum.resp, resp)
		}
	}
}
