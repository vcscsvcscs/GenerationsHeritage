package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) CreatePerson(c *gin.Context) {}

func (srv *server) SoftDeletePerson(c *gin.Context, id int, params api.SoftDeletePersonParams) {}

func (srv *server) UpdatePerson(c *gin.Context, id int, params api.UpdatePersonParams) {}

func (srv *server) HardDeletePerson(c *gin.Context, id int, params api.HardDeletePersonParams) {}
