package handlers

import (
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
	"testing"
)

//var a main.App

var dataStore store.Storage

func TestMain(m *testing.M) {

	//cmdString := `docker run -d -e "ACCEPT_EULA=Y" -e "MSSQL_SA_PASSWORD=___Aa123" --name sqlserver_test --hostname sqlserver_test --rm mcr.microsoft.com/mssql/server:2019-latest`

	//dataStore.Init()

	//store.DSN = fmt.Sprintf("user id=%s;password=%s;port=%s;database=%s", store.User, store.Password, store.Port, store.Database)

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path) // for example /home/user

	cmdString := `./ctr.sh`
	cmd := exec.Command(cmdString)

	err = cmd.Run()

	if err != nil {
		log.Fatalf("unable to start test container %v", err)

	}

	dataStore = store.GetStorage("sqlserver", "user id=SA;password=___Aa123;port=1434;database=master")

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

	h := res.Header.Get("location")
	path := h[11:]

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

	h2 := uRes.Header.Get("location")
	path2 := h2[11:]

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

	h := res.Header.Get("location")
	path := h[11:]
	req = httptest.NewRequest(http.MethodGet, path, nil)

	w2 := httptest.NewRecorder()

	ReadMacroHandler(dataStore, w2, req)

	readRes := w2.Result()
	defer readRes.Body.Close()

	//	var fileInfo stagingfile

	data, err := io.ReadAll(readRes.Body)

	//			err := json.Unmarshal(r.Body, fileInfo)
	//	err := json.NewDecoder(r.Body).Decode(&fileInfo)

	if err == nil {
		fmt.Println(string(data))
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
