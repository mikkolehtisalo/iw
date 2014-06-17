package controllers

import (
    "github.com/revel/revel"
    "iw/app/models"
    "github.com/mikkolehtisalo/revel/common"
)

type FavoriteWikis struct {
    *revel.Controller
}

// CREATE
func (c FavoriteWikis) Create(wiki string) revel.Result {
    revel.TRACE.Printf("FavoriteWikis Create(): %+v", wiki)

    // Make sure the id at least looks like one
    if !common.IsUUID(wiki) {
        revel.ERROR.Printf("Garbage favorite %+v create received from %+v", wiki, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Already exists?
    if models.IsFavoriteWiki(wiki, c.Session["username"]) {
        revel.ERROR.Printf("Wiki %+v already favorite!", wiki)
        return c.RenderText("{\"success\":false}")
    }

    fav := models.FavoriteWiki{}
    fav.Username = c.Session["username"]
    fav.Wiki_id = wiki
    fav.Status = "ACTIVE"
    fav.Save()

    return c.RenderText("{\"success\":true}")
}

// READ
func (c FavoriteWikis) Read() revel.Result {
    revel.TRACE.Printf("FavoriteWikis List()")
    return c.RenderJson(models.ListFavoriteWikis(c.Session["username"]))
}


// DELETE
func (c FavoriteWikis) Delete(wiki string) revel.Result {
    revel.TRACE.Printf("FavoriteWikis Delete(): %+v", wiki)

    // Make sure the id at least looks like one
    if !common.IsUUID(wiki) {
        revel.ERROR.Printf("Garbage favorite %+v delete received from %+v", wiki, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Already exists?
    if !models.IsFavoriteWiki(wiki, c.Session["username"]) {
        revel.ERROR.Printf("Wiki %+v not favorite!", wiki)
        return c.RenderText("{\"success\":false}")
    }

    fav := models.FavoriteWiki{}
    fav.Username = c.Session["username"]
    fav.Wiki_id = wiki
    fav.Status = "DELETED"
    fav.Save()

    return c.RenderText("{\"success\":true}")
}

