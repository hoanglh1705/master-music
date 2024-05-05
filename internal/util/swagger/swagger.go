package swagger

import (
	httputil "music-master/internal/util/http"
	"music-master/internal/util/server"
)

// Success empty response
// swagger:response ok
type swaggOKResp struct{}

// Error empty response
// swagger:response err
type swaggErrResp struct{}

// Error response with details
// swagger:response errDetails
type swaggErrDetailsResp struct {
	//in: body
	Body server.ErrorResponse
}

// ListRequest holds data of listing request
// swagger:parameters customerMusicTrackSearch customerPlaylistSearch
type ListRequest struct {
	httputil.ListRequest
}
