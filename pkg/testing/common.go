package saasy_testing

import (
	"io/ioutil"
	"os"
	"path"
)

// CreateJSONSpecFileAndOutdir creates a JSON spec file from a string and
// prepares a temporary directory for the results. Caller is responsible
// with deleting the temporary spec file and directory in the end.
func CreateJSONSpecFileAndOutdir(spec, dir, pattern string) (pSpec, pOutdir string, err error) {
	content := []byte(spec)
	f, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		return
	}
	if _, err = f.Write(content); err != nil {
		return
	}

	pSpec = path.Join(dir, f.Name())
	err = f.Close()

	pOutdir, err = ioutil.TempDir("", "testSaasy")
	if err != nil {
		os.Remove(f.Name())
	}
	return
}
