package notifier

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
)

type Input struct {
	file *os.File
}

func NewInput(filename string) (*Input, error) {
	f, err := inputFile(filename)
	if err != nil {
		return nil, err
	}

	return &Input{
		file: f,
	}, nil
}

func (i *Input) Bulk() (string, error) {
	b, err := ioutil.ReadAll(i.file)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (i *Input) Line() <-chan string {
	br := bufio.NewScanner(i.file)

	line := make(chan string)

	go func() {
		for br.Scan() {
			line <- br.Text()
		}
		close(line)
	}()

	return line
}

func (i *Input) Close() error {
	return i.file.Close()
}

func inputFile(filename string) (*os.File, error) {
	switch {
	case filename != "":
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}

		return f, nil
	case hasStdin():
		return os.Stdin, nil
	default:
		return nil, errors.New("no input data")
	}
}

func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	isPipedFromChrDev := (stat.Mode() & os.ModeCharDevice) == 0
	isPipedFromFIFO := (stat.Mode() & os.ModeNamedPipe) != 0

	return isPipedFromChrDev || isPipedFromFIFO
}
