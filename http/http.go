package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	pa "github.com/Lambels/patrickarvatu.com"
)

// error code -> http error code.
var httpCodes = map[string]int{
	pa.ENOTFOUND:       http.StatusNotFound,
	pa.ECONFLICT:       http.StatusConflict,
	pa.EINTERNAL:       http.StatusInternalServerError,
	pa.EINVALID:        http.StatusBadRequest,
	pa.ENOTIMPLEMENTED: http.StatusNotImplemented,
	pa.EUNAUTHORIZED:   http.StatusUnauthorized,
}

// getErrorCode maps the code to an http code if possible or returns 500.
func getErrorCode(code string) int {
	v, ok := httpCodes[code]
	if !ok {
		return http.StatusInternalServerError
	}
	return v
}

// errorResponse represents an http error response used by send error.
type errorResponse struct {
	Code  int    `json:"code"`
	Trace string `json:"trace"`
}

// SendError sends an error over http.
// maps system error codes to http error codes.
func SendError(w http.ResponseWriter, r *http.Request, err error) {
	code, message := pa.ErrorCode(err), pa.ErrorMessage(err)

	if code == pa.EINTERNAL {
		LogError(r, err)
	}

	v := getErrorCode(code)
	w.WriteHeader(v)
	SendJSON(w, errorResponse{
		Code:  v,
		Trace: message,
	})
}

// LogError logs an error with format '[HTTP] error: {method} {path}: {error message}'
func LogError(r *http.Request, err error) {
	log.Printf("[HTTP] error: %s %s: %s\n", r.URL.Path, r.Method, err)
}

// SendJSON sends json: data over http.
func SendJSON(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}

// response and request types ----------------------------

type getMeResponse struct {
	User   *pa.User `json:"user"`
	PfpURL string   `json:"pfpUrl"`
}

type getOtherUserResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type getMyProfileResponse struct {
	User          *pa.User                 `json:"user"`
	Comments      getCommentsResponse      `json:"comments"`
	Subscriptions getSubscriptionsResponse `json:"subscriptions"`
}

type getOtherUserProfileResponse struct {
	User     getOtherUserResponse `json:"user"`
	Comments getCommentsResponse  `json:"comments"`
}

type getCommentsResponse struct {
	N        int           `json:"n"`
	Comments []*pa.Comment `json:"comments"`
}

type getSubscriptionsResponse struct {
	N             int `json:"n"`
	Subscriptions []struct {
		ID     int `json:"id"`
		UserID int `json:"userID"`
		On     int `json:"on"`
	} `json:"subscriptions"`
}

// TODO: test
func (s *getSubscriptionsResponse) serializeIn(subs ...*pa.Subscription) {
	// loop over each subscription.
	for _, sub := range subs {
		// filter each subscription by type so we can identify each payload
		switch sub.Topic {
		case pa.EventTopicNewSubBlog:
			payload := sub.Payload.(pa.SubBlogPayload)
			s.Subscriptions = append(s.Subscriptions, struct {
				ID     int `json:"id"`
				UserID int `json:"userID"`
				On     int `json:"on"`
			}{
				ID:     sub.ID,
				UserID: sub.ID,
				On:     payload.BlogID,
			})

		case pa.EventTopicNewComment:
			payload := sub.Payload.(pa.CommentPayload)
			s.Subscriptions = append(s.Subscriptions, struct {
				ID     int `json:"id"`
				UserID int `json:"userID"`
				On     int `json:"on"`
			}{
				ID:     sub.ID,
				UserID: sub.ID,
				On:     payload.SubBlogID,
			})

		}
	}
}

type getSubBlogsResponse struct {
	N        int           `json:"n"`
	SubBlogs []*pa.SubBlog `json:"subBlogs"`
}

type getBlogsResponse struct {
	N     int        `json:"n"`
	Blogs []*pa.Blog `json:"blogs"`
}
