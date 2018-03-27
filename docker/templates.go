// Code generated by go-bindata.
// sources:
// static/fix.sh.template
// DO NOT EDIT!

package docker

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _staticFixShTemplate = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x91\x4f\x6f\x13\x31\x10\xc5\xef\xfe\x14\x2f\x4e\xd4\xd0\x4a\xee\x8a\x2b\x51\x4e\x14\x45\x9c\x40\xa9\x72\xaa\x10\x75\xec\xe9\xae\xa5\xc4\x5e\xec\x71\x29\x54\xf9\xee\x68\x67\x97\xe5\xcf\xc9\xd2\xf3\x6f\xde\xfc\x79\xcb\x45\x73\x0c\xb1\x29\x9d\x52\x4b\xec\x88\x51\xd8\x46\x6f\xb3\x87\xb3\xa7\x80\xc3\xfd\x87\xfd\xd7\x8f\x77\x78\xb6\x39\xd8\xe3\x89\xd4\x24\x6c\x5f\x5f\x6f\x0f\xc1\x5f\x2e\x6a\xb7\xff\x74\xf8\x3c\x29\x3b\x51\x94\x0f\x85\x73\xda\x3e\x8a\xb3\xb3\x8c\xa6\xcf\xc9\x35\xcf\x94\x4b\x48\xf1\x51\x29\x67\x0b\x61\x35\x62\x08\x51\x01\x37\x87\x63\x8d\x5c\x6f\xae\x15\x00\xd4\x42\xd9\x7a\x0f\x73\x2e\x98\xc6\x13\x4d\xc9\xef\x12\xef\x3b\x1b\x5b\xc2\x7a\xd0\xd6\xa8\xc1\x83\x13\xba\x54\x58\xa8\x75\x19\x24\x41\xc3\x13\x1e\xb0\x80\xf9\x09\xbd\x9a\x26\xd7\xf8\x82\xab\x2b\x3c\x40\xaf\xde\x04\x0f\x53\xa5\xe6\x5a\x63\xb1\xfd\x07\xda\x80\x3b\x8a\xe2\x22\x2d\x33\x59\xa6\x41\x13\x1e\x6d\x4e\xb5\x1f\xfc\x03\xc3\x27\x2a\x88\x89\x41\x2f\xa1\xf0\x54\x22\x80\x2c\x61\x62\x8a\xa6\xc6\xf0\xad\x12\x4c\x0b\xbd\xfa\x7d\x32\x3d\x42\x6a\x6e\x72\x4f\x3c\x77\x18\xb7\x80\x8d\x1e\xed\xb4\xcd\x78\x98\x73\xfa\xdf\xd3\x0c\xe0\x5f\xb3\x1b\xd3\x8a\xf0\xa7\x8f\x1c\x6f\x30\x78\x0a\xf2\xb8\x2e\x7d\x8f\x30\x7b\xf9\x78\x87\xa6\x4b\x67\x6a\x66\x68\xb3\x19\x12\x19\xa3\x20\xd7\x25\xe8\xbb\x31\xaa\x39\xb2\x02\x57\x73\xa6\xc8\xa7\x1f\xb2\x78\xa9\x7d\x9f\x32\x93\xbf\xd5\x63\xd5\x4b\x60\xbc\x55\x54\xac\xfb\x15\x00\x00\xff\xff\x8e\x29\xa1\xb4\x61\x02\x00\x00")

func staticFixShTemplateBytes() ([]byte, error) {
	return bindataRead(
		_staticFixShTemplate,
		"static/fix.sh.template",
	)
}

func staticFixShTemplate() (*asset, error) {
	bytes, err := staticFixShTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "static/fix.sh.template", size: 609, mode: os.FileMode(436), modTime: time.Unix(1522156906, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"static/fix.sh.template": staticFixShTemplate,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"static": &bintree{nil, map[string]*bintree{
		"fix.sh.template": &bintree{staticFixShTemplate, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
