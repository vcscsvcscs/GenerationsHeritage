package api

import "github.com/gin-gonic/gin"

func (srv *server) GetPersonByGoogleId(c *gin.Context, googleId string) {}

func (srv *server) CreatePersonByGoogleIdAndInviteCode(c *gin.Context, googleId string) {}

func (srv *server) CreatePersonByGoogleId(c *gin.Context, googleId string) {}
