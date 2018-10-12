package fix

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var (
	// Gob is a Codec that uses the gob package.
	Gob = Codec{gobMarshal, gobCompare}
	// JSON is a Codec that uses the json package.
	JSON = Codec{jsonMarshal, jsonCompare}
)

// Codec is funcstions to store structure to file.
type Codec struct {
	Marshal func(interface{}) ([]byte, error)
	Compare func(old, new []byte) error
}

// Fix uses JSON codec as default
func Fix(target interface{}, additional ...string) error {
	// Fix() caller's func name
	pt, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pt).Name()

	return JSON.fix(funcName, target, additional...)
}

func jsonMarshal(target interface{}) ([]byte, error) {
	return json.MarshalIndent(target, "", "  ")
}

func jsonCompare(old, new []byte) error {
	// If equal, target still have the same result.
	if bytes.Equal(new, old) {
		return nil
	}

	return fmt.Errorf("Diff: %s", lineDiff(string(old), string(new)))
}

func gobMarshal(target interface{}) ([]byte, error) {
	gob.Register(target)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(target)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func gobUnmarshal(raw []byte, target interface{}) error {
	return gob.NewDecoder(bytes.NewReader(raw)).Decode(target)
}

func gobCompare(old, new []byte) error {
	// If equal, target still have the same result.
	if bytes.Equal(new, old) {
		return nil
	}

	// If not, try decode and show diff as error
	var decodedOld interface{}
	err := gobUnmarshal(old, &decodedOld)
	if err != nil {
		return errors.Wrap(err, "File cannot decode")
	}
	var decodedNew interface{}
	err = gobUnmarshal(new, &decodedNew)
	if err != nil {
		return errors.Wrap(err, "File cannot decode")
	}

	// If decoded results are equal, It may be OK.
	if cmp.Equal(decodedOld, decodedNew, cmpopts.EquateEmpty()) {
		return nil
	}
	return fmt.Errorf("Diff: %s", cmp.Diff(decodedOld, decodedNew, cmpopts.EquateEmpty()))
}

// DefaultOutputPath is ./testdata/<caller_func_name>
func DefaultOutputPath(funcName string, additional ...string) string {
	f := []string{}
	f = append(f, funcName)
	f = append(f, additional...)

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

// Fix target pointer as codec encoded data file.
// if file doesn't exist, just write data to file and return error.
// if file exist, test if the target is the same as file's data.
func (c *Codec) Fix(target interface{}, additional ...string) error {
	// Fix() caller's func name
	pt, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pt).Name()

	return c.fix(funcName, target, additional...)
}

func (c *Codec) fix(funcName string, target interface{}, additional ...string) error {
	// Generate path
	path := outputPath(funcName, additional...)

	new, err := c.Marshal(&target)
	if err != nil {
		return errors.Wrap(err, "Target cannot encode")
	}

	// file is not exist, write and exit.
	if _, err = os.Stat(path); err != nil {
		err = os.MkdirAll(filepath.Dir(path), 0777)
		if err != nil {
			return errors.Wrap(err, "Cannot make dir")
		}

		err = ioutil.WriteFile(path, new, 0666)
		if err != nil {
			return errors.Wrap(err, "Failed to write file")
		}
		return errors.New("Write to file")
	}

	// if file exist..
	old, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "File cannot read: "+path)
	}

	// Compare and PrintDiff
	return c.Compare(old, new)
}

func lineDiff(src1, src2 string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(src1, src2, false)
	return dmp.DiffPrettyText(diffs)
}
