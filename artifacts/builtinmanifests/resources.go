// Code generated by "esc -pkg builtinmanifests -prefix  -ignore  -include  -o resources.go manifests"; DO NOT EDIT.

package builtinmanifests

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isDir {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is not directory", f.name)
	}

	fis, ok := _escDirs[f.local]
	if !ok {
		return nil, fmt.Errorf(" escFile.Readdir: '%s' is directory, but we have no info about content of this dir, local=%s", f.name, f.local)
	}
	limit := count
	if count <= 0 || limit > len(fis) {
		limit = len(fis)
	}

	if len(fis) == 0 && count > 0 {
		return nil, io.EOF
	}

	return []os.FileInfo(fis[0:limit]), nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/manifests/large_php.yml": {
		name:    "large_php.yml",
		local:   "manifests/large_php.yml",
		size:    3838,
		modtime: 1543322378,
		compressed: `
H4sIAAAAAAAC/+xWTYvkOA++51eI7svUS1V1VzMML7kt+8EODCxM97CHZQmOo0o87dhZS051za9f7LiS
1Ef3fLCHXRgamrIkS/KjR4pWq1VGssFW5NDfrjdZhqYvjGgxh0dvtDeZQ7LeSSxqZ32XdK9eXRAvFpm2
UrCyJlgcfi8WWaaEoBzEJ+8QroEbRSCFgRKTzDog2yJYbtBBpwVvrWuzrG8H/5RnACsYgn/0bVfapwwA
oEUWefwFMGTC++7ERlpvOIdNPNCjz+GehamEq4qf7jdFfxcVw72+Pcmvb0NyfUu0BOmdQ8N6D9bofVAp
AvJdZx1jFb1YKjpnt0rjIStRtcqMsE2nxSIZaGX8UyGt2arauwG+pAIgaqYDQOdLrWTxiHuaiwM2B99E
zXoyS1GIrRP1mJNqZwcAu92iy+FD6Q37e3Q9uuwoIjVB/6Mw1igp9JRdAHPzZn37evXu4X4U9+gokkAL
RuIkt1RUih6nsK0wosYqSlPZxsK8ez+5k0I2ytQ5vEdR/e4U46RyKBgL2w2k+8XZ9m14WzKoBIvofobW
6ovjvhj5LPbPbcf7mTp6J/UJi7rMYXN7O+quE5NDzXs0A1Teq2pGvRKj5CZyeglqG49LsGZQwM56XQWC
1gadYKxAUPQawxjknXXTs1dAvjTIiYdkkFebMaFkXRBK7xTvj1qdqCkM1dkJBVWXA7FgJaPGWY2zaPMJ
sg5VCGhksx7eYbmiyDT6hh793mrfW+2/2WraiqoohRZGoksmx7LnbEshH9FUhagqh0RFZ61ODi6qvrC9
d1gOfTg2+dDKB4AUgSesgC3UOLwfuEGosNN236JhIOlUx/SZ5r+sjV6GTULaaqp5L9xR23XOfkTJRVo6
gm3RCW5Coz+nG9t+7qBWXDjsFaUt5ZJ8vHiSbdd0WdYb5KN1pJ8qfCgAdUJiYOI6/t1s3gyjIRLjbEye
UMQJU88v370eVbVg3In9qNxk2RFJjvK6RKujyRoll7hzluELBGtQaG7CJC7PPwANczeoslkd2Eqrc3iQ
HVzD0NNCQy+0D7z740F2S/iVOf2nP6e71nEO/5+3t8O/PBInLlzdXCWZcliFVj6Ei5sacqBxcBp2uuh8
Db+FlXOnCJegOJgZyyC0tjus1qeVoual55wkejd8tJzXzyATNF8LzIeqW8IPWk+obJ01HKpzBs+hbM8q
vn6WzMtdnBU5uzxqjmh5PnDOLl1ELJZkFYLNIFM2XMjhbjbzK+VQDh+Kt6a03lSjSkiJRHnAz+4+U8M0
TQJ2RWrJq/9dTVGQWJm4QhzZzFBOHkYkHW7V0/NeLtvNqfctgBE1F/Da/GvwSk3yD+GV9cIpUSY0RiQO
29psiQynvwMAAP//FAGNRf4OAAA=
`,
	},

	"/manifests/maximum_php.yml": {
		name:    "maximum_php.yml",
		local:   "manifests/maximum_php.yml",
		size:    3838,
		modtime: 1543322378,
		compressed: `
H4sIAAAAAAAC/+xWTYvkOA++51eI7svUS1V1VzMML7kt+8EODCxM97CHZQmOo0o87dhZS051za9f7LiS
1Ef3fLCHXRgamrIkS/KjR4pWq1VGssFW5NDfrjdZhqYvjGgxh0dvtDeZQ7LeSSxqZ32XdK9eXRAvFpm2
UrCyJlgcfi8WWaaEoBzEJ+8QroEbRSCFgRKTzDog2yJYbtBBpwVvrWuzrG8H/5RnACsYgn/0bVfapwwA
oEUWefwFMGTC++7ERlpvOIdNPNCjz+GehamEq4qf7jdFfxcVw72+Pcmvb0NyfUu0BOmdQ8N6D9bofVAp
AvJdZx1jFb1YKjpnt0rjIStRtcqMsE2nxSIZaGX8UyGt2arauwG+pAIgaqYDQOdLrWTxiHuaiwM2B99E
zXoyS1GIrRP1mJNqZwcAu92iy+FD6Q37e3Q9uuwoIjVB/6Mw1igp9JRdAHPzZn37evXu4X4U9+gokkAL
RuIkt1RUih6nsK0wosYqSlPZxsK8ez+5k0I2ytQ5vEdR/e4U46RyKBgL2w2k+8XZ9m14WzKoBIvofobW
6ovjvhj5LPbPbcf7mTp6J/UJi7rMYXN7O+quE5NDzXs0A1Teq2pGvRKj5CZyeglqG49LsGZQwM56XQWC
1gadYKxAUPQawxjknXXTs1dAvjTIiYdkkFebMaFkXRBK7xTvj1qdqCkM1dkJBVWXA7FgJaPGWY2zaPMJ
sg5VCGhksx7eYbmiyDT6hh793mrfW+2/2WraiqoohRZGoksmx7LnbEshH9FUhagqh0RFZ61ODi6qvrC9
d1gOfTg2+dDKB4AUgSesgC3UOLwfuEGosNN236JhIOlUx/SZ5r+sjV6GTULaaqp5L9xR23XOfkTJRVo6
gm3RCW5Coz+nG9t+7qBWXDjsFaUt5ZJ8vHiSbdd0WdYb5KN1pJ8qfCgAdUJiYOI6/t1s3gyjIRLjbEye
UMQJU88v370eVbVg3In9qNxk2RFJjvK6RKujyRoll7hzluELBGtQaG7CJC7PPwANczeoslkd2Eqrc3iQ
HVzD0NNCQy+0D7z740F2S/iVOf2nP6e71nEO/5+3t8O/PBInLlzdXCWZcliFVj6Ei5sacqBxcBp2uuh8
Db+FlXOnCJegOJgZyyC0tjus1qeVoual55wkejd8tJzXzyATNF8LzIeqW8IPWk+obJ01HKpzBs+hbM8q
vn6WzMtdnBU5uzxqjmh5PnDOLl1ELJZkFYLNIFM2XMjhbjbzK+VQDh+Kt6a03lSjSkiJRHnAz+4+U8M0
TQJ2RWrJq/9dTVGQWJm4QhzZzFBOHkYkHW7V0/NeLtvNqfctgBE1F/Da/GvwSk3yD+GV9cIpUSY0RiQO
29psiQynvwMAAP//FAGNRf4OAAA=
`,
	},

	"/manifests/medium_php.yml": {
		name:    "medium_php.yml",
		local:   "manifests/medium_php.yml",
		size:    3838,
		modtime: 1543322378,
		compressed: `
H4sIAAAAAAAC/+xWTYvkOA++51eI7svUS1V1VzMML7kt+8EODCxM97CHZQmOo0o87dhZS051za9f7LiS
1Ef3fLCHXRgamrIkS/KjR4pWq1VGssFW5NDfrjdZhqYvjGgxh0dvtDeZQ7LeSSxqZ32XdK9eXRAvFpm2
UrCyJlgcfi8WWaaEoBzEJ+8QroEbRSCFgRKTzDog2yJYbtBBpwVvrWuzrG8H/5RnACsYgn/0bVfapwwA
oEUWefwFMGTC++7ERlpvOIdNPNCjz+GehamEq4qf7jdFfxcVw72+Pcmvb0NyfUu0BOmdQ8N6D9bofVAp
AvJdZx1jFb1YKjpnt0rjIStRtcqMsE2nxSIZaGX8UyGt2arauwG+pAIgaqYDQOdLrWTxiHuaiwM2B99E
zXoyS1GIrRP1mJNqZwcAu92iy+FD6Q37e3Q9uuwoIjVB/6Mw1igp9JRdAHPzZn37evXu4X4U9+gokkAL
RuIkt1RUih6nsK0wosYqSlPZxsK8ez+5k0I2ytQ5vEdR/e4U46RyKBgL2w2k+8XZ9m14WzKoBIvofobW
6ovjvhj5LPbPbcf7mTp6J/UJi7rMYXN7O+quE5NDzXs0A1Teq2pGvRKj5CZyeglqG49LsGZQwM56XQWC
1gadYKxAUPQawxjknXXTs1dAvjTIiYdkkFebMaFkXRBK7xTvj1qdqCkM1dkJBVWXA7FgJaPGWY2zaPMJ
sg5VCGhksx7eYbmiyDT6hh793mrfW+2/2WraiqoohRZGoksmx7LnbEshH9FUhagqh0RFZ61ODi6qvrC9
d1gOfTg2+dDKB4AUgSesgC3UOLwfuEGosNN236JhIOlUx/SZ5r+sjV6GTULaaqp5L9xR23XOfkTJRVo6
gm3RCW5Coz+nG9t+7qBWXDjsFaUt5ZJ8vHiSbdd0WdYb5KN1pJ8qfCgAdUJiYOI6/t1s3gyjIRLjbEye
UMQJU88v370eVbVg3In9qNxk2RFJjvK6RKujyRoll7hzluELBGtQaG7CJC7PPwANczeoslkd2Eqrc3iQ
HVzD0NNCQy+0D7z740F2S/iVOf2nP6e71nEO/5+3t8O/PBInLlzdXCWZcliFVj6Ei5sacqBxcBp2uuh8
Db+FlXOnCJegOJgZyyC0tjus1qeVoual55wkejd8tJzXzyATNF8LzIeqW8IPWk+obJ01HKpzBs+hbM8q
vn6WzMtdnBU5uzxqjmh5PnDOLl1ELJZkFYLNIFM2XMjhbjbzK+VQDh+Kt6a03lSjSkiJRHnAz+4+U8M0
TQJ2RWrJq/9dTVGQWJm4QhzZzFBOHkYkHW7V0/NeLtvNqfctgBE1F/Da/GvwSk3yD+GV9cIpUSY0RiQO
29psiQynvwMAAP//FAGNRf4OAAA=
`,
	},

	"/manifests/small_jmeter.yml": {
		name:    "small_jmeter.yml",
		local:   "manifests/small_jmeter.yml",
		size:    3043,
		modtime: 1543323298,
		compressed: `
H4sIAAAAAAAC/+xWTY/bNhC981cM0ktcVK5lbAJEt6IfQICcsil6FEbUWGJXIoUZUrvOry8oyZL8seki
6KEIAl/EeaM3o8d5pJMkUaJrajGDfrdNlSLb5xZbyuAh2CZYxSQusKa8Yhe6CXv9+kZ4s1GN0+iNszHj
9LzZKGUQJQP8HJjgB/C1EdBooaAp5hjEtQTO18TQNegPjlul+nbkl0wBJDAW/zu0XeGeFABASx6z4Qlg
7MQfu4sc7YL1GaTDQh5CBvcebYlc5r/dp3m/H4Dxvb696K9vY3N9K/IT6MBM1jdHcLY5RsgISOg6x57K
gcVJ3rE7mIZOXWHZGjvLtqw2mymhMTY85drZg6kCj/JNEIBIvSwAulA0RucPdJR1OGpz4hapt0vaVEW8
Y6zmnky7WgC4w4E4gz+LYH24J+6J1VlFqSP+K1pnjcZm6S6Kmb7d7u6SD5/u53BPLMMQNOhJ/BR3kpdG
HpayLVqsqByi07bNG/Ph40KnUdfGVhl8JCz/YuNpgZjQU+66cej+YNe+j982JZTocaBfqZW8uO4XK1/V
/r3t/HEFD+xiPlNeFRmku92AWfKPjpd+EpBQWPLTgIgln6Qzy5SdC+nAxh/PPChS51YqdTEbpstAPHqj
B4RdQ6tqa2tvozyxTbU2V0ueOIllieUr/PPdBt9t8KwN7l7ugpVk/2KEcWSjFwCedwO83A430bGKUr0l
f3Yd9YtnsSyZRHLpUFM0/Xb4/Zy+Hcdv+Mir8hemZ7TV+uX93QxV6OkRjzOY3mI6E25Ntn/z5kt0EU5P
kFKNwzIvsEGriSVT6vYenEmxPpKu0jnckB6bxj0mIvWycWxczF/OzHGMmPQ4Ye9t4YItZwi1JpEMfolU
KxrnnXZNBp90tzh1/M8SL+x8kubVj6+WKiTe2OHsOcvZ7y8ZTjvdMR3M0/Mst/POD9uvlSwpqb8tW/qN
6nYtw+TKmyrs/ycqpLt37/5LHVSPbLCYRmPW43TtrW7juPonAAD//0GbVcDjCwAA
`,
	},

	"/manifests/small_kubernetes.yml": {
		name:    "small_kubernetes.yml",
		local:   "manifests/small_kubernetes.yml",
		size:    3924,
		modtime: 1543231989,
		compressed: `
H4sIAAAAAAAC/+xWW6vbuBN/96cYTl56/iRpEkopfvuzF7ZQWOhp2YdlMbI8sdUjS17NyGn66RfJiu1c
Ti/LPuxCOXCI5jeaGf/molmtVhnJBluRQ79Zb7MMTV8Y0WIOj95obzKHZL2TWNTO+i5hz57dEN/fZ9pK
wcqaoHH6fX+fZUoIykF88g5hAdwoAikMlJhk1gHZFsFygw46LXhvXZtlfTvYpzwDWMHg/INvu9J+zAAA
WmSRx18AQyR87C50pPWGc9jGAz36HB5YmEq4qvjxYVv0uwgM9/r2Ir6+DcH1LdESpHcODesjWKOPAVIE
5LvOOsYqWrFUdM7ulcZTVKJqlRlpm07390lBK+M/FtKavaq9G+hLEABRMx0AOl9qJYtHPNJcHLg52SZq
1pNa8kJsnajHmFQ7OwDY/R5dDu9Lb9g/oOvRZWceqQn4D8JYo6TQU3SBzO3L9ebF6s27h1Hco6NYBFow
Eie5paJS9Di5bYURNVZRmtI2JubN28mcFLJRps7hLYrqN6cYJ8ihYCxsNxTdz862r8O3JYVKsIjmZ2yt
vtrvZz1f+f6p7fg4g6N1Up+wqMsctpvNiC1SJYec92gGqrxX1az0SoyS57Gml6D28bgEawYADtbrKhRo
bdAJxgoERavRjUE+WDd99grIlwY51SEZ5NV2DChpF4TSO8XHs1YnagpDdXZRgqrLgViwkhFxVuPM23yC
rEMWAhtZNodKdAYZ6csturtq0e+t9r3V/puttgBtRVWUQgsj0SWlc9nT2qWQj2iqQlSVQ6Kis1YnEzeh
r2zxqRnHTh/6+cSSIvCEFbCFGgcSgBuECjttjy0aBpJOdUxfmAA30YtZMGVtuh/9DAuHtBWOegC9cCP1
i6GNnP2Akou0nwT9ohPchJnwFDZOiHMTteLCYa8orTS35LOrF1F3TZdlvUE+2176qSBOuaJOSAyFu45/
z7cvh0kS6+hqql5UlBOmnl/evRihWjAexHEEt1l2Vk/R9BT2dRUuLqdxkt0qtmRsHumTNRnUGhSamzDC
S7xxt2HuBnCEYl7YSqtzeCc7WMAwEISGXmgf6vX3d7Jbwi/M6T/9Mb9tHefwajMTLcDhnx6JU4XcPb9L
MuWwCrPg5DKuesihBYLhsBRGB2v4NeysB0W4BMVBzVgGobU9YLW++q7wnH72s64C3u2SyHn9JFMB+3ai
3lfdEv6v9ZylvbOGQ9Zu0HVK6GegbxlNp7vzUiiu0p/dnl1nTXVjgl3dGumbkxfztArepmfYKRsu5LCb
vSSVciiH5+e1Ka031QgJKZEoD0Taw8zMeVLTMz4MnkBgkTr37n93kxckViYuJmc6rzaXFkYyHe7Vx6et
3NabV+PfIYyoucHX9l/D1273T/KV9cIpUSY2RiZOO+BsNw2nvwIAAP//Ps/0l1QPAAA=
`,
	},

	"/manifests/small_php.yml": {
		name:    "small_php.yml",
		local:   "manifests/small_php.yml",
		size:    3839,
		modtime: 1543321722,
		compressed: `
H4sIAAAAAAAC/+xWTYvkOA++51eI7svUS1V1VzMML7kt+8EODCxM97CHZQmOo0o87dhZS051za9f7LiS
1Ef3fLCHXRgamrIkS/KjR4pWq1VGssFW5NDfrjdZhqYvjGgxh0dvtDeZQ7LeSSxqZ32XdK9eXRAvFpm2
UrCyJlgcfi8WWaaEoBzEJ+8QroEbRSCFgRKTzDog2yJYbtBBpwVvrWuzrG8H/5RnACsYgn/0bVfapwwA
oEUWefwFMGTC++7ERlpvOIdNPNCjz+GehamEq4qf7jdFfxcVw72+Pcmvb0NyfUu0BOmdQ8N6D9bofVAp
AvJdZx1jFb1YKjpnt0rjIStRtcqMsE2nxSIZaGX8UyGt2arauwG+pAIgaqYDQOdLrWTxiHuaiwM2B99E
zXoyS1GIrRP1mJNqZwcAu92iy+FD6Q37e3Q9uuwoIjVB/6Mw1igp9JRdAHPzZn37evXu4X4U9+gokkAL
RuIkt1RUih6nsK0wosYqSlPZxsK8ez+5k0I2ytQ5vEdR/e4U46RyKBgL2w2k+8XZ9m14WzKoBIvofobW
6ovjvhj5LPbPbcf7mTp6J/UJi7rMYXN7O+quE5NDzXs0A1Teq2pGvRKj5CZyeglqG49LsGZQwM56XQWC
1gadYKxAUPQawxjknXXTs1dAvjTIiYdkkFebMaFkXRBK7xTvj1qdqCkM1dkJBVWXA7FgJaPGWY2zaPMJ
sg5VCGhksx7eYbmiyDT6lh793mvfe+2/2WvaiqoohRZGoksmx7LnbEshH9FUhagqh0RFZ61ODi6qvrC/
d1gOjTh2+dDLB4AUgSesgC3UOLwfuEGosNN236JhIOlUx/SZ7r+sjV6GVULaaqp5L9xR23XOfkTJRdo6
gm3RCW5Coz+nG9t+7qBWXDjsFaU15ZJ8vHiSbdd0WdYb5KN9pJ8qfCgAdUJiYOI6/t1s3gyjIRLjbE6e
UMQJU88v370eVbVg3In9qNxk2RFJjvK6RKuj0Roll7hzluELBGtQaG7CJC7PvwANczeoslkd2Eqrc3iQ
HVzD0NNCQy+0D7z740F2S/iVOf2nP6e71nEO/5+3t8O/PBInLlzdXCWZcliFVj6Ei6sacqBxcBqWuuh8
Db+FnXOnCJegOJgZyyC0tjus1qeVoual55wkejd8tZzXzyATNF8LzIeqW8IPWk+obJ01HKpzBs+hbM8q
vn6WzMtdnBU5uzxqjmh5PnDOLl1ELJZkFYLNIFM2XMjhbjbzK+VQDh+Kt6a03lSjSkiJRHnAz+4+U8M0
TQJ2RWrJq/9dTVGQWJm4QhzZzFBOHkYkHW7V0/NeLtvNqfctgBE1F/Da/GvwSk3yD+GV9cIpUSY0RiQO
29psiwynvwMAAP//GjIHz/8OAAA=
`,
	},

	"/manifests": {
		name:  "manifests",
		local: `manifests`,
		isDir: true,
	},
}

var _escDirs = map[string][]os.FileInfo{

	"manifests": {
		_escData["/manifests/large_php.yml"],
		_escData["/manifests/maximum_php.yml"],
		_escData["/manifests/medium_php.yml"],
		_escData["/manifests/small_jmeter.yml"],
		_escData["/manifests/small_kubernetes.yml"],
		_escData["/manifests/small_php.yml"],
	},
}
