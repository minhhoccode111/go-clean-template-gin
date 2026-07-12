package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minhhoccode111/go-clean-template-gin/internal/controller/restapi/v1/request"
	"github.com/minhhoccode111/go-clean-template-gin/internal/entity"
)

// @Summary     Show history
// @Description Show all translation history for current user
// @ID          history
// @Tags        translation
// @Produce     json
// @Success     200 {object} entity.TranslationHistory
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /translation/history [get]
func (r *V1) history(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")

		return
	}

	translationHistory, err := r.t.History(c.Request.Context(), userID.(string))
	if err != nil {
		r.l.Error(err, "restapi - v1 - history")
		errorResponse(c, http.StatusInternalServerError, "database problems")

		return
	}

	c.JSON(http.StatusOK, translationHistory)
}

// @Summary     Translate
// @Description Translate a text
// @ID          do-translate
// @Tags        translation
// @Accept      json
// @Produce     json
// @Param       request body     request.Translate true "Set up translation"
// @Success     200     {object} entity.Translation
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /translation/do-translate [post]
func (r *V1) doTranslate(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")

		return
	}

	var body request.Translate

	if err := c.ShouldBindJSON(&body); err != nil {
		r.l.Error(err, "restapi - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "restapi - v1 - doTranslate")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	translation, err := r.t.Translate(
		c.Request.Context(),
		userID.(string),
		entity.Translation{
			Source:      body.Source,
			Destination: body.Destination,
			Original:    body.Original,
		},
	)
	if err != nil {
		r.l.Error(err, "restapi - v1 - doTranslate")
		errorResponse(c, http.StatusInternalServerError, "translation service problems")

		return
	}

	c.JSON(http.StatusOK, translation)
}
