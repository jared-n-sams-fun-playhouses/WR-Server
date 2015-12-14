package main

import (
	//"http/template"
	"io/ioutil"
)

type Page struct {
	Location string
	Body []byte
}

func(p *Page) CreateView(filepath string) (error) {
	fileContents, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}

	p.Location = filepath
	p.Body  = fileContents

	return nil
}

func(p *Page) BuildPage() {

}

func(p *Page) ServePage() ([]byte) {
	return p.Body
}
