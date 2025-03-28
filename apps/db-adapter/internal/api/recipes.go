package api

import (
	"github.com/gin-gonic/gin"
	"github.com/vcscsvcscs/GenerationsHeritage/apps/db-adapter/pkg/api"
)

func (srv *server) SoftDeleteRecipe(c *gin.Context, id int, params api.SoftDeleteRecipeParams) {}

func (srv *server) UpdateRecipe(c *gin.Context, id int, params api.UpdateRecipeParams) {}

func (srv *server) HardDeleteRecipe(c *gin.Context, id int, params api.HardDeleteRecipeParams) {}
