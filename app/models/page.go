package models

import (
    "time"
    "github.com/revel/revel"
    "github.com/mikkolehtisalo/revel/acl"
    "regexp"
    "html"
    "strings"
    "fmt"
    //. "github.com/mikkolehtisalo/revel/common"
)

type Page struct {
    Page_id string
    Wiki_id string
    Path string
    Title string
    Create_user string
    Readacl string
    Writeacl string
    Adminacl string
    Stopinheritation bool
    Index int
    Depth int
    Status string
    Modified time.Time
    MatchedPermissions []string
    Loaded bool `json:"loaded"` 
}

// Sets modified time
// Updates depth automatically
func (p *Page) Save(save_activity bool) {
    revel.TRACE.Printf("Page Save(): %+v", p)
    db := get_db()
    defer db.Close()

    // Update depth
    update_depth(p)
    p.Modified = time.Now()
    
    _, err := db.Exec("insert into pages(page_id, wiki_id, path, title, create_user, readacl, writeacl, adminacl, stopinheritation, index, depth, modified, status) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
        p.Page_id, p.Wiki_id, p.Path, p.Title, p.Create_user, p.Readacl, p.Writeacl, p.Adminacl, p.Stopinheritation, p.Index, p.Depth, p.Modified, p.Status )

    if err != nil {
        revel.ERROR.Printf("Page Save(): failed with %+v", err)
    }

    if save_activity {
        SaveActivity(p)
    }
}

// Wiki id, one page id, path
// The pages in path must either exist in database or match given single id (in case target is not already in database)
func valid_path(w string, me string, p string) bool {
    valid := true
    split_path := strings.Split(p, "/")

    for _, part := range split_path {
        x := GetPage(w, part)
        if part != me && x.Page_id != part {
            valid = false
        }
    }

    return valid
}

