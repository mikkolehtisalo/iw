package controllers

import (
    "github.com/revel/revel"
    "iw/app/models"
    "encoding/json"
    . "github.com/mikkolehtisalo/revel/common"
     "github.com/mikkolehtisalo/revel/acl"
)

type Wikis struct {
    *revel.Controller
}

// CREATE
func (c Wikis) Create(wiki string) revel.Result {
    revel.TRACE.Printf("Wikis Create(): %+v", wiki)

    // Make sure the id at least looks like one
    if !IsUUID(wiki) {
        revel.ERROR.Printf("Garbage wiki %+v received from %+v", wiki, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Make sure the wiki doesn't pre-exist
    exists_test := models.GetWiki(wiki)
    if exists_test.Wiki_id == wiki {
        revel.ERROR.Printf("Attempt to rewrite pre-existing wiki %+v by user %+v", wiki, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Decode wiki from input json
    var new_wiki models.Wiki
    decoder := json.NewDecoder(c.Request.Body)
    err := decoder.Decode(&new_wiki)
    if err != nil {
        revel.ERROR.Printf("Unable to parse wiki %+v: %+v", wiki, err)
        return c.RenderText("{\"success\":false}")
    }

    // ID fields must match!
    if new_wiki.Wiki_id != wiki {
        revel.ERROR.Printf("Wiki id mismatch %+v != %+v", new_wiki.Wiki_id, wiki)
        return c.RenderText("{\"success\":false}")
    }

    // Validate fields
    new_wiki.Validate(c.Validation)
    if c.Validation.HasErrors() {
        revel.ERROR.Printf("Validation errors parsing wiki %+v: %+v", wiki, c.Validation.ErrorMap())
        return c.RenderText("{\"success\":false}")
    }

    // Make user the author has admin access right by default
    AddUserToACLList(c.Session["username"], &new_wiki.Adminacl)

    // Save the wiki
    new_wiki.Create_user = c.Session["username"]
    new_wiki.Status = "ACTIVE"
    new_wiki.Save(true)

    revel.INFO.Printf("User %+v created wiki %+v: %+v", c.Session["username"], new_wiki.Wiki_id, new_wiki.Title )
    return c.RenderText("{\"success\":true}")
}

// READ
func (c Wikis) Read() revel.Result {
    revel.TRACE.Printf("Wikis Read()")

    wikis := models.ListWikis()
    filtered := acl.Filter(c.Args, []string{"read","write","admin"}, wikis, false)

    revel.TRACE.Printf("Wikis Read() returning: %+v", filtered)
    return c.RenderJson(filtered)
}


// UPDATE
func (c Wikis) Update(wiki string) revel.Result {
    revel.TRACE.Printf("Wikis Update(): %s", wiki)

    // Make sure the id at least looks like one
    if !IsUUID(wiki) {
        revel.ERROR.Printf("Garbage wiki %+v received from %+v", wiki, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Make sure the wiki exists
    exists_test := models.GetWiki(wiki)
    if exists_test.Wiki_id != wiki {
        revel.ERROR.Printf("Attempt to update non-existing wiki %+v by user %+v", wiki, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Decode the wiki from input
    var new_wiki models.Wiki
    decoder := json.NewDecoder(c.Request.Body)
    err := decoder.Decode(&new_wiki)
    if err != nil {
        revel.ERROR.Printf("Unable to parse wiki %+v: %+v", wiki, err)
        return c.RenderText("{\"success\":false}")
    }

    // ID fields must match!
    if new_wiki.Wiki_id != wiki {
        revel.ERROR.Printf("Wiki id mismatch %+v != %+v", new_wiki.Wiki_id, wiki)
        return c.RenderText("{\"success\":false}")
    }

    // Validate fields
    new_wiki.Validate(c.Validation)
    if c.Validation.HasErrors() {
        revel.ERROR.Printf("Validation errors: %+v", c.Validation.ErrorMap()) 
        return c.RenderText("{\"success\":false}")
    }

    // Make sure the user has rights to modify the wiki
    filtered := acl.Filter(c.Args, []string{"admin","write"}, []models.Wiki{exists_test}, false)
    if len(filtered) != 1 {
        revel.ERROR.Printf("Attempt to update wiki without access rights: %+v: %+v, user: %+v", exists_test.Wiki_id, exists_test.Title, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    new_wiki.Status = "ACTIVE"
    new_wiki.Create_user = c.Session["username"]
    new_wiki.Save(true)

    revel.INFO.Printf("User %+v updated wiki %+v: %+v", c.Session["username"], new_wiki.Wiki_id, new_wiki.Title )

    return c.RenderText("{\"success\":true}")
}

// DELETE
func (c Wikis) Delete(wiki string) revel.Result {
    revel.TRACE.Printf("Wikis Delete(): %s", wiki)
    wi := models.GetWiki(wiki)
    filtered := acl.Filter(c.Args, []string{"admin"}, []models.Wiki{wi}, false)

    // Delete everything that survived filtering
    for _, w := range filtered {
        // Will also do other house cleaning
        w.(models.Wiki).Delete(c.Session["username"])
        revel.INFO.Printf("User %+v deleted wiki %+v: %+v", c.Session["username"], w.(models.Wiki).Wiki_id, w.(models.Wiki).Title)
    }

    return c.RenderText("{\"success\":true}")
}
