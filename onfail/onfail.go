package onfail

import (
	"errors"
	"log"
)

type Func func(string, error)

func Fatal(note string, err error) {
	log.Fatalln(convNoteErrToString(note, err))
}

func Log(note string, err error) {
	log.Println(convNoteErrToString(note, err))
}

func Panic(note string, err error) {
	panic(convNoteErrToError(note, err))
}

func convNoteErrToError(note string, err error) error {
	if len(note) < 1 {
		return err
	}
	return errors.New(note + ":\t" + err.Error())
}

func convNoteErrToString(note string, err error) string {
	if len(note) < 1 {
		return err.Error()
	}
	return note + ":\t" + err.Error()
}
