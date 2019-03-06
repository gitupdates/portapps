package portapps

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

// SetFileAttributes set attributes to a file
func SetFileAttributes(path string, attrs uint32) error {
	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	return syscall.SetFileAttributes(pointer, attrs)
}

// CopyFile copy a file
func CopyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// CopyFolder copy a folder
func CopyFolder(source string, dest string) (err error) {
	err = os.MkdirAll(dest, 777)
	if err != nil {
		return err
	}

	folder, _ := os.Open(source)
	objects, err := folder.Readdir(-1)
	for _, object := range objects {
		sourceFile := path.Join(source, object.Name())
		destFile := path.Join(dest, object.Name())
		if object.IsDir() {
			err = CopyFolder(sourceFile, destFile)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(sourceFile, destFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// RemoveContents remove contents of a specified directory
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateFolderCheck to create a folder and get its path and return error
func CreateFolderCheck(path string) (string, error) {
	if err := os.MkdirAll(path, 777); err != nil {
		return "", err
	}
	return path, nil
}

// CreateFolder to create a folder and get its path
func CreateFolder(path string) string {
	Log.Infof("Create folder %s...", path)
	if _, err := CreateFolderCheck(path); err != nil {
		Log.Errorf("Cannot create folder: %v", err)
	}
	return path
}

// CreateFile creates / overwrites a file with content
func CreateFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err = file.Sync(); err != nil {
		return err
	}
	return nil
}

// PathJoin to join paths
func PathJoin(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return strings.Join(elem[i:], `\`)
		}
	}
	return ""
}

// AppPathJoin to join paths from Papp.Path
func AppPathJoin(elem ...string) string {
	return PathJoin(append([]string{Papp.Path}, elem...)...)
}

// FormatUnixPath to format a path for unix
func FormatUnixPath(path string) string {
	return strings.Replace(path, `\`, `/`, -1)
}

// FormatWindowsPath to format a path for windows
func FormatWindowsPath(path string) string {
	return strings.Replace(path, `/`, `\`, -1)
}

// Exists reports whether the named file or directory exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// WriteToFile reports writes content to a file
func WriteToFile(name string, content string) error {
	fo, err := os.Create(name)
	defer fo.Close()
	if err != nil {
		return err
	}
	if _, err = io.Copy(fo, strings.NewReader(content)); err != nil {
		return err
	}
	return nil
}

// ReplaceByPrefix replaces line in file starting with a specific prefix
func ReplaceByPrefix(filename string, prefix string, replace string) error {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, prefix) {
			lines[i] = replace
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filename, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}

// IsDirEmpty determines if directory is empty
func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	if _, err = f.Readdir(1); err == io.EOF {
		return true, nil
	}

	return false, err
}
