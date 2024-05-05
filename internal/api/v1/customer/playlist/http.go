package playlist

import (
	"context"
	"fmt"
	"music-master/internal/model"
	"net/http"

	httputil "music-master/internal/util/http"

	"github.com/labstack/echo/v4"
)

// HTTP represents playlist http service
type HTTP struct {
	svc Service
}

// Service represents playlist application interface
type Service interface {
	Create(ctx context.Context, authUsr *model.AuthUser, data CreationData) (*model.Playlist, error)
	View(ctx context.Context, authUsr *model.AuthUser, id string) (*model.Playlist, error)
	// List(ctx context.Context, authUsr *model.AuthUser, lq *dbutil.ListQueryCondition, count *int64) (*ListLateFeeResp, error)
	Update(ctx context.Context, authUsr *model.AuthUser, id string, data UpdateData) (*model.Playlist, error)
	Delete(ctx context.Context, authUsr *model.AuthUser, id string) error
	DeleteMusicTrack(ctx context.Context, authUsr *model.AuthUser, id string, data DeleteMusicTrack) error
	Search(ctx context.Context, authUsr *model.AuthUser, lq *httputil.ListRequest) ([]*model.Playlist, error)
}

// NewHTTP creates new playlist http service
func NewHTTP(svc Service, auth model.Auth, eg *echo.Group) {
	h := HTTP{svc}

	// swagger:operation POST /v1/customer/playlists customer-playlists customerPlaylistCreate
	// ---
	// summary: Creates new playlist
	// parameters:
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CustomerPlaylistCreationData"
	// responses:
	//   "200":
	//     description: The new playlist
	//     schema:
	//       "$ref": "#/definitions/Playlist"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	eg.POST("", h.create)

	// swagger:operation GET /v1/customer/playlists/{id} customer-playlists customerPlaylistView
	// ---
	// summary: Returns a single playlist
	// parameters:
	// - name: id
	//   in: path
	//   description: id of playlist
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: The playlist
	//     schema:
	//       "$ref": "#/definitions/Playlist"
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

	// swagger:operation GET /v1/customer/playlists customer-playlists customerPlaylistSearch
	// ---
	// summary: Returns a list playlist
	// responses:
	//   "200":
	//     description: The playlists
	//     schema:
	//       "$ref": "#/definitions/CustomerPlaylistListResp"
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

	// swagger:operation PATCH /v1/customer/playlists/{id} customer-playlists customerPlaylistUpdate
	// ---
	// summary: Update playlist's information
	// parameters:
	// - name: id
	//   in: path
	//   description: id of playlist
	//   type: string
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CustomerPlaylistUpdateData"
	// responses:
	//   "200":
	//     description: The updated playlist
	//     schema:
	//       "$ref": "#/definitions/Playlist"
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

	// swagger:operation DELETE /v1/customer/playlists/{id} customer-playlists customerPlaylistDelete
	// ---
	// summary: Deletes a playlist
	// parameters:
	// - name: id
	//   in: path
	//   description: id of playlist
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

	// swagger:operation DELETE /v1/customer/playlists/music-tracks/{id} customer-playlists customerPlaylistMusicTrackDelete
	// ---
	// summary: Deletes a music track in playlist
	// parameters:
	// - name: id
	//   in: path
	//   description: id of playlist
	//   type: string
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/CustomerPlaylistDeleteMusicTrack"
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
	eg.DELETE("/music-tracks/:id", h.deleteMusicTrack)
}

// CreationData contains playlist data from json request
// swagger:model CustomerPlaylistCreationData
type CreationData struct {
	// example: My playlist
	Name string `json:"name"`
	// example: [{"id": "661ffc6c12e6a410902997b0", "title": "Em của ngày hôm qua","artist": "Sơn Tùng MTP",  "album": "Em của ngày hôm qua", "genre": "Ballad", "release_year": 2017, "duration": 189, "mp3_file": "AA=="}]
	Tracks []*model.MusicTrack `json:"tracks"`
}

// UpdateData contains playlist data from json request
// swagger:model CustomerPlaylistUpdateData
type UpdateData struct {
	Name   *string              `json:"name"`
	Tracks *[]*model.MusicTrack `json:"tracks"`
}

// DeleteMusicTrack contains playlist data from json request
// swagger:model CustomerPlaylistDeleteMusicTrack
type DeleteMusicTrack struct {
	MusicTrackID string `json:"music_track_id" validate:"required"`
}

// ListResp contains list of playlist and current page number response
// swagger:model CustomerPlaylistListResp
type ListResp struct {
	Data       []*model.Playlist `json:"data"`
	TotalCount int64             `json:"total_count"`
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

func (h *HTTP) deleteMusicTrack(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	fmt.Println("playlist id", id)
	r := DeleteMusicTrack{}
	if err := c.Bind(&r); err != nil {
		return err
	}

	if err := h.svc.DeleteMusicTrack(c.Request().Context(), nil, id, r); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
