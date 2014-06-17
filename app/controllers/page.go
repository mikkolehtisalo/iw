package controllers

import (
    "github.com/revel/revel"
    "iw/app/models"
    "encoding/json"
    "strings"
    "regexp"
    . "github.com/mikkolehtisalo/revel/common"
    "github.com/mikkolehtisalo/revel/acl"
)

type Pages struct {
    *revel.Controller
}

// CREATE
func (c Pages) Create(wiki string, page string) revel.Result {
    revel.TRACE.Printf("Pages Create() wiki: %+v, page: %+v", wiki, page)

    // Make sure the ids at least look like one
    if (!IsUUID(wiki)) || (!IsUUID(page)) {
        revel.ERROR.Printf("Garbage page %+v/%+v received from %+v", wiki, page, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Make sure the page doesn't pre-exist
    exists_test := models.GetPage(wiki, page)
    if exists_test.Wiki_id == wiki && exists_test.Page_id == page {
        revel.ERROR.Printf("Attempt to rewrite pre-existing page %+v/%+v by user %+v", wiki, page, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Make sure the wiki exists!
    wiki_exists := models.GetWiki(wiki)
    if wiki_exists.Wiki_id != wiki {
        revel.ERROR.Printf("Attempt to add page to non-existent wiki: %+v/%+v", wiki, page)
        return c.RenderText("{\"success\":false}")
    }

    // Decode page from json input
    var new_page models.Page
    decoder := json.NewDecoder(c.Request.Body)
    err := decoder.Decode(&new_page)
    if err != nil {
        revel.ERROR.Printf("Unable to parse page %+v/%+v: %+v", wiki, page, err)
        return c.RenderText("{\"success\":false}")
    }

    // ID fields must match!
    if (new_page.Wiki_id != wiki) || (new_page.Page_id != page) {
        revel.ERROR.Printf("Page id mismatch %+v/%+v != %+v/%+v from user %+v", wiki, page, new_page.Wiki_id, new_page.Page_id, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Validate fields
    new_page.Validate(c.Validation)
    if c.Validation.HasErrors() {
        revel.ERROR.Printf("Validation errors: %+v", c.Validation.ErrorMap())
        return c.RenderText("{\"success\":false}")
    }

    // Make sure the user has rights to create page here
    filtered := acl.Filter(c.Args, []string{"admin","write"}, []models.Page{new_page}, true)
    if len(filtered) != 1 {
        revel.ERROR.Printf("Attempt to create page %+v/%+v without access rights: %+v, user: %+v", wiki, page, new_page.Title, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    new_page.Create_user = c.Session["username"]
    new_page.Status = "ACTIVE"
    new_page.Save(true)

    // Create corresponding contentfield!
    cf := models.ContentField{}
    cf.Contentfield_id = new_page.Page_id
    cf.Wiki_id = new_page.Wiki_id
    cf.Content = "<p>Add new content here!</p>"
    cf.Status = "ACTIVE"
    cf.Create_user = c.Session["username"]
    cf.Save(false)

    return c.RenderText("{\"success\":true}")
}

// READ
func (c Pages) Read() revel.Result {
    node := c.Params.Values.Get("node")
    revel.TRACE.Printf("Pages Read() node: %+v", node)

    // Check that the input parameter is at least roughly valid
    re := regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})*")
    if !re.MatchString(node) {
        revel.TRACE.Printf("Pages List() invalid node: %+v", node)
        return c.RenderText("{\"success\":false}")
    }

    // This won't blow up thanks to previous
    wiki := strings.Split(node, "/")[0]
    page := strings.Split(node, "/")[1]

    // Get the list of pages
    pages := models.ListPages(wiki, page)
    // Filter using the ACL system
    filtered := acl.Filter(c.Args, []string{"read","write","admin"}, pages, true)
    revel.TRACE.Printf("Returning: %+v", filtered)

    return c.RenderJson(filtered)
}

// UPDATE
// Not DRY, but this logic might still change: revise later...
func (c Pages) Update(wiki string, page string) revel.Result {
    revel.TRACE.Printf("Pages Update() wiki: %+v, page: %+v", wiki, page)

    // Make sure the ids at least look like one
    if (!IsUUID(wiki)) || (!IsUUID(page)) {
        revel.ERROR.Printf("Garbage page %+v/%+v received from %+v", wiki, page, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Make sure the page pre-exists
    // Also the wiki probably exists if the page exists
    old_page := models.GetPage(wiki, page)
    if old_page.Wiki_id != wiki && old_page.Page_id != page {
        revel.ERROR.Printf("Attempt to update non-existing page %+v/%+v by user %+v", wiki, page, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Decode page from json input
    var new_page models.Page
    decoder := json.NewDecoder(c.Request.Body)
    err := decoder.Decode(&new_page)
    if err != nil {
        revel.ERROR.Printf("Unable to parse page %+v/%+v: %+v", wiki, page, err)
        return c.RenderText("{\"success\":false}")
    }

    // ID fields must match!
    if (new_page.Wiki_id != wiki) || (new_page.Page_id != page) {
        revel.ERROR.Printf("Page id mismatch %+v/%+v != %+v/%+v from user %+v", wiki, page, new_page.Wiki_id, new_page.Page_id, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    // Validate fields
    new_page.Validate(c.Validation)
    if c.Validation.HasErrors() {
        revel.ERROR.Printf("Validation errors: %+v", c.Validation.ErrorMap())
        return c.RenderText("{\"success\":false}")
    }

    // Make sure the user has rights to the page
    filtered := acl.Filter(c.Args, []string{"admin","write"}, []models.Page{new_page}, true)
    if len(filtered) != 1 {
        revel.ERROR.Printf("Attempt to create page %+v/%+v without access rights: %+v, user: %+v", wiki, page, new_page.Title, c.Session["username"])
        return c.RenderText("{\"success\":false}")
    }

    new_page.Create_user = c.Session["username"]
    new_page.Status = "ACTIVE"
    new_page.Save(true)

    return c.RenderText("{\"success\":true}")
}


// DELETE
func (c Pages) Delete(wiki string, page string) revel.Result {
    revel.TRACE.Printf("Pages Delete() wiki: %+v, page: %+v", wiki, page)

    // If the parameters don't seem valid, bail out
    if (!IsUUID(wiki)) || (!IsUUID(page)) {
        revel.ERROR.Printf("Pages Delete() invalid IDs! wiki: %+v, page: %+v", wiki, page)
        return c.RenderText("{\"success\":false}")
    }

    p := models.GetPage(wiki, page)
    if p.Wiki_id == "" {
        revel.ERROR.Printf("Pages Delete() Attempt to delete non-existing page! user: %+v wiki %+v page %+v", c.Session["username"], wiki, page)
        return c.RenderText("{\"success\":false}")
    }

    // Write or admin is enough to delete
    filtered := acl.Filter(c.Args, []string{"write", "admin"}, []models.Page{p}, true)
    if len(filtered) < 1 {
        revel.ERROR.Printf("Pages Delete() insufficient rights! user: %+v wiki: %+v, page: %+v", c.Session["username"], wiki, page)
        return c.RenderText("{\"success\":false}")
    }

    // Also handles children and contentfields
    models.DeletePage(wiki, page, c.Session["username"], true)

    return c.RenderText("{\"success\":true}")
}