package main

type properties map[string]string

type itemData struct {
	ConstName string
	Value     string
}

type keyFileData struct {
	Items []itemData
}

type dictFileData struct {
	FuncName string
	LangTag  string
	Keys     []itemData
}
