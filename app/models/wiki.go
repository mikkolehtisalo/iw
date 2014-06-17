package models

import (
    "time"
    "github.com/mikkolehtisalo/revel/acl"
    "strings"
    "github.com/revel/revel"
    "regexp"
    "html"
)

type Wiki struct {
    Wiki_id string
    Title string
    Description string
    Create_user string
    Readacl string
    Writeacl string 
    Adminacl string
    Status string
    Modified time.Time
    MatchedPermissions []string
    Favorite bool
}

// Sets the modified timestamp
func (w *Wiki) Save(save_activity bool) {
    revel.TRACE.Printf("Wiki Save(): %+v", w)

    db := get_db()
    defer db.Close()
    w.Modified = time.Now()

    _, err := db.Exec("insert into wikis(wiki_id, title, description, create_user, readacl, writeacl, adminacl, status, modified) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
        w.Wiki_id, w.Title, w.Description, w.Create_user, w.Readacl, w.Writeacl, w.Adminacl, w.Status, w.Modified)

    if err != nil {
        revel.ERROR.Printf("Wiki Save(): error %+v", err)
    }

    if save_activity {
        SaveActivity(w)
    }
}

func (w *Wiki) Validate(v *revel.Validation) {
    // Required fields
    v.Required(w.Wiki_id)
    v.Required(w.Title)

    // Match against regexp patterns
    v.Match(w.Wiki_id, regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")).Message("Not UUID?")

    // Escape HTML from the fields that might be rendered to users
    w.Description = html.EscapeString(w.Description)
    w.Title = html.EscapeString(w.Title)
    w.Readacl = html.EscapeString(w.Readacl)
    w.Writeacl = html.EscapeString(w.Writeacl)
    w.Adminacl = html.EscapeString(w.Adminacl)
}

//All ACTIVE Wikis
func ListWikis() []Wiki {
    revel.TRACE.Printf("ListWikis()")

    wikis := []Wiki{}
    db := get_db()
    defer db.Close()

    err := db.Select(&wikis, "select * from wikis w1 where not exists (select * from wikis w2 where w2.modified>w1.modified and w1.wiki_id=w2.wiki_id) and status='ACTIVE'")
    if err != nil {
        revel.ERROR.Printf("ListWikis() err: %+v", err)
    }

    revel.TRACE.Printf("ListWikis() returning %+v", wikis)
    return wikis
}

//Newest ACTIVE version of wiki
func GetWiki(id string) Wiki {
    revel.TRACE.Printf("GetWiki() %+v", id)
    wikis := []Wiki{}
    wiki := Wiki{}
    db := get_db()
    defer db.Close()

    err := db.Select(&wikis, "select * from wikis w1 where w1.wiki_id=uuid_in($1) and w1.status='ACTIVE' and not exists (select * from wikis w2 where w1.wiki_id=w2.wiki_id and w2.modified>w1.modified)", id)
    if err != nil {
        revel.ERROR.Printf("GetWiki(): error %+v", err)
    }

    if len(wikis)>0 {
        wiki = wikis[0]
    }

    revel.TRACE.Printf("GetWiki() returning %+v", wiki)
    return wiki
}

// Goodbye, wiki!
func (w Wiki) Delete(user string) {
    revel.TRACE.Printf("Wiki Delete() user: %+v", user)
    w.Status = "DELETED"
    w.Create_user = user
    w.Save(true)
    DeletePages(w.Wiki_id, user)
    DeleteContentFields(w.Wiki_id, user)
    DeleteAttachments(w.Wiki_id, user)
    DeleteFavorites(w.Wiki_id)
}

// ACL stuff
// ---------

// Build ACL entry for reference
func (w Wiki) BuildACLEntry(reference string) acl.ACLEntry {
    entry := acl.ACLEntry{}

    tgt := w
    if reference != ("wiki:"+w.Wiki_id) {
        // We are not working on this copy, get from database
        re := regexp.MustCompile("wiki:(.*)")
        ref := re.FindStringSubmatch(reference)[1]
        tgt = GetWiki(ref)
    }

    // Build the ACL from tgt
    read_acl := acl.BuildPermissionACLs("read", strings.Split(tgt.Readacl, ","))
    write_acl := acl.BuildPermissionACLs("write", strings.Split(tgt.Writeacl, ","))
    admin_acl := acl.BuildPermissionACLs("admin", strings.Split(tgt.Adminacl, ","))
    acls := append(read_acl, write_acl...)
    acls = append(acls, admin_acl...)
    entry.ObjReference = reference
    entry.ACLs = acls
    entry.Inheritation = tgt.BuildACLInheritation()
    entry.Parent = tgt.BuildACLParent()

    return entry
}

// Set the matched permissions to a variable
func (w Wiki) SetMatched(permissions []string) interface{} {
    w.MatchedPermissions = permissions
    return w
}

// Just append type to the id
func (w Wiki) BuildACLReference() string {
    return "wiki:"+w.Wiki_id
}

// No wiki inherits ACL
func (w Wiki) BuildACLInheritation() bool {
    return false
}

// No wiki has parent
func (w Wiki) BuildACLParent() string {
    return ""
}