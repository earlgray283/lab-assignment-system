package server

import (
	"lab-assignment-system-backend/repository"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
)

type GradeJson struct {
	UnitSum       int          `json:"unitSum,omitempty"`
	GpSum         int          `json:"gpSum,omitempty"`
	Gpa           float64      `json:"gpa,omitempty"`
	Grades        []ChildGrade `json:"grades,omitempty"`
	StudentName   string       `json:"studentName,omitempty"`
	StudentNumber int          `json:"studentNumber,omitempty"`
}

type ChildGrade struct {
	UnitNum    int     `json:"unitNum,omitempty"`
	Gp         float64 `json:"gp,omitempty"`
	ReportedAt string  `json:"reportedAt,omitempty"`
}

func (srv *Server) GradesRouter() {
	gradesRouter := srv.r.Group("/grades")
	{
		gradesRouter.POST("/", srv.HandlePostGrade())
		gradesRouter.GET("/", func(ctx *gin.Context) {
			// TODO
		})
	}
}

func (srv *Server) HandlePostGrade() gin.HandlerFunc {
	return func(c *gin.Context) {
		var gradeJson GradeJson
		if err := c.BindJSON(&gradeJson); err != nil {
			srv.logger.Println(err)
			return
		}
		key := repository.NewGradesKey(gradeJson.StudentNumber)
		var grade repository.Grade
		err := srv.dc.Get(c, key, &grade)
		if err == nil {
			// 既に成績が存在する場合は処理を行わない
			// 上書きしたい場合は delete リクエストを送ってもらう
			srv.logger.Println("Duplicate post grades request")
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		if err != datastore.ErrNoSuchEntity {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if _, err := srv.dc.Put(c, key, &grade); err != nil {
			srv.logger.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
