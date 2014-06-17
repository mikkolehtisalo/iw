package controllers

import (
    "github.com/revel/revel"
    "iw/app/models"
    "github.com/mikkolehtisalo/revel/acl"
)

type Activities struct {
    *revel.Controller
}

// READ
func (a Activities) Read() revel.Result {
    user := a.Session["username"]
    revel.TRACE.Printf("Activities Read() user: %+v", user)

    activities := models.GetActivities()
    filtered := acl.Filter(a.Args, []string{"read", "admin","write"}, activities, false)

    revel.TRACE.Printf("Activities Read() returning: %+v", filtered)
    return a.RenderJson(filtered)
}
