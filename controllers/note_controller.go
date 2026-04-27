package controllers

import (
	"net/http"
	"note-app/models"
	"note-app/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NoteController struct {
	Repo *repository.NoteRepository
}

func (ctrl *NoteController) Index(c *gin.Context) {
	notes, _ := ctrl.Repo.GetAll()
	c.HTML(http.StatusOK, "index.html", gin.H{"notes": notes})
}

func (ctrl *NoteController) Create(c *gin.Context) {
	var note models.Note
	// Title is now optional/empty based on your previous request
	note.Title = c.PostForm("title") 
	note.Content = c.PostForm("content")
	
	_ = ctrl.Repo.Create(&note)
	c.Redirect(http.StatusSeeOther, "/")
}

// ADD THIS: Method to show the edit page
func (ctrl *NoteController) Edit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	note, err := ctrl.Repo.GetByID(uint(id))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	c.HTML(http.StatusOK, "edit.html", gin.H{"note": note})
}

// ADD THIS: Method to handle the update POST request
func (ctrl *NoteController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	note, _ := ctrl.Repo.GetByID(uint(id))
	
	note.Content = c.PostForm("content")
	
	_ = ctrl.Repo.Update(&note)
	c.Redirect(http.StatusSeeOther, "/")
}

func (ctrl *NoteController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_ = ctrl.Repo.Delete(uint(id))
	c.Redirect(http.StatusSeeOther, "/")
}
