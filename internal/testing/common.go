package saasytesting

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

var generateReferenceFiles = false

func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// GetTestingCommonDirectory returns the directory of the common testing package.
func GetTestingCommonDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}

// CreateJSONSpecFile creates a JSON spec file from a string.
// Caller is responsible with deleting the temporary spec file in the end.
func CreateJSONSpecFile(spec, dir, pattern string) (pSpec string, err error) {
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
	return
}

// CreateOutdir creates a temporary directory for the results.
// Caller is responsible with deleting the temporary directory in the end.
func CreateOutdir() (pOutdir string, err error) {
	pOutdir, err = ioutil.TempDir("", "testSaasy")
	return
}

// CheckFilesInDirsEqual checks that files with the same name from the two
// given directories are equal in content
func CheckFilesInDirsEqual(t *testing.T, outDir string, referenceDir string, filenames []string) {
	for _, fname := range filenames {
		actualFileName := path.Join(outDir, fname)
		referenceFileName := path.Join(referenceDir, fname+".reference")

		if !generateReferenceFiles {
			bActual, err := ioutil.ReadFile(actualFileName)
			require.NoError(t, err)

			bExpected, err := ioutil.ReadFile(referenceFileName)
			require.NoError(t, err)
			require.Equal(t, bExpected, bActual)
		} else {
			err := copyFileContents(actualFileName, referenceFileName)
			if err != nil {
				require.FailNowf(t, "could not write reference file", "error - %v", err)
			}
		}
	}
}
