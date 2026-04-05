package domain

import (
	"errors"
	"time"
)

// Все ID uuid7

type CreateClass struct {
	ClassName string
	// Лет обучается
	YearOfStudy uint32
	// Год окончания, если не установлен то будет по считан:
	// year_of_study + текущий год
	GraduationYear *uint32
}

func (cc *CreateClass) Validate() error {
	if cc.ClassName == "" {
		return errors.New("The class name cannot be empty.")
	}

	if cc.GraduationYear != nil && *cc.GraduationYear <= uint32(time.Now().Year()) {
		return errors.New("Not a possible graduation year")
	}
	return nil
}

type Class struct {
	ID             ID
	ClassName      string
	YearOfStudy    uint32
	GraduationYear uint32
}

func (c *Class) Validate() error {
	if err := c.ID.Validate(); err != nil {
		return err
	}
	if c.ClassName == "" {
		return errors.New("The class name cannot be empty.")
	}
	if c.GraduationYear <= uint32(time.Now().Year()) {
		return errors.New("Not a possible graduation year")
	}
	return nil
}

type ListTeacherClasses struct {
	TeacherID ID
	Limit     uint32
	Offset    uint32
}

func (ltc *ListTeacherClasses) Validate() error {
	if err := ltc.TeacherID.Validate(); err != nil {
		return err
	}
	return nil
}

type ListClasses struct {
	TotalCount uint32
	Classes    []Class
}

type UpdateClass struct {
	ID             ID
	ClassName      *string
	YearOfStudy    *uint32
	GraduationYear *uint32
}
