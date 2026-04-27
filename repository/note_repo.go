package repository

import (
	"note-app/models"
	"gorm.io/gorm"
)

type NoteRepository struct {
	DB *gorm.DB
}

func (r *NoteRepository) GetAll() ([]models.Note, error) {
	var notes []models.Note
	err := r.DB.Order("id desc").Find(&notes).Error
	return notes, err
}

func (r *NoteRepository) Create(note *models.Note) error {
	return r.DB.Create(note).Error
}

func (r *NoteRepository) Update(note *models.Note) error {
	return r.DB.Save(note).Error
}

func (r *NoteRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Note{}, id).Error
}

func (r *NoteRepository) GetByID(id uint) (models.Note, error) {
	var note models.Note
	err := r.DB.First(&note, id).Error
	return note, err
}
