package main

import (
	"fmt"
	"log"
	"strings"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, strings.ToUpper(value))
	return nil
}

func cliRun(uniprotID string, pdbFlags arrayFlags) {
	msgs := make(chan string, 100)
	go func() {
		for m := range msgs {
			fmt.Println(m)
		}
	}()
	p, err := NewPipeline(uniprotID, pdbFlags, msgs)
	if err != nil {
		log.Fatal(err)
	}

	err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
	close(msgs)
}
