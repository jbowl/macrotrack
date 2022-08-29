package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"macrotrack/internal/pkg/store"
	"macrotrack/internal/pkg/types"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
)

//var a main.App

var dataStore store.Storage

func DSNFromFile(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	// find string after DSN=
	for scanner.Scan() {

		ss := strings.SplitAfter(scanner.Text(), "DSN=")

		if len(ss) > 1 {

			return ss[1]

		}

		fmt.Println(scanner.Text())
	}
	return ""
}

func TestMain(m *testing.M) {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path) // for example /home/user

	cmdString := `./ctr.sh`
	//cmdString := `./mysqldb.sh`

	// DSN="..."
	DSN := DSNFromFile(cmdString)
	if len(DSN) < 1 {
		log.Fatalf("DSN string not found in container script file")
	}

	// startup contaner
	cmd := exec.Command(cmdString)
	err = cmd.Run()

	if err != nil {
		log.Fatalf("unable to start test container %v", err)
	}

	//dataStore = store.GetStorage("sqlserver", "user id=SA;password=___Aa123;port=1434;database=master")
	//dataStore = store.GetStorage("mysql", DSN)
	dataStore = store.GetStorage("sqlserver", DSN)

	//s.DSN = "gouser:gopwd@tcp(127.0.0.1:3306)/macros"

	dataStore.Init()

	exitVal := m.Run()

	os.Exit(exitVal)

}

func Create(c types.Macro) *http.Response {
	body, _ := json.Marshal(c)

	r := httptest.NewRequest(http.MethodPost, "/macros", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateMacroHandler(dataStore, w, r)

	return w.Result()
}

func Read(target string) *http.Response {

	r := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()

	ReadMacroHandler(dataStore, w, r)
	return w.Result()

}

func Update(target string, u types.Macro) *http.Response {

	body, _ := json.Marshal(u)

	r := httptest.NewRequest(http.MethodPut, target, bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	//	res := Read(target)

	UpdateMacroHandler(dataStore, w, r)

	return w.Result()
}

func TestCRUD(t *testing.T) {

	createData := types.Macro{Carbs: 100, Protein: 99, Fat: 98, Alcohol: 97}
	res := Create(createData)

	if res.StatusCode != 201 {
		res.Body.Close()
		t.Errorf("expected Status == 201  error got %v", res.StatusCode)
	}

	res.Body.Close()

	path := res.Header.Get("location")

	readRes := Read(path)

	defer readRes.Body.Close()

	//	var fileInfo stagingfile
	//	data, err := io.ReadAll(readRes.Body)

	var readData types.Macro

	//err := json.Unmarshal(readRes.Body, &readData)
	err := json.NewDecoder(readRes.Body).Decode(&readData)

	if err != nil {
		fmt.Println(err)
	}

	// check header

	if readRes.StatusCode != http.StatusOK {
		fmt.Println(readRes.StatusCode)
	}

	// compare all fields
	if createData.Carbs != readData.Carbs {

	}

	u := types.Macro{Carbs: 50, Protein: 50, Fat: 98, Alcohol: 97}

	uRes := Update(path, u)

	var updateData types.Macro

	//err := json.Unmarshal(readRes.Body, &readData)
	err = json.NewDecoder(uRes.Body).Decode(&updateData)

	if err != nil {
		fmt.Println(err)
	}

	path2 := uRes.Header.Get("location")
	//	path2 := h2[11:]

	readRes2 := Read(path2)

	defer readRes2.Body.Close()

	//	var fileInfo stagingfile
	//	data, err := io.ReadAll(readRes.Body)

	var updatedData types.Macro

	//err := json.Unmarshal(readRes.Body, &readData)
	err = json.NewDecoder(readRes2.Body).Decode(&updatedData)

	if err != nil {
		fmt.Println(err)
	}

	// now update

}

func TestCreateMacro(t *testing.T) {

	//	m := types.Macro{Carbs: 100, Protein: 99, Fat: 98, Alcohol: 97}

	//	json.Marshal(m)

	var jsonStr = []byte(`{"carbs":100, "protein": 99, "fat":98, "alcohol":97}`)

	req := httptest.NewRequest(http.MethodPost, "/macros", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	CreateMacroHandler(dataStore, w, req)

	res := w.Result()

	//defer res.Body.Close()

	if res.StatusCode != 201 {
		res.Body.Close()
		t.Errorf("expected Status == 201  error got %v", res.StatusCode)
	}

	res.Body.Close()

	//	fmt.Println(res.Body)

	//	data, err := io.ReadAll(res.Body)

	//	if err != nil {
	//		t.Errorf("expected nil error got %v", err)
	//	}

	//	if len(data) > 0 {

	//	}

	loc := res.Header.Get("location")
	req = httptest.NewRequest(http.MethodGet, loc, nil)

	w2 := httptest.NewRecorder()

	ReadMacroHandler(dataStore, w2, req)

	readRes := w2.Result()
	defer readRes.Body.Close()

	//	var fileInfo stagingfile

	data, err := io.ReadAll(readRes.Body)
	if err == nil {
		fmt.Println(err)
	}

	var m types.Macro

	err = json.Unmarshal(data, &m) // unmarshall byte to JSON
	//bytes, err =  json.Marshal(m)            // marshall json to byte
	//	err := json.NewDecoder(r.Body).Decode(&fileInfo)

	if err == nil {
		fmt.Println(err)
		//	fmt.Println(string(data))
	}

}

func TestReadMacro(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/macros/7D805D59-A2E2-464B-BCD1-8F07ABC84C10", nil)

	w := httptest.NewRecorder()

	//store := store.GetStorage("foo")
	//store.Init()

	ReadMacroHandler(dataStore, w, req)

	res := w.Result()

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		t.Errorf("expected nil error got %v", err)
	}

	if len(data) > 0 {

	}

}
