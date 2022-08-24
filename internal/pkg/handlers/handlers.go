package handlers

import (
	"encoding/json"
	"fmt"
	"macrotrack/internal/pkg/store"
	"macrotrack/internal/pkg/types"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Create
func CreateMacro(store store.Storage) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			CreateMacroHandler(store, w, r)
		})
}

func CreateMacroHandler(store store.Storage, w http.ResponseWriter, r *http.Request) {
	///////////   grab JSON from request body
	var macro types.Macro
	if err := json.NewDecoder(r.Body).Decode(&macro); err != nil {
		detailsResponse(w, types.ProblemDetails{
			Detail: "unable to decode request body",
			Title:  http.StatusText(http.StatusBadRequest),
			Status: http.StatusBadRequest,
		})
		return
	}

	if err := macro.Validate(); err != nil {
		detailsResponse(w, types.ProblemDetails{
			Detail: "unable to decode request body",
			Title:  http.StatusText(http.StatusBadRequest),
			Status: http.StatusBadRequest,
		})
		return
	}

	// location header

	uuid, err := store.Create(macro)

	if err != nil {
		detailsResponse(w, types.ProblemDetails{
			Detail: "store.Create",
			Title:  http.StatusText(http.StatusBadRequest),
			Status: http.StatusBadRequest,
		})
		return
	}

	//w.Header().Set("location", fmt.Sprintf("%s%s/%s", r.Host, r.URL.Path, uuid.UUID.String()))
	w.Header().Set("location", fmt.Sprintf("%s%s/%s", r.Host, r.URL.Path, uuid))
	w.WriteHeader(http.StatusCreated)
}

// Read
func ReadAllMacro(store store.Storage) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ReadAllMacroHandler(store, w, r)
		})
}

func ReadAllMacroHandler(store store.Storage, w http.ResponseWriter, r *http.Request) {
	macros, err := store.ReadAll()

	if err != nil {
		//if err == sql.ErrNoRows
		detailsResponse(w, types.ProblemDetails{
			Detail: err.Error(),
			Title:  http.StatusText(http.StatusBadRequest),
			Status: http.StatusBadRequest,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(macros)

}

func ReadMacro(store store.Storage) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			ReadMacroHandler(store, w, r)

		})
}

func ReadMacroHandler(store store.Storage, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_uuid := vars["uuid"]

	// test cases
	if len(_uuid) == 0 && len(r.URL.Path) > 32 {
		// strip of /macros from testcase
		_uuid = r.URL.Path[8:]

	}
	// TODO validate
	//u := uuid.UUID{}
	u, err := uuid.Parse(_uuid)

	if err != nil {
		detailsResponse(w, types.ProblemDetails{
			Detail: "uuid.Parse",
			Title:  http.StatusText(http.StatusBadRequest),
			Status: http.StatusBadRequest,
		})
		return
	}

	macro, err := store.Read(u)

	if err != nil {
		//if err == sql.ErrNoRows
		detailsResponse(w, types.ProblemDetails{
			Detail: err.Error(),
			Title:  http.StatusText(http.StatusBadRequest),
			Status: http.StatusBadRequest,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(macro)
}

func UpdateMacro(store store.Storage) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			vars := mux.Vars(r)

			u := vars["uuid"]

			macro_uuid, err := uuid.Parse(u)

			if err != nil {
				detailsResponse(w, types.ProblemDetails{
					Detail: "unable to decode request body",
					Title:  http.StatusText(http.StatusBadRequest),
					Status: http.StatusBadRequest,
				})
				return
			}

			fmt.Println(macro_uuid)

			///////////   grab JSON from request body
			var macro types.Macro
			if err := json.NewDecoder(r.Body).Decode(&macro); err != nil {
				detailsResponse(w, types.ProblemDetails{
					Detail: "unable to decode request body",
					Title:  http.StatusText(http.StatusBadRequest),
					Status: http.StatusBadRequest,
				})
				return
			}

			if err := macro.Validate(); err != nil {
				detailsResponse(w, types.ProblemDetails{
					Detail: "unable to decode request body",
					Title:  http.StatusText(http.StatusBadRequest),
					Status: http.StatusBadRequest,
				})
				return
			}

			// location header

			if err := store.Update(macro_uuid, macro); err != nil {
				detailsResponse(w, types.ProblemDetails{
					Detail: "store.Create",
					Title:  http.StatusText(http.StatusBadRequest),
					Status: http.StatusBadRequest,
				})
				return
			}

			//w.Header().Set("location", fmt.Sprintf("%s%s/%s", r.Host, r.URL.Path, uuid.UUID.String()))
			//w.Header().Set("location", fmt.Sprintf("%s%s/%s", r.Host, r.URL.Path, uuid))
			w.WriteHeader(http.StatusOK)

		})
}

func DeleteMacro(store store.Storage) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			vars := mux.Vars(r)
			_uuid := vars["uuid"]

			// TODO validate

			u, err := uuid.Parse(_uuid)
			if err != nil {
				detailsResponse(w, types.ProblemDetails{
					Detail: "uuid.Parse",
					Title:  http.StatusText(http.StatusBadRequest),
					Status: http.StatusBadRequest,
				})
				return
			}

			store.Delete(u)
		})
}

func detailsResponse(w http.ResponseWriter, details types.ProblemDetails) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(details.Status)

	json.NewEncoder(w).Encode(&details)
}

func Custom404Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		detailsResponse(w, types.ProblemDetails{
			Title:  http.StatusText(http.StatusNotFound),
			Status: http.StatusNotFound,
		})
		//	JSONError(w, "Route does not exist", http.StatusNotFound)
	})
}

func Custom405Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		detailsResponse(w, types.ProblemDetails{
			Title:  http.StatusText(http.StatusNotFound),
			Status: http.StatusNotFound,
		})

	})
}
