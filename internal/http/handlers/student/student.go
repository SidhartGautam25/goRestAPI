package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/SidhartGautam25/goRestAPI/internal/storage"
	types "github.com/SidhartGautam25/goRestAPI/internal/types/student"
	"github.com/SidhartGautam25/goRestAPI/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("craeting a student ")
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation (0 trust policy)
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully")
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int{"id": lastId})

	}

}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student ", slog.String("id", id))
		intid, err := strconv.Atoi(id)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}
		student, err := storage.GetStudentById(intid)
		if err != nil {
			slog.Error("error occured while getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)

	}
}
