package controllers

import (
    "github.com/revel/revel"
    . "github.com/mikkolehtisalo/revel/common"
    "iw/app/models"
    "github.com/mikkolehtisalo/revel/acl"
    "encoding/json"
)

type ContentFields struct {
    *revel.Controller
}

// READ
func (c ContentFields) Read(wiki string, page string) revel.Result {
    revel.TRACE.Printf("ContentFields Read() wiki: %+v, page: %+v", wiki, page)

    // Make sure the ids at least look like one
    if (!IsUUID(wiki)) || (!IsUUID(page)) {
        revel.ERROR.Printf("Garbage contentfield %+v/%+v received from %+v", wiki, page, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Get the contentfield
    cf := models.GetContent(wiki, page, c.Args)
    // Check the ACL
    filtered := acl.Filter(c.Args, []string{"read", "write", "admin"}, []models.ContentField{cf}, true)

    if len(filtered) < 1 {
        revel.ERROR.Printf("Unable to read content field! wiki: %+v, page: %+v, user: %+v", wiki, page, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    return c.RenderJson(filtered[0])
}

// UPDATE
func (c ContentFields) Update(wiki string, page string) revel.Result {

    revel.TRACE.Printf("ContentFields Update() wiki: %+v, page: %+v", wiki, page)

    // Make sure the ids at least look like one
    if (!IsUUID(wiki)) || (!IsUUID(page)) {
        revel.ERROR.Printf("Garbage contentfield %+v/%+v received from %+v", wiki, page, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Make sure the contentfield pre-exists
    old_cf := models.GetContent(wiki, page, c.Args)
    if old_cf.Wiki_id != wiki && old_cf.Contentfield_id != page {
        revel.ERROR.Printf("Attempt to update non-existing contentfield %+v/%+v by user %+v", wiki, page, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Decode contentfield from json input
    var new_cf models.ContentField
    decoder := json.NewDecoder(c.Request.Body)
    err := decoder.Decode(&new_cf)
    if err != nil {
        revel.ERROR.Printf("Unable to parse contentfield %+v/%+v: %+v", wiki, page, err)
        return c.RenderText("{\"success\":false}")
    }

    // ID fields must match!
    if (new_cf.Wiki_id != wiki) || (new_cf.Contentfield_id != page) {
        revel.ERROR.Printf("Contentfield id mismatch %+v/%+v != %+v/%+v from user %+v", wiki, page, new_cf.Wiki_id, new_cf.Contentfield_id, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Validate fields
    new_cf.Validate(c.Validation)
    if c.Validation.HasErrors() {
        revel.ERROR.Printf("Validation errors: %+v", c.Validation.ErrorMap())
        return c.RenderText("{\"success\":false}")
    }

    // Make sure the user has rights to the page
    filtered := acl.Filter(c.Args, []string{"admin","write"}, []models.ContentField{new_cf}, true)
    if len(filtered) != 1 {
        revel.ERROR.Printf("Attempt to update contentfield %+v/%+v without access rights! user: %+v", wiki, page, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    new_cf.Create_user = c.Session["username"]
    new_cf.Status = "ACTIVE"
    new_cf.Save(true)

    return c.RenderText("{\"success\":true}")
}

