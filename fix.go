package fix

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
)

func encode(target interface{}) ([]byte, error) {
	gob.Register(target)

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&target)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decode(raw []byte, target *interface{}) error {
	dec := gob.NewDecoder(bytes.NewReader(raw))
	err := dec.Decode(target)
	if err != nil {
		return err
	}
	return nil
}

// DefaultOutputPath is ./testdata/<caller_func_name>
func DefaultOutputPath(funcName string, additional ...string) string {
	f := []string{}
	f = append(f, funcName)
	for _, a := range additional {
		f = append(f, a)
	}
	baseDir := "testdata/"
	return filepath.Join(baseDir,
		strings.Replace(strings.Join(f, "_"), "/", "_", -1),
	)
}

var outputPath = DefaultOutputPath

// SetOutputPathFunc overwrite DefaultOutputPath
func SetOutputPathFunc(in func(funcName string, additional ...string) string) {
	outputPath = in
}

// Fix target pointer as gob encoded data file.
// if file doesn't exist, just write data to file and return error.
// if file exist, test if the target is the same as file's data.
func Fix(target interface{}, additional ...string) error {
	pt, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pt).Name()
	path := outputPath(funcName, additional...)

	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return errors.Wrap(err, "Cannot make dir")
	}

	ret, err := encode(target)
	if err != nil {
		return errors.Wrap(err, "Target cannot encode")
	}

	_, err = os.Stat(path)
	if err != nil { // file is not exist
		ioutil.WriteFile(path, ret, 0666)
		return errors.New("Write to file")
	}

	// file exist
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "File cannot read: "+path)
	}

	// If equal, target still have the same result.
	if bytes.Equal(ret, raw) {
		return nil
	}

	// If not, try decode and show diff as error
	var valid interface{}
	err = decode(raw, &valid)
	if err != nil {
		return errors.Wrap(err, "File cannot decode: "+path)
	}

	// If decoded results are equal, It may be OK.
	if cmp.Equal(valid, target, cmpopts.EquateEmpty()) {
		return nil
	}

	return fmt.Errorf("Diff: %s", cmp.Diff(valid, target, cmpopts.EquateEmpty()))
}
