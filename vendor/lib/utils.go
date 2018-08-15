package lib

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strings"

	"github.com/simplejia/namecli/api"
)

func TestPost(h http.HandlerFunc, params interface{}) (body []byte, err error) {
	v, err := json.Marshal(params)
	if err != nil {
		return
	}
	r, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(v))
	if err != nil {
		return
	}
	w := httptest.NewRecorder()
	h(w, r)
	body = w.Body.Bytes()
	if g, e := w.Code, http.StatusOK; g != e {
		err = fmt.Errorf("http resp status not ok: %s", http.StatusText(g))
		return
	}
	return
}

func DeduplicateInt64s(a []int64) (result []int64) {
	if len(a) == 0 {
		return
	}

	exists := map[int64]bool{}

	result = a[:0]
	for _, e := range a {
		if exists[e] {
			continue
		}
		exists[e] = true

		result = append(result, e)
	}

	return
}

func NameWrap(name string) (addr string, err error) {
	if strings.HasSuffix(name, ".ns") {
		return api.Name(name)
	}

	return name, nil
}

func SearchInt64s(a []int64, x int64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func BingoDisorderInt64s(a []int64, x int64) bool {
	for _, e := range a {
		if e == x {
			return true
		}
	}

	return false
}

func Int64s(a []int64) {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
}

func ZipInt64s(a []int64) (result []byte, err error) {
	if len(a) == 0 {
		return
	}

	var b bytes.Buffer
	zw, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return
	}

	err = json.NewEncoder(zw).Encode(a)
	if err != nil {
		return
	}
	zw.Close()

	result = b.Bytes()
	return
}

func UnzipInt64s(a []byte) (result []int64, err error) {
	if len(a) == 0 {
		return
	}

	zr := flate.NewReader(bytes.NewReader(a))
	err = json.NewDecoder(zr).Decode(&result)
	if err != nil {
		return
	}
	return
}

func ZipBytes(a []byte) (result []byte, err error) {
	if len(a) == 0 {
		return
	}

	var b bytes.Buffer
	zw, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return
	}

	zw.Write(a)
	zw.Close()

	result = b.Bytes()
	return
}

func UnzipBytes(a []byte) (result []byte, err error) {
	if len(a) == 0 {
		return
	}

	zr := flate.NewReader(bytes.NewReader(a))
	bs, err := ioutil.ReadAll(zr)
	if err != nil {
		return
	}

	result = bs
	return
}

func TrimDataURL(data string) (ret string) {
	if data == "" {
		return
	}

	return regexp.MustCompile(`^data:[0-9a-zA-Z/]+?;base64,`).ReplaceAllString(data, "")
}

func TruncateWithSuffix(data string, length int, suffix string) (ret string) {
	rdata := []rune(data)
	if len(rdata) > length {
		ret = string(rdata[:length]) + suffix
	} else {
		ret = data
	}

	return
}
