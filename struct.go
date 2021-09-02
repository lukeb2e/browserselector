package main

type configuration struct {
	Domain  []domain
	Browser map[string]*browser
	Debug   bool
}

type domain struct {
	Browser  string
	Regex    string
	Priority uint
}

type browser struct {
	Exec   string
	Script string
}
