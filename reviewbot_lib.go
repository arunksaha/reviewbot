package main

import (
	"strconv"
	"strings"
)

type UserId int64

// State of a single user.
type UserInfo struct {
	firstName string

	// User is in the middle of reviewing.
	reviewInProgress bool

	// All reviews of this user.
	reviews []string

	// Content of the current ongoing review.
	messages []string
}

// State of the entire chatbot.
type ReviewBot struct {
	userMap map[UserId]*UserInfo

	commandMap map[string]string

	menuText string
}

func NewReviewBot() *ReviewBot {
	var rb ReviewBot
	rb.userMap = make(map[UserId]*UserInfo)
	rb.commandMap = map[string]string{
		"/menu":       "Show the list of commands",
		"/postreview": "Post a review",
		"/myreviews":  "Retrieve all my reviews",
	}
	rb.menuText = composeMenu(rb.commandMap)
	return &rb
}

func composeMenu(cmdMap map[string]string) string {
	menuText := "You can control me by sending these commands:\n\n"
	for cmd, desc := range cmdMap {
		menuText += cmd + " - " + desc + "\n"
	}
	return menuText
}

func (rb ReviewBot) Menu() string {
	return rb.menuText
}

// Find user by 'userId' if existing, else add one.
func (rb *ReviewBot) UserFindOrAdd(userId UserId) *UserInfo {
	_, ok := rb.userMap[userId]
	if !ok {
		rb.userMap[userId] = &UserInfo{}
	}
	userInfo := rb.userMap[userId]
	return userInfo
}

// If 'text' starts with a comamnd, return (<command>, true, true);
// else return ("", <start-with-slash>, false).
func (rb ReviewBot) IsCommand(text string) (string, bool, bool) {
	words := strings.Split(text, " ")

	// The first word might be a command.
	cmd := words[0]
	startWithSlash := text[0] == '/'

	_, knownCmd := rb.commandMap[cmd]
	validCmd := knownCmd || (cmd == "/start")

	return cmd, startWithSlash, validCmd
}

func (rb *ReviewBot) HandleText(text string, userId UserId, firstName string) string {
	cmd, startWithSlash, validCmd := (*rb).IsCommand(text)

	var resp string

	if startWithSlash && !validCmd {
		resp = cmd + " is not a valid command."
		return resp
	}

	userInfo := rb.UserFindOrAdd(userId)
	userInfo.PutFirstName(firstName)

	if validCmd {
		resp = userInfo.HandleCommand(*rb, cmd)
	} else {
		userInfo.HandleMessage(text)
	}

	return resp
}

func (userInfo *UserInfo) HandleCommand(rb ReviewBot, cmd string) string {
	var resp string

	// If we were already in a review, save that first
	// before processing a new command.
	userInfo.SaveReview()

	if cmd == "/start" {
		resp = "Welcome " + userInfo.firstName + "!\n"
		resp += rb.Menu()
	} else if cmd == "/menu" {
		resp = rb.Menu()
	} else if cmd == "/postreview" {
		userInfo.StartReview()
	} else if cmd == "/myreviews" {
		resp = userInfo.GetAllReviews()
	}

	return resp
}

func (userInfo *UserInfo) StartReview() {
	userInfo.reviewInProgress = true
}

func (userInfo *UserInfo) StopReview() {
	userInfo.reviewInProgress = false
}

func (userInfo UserInfo) ReviewInProgress() bool {
	return userInfo.reviewInProgress
}

// If user has an ongoing review, then save it.
func (userInfo *UserInfo) SaveReview() {
	if !userInfo.reviewInProgress {
		return
	}

	var review string
	for _, mesg := range userInfo.messages {
		review += mesg + " "
	}

	userInfo.reviews = append(userInfo.reviews, review)

	userInfo.messages = make([]string, 0)

	userInfo.StopReview()
}

func (userInfo UserInfo) GetAllReviews() string {
	var reviews string

	for idx, review := range userInfo.reviews {
		reviews += "Review " + strconv.Itoa(idx+1) + ": \n"
		reviews += review
		reviews += "\n\n"
	}

	return reviews
}

func (userInfo *UserInfo) HandleMessage(mesg string) {
	if userInfo.reviewInProgress {
		userInfo.messages = append(userInfo.messages, mesg)
	}
}

func (userInfo *UserInfo) PutFirstName(firstName string) {
	userInfo.firstName = firstName
}
