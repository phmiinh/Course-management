package controller

import (
	"first-app/database"
	"first-app/models"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LessonController(c *fiber.Ctx) error {
	courseID := c.Params("id")

	var lessons []models.Lesson

	if err := database.DB.Where("course_id = ?", courseID).Find(&lessons).Error; err != nil {
		log.Println(err)
	}

	data := fiber.Map{
		"Ctx":      c,
		"CourseID": courseID,
		"Lessons":  lessons,
	}
	return c.Render("lesson/index", data, "layouts/main")
}

func CreateLessonPostController(c *fiber.Ctx) error {
	var lesson models.Lesson

	c.BodyParser(&lesson)

	courseID, _ := strconv.Atoi(c.Params("id"))
	user := GetSessionUser(c)

	lesson.CourseID = courseID
	lesson.UserID = user.UserID

	date, err := time.Parse("2006-01-02", c.FormValue("day_start"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid date format")
	}

	newLesson := models.Lesson{
		LessonTitle:       lesson.LessonTitle,
		LessonDescription: lesson.LessonDescription,
		UserID:            user.UserID,
		CourseID:          courseID,
		StartAt:           date,
	}

	database.DB.Create(&newLesson)
	var course models.Course

	database.DB.Where("course_id = ?", courseID).First(&course)

	data := fiber.Map{
		"Ctx":               c,
		"CourseID":          courseID,
		"CourseTitle":       course.CourseTitle,
		"CourseDescription": course.CourseDescription,
		"InstructorName":    user.Name,
		"StartAt":           lesson.StartAt,
		"LessonID":          lesson.LessonID,
	}

	return c.Render("course/instructor", data, "layouts/main")
}

func LessonDetailController(c *fiber.Ctx) error {
	lessonID := c.Params("lessonID")

	var lesson models.Lesson
	var assignments []models.Assignment

	if err := database.DB.Where("lesson_id = ?", lessonID).First(&lesson).Error; err != nil {
		log.Println(err)
	}

	if err := database.DB.Where("lesson_id = ?", lessonID).Find(&assignments).Error; err != nil {
		log.Println(err)
	}

	data := fiber.Map{
		"Ctx":               c,
		"CourseID":          lesson.CourseID,
		"LessonID":          lesson.LessonID,
		"LessonTitle":       lesson.LessonTitle,
		"StartAt":           lesson.StartAt.Format("2006-01-02"),
		"LessonDescription": lesson.LessonDescription,
		"Assignments":       assignments,
	}

	return c.Render("lesson/detail", data, "layouts/main")
}

func LessonDeleteController(c *fiber.Ctx) error {
	lessonID := c.Params("lessonID")

	var lesson models.Lesson
	if err := database.DB.Where("lesson_id = ?", lessonID).First(&lesson).Error; err != nil {
		log.Println(err)
	}

	if err := database.DB.Delete(&lesson).Error; err != nil {
		log.Println(err)
	}

	var lessons []models.Lesson
	if err := database.DB.Where("course_id = ?", lesson.CourseID).Find(&lessons).Error; err != nil {
		log.Println(err)
	}

	data := fiber.Map{
		"Ctx":      c,
		"CourseID": lesson.CourseID,
		"Lessons":  lessons,
	}
	return c.Render("lesson/index", data, "layouts/main")
}