func (p *Page) Validate(v *revel.Validation) {
    // Required fields
    v.Required(p.Page_id)
    v.Required(p.Wiki_id)
    v.Required(p.Path)
    v.Required(p.Title)

    // Match against regexp patterns
    v.Match(p.Wiki_id, regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")).Message("Wiki_id not UUID?")
    v.Match(p.Page_id, regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")).Message("Page_id not UUID?")

    // Escape HTML from the fields that might be rendered to users
    p.Title = html.EscapeString(p.Title)

    // Validate Path
    v.Required(valid_path(p.Wiki_id, p.Page_id, p.Path)).Message(fmt.Sprintf("Path is probably invalid: %+v", p.Path))
}

func has_children(page Page) bool {
    result := false
    pages := Get_children(page)
    if len(pages) > 0 {
        result = true
    }

    return result
}

// Only the next level
func Get_children(page Page) []Page {
    revel.TRACE.Printf("Page Get_children(): %+v", page)
    pages := []Page{}
    db := get_db()
    defer db.Close()
    
    path := page.Path + "/%"
    db.Select(&pages, "select * from pages p1 where p1.status='ACTIVE' and not exists (select * from pages p2 where p1.wiki_id=p2.wiki_id and p1.page_id=p2.page_id and p2.modified>p1.modified) and p1.path like $1 and p1.depth=$2 order by index", path, page.Depth + 1)

    return pages
}

// Used only from Wiki.Delete()
func DeletePages(wiki string, user string) {
    revel.TRACE.Printf("Page DeletePages(): wiki %+v, user %+v", wiki, user)
    pages := []Page{}
    db := get_db()
    defer db.Close()

    db.Select(&pages, "select * from pages p1 where p1.status='ACTIVE' and p1.wiki_id=uuid_in($1) and not exists (select * from pages p2 where p1.wiki_id=p2.wiki_id and p1.page_id=p2.page_id and p2.modified>p1.modified)",
        wiki)

    for _, page := range pages {
        DeletePage(wiki, page.Page_id, user, false)
    }
}

func ListPages(wiki_id string, node string) []Page {
    revel.TRACE.Printf("ListPages() wiki_id: %+v, node: %+v", wiki_id, node)
    pages := []Page{}
    db := get_db()
    defer db.Close()

    if len(node) > 0 {
        // Get the parent, and find what's under it - if anything
        p := GetPage(wiki_id, node)
        path := p.Path + "/%"
        depth := p.Depth + 1
        db.Select(&pages, "select * from pages p1 where p1.status='ACTIVE' and not exists (select * from pages p2 where p1.wiki_id=p2.wiki_id and p1.page_id=p2.page_id and p2.modified>p1.modified) and p1.path like $1 and p1.depth=$2", path, depth)
    } else {
        db.Select(&pages, "select * from pages p1 where p1.status='ACTIVE' and p1.wiki_id=uuid_in($1) and p1.depth=0 and not exists (select * from pages p2 where p1.wiki_id=p2.wiki_id and p1.page_id=p2.page_id and p2.modified>p1.modified)", wiki_id)
    }

    // Update children status
    for x, _ := range pages {
        pages[x].Loaded = !has_children(pages[x])
        revel.TRACE.Printf("Page %+v children: %+v", pages[x], has_children(pages[x]))
    }

    revel.TRACE.Printf("ListPages returning %+v", pages)
    return pages
}

func DeletePage(wiki_id string, page_id string, user string, save_activity bool) {
    revel.TRACE.Printf("DeletePage() wiki: %+v, page: %+v, user: %+v", wiki_id, page_id, user)

    page := GetPage(wiki_id, page_id)
    pages := Get_children(page)
    // For looping
    pages = append(pages, page) 

    for _, p := range pages {
        if p.Wiki_id != "" {
            p.Status = "DELETED"
            p.Create_user = user
            p.Save(save_activity)

            // Delete ContentField too
            cf := GetContent(p.Wiki_id, p.Page_id, nil)
            cf.Status = "DELETED"
            cf.Create_user = user
            cf.Save(save_activity)
        }
    }
}

func GetPage(wiki_id string, page_id string) Page {
    revel.TRACE.Printf("Page GetPage(): page: %+v wiki: %+v", page_id, wiki_id)
    pages := []Page{}
    page := Page{}
    db := get_db()
    defer db.Close()

    db.Select(&pages, "select * from pages p1 where p1.wiki_id=uuid_in($1) and p1.page_id=uuid_in($2) and p1.status='ACTIVE' and not exists (select * from pages p2 where p1.wiki_id=p2.wiki_id and p1.page_id=p2.page_id and p2.modified>p1.modified)",
        wiki_id, page_id)
    if len(pages)>0 {
        page = pages[0]
    }
    return page
}

func GetPageAllStatuses(wiki_id string, page_id string) Page {
    revel.TRACE.Printf("Page GetPageAllStatuses(): page: %+v wiki: %+v", page_id, wiki_id)
    pages := []Page{}
    page := Page{}
    db := get_db()
    defer db.Close()

    db.Select(&pages, "select * from pages p1 where p1.wiki_id=uuid_in($1) and p1.page_id=uuid_in($2) and not exists (select * from pages p2 where p1.wiki_id=p2.wiki_id and p1.page_id=p2.page_id and p2.modified>p1.modified)",
        wiki_id, page_id)
    if len(pages)>0 {
        page = pages[0]
    }
    return page
}

func update_depth(page *Page) {
    split := strings.Split(page.Path, "/")
    page.Depth = len(split) - 1
}

// ACL stuff

// Build ACL entry for reference
func (p Page) BuildACLEntry(reference string) acl.ACLEntry {
    revel.TRACE.Printf("BuildACLEntry() %+v", reference)
    var entry acl.ACLEntry

    if reference != ("page:"+p.Wiki_id+"/"+p.Page_id) {
        if strings.Index(reference, "page") == 0 {
            // We are not working on this copy, get from database
            re := regexp.MustCompile("page:([^/]*)/(.*)")
            m := re.FindStringSubmatch(reference)
            wref := m[1]
            pref := m[2]
            pa := GetPage(wref, pref)
            entry = entry_helper(pa.Readacl, pa.Writeacl, pa.Adminacl, reference, pa)
        } else {
            // This must be wiki!
            re := regexp.MustCompile("wiki:(.*)")
            ref := re.FindStringSubmatch(reference)[1]
            wi := GetWiki(ref)
            entry = entry_helper(wi.Readacl, wi.Writeacl, wi.Adminacl, reference, wi)
        }
    } else {
        // It's exactly the originating item!
        entry = entry_helper(p.Readacl, p.Writeacl, p.Adminacl, reference, p)
    }

    return entry
}

// Set the matched permissions to a variable
func (p Page) SetMatched(permissions []string) interface{} {
    p.MatchedPermissions = permissions
    return p
}

// Building parent information
func (p Page) BuildACLParent() string {
    if p.Depth==0 {
        return "wiki:"+p.Wiki_id
    } else {
        pslice := strings.Split(p.Path, "/")
        return "page:"+p.Wiki_id + "/" + pslice[len(pslice)-2]
    }
}

// Wiki+page ids...
func (p Page) BuildACLReference() string {
    return "page:"+p.Wiki_id+"/"+p.Page_id
}

// Set on data
func (p Page) BuildACLInheritation() bool {
    return !p.Stopinheritation
}
