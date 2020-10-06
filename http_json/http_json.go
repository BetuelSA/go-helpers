package http_json

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/BetuelSA/go-helpers/errors"
	"github.com/google/jsonapi"
)

const (
	headerContentType = "Content-Type"
)

func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func WriteErrJSON(w http.ResponseWriter, err error) {
	var code int
	errorType := errors.GetType(err)

	switch errorType {
	case errors.BadRequest:
		code = http.StatusBadRequest
	case errors.Unauthorized:
		code = http.StatusUnauthorized
	case errors.Forbidden:
		code = http.StatusForbidden
	case errors.NotFound:
		code = http.StatusNotFound
	case errors.MethodNotAllowed:
		code = http.StatusMethodNotAllowed
	case errors.PreconditionFailed:
		code = http.StatusPreconditionFailed
	case errors.UnsupportedMediaType:
		code = http.StatusUnsupportedMediaType
	case errors.InternalServerError:
		code = http.StatusInternalServerError
	case errors.NotImplemented:
		code = http.StatusNotImplemented
	case errors.ServiceUnavailable:
		code = http.StatusServiceUnavailable
	default:
		code = http.StatusInternalServerError
	}

	if errorType == errors.NoType {
		log.Printf(err.Error())
	}

	var errorMap = map[string]interface{}{
		"error":   code,
		"message": errors.GetErrorMessage(err),
		"detail":  errors.GetErrorDetail(err),
	}

	errorContext := errors.GetErrorContext(err)
	if errorContext != nil {
		errorMap["context"] = errorContext
	}

	WriteJSON(w, code, errorMap)
}

func WriteJSONAPI(w http.ResponseWriter, code int, data interface{}) {
	jsonapiRuntime := jsonapi.NewRuntime()
	w.Header().Set(headerContentType, jsonapi.MediaType)
	w.WriteHeader(code)

	err := jsonapiRuntime.MarshalPayload(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteErrJSONAPI(w http.ResponseWriter, err error) {
	var code int
	errorType := errors.GetType(err)

	switch errorType {
	case errors.BadRequest:
		code = http.StatusBadRequest
	case errors.Unauthorized:
		code = http.StatusUnauthorized
	case errors.Forbidden:
		code = http.StatusForbidden
	case errors.NotFound:
		code = http.StatusNotFound
	case errors.MethodNotAllowed:
		code = http.StatusMethodNotAllowed
	case errors.PreconditionFailed:
		code = http.StatusPreconditionFailed
	case errors.UnsupportedMediaType:
		code = http.StatusUnsupportedMediaType
	case errors.InternalServerError:
		code = http.StatusInternalServerError
	case errors.NotImplemented:
		code = http.StatusNotImplemented
	case errors.ServiceUnavailable:
		code = http.StatusServiceUnavailable
	default:
		code = http.StatusInternalServerError
	}

	if errorType == errors.NoType {
		log.Printf(err.Error())
	}

	var errorObject = jsonapi.ErrorObject{
		Status: strconv.Itoa(code),
		Title:  errors.GetErrorMessage(err),
		Detail: errors.GetErrorDetail(err),
	}

	errorObjects := []*jsonapi.ErrorObject{&errorObject}

	w.Header().Set(headerContentType, jsonapi.MediaType)
	w.WriteHeader(code)

	err = jsonapi.MarshalErrors(w, errorObjects)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
