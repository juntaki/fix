package fix

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var (
	// JSON is a Codec that uses the json package.
	JSON = Codec{jsonMarshal, jsonCompare}
	// PP is a Codec that uses the pp package.
	PP = Codec{ppMarshal, ppCompare}
)

// Codec is functions to store structure to file.
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

// DefaultOutputPath is ./testdata/<caller_func_name>
func DefaultOutputPath(funcName string, additional ...string) string {
	var f []string
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
