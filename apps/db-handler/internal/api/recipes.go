package api

import "github.com/gin-gonic/gin"

func (srv *server) SoftDeleteRecipe(c *gin.Context, id int, params SoftDeleteRecipeParams) {}

func (srv *server) UpdateRecipe(c *gin.Context, id int, params UpdateRecipeParams) {}

func (srv *server) HardDeleteRecipe(c *gin.Context, id int, params HardDeleteRecipeParams) {}
