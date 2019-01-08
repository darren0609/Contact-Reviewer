# Contact Reviewer (previously: msgraph-go-contacts)

[![godoc](https://godoc.org/github.com/darren0609/msgraph-go-contacts?status.svg)](https://godoc.org/github.com/darren0609/msgraph-go-contacts)
[![Build Status](https://travis-ci.org/darren0609/msgraph-go-contacts.svg?branch=master)](https://travis-ci.org/darren0609/msgraph-go-contacts)
[![Go Report Card](https://goreportcard.com/badge/github.com/darren0609/msgraph-go-contacts)](https://goreportcard.com/report/github.com/darren0609/msgraph-go-contacts)
[![Github Pages](https://img.shields.io/badge/Github%20Pages-URL-blue.svg)](https://darren0609.github.io/Contact-Reviewier/)

Overall intent here is to build a solution to allow you to consolidate all of your contact information from various source infromation systems and allow them to populate and be edited. 

Ideal future functionality; 
* Consolidate all contacts to a source system
* Remove duplicates
* Make decisions around duplicated information - where different email, phone, address details may exist

## connect.go 

- Core of program, fundamentally creates authentication to MSGraph using details found in PRIVATE.TXT


## quickstart.go 

- Taken from https://developers.google.com/people/quickstart/go 
- Code to pull Google Contacts from Google Api

## tpl 

- Folder holds templates to process HTML

