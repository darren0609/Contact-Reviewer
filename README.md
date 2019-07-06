# Golang Contacts App

[![godoc](https://godoc.org/github.com/darren0609/Contact-Reviewer?status.svg)](https://godoc.org/github.com/darren0609/Contact-Reviewer)
[![Build Status](https://travis-ci.org/darren0609/Contact-Reviewer.svg?branch=Refactoring)](https://travis-ci.org/darren0609/Contact-Reviewer)
[![Go Report Card](https://goreportcard.com/badge/github.com/darren0609/Contact-Reviewer)](https://goreportcard.com/report/github.com/darren0609/Contact-Reviewer)
[![codecov](https://codecov.io/gh/darren0609/Contact-Reviewer/branch/master/graph/badge.svg)](https://codecov.io/gh/darren0609/Contact-Reviewer)

Overall intent here is to build a solution to allow you to consolidate all of your contact information from various source infromation systems and allow them to populate and be edited. 

Initially focus will be based on the Microsoft Graph API pulling data directly from the Office 365 Contacts listing of a registered user. 

Ideal future functionality; 
* Consolidate all contacts to a source system of the users choice
* Remove duplicates across the different source systems
* Make decisions around duplicated information - where different email, phone, address details may exist
* Create in a way that can be called either via a UI, or via an API call

## main.go 

- You will need to create a folder called `init` where you will need to place the PRIVATE.TXT
- To create your App ID and Secret - start here  [Microsoft - Application Registration Portal](https://apps.dev.microsoft.com)
- Core of program, fundamentally creates authentication to MSGraph using details found in PRIVATE.TXT
- Once you have your PRIVATE.TXT file housed within your init folder, you should be able to run the GO MAIN.GO 

- It will start a server, and you will need to view at `http://localhost:8080/` 

