# msgraph-go-contacts

[![godoc](https://godoc.org/github.com/darren0609/msgraph-go-contacts?status.svg)](https://godoc.org/github.com/darren0609/msgraph-go-contacts)
[![Build Status](https://travis-ci.org/darren0609/msgraph-go-contacts.svg?branch=master)](https://travis-ci.org/darren0609/msgraph-go-contacts)
[![Go Report Card](https://goreportcard.com/badge/github.com/darren0609/msgraph-go-contacts)](https://goreportcard.com/report/github.com/darren0609/msgraph-go-contacts)

Overall intent here is to build a solution to allow you to consolidate all of your contact information from various source infromation systems and allow them to populate and be edited. 

Ideal future functionality; 
* Consolidate all contacts to a source system of the users choice
* Remove duplicates across the different source systems
* Make decisions around duplicated information - where different email, phone, address details may exist
* Create in a way that can be called either via a UI, or via an API call

## main.go 

- You will need to create a folder called `init` where you will need to place the PRIVATE.TXT
- To create your App ID and Secret - start here ![](https://assets.onestore.ms/cdnfiles/external/uhf/long/9a49a7e9d8e881327e81b9eb43dabc01de70a9bb/images/microsoft-white.png) [Microsoft - Application Registration Portal](https://apps.dev.microsoft.com)
- Core of program, fundamentally creates authentication to MSGraph using details found in PRIVATE.TXT


## quickstart.go 

- Taken from https://developers.google.com/people/quickstart/go 
- Code to pull Google Contacts from Google Api

## tpl 

- Folder holds templates to process HTML

