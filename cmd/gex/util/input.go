package util

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/afritzler/garden-examiner/cmd/gex/cleanup"
)

type FileInput interface {
	InheritedFiles(f []*os.File) ([]*os.File, string)
	CleanupFunction() func()
}

type FileInputFactory func(input []byte) (FileInput, error)

/////////////////////////////////////////////////////////////////////////////////
type StreamedFileInput struct {
	r *os.File
}

var _ FileInput = &StreamedFileInput{}

func NewStreamedFileInput(input []byte) (FileInput, error) {
	fmt.Printf("USING stream\n")
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	go func() {
		w.Write([]byte(input))
		w.Close()
	}()
	fi := &StreamedFileInput{r: r}
	return fi, nil
}

func (this *StreamedFileInput) InheritedFiles(f []*os.File) ([]*os.File, string) {
	if f == nil {
		return []*os.File{this.r}, "/dev/fd/3"
	}
	return append(f, this.r), fmt.Sprintf("/dev/fd/%d", 3+len(f))
}

func (this *StreamedFileInput) CleanupFunction() func() {
	return this.Close
}

func (this *StreamedFileInput) Close() {
	this.r.Close()
}

/////////////////////////////////////////////////////////////////////////////////

type TempFileInput struct {
	tmp *os.File
}

var _ FileInput = &TempFileInput{}

func NewTempFileInput(input []byte) (FileInput, error) {
	tmpfile, err := ioutil.TempFile("/tmp", "input")
	if err != nil {
		return nil, fmt.Errorf("cannot get temporary file: %s", err)
	}
	inp := &TempFileInput{tmp: tmpfile}

	if _, err := tmpfile.Write(input); err != nil {
		return nil, fmt.Errorf("cannot write temporary file '%s': %s", tmpfile.Name, err)
	}
	if err := tmpfile.Close(); err != nil {
		inp.Close()
		return nil, err
	}
	return inp, nil
}

func (this *TempFileInput) InheritedFiles(f []*os.File) ([]*os.File, string) {
	return f, this.tmp.Name()
}

func (this *TempFileInput) CleanupFunction() func() {
	return Cleanup(this.Close)
}

func (this *TempFileInput) Close() {
	os.Remove(this.tmp.Name())
}
