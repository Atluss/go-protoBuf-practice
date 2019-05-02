package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	pblist "github.com/Atluss/protoBufPractice/proto"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"strings"
)

const (
	sizeOfLength = 8
	dbPath       = "mydb.pb"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: list or add")
		os.Exit(1)
	}

	var err error
	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list()
	case "add":
		err = add(strings.Join(flag.Args()[1:], ""))
	default:
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func add(text string) error {

	task := &pblist.Task{
		Text: text,
		Done: false,
	}

	var b, err = proto.Marshal(task)
	if err != nil {
		return fmt.Errorf("could not encode task: %v", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", dbPath, err)
	}

	if err := binary.Write(f, binary.LittleEndian, int64(len(b))); err != nil {
		return fmt.Errorf("could not encode length of message: %v", err)
	}
	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could not write task to file: %v", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close file %s: %v", dbPath, err)
	}
	return nil
}

func list() error {
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("could not read %s: %v", dbPath, err)
	}

	for {
		if len(b) == 0 {
			return nil
		} else if len(b) < sizeOfLength {
			return fmt.Errorf("remaining odd %d bytes, what to do", len(b))
		}

		var l int64
		if err := binary.Read(bytes.NewReader(b[:sizeOfLength]), binary.LittleEndian, &l); err != nil {
			return fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[sizeOfLength:]

		var task pblist.Task
		if err := proto.Unmarshal(b[:l], &task); err != nil {
			return fmt.Errorf("could not read task: %v", err)
		}
		b = b[l:]

		if task.Done {
			fmt.Printf("👍")
		} else {
			fmt.Printf("😱")
		}
		fmt.Printf(" %s\n", task.Text)
	}
}
