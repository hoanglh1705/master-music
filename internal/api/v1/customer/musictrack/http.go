package musictrack

import (
	"context"
	"music-master/internal/model"
	"net/http"

	httputil "music-master/internal/util/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents music track http service
type HTTP struct {
	svc Service
}

// Service represents music track application interface
type Service interface {
	Create(ctx context.Context, authUsr *model.AuthUser, data CreationData) (*model.MusicTrack, error)
	View(ctx context.Context, authUsr *model.AuthUser, id string) (*model.MusicTrack, error)
	Search(ctx context.Context, authUsr *model.AuthUser, lq *httputil.ListRequest) ([]*model.MusicTrack, error)
	// List(ctx context.Context, authUsr *model.AuthUser, lq *dbutil.ListQueryCondition, count *int64) (*ListLateFeeResp, error)
	Update(ctx context.Context, authUsr *model.AuthUser, id string, data UpdateData) (*model.MusicTrack, error)
	Delete(ctx context.Context, authUsr *model.AuthUser, id string) error
}

// NewHTTP creates new music track http service
func NewHTTP(svc Service, auth model.Auth, eg *echo.Group) {
	h := HTTP{svc}

	// swagger:operation POST /v1/customer/music-tracks customer-musictracks customerMusicTrackCreate
	// ---
	// summary: Creates new music track
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CustomerMusicTrackCreationData"
	// responses:
	//   "200":
	//     description: The new music track
	//     schema:
	//       "$ref": "#/definitions/MusicTrack"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.POST("", h.create)

	// swagger:operation GET /v1/customer/music-tracks/{id} customer-musictracks customerMusicTrackView
	// ---
	// summary: Returns a single music track
	// parameters:
	// - name: id
	//   in: path
	//   description: id of music track
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: The music track
	//     schema:
	//       "$ref": "#/definitions/MusicTrack"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "404":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.GET("/:id", h.view)

	// swagger:operation GET /v1/customer/music-tracks customer-musictracks customerMusicTrackSearch
	// ---
	// summary: Returns a list music track
	// responses:
	//   "200":
	//     description: The music track
	//     schema:
	//       "$ref": "#/definitions/CustomerMusicTrackListResp"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "404":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.GET("", h.search)

	// swagger:operation PATCH /v1/customer/music-tracks/{id} customer-musictracks customerMusicTrackUpdate
	// ---
	// summary: Update music track's information
	// parameters:
	// - name: id
	//   in: path
	//   description: id of music track
	//   type: string
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CustomerMusicTrackUpdateData"
	// responses:
	//   "200":
	//     description: The updated music track
	//     schema:
	//       "$ref": "#/definitions/MusicTrack"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "404":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.PATCH("/:id", h.update)

	// swagger:operation DELETE /v1/customer/music-tracks/{id} customer-musictracks customerMusicTrackDelete
	// ---
	// summary: Deletes a music track
	// parameters:
	// - name: id
	//   in: path
	//   description: id of music track
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     "$ref": "#/responses/ok"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "404":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.DELETE("/:id", h.delete)
}

// CreationData contains music track data from json request
// swagger:model CustomerMusicTrackCreationData
type CreationData struct {
	// example: Em của ngày hôm qua
	Title string ` json:"title" validate:"required"`
	// example: Sơn Tùng MTP
	Artist string `json:"artist" validate:"required"`
	// example: Em của ngày hôm qua
	Album string `json:"album" validate:"required"`
	// example: Ballad
	Genre string `json:"genre"`
	// example: 2017
	ReleaseYear int `json:"release_year"`
	// example: 189
	Duration int    `json:"duration"` // Duration in seconds
	MP3File  []byte `json:"mp3_file"` // Binary data of the MP3 file
}

// UpdateData contains music track data from json request
// swagger:model CustomerMusicTrackUpdateData
type UpdateData struct {
	Title       *string `json:"title,omitempty"`
	Artist      *string `json:"artist,omitempty"`
	Album       *string `json:"album,omitempty"`
	Genre       *string `json:"genre,omitempty"`
	ReleaseYear *int    `json:"release_year,omitempty"`
	Duration    *int    `json:"duration,omitempty"`
	MP3File     *[]byte `json:"mp3_file,omitempty"`
}

// ListResp contains list of music track and current page number response
// swagger:model CustomerMusicTrackListResp
type ListResp struct {
	Data       []*model.MusicTrack `json:"data"`
	TotalCount int64               `json:"total_count"`
}

func (h *HTTP) create(c echo.Context) error {
	r := CreationData{}
	if err := c.Bind(&r); err != nil {
		return err
	}

	resp, err := h.svc.Create(c.Request().Context(), nil, r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) view(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	resp, err := h.svc.View(c.Request().Context(), nil, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) search(c echo.Context) error {
	lr := &httputil.ListRequest{}
	if err := c.Bind(lr); err != nil {
		return err
	}

	resp, err := h.svc.Search(c.Request().Context(), nil, lr)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ListResp{
		Data: resp,
	})
}

func (h *HTTP) update(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	u := UpdateData{}
	if err := c.Bind(&u); err != nil {
		return err
	}

	usr, err := h.svc.Update(c.Request().Context(), nil, id, u)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	if err := h.svc.Delete(c.Request().Context(), nil, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
