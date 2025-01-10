package storage

import types "github.com/SidhartGautam25/goRestAPI/internal/types/student"

type Storage interface {
	CreateStudent(name string, email string, age int) (int, error)
	GetStudentById(id int) (types.Student, error)
}
