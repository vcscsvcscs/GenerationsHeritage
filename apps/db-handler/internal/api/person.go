package api

import "github.com/gin-gonic/gin"

func (srv *server) CreatePerson(c *gin.Context) {}

func (srv *server) SoftDeletePerson(c *gin.Context, id int, params SoftDeletePersonParams) {}

func (srv *server) UpdatePerson(c *gin.Context, id int, params UpdatePersonParams) {}

func (srv *server) HardDeletePerson(c *gin.Context, id int, params HardDeletePersonParams) {}
