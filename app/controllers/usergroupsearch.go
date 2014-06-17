package controllers

import "github.com/revel/revel"
import "iw/app/models"
//import "fmt"

type UserGroupSearch struct {
    *revel.Controller
}

func (c UserGroupSearch) List() revel.Result {
    query := c.Params.Query["query"][0]
    return c.RenderJson(models.ListUserGroupSearchItems(query))
}

