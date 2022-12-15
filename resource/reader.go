package resource

import (
	"bytes"
	"debug/pe"
	"expvar"
	"io"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/unicode"
)

var (
	peReaderTimeouts = expvar.NewInt("pe.reader.timeouts")
)

// Reader is the interface for PE (Portable Executable) format metadata parsing. The stdlib debug/pe package underpins
// the core functionality of the reader, but additionally, it provides numerous methods for reading resources, strings,
// IAT directories and other information that is not offered by the standard library package.
type Reader interface {
	// Read is the main method that reads the PE metadata for the specified image file.
	Read(filename string) (VersionResources, error)
	// FindSectionByRVA gets the section containing the given address.
	FindSectionByRVA(rva uint32) (*pe.Section, error)
	// FindOffsetByRVA returns the file offset that maps to the given RVA.
	FindOffsetByRVA(rva uint32) (int64, error)
}

type reader struct {
	f        *os.File
	sections []*pe.Section
	oh       interface{}
}

// NewReader builds a new instance of the PE reader.
func NewReader() Reader {
	return &reader{}
}

func (r *reader) Read(filename string) (VersionResources, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	r.f = f
	defer r.f.Close()
	pefile, err := pe.NewFile(f)
	if err != nil {
		return nil, err
	}
	r.sections = pefile.Sections
	r.oh = pefile.OptionalHeader

	p := make(map[string]string)

	var resDir pe.DataDirectory
	switch hdr := r.oh.(type) {
	case *pe.OptionalHeader32:
		resDir = hdr.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_RESOURCE]
	case *pe.OptionalHeader64:
		resDir = hdr.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_RESOURCE]
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		p, err = r.readResources(resDir.VirtualAddress)
		if err != nil {
			log.Warnf("fail to read %q PE resources: %v", filename, err)
		}
	}(&wg)

	// ensure this method terminates in a timely manner
	done := make(chan struct{})

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return p, nil
	case <-time.After(time.Second):
		log.Warn("wait timeout reached during PE metadata parsing")
		peReaderTimeouts.Add(1)
		return p, nil
	}
}

// readUTF16String reads an UTF16 string at the specified RVA.
func (r *reader) readUTF16String(rva uint32) (string, error) {
	data := make([]byte, 1024)
	offset, err := r.FindOffsetByRVA(rva)
	if err != nil {
		return "", err
	}
	n, err := r.f.ReadAt(data, offset)
	if err != nil {
		if err == io.EOF {
			return "", nil
		}
		return "", err
	}
	idx := bytes.Index(data[:n], []byte{0, 0})
	if idx < 0 {
		idx = n - 1
	}
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	utf8, err := decoder.Bytes(data[0 : idx+1])
	if err != nil {
		return "", err
	}
	return string(utf8), nil
}

func dwordAlign(offset, base int64) int64 {
	return ((offset + base + 3) & 0xfffffffc) - (base & 0xfffffffc)
}
