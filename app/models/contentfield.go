package models

import (
    "github.com/revel/revel"
    "github.com/mikkolehtisalo/revel/deXSS"
    "time"
    "regexp"
    "github.com/mikkolehtisalo/revel/acl"
    "strings"
    "fmt"
)

var (
    allowed map[string]string
)

func init() {
    allowed = make(map[string]string)
    // This is actually what most basic editing functions of CKEditor require
    allowed["p"] = "class,id"
    allowed["div"] = "class,id"
    allowed["h1"] = "class,id"
    allowed["h2"] = "class,id"
    allowed["h3"] = "class,id"
    allowed["ul"] = "class,id"
    allowed["li"] = "class,id"
    allowed["a"] = "class,id,href,rel"
    allowed["img"] = "class,id,src,alt,hspace,vspace,width,height"
    allowed["span"] = "class,id,style"
}

type ContentField struct {
    Contentfield_id string
    Wiki_id string
    Content string
    Contentwithmacros string
    Modified time.Time
    Status string
    Create_user string
    MatchedPermissions []string
    // Not really used, but required for handling the acls
    Readacl string
    Writeacl string
    Adminacl string
}

func (c *ContentField) Validate(v *revel.Validation) {
    // Required fields
    v.Required(c.Contentfield_id)
    v.Required(c.Wiki_id)

    // Match against regexp patterns
    v.Match(c.Contentfield_id, regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")).Message("Contentfield_id not UUID?")
    v.Match(c.Wiki_id, regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")).Message("Wiki_id not UUID?")

    // Clean the HTML
    c.Content = deXSS.FilterHTML(c.Content, allowed, true)
}

func toc_page(p Page, mode string, co map[string]interface {}) string {
    toc := ""

    level := p.Depth + 1
    children := Get_children(p)
    filtered := acl.Filter(co, []string{"read", "write", "admin"}, children, true)
    var link string

    for _, child := range filtered {
        x := child.(Page)
        if mode == "anchor" {  
            link = "<a href='#" + x.Page_id + "' class='extpageanchor'>" + x.Title + "</a>"
        } else {
            parent := ""
            if x.Depth > 0 {
                split := strings.Split(x.Path, "/")
                parent = split[len(split)-2]
            }
            link = "<a href='#' class='extpagelink' id='" + x.Wiki_id + "/" + parent + "/" + x.Page_id + "'>" + x.Title + "</a>"
        }

        toc = toc + "<li class='level-" + fmt.Sprintf("%v", level) + "'>" + link + "</li>"

        if len(Get_children(x)) > 0 {
            toc = toc + toc_page(x, mode, co)
        }
    }

    return toc
}

// Mode: either "anchor" or "link"
func (c *ContentField) build_toc(mode string, co map[string]interface {}) string {
    toc := "<div class='toc'><h1>Table of contents</h1><ul>\n"
    page := GetPage(c.Wiki_id, c.Contentfield_id)
    toc = toc + toc_page(page, mode, co)
    toc = toc + "</ul></div>"
    return toc
}

func include_page(p Page, co map[string]interface {}) string {
    include := ""

    level := p.Depth
    children := Get_children(p)
    filtered := acl.Filter(co, []string{"read", "write", "admin"}, children, true)

    for _, child := range filtered {
        p := child.(Page)
        content := GetContent(p.Wiki_id, p.Page_id, co)
        item := "<h" + fmt.Sprintf("%v", level + 1) + " class='extpagetarget' id='" + p.Page_id + "'>" + p.Title +"</h" + fmt.Sprintf("%v", level + 1) + ">\n" + content.Content + "\n"
        include = include + item
        if len(Get_children(p)) > 0 {
            include = include + include_page(p, co)
        }
    }

    return include
}

// Sets up the recursion and fires away
func (c *ContentField) build_include(co map[string]interface {}) string {
    include := ""
    page := GetPage(c.Wiki_id, c.Contentfield_id)
    include = include + include_page(page, co)
    return include
}

func (c *ContentField) run_macros(co map[string]interface {}) {
    revel.TRACE.Printf("ContentField run_macros()")
    c.Contentwithmacros = c.Content

    toc_anchor_regexp := regexp.MustCompile("(?m)^<p>::toc_anchor(.*)$")
    toc_link_regexp := regexp.MustCompile("(?m)^<p>::toc_link(.*)$")
    children_regexp := regexp.MustCompile("(?m)^<p>::children(.*)$")

    if toc_anchor_regexp.MatchString(c.Content) {
        revel.TRACE.Printf("TOC anchor macro found!")
        toc := c.build_toc("anchor", co)
        c.Contentwithmacros = toc_anchor_regexp.ReplaceAllString(c.Contentwithmacros, toc)
    }

    if toc_link_regexp.MatchString(c.Content) {
        revel.TRACE.Printf("TOC link macro found!")
        toc := c.build_toc("link", co)
        c.Contentwithmacros = toc_link_regexp.ReplaceAllString(c.Contentwithmacros, toc)
    }

    if children_regexp.MatchString(c.Content) {
        revel.TRACE.Printf("Include children macro found!")
        include := c.build_include(co)
        c.Contentwithmacros = children_regexp.ReplaceAllString(c.Contentwithmacros, include)
    }
}

// Updates modified time
func (c *ContentField) Save(save_activity bool) {
    revel.TRACE.Printf("ContentField Save(): %+v", c)
    db := get_db()
    defer db.Close()

    c.Modified = time.Now()

    _, err := db.Exec("insert into contentfields(contentfield_id, wiki_id, content, modified, status, create_user) values ($1, $2, $3, $4, $5, $6)",
            c.Contentfield_id, c.Wiki_id, c.Content, c.Modified, c.Status, c.Create_user )

    if err != nil {
        revel.ERROR.Printf("ContentField Save(): failed with %+v", err)
    }

    if save_activity {
        // Save the activity just if this was page anyways
        p := GetPageAllStatuses(c.Wiki_id, c.Contentfield_id)
        SaveActivity(&p)
    }
}

func DeleteContentFields(wiki_id string, user string) {
    revel.TRACE.Printf("DeleteContentFields(): wiki: %+v, user: %+v", wiki_id, user)

    contents := []ContentField{}
    db := get_db()
    defer db.Close()

    err := db.Select(&contents, "select * from contentfields c1 where wiki_id=uuid_in($1) and status='ACTIVE' and not exists (select * from contentfields c2 where c1.contentfield_id=c2.contentfield_id and c1.wiki_id=c2.wiki_id and c2.modified>c1.modified)",
        wiki_id)

    if err != nil {
        revel.ERROR.Printf("Unable to get contents %+v: %+v", wiki_id,  err)
    }

    for _, item := range contents {
        item.Status = "DELETED"
        item.Create_user = user
        item.Save(false)
    }
}

func GetContent(wiki_id string, content_id string, c map[string]interface {}) ContentField {
    revel.TRACE.Printf("GetContent() wiki: %+v, content: %+v", wiki_id, content_id)
    contents := []ContentField{}
    content := ContentField{}
    db := get_db()
    defer db.Close()

    err := db.Select(&contents, "select * from contentfields c1 where contentfield_id=uuid_in($1) and wiki_id=uuid_in($2) and status='ACTIVE' and not exists (select * from contentfields c2 where c1.contentfield_id=c2.contentfield_id and c1.wiki_id=c2.wiki_id and c2.modified>c1.modified)",
        content_id, wiki_id)

    if err != nil {
        revel.ERROR.Printf("Unable to get content %+v/%+v: %+v", wiki_id, content_id, err)
    }

    if len(contents)>0 {
        content = contents[0]
    }

    content.run_macros(c)
    
    revel.TRACE.Printf("GetContent() returning: %+v", content)
    return content
}


// ACL stuff

// Build ACL entry for reference
func (c ContentField) BuildACLEntry(reference string) acl.ACLEntry {
    revel.TRACE.Printf("BuildACLEntry() %+v", reference)
    var entry acl.ACLEntry

    if reference != ("cf:"+c.Wiki_id+"/"+c.Contentfield_id) {

        if strings.Index(reference, "page") == 0 {
            // We are not working on this copy, get from database
            re := regexp.MustCompile("page:([^/]*)/(.*)")
            m := re.FindStringSubmatch(reference)
            wref := m[1]
            pref := m[2]
            pa := GetPage(wref, pref)
            entry = entry_helper(pa.Readacl, pa.Writeacl, pa.Adminacl, reference, pa)
        } 

        if strings.Index(reference, "wiki") == 0 {
            // This must be wiki!
            re := regexp.MustCompile("wiki:(.*)")
            ref := re.FindStringSubmatch(reference)[1]
            wi := GetWiki(ref)
            entry = entry_helper(wi.Readacl, wi.Writeacl, wi.Adminacl, reference, wi)
        }

    } else {
        // It's exactly the originating item!
        entry = entry_helper(c.Readacl, c.Writeacl, c.Adminacl, reference, c)
    }

    revel.TRACE.Printf("BuildACLEntry() returning %+v", entry)

    return entry
}

// Set the matched permissions to a variable
func (c ContentField) SetMatched(permissions []string) interface{} {
    c.MatchedPermissions = permissions
    return c
}

// Building parent information
func (c ContentField) BuildACLParent() string {
    return "page:" + c.Wiki_id + "/" + c.Contentfield_id
}

// Wiki+contentfield ids...
func (c ContentField) BuildACLReference() string {
    return "cf:"+c.Wiki_id+"/"+c.Contentfield_id
}

// Always inherit the parent!
func (c ContentField) BuildACLInheritation() bool {
    return true
}
