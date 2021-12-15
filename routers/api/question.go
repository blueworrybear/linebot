package api

import (
	"strconv"

	"github.com/blueworrybear/lineBot/core"
	"github.com/gin-gonic/gin"
)

func HandlePostQuestions(questionStore core.QuestionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		question := &core.Question{}
		if err := c.BindJSON(question); err != nil {
			c.String(400, err.Error())
			return
		}
		if err := questionStore.Create(question); err != nil {
			c.String(400, err.Error())
			return
		}
		c.JSON(200, question)
	}
}

func HandleGetQuestionsAsk(userStore core.UserStore, questionStore core.QuestionStore, questionService core.QuestionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(400, err.Error())
			return
		}
		q, err := questionStore.Find(&core.Question{ID: id})
		if err != nil {
			c.String(400, err.Error())
			return
		}
		bot := MustGetLineBotFrom(c)
		users, err := userStore.VIPs()
		if err != nil {
			c.String(400, err.Error())
			return
		}
		for _, u := range users {
			if err := questionService.Ask(bot, u, q); err != nil {
				c.String(400, err.Error())
				return
			}
		}
	}
}
