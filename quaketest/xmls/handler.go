// fortune_server.go implements the 'handler' and other functions to return a result of 'Omikuji'.
package xmls

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// var resOmikuji = [4]string{"大吉", "中吉", "小吉", "凶"}
// var retXml = []string{}

// Res is a struct for json has one field 'result'.
type Res struct {
	Result string
}

// WeatherHandler implements the 'handler' for the  server.
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/xml")

	res, code, err := result()
	if err != nil {
		log.Printf("handler error: %v\ncode: %d\n", err, code)
		http.Error(w, err.Error(), code)
	}
	w.Header().Add("Status Code", "200 OK")
	if _, err := fmt.Fprint(w, res); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func globQuakeXmls(dir string) ([]string, error) {
	return filepath.Glob(dir + "/32-35_??_??_100*.xml")
}

// result function returns the contents in a xml file.
func result() (string, int, error) {
	xmlPaths, err := globQuakeXmls("./jmaxml_20210730_Samples")
	if err != nil {
		return "", 404, err
	}
	idx := rand.Intn(len(xmlPaths))
	f, err := os.OpenFile(xmlPaths[idx], os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", 404, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("open:", xmlPaths[idx])

	res, err := readFile(f)
	if err != nil {
		return "", 500, err
	}
	println(utf8.ValidString(res))

	return res, 200, nil
}

func readFile(f *os.File) (string, error) {
	r := &bytes.Buffer{}
	n, err := io.Copy(r, f)
	if err != nil {
		return "", err
	}
	if n == 0 {
		return "", fmt.Errorf("the xml file is empty")
	}

	return strings.Replace(r.String(), "\t", "", -1), nil
}
