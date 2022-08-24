package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"macrotrack/internal/pkg/store"
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

	cmdString := `./ctr.sh`
	cmd := exec.Command(cmdString)

	err := cmd.Run()

	if err != nil {
		log.Fatalf("unable to start test container %v", err)

	}

	dataStore = store.GetStorage("sqlserver", "user id=SA;password=___Aa123;port=1434;database=master")

	dataStore.Init()

	exitVal := m.Run()

	os.Exit(exitVal)

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
