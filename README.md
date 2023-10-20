# Problem Statement

Develop a rudimentary chatbot based on Telegram where users can submit product reviews and view them later.

# Functional Specification

  - The chatbot registers to Telegram with its token.
  - It can receive messages sent to it via Telegram.
  - It can send messages to the user via Telegram.
  - The chatbot supports the following commands:
    - /menu       - Show the list of commands
    - /postreview - Post a review
    - /myreviews  - Retrieve all reviews of the user
  - The user chooses the command option from the provided menu in the usual Telegram way, e.g., "/<command>".
  - The user can send one or more text messages that is processed in context of the last command.
  - For example, user can choose "/postreview" and start posting a review.
  - The content of the review can be sent over multiple (successive) messages. 
  - At any time, the user can choose a differet command and switch the context. 
  - "/myreviews" shows all reviews submitted by the user in the chronological order.

# Design

  - The chatbot is designed as a daemon process.
  - The chatbot application is split into the following parts:
    - reviewbot_main.go: 
      - As the name suggests, this contains the "main" function.
      - Reads the token from the environment.
      - Connection and registration with Telegram.
      - Receiving messages from and sending messages to Telegram.
    - reviewbot_lib.go
      - This is the reviewbot library with two main structures:
        - ReviewBot - the commands and the entire user information collection
        - UserInfo - user specific information

# Build

## Local
```
  $ make build
  go build -o reviewbot
```
## Docker

```
  $ make docker-build
  docker build -t reviewbot .
  <long output trimmed>
```

# Execution

In the following, <token> refers to the Telegram token for the reviewbot account.

## Local

```
  $ TOKEN=<token> ./reviewbot
```

## Docker

```
  $ docker run -e TOKEN='<token>' reviewbot
```

### DockerHub

A version of the image is available at dockerhub as per the following.

  arunksaha/reviewbot

  https://hub.docker.com/repository/docker/arunksaha/reviewbot/general

# Test

There are a bunch of unit tests at reviewbot_lib_test.go.

```
  $ make test # go test -v
  === RUN   TestUserSaveReview1
  --- PASS: TestUserSaveReview1 (0.00s)
  === RUN   TestUserSaveReview2
  --- PASS: TestUserSaveReview2 (0.00s)
  === RUN   TestUserHandleCommandStart
  --- PASS: TestUserHandleCommandStart (0.00s)
  === RUN   TestUserHandleCommandMenu
  --- PASS: TestUserHandleCommandMenu (0.00s)
  === RUN   TestUserHandleCommandPostReview
  --- PASS: TestUserHandleCommandPostReview (0.00s)
  === RUN   TestUserHandleCommandMyReviews
  --- PASS: TestUserHandleCommandMyReviews (0.00s)
  === RUN   TestRbUserFindOrAdd
  --- PASS: TestRbUserFindOrAdd (0.00s)
  === RUN   TestRbIsCommand
  --- PASS: TestRbIsCommand (0.00s)
  === RUN   TestRbHandleText
  --- PASS: TestRbHandleText (0.00s)
  PASS
```

## Test Coverage

Collect the test coverage.

```
  $ go test -v -coverprofile=coverage.out
```

View the test coverage

```
  $ go tool cover -func=coverage.out 
  reviewbot/reviewbot_lib.go:33:	NewReviewBot		100.0%
  reviewbot/reviewbot_lib.go:45:	composeMenu		100.0%
  reviewbot/reviewbot_lib.go:53:	Menu			100.0%
  reviewbot/reviewbot_lib.go:58:	UserFindOrAdd		100.0%
  reviewbot/reviewbot_lib.go:69:	IsCommand		100.0%
  reviewbot/reviewbot_lib.go:82:	HandleText		100.0%
  reviewbot/reviewbot_lib.go:104:	HandleCommand		100.0%
  reviewbot/reviewbot_lib.go:125:	StartReview		100.0%
  reviewbot/reviewbot_lib.go:129:	StopReview		100.0%
  reviewbot/reviewbot_lib.go:133:	ReviewInProgress	100.0%
  reviewbot/reviewbot_lib.go:138:	SaveReview		100.0%
  reviewbot/reviewbot_lib.go:155:	GetAllReviews		100.0%
  reviewbot/reviewbot_lib.go:167:	HandleMessage		100.0%
  reviewbot/reviewbot_lib.go:173:	PutFirstName		100.0%
  reviewbot/reviewbot_main.go:15:	main			0.0%
  reviewbot/reviewbot_main.go:37:	getUpdatesChan		0.0%
  reviewbot/reviewbot_main.go:56:	receiveUpdates		0.0%
  reviewbot/reviewbot_main.go:64:	handleUpdate		0.0%
  reviewbot/reviewbot_main.go:70:	recvMessage		0.0%
  reviewbot/reviewbot_main.go:87:	sendMessage		0.0%
  total:							(statements)		61.0%
```

The test coverage can also viewed in a browser as the following.

```
  $ go tool cover -html=coverage.out
```

The collection and command-line viewing can be done together as the following.

```
  $ make test-coverage
```

# Static Analysis

```
  $ make check
```

# Future Work
  
  - Make the chatbot restartable.
    - Save the state persistently.
    - Load the previously saved state (if any) at the startup.

  - Add a mock Telegram server

  - Add tests for the functions in reviewbot_main.go, i.e., the functionality of receiving messages from and sending messages to Telegram.

  - More features of reviewing, e.g.
    - Choosing a product to review.
    - Numerical star rating, e.g., the conventional rating from 1.0 to 5.0.
    - Aggregating reviews by product.
      - Average rating of a product.
      - Bucketing reviews by rating.
        - Percentage (histogram-ish) ratings of a product.

