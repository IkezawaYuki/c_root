package session

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
)

const UserSession = "user_session"

const (
	Name = "c_root_id"
)

func GetLoginSession(ctx *gin.Context, store sessions.Store) (*sessions.Session, error) {
	sess, err := store.Get(ctx.Request, Name)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func SetLoginSession(ctx *gin.Context, store sessions.Store, uid string) {
	sess, _ := store.Get(ctx.Request, Name)
	sess.Values["uid"] = uid
	err := sess.Save(ctx.Request, ctx.Writer)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
