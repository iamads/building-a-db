This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.


# Essential commands

build: go build -v
test: go test -v

# Interaction Style
Be direct and honest.
Skip unnecessary acknowledgments.
Correct me when I'm wrong and explain why.
Suggest better alternatives if my ideas can be improved.
Focus on accuracy and efficiency.
Challenge my assumptions when needed.
Prioritize quality information and directness.


# Test Case Flow

- Analyse the user request and identify codde portions you have to create tests for.
- If you identify any bugs, Report to user.
- Document tests plan: cases that you will be testing for
  - In test function you shhould have sections for: Happy path scenario, basic use cases , edge cases
  - If you find a bug, create a failing case for it
- You are only allowed to make changes to test files i.e. *_test.go
- Once plan is complete share it to user and wait for approval
- Once Approval is granted, create a todo and start working on the test cases