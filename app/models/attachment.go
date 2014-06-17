package models

import (
    "time"
    "github.com/revel/revel"
    "github.com/mikkolehtisalo/revel/acl"
    "regexp"
    "strings"
)

type Attachment struct {
    Attachment_id string
    Wiki_id string
    // pg driver works probably better with BASE64 encoding instead of handling bytea hex string
    //Attachment []byte `json:",omitempty"`
    Attachment string
    Mime string
    Filename string
    Create_user string
    Modified time.Time
    MatchedPermissions []string
    Readacl string
    Writeacl string 
    Adminacl string
    Status string
}

func (a *Attachment) Validate(v *revel.Validation) {
    // Required fields
    v.Required(a.Attachment_id)
    v.Required(a.Wiki_id)
    v.Required(a.Attachment)
    v.Required(a.Filename)
    v.Required(a.Modified)

    // Match against regexp patterns
    v.Match(a.Wiki_id, regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")).Message("Wiki_id not UUID?")
    v.Match(a.Attachment_id, regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")).Message("Wiki_id not UUID?")

    // Validate the Attachment field
    re := regexp.MustCompile("^data:([^;]*);base64,(.*)$")
    v.Match(a.Attachment, re).Message("Invalid attachment data!")

    // Fix the content & mime filed
    if re.MatchString(a.Attachment) {
        parts := re.FindStringSubmatch(a.Attachment)
        a.Mime = parts[1]
        a.Attachment = parts[2]

    }
}

// Sets the modified timestamp
func (a *Attachment) Save(save_activity bool) {
    revel.TRACE.Printf("Attachment Save(): %+v", a)

    db := get_db()
    defer db.Close()
    a.Modified = time.Now()

    _, err := db.Exec("insert into attachments(wiki_id, attachment_id, attachment, mime, filename, create_user, modified, status) values ($1, $2, $3, $4, $5, $6, $7, $8)",
        a.Wiki_id, a.Attachment_id, a.Attachment, a.Mime, a.Filename, a.Create_user, a.Modified, a.Status)

    if err != nil {
        revel.ERROR.Printf("Attachment Save(): error %+v", err)
    }

    if save_activity {
        SaveActivity(a)
    }
}

func GetAttachment(wiki string, attachment string) Attachment {
    revel.TRACE.Printf("GetAttachment() wiki:%+v, attachment:%+v", wiki, attachment)
    attachments := []Attachment{}
    att := Attachment{}
    db := get_db()
    defer db.Close()

    err := db.Select(&attachments, "select * from attachments a1 where a1.wiki_id=uuid_in($1) and a1.status='ACTIVE' and a1.attachment_id=uuid_in($2) and not exists (select * from attachments a2 where a1.wiki_id=a2.wiki_id and a1.attachment_id=a2.attachment_id and a2.modified>a1.modified)", wiki, attachment)

    if err != nil {
        revel.ERROR.Printf("GetAttachment(): error %+v", err)
    }

    if len(attachments)>0 {
        att = attachments[0]
    }

    revel.TRACE.Printf("GetAttachment() returning %+v", att)
    return att
}

func GetAttachments(wiki string) []Attachment {
    revel.TRACE.Printf("GetAttachments() wiki: %+v", wiki)
    attachments := []Attachment{}
    db := get_db()
    defer db.Close()

    err := db.Select(&attachments, "select * from attachments a1 where a1.wiki_id=uuid_in($1) and a1.status='ACTIVE' and not exists (select * from attachments a2 where a1.wiki_id=a2.wiki_id and a1.attachment_id=a2.attachment_id and a2.modified>a1.modified)", wiki)

    if err != nil {
        revel.ERROR.Printf("GetAttachments(): error %+v", err)
    }

    revel.TRACE.Printf("GetAttachments() returning: %+v", attachments)
    return attachments
}

func DeleteAttachments(wiki string, user string) {
    revel.TRACE.Printf("DeleteAttachments() wiki: %+v, user: %+v", wiki, user)
    attachments := []Attachment{}
    db := get_db()
    defer db.Close()

    err := db.Select(&attachments, "select * from attachments a1 where a1.wiki_id=uuid_in($1) and a1.status='ACTIVE' and not exists (select * from attachments a2 where a1.wiki_id=a2.wiki_id and a1.attachment_id=a2.attachment_id and a2.modified>a1.modified)", wiki)

    if err != nil {
        revel.ERROR.Printf("DeleteAttachments(): error %+v", err)
    }

    for _, item := range attachments {
        item.Create_user = user
        item.Status = "DELETED"
        item.Save(false)
    }
}


// ACL stuff

// Build ACL entry for reference
func (a Attachment) BuildACLEntry(reference string) acl.ACLEntry {
    revel.TRACE.Printf("BuildACLEntry() %+v", reference)
    var entry acl.ACLEntry

    if reference != ("attachment:"+a.Wiki_id+"/"+a.Attachment_id) {
        if strings.Index(reference, "attachment") == 0 {
            // We are not working on this copy, get from database
            re := regexp.MustCompile("attachment:([^/]*)/(.*)")
            m := re.FindStringSubmatch(reference)
            wref := m[1]
            aref := m[2]
            aa := GetAttachment(wref, aref)
            entry = entry_helper(aa.Readacl, aa.Writeacl, aa.Adminacl, reference, aa)
        } else {
            // This must be wiki!
            re := regexp.MustCompile("wiki:(.*)")
            ref := re.FindStringSubmatch(reference)[1]
            wi := GetWiki(ref)
            entry = entry_helper(wi.Readacl, wi.Writeacl, wi.Adminacl, reference, wi)
        }
    } else {
        // It's exactly the originating item!
        entry = entry_helper(a.Readacl, a.Writeacl, a.Adminacl, reference, a)
    }

    return entry
}

// Set the matched permissions to a variable
func (a Attachment) SetMatched(permissions []string) interface{} {
    a.MatchedPermissions = permissions
    return a
}

// Building parent information
func (a Attachment) BuildACLParent() string {
    return "wiki:"+a.Wiki_id
}

// Wiki+attachment ids...
func (a Attachment) BuildACLReference() string {
    return "attachment:"+a.Wiki_id+"/"+a.Attachment_id
}

// Set on data
func (a Attachment) BuildACLInheritation() bool {
    return true
}
