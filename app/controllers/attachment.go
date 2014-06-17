package controllers

import (
    "github.com/revel/revel"
    . "github.com/mikkolehtisalo/revel/common"
    "iw/app/models"
    "encoding/json"
    "github.com/mikkolehtisalo/revel/acl"
    "encoding/base64"
    //"fmt"
    "time"
    "bytes"
)

type Attachments struct {
    *revel.Controller
}

// CREATE
func (a Attachments) Create(wiki string, attachment string) revel.Result {
    revel.TRACE.Printf("Attachments Create() wiki: %+v, attachment: %+v", wiki, attachment)

    // Make sure the ids at least look like one
    if (!IsUUID(wiki)) || (!IsUUID(attachment)) {
        revel.ERROR.Printf("Garbage attachment %+v/%+v received from %+v", wiki, attachment, a.Session["username"])
        return a.RenderText("{\"success\":false}")
    }

    // Make sure the attachment doesn't pre-exist
    exists_test := models.GetAttachment(wiki, attachment)
    if exists_test.Wiki_id == wiki && exists_test.Attachment_id == attachment {
        revel.ERROR.Printf("Attempt to rewrite pre-existing attachment %+v/%+v by user %+v", wiki, attachment, a.Session["username"])
        return a.RenderText("{\"success\":false}")
    }

    // Make sure the wiki exists!
    wiki_exists := models.GetWiki(wiki)
    if wiki_exists.Wiki_id != wiki {
        revel.ERROR.Printf("Attempt to add attachment to non-existent wiki: %+v/%+v", wiki, attachment)
        return a.RenderText("{\"success\":false}")
    }

    // Decode attachment from json input
    var new_att models.Attachment
    decoder := json.NewDecoder(a.Request.Body)
    err := decoder.Decode(&new_att)
    if err != nil {
        revel.ERROR.Printf("Unable to parse attachment %+v/%+v: %+v", wiki, attachment, err)
        return a.RenderText("{\"success\":false}")
    }

    // ID fields must match!
    if (new_att.Wiki_id != wiki) || (new_att.Attachment_id != attachment) {
        revel.ERROR.Printf("Attachment id mismatch %+v/%+v != %+v/%+v from user %+v", wiki, attachment, new_att.Wiki_id, new_att.Attachment_id, a.Session["username"])
        return a.RenderText("{\"success\":false}")
    }

    // Validate fields
    new_att.Validate(a.Validation)
    if a.Validation.HasErrors() {
        revel.ERROR.Printf("Validation errors: %+v", a.Validation.ErrorMap())
        return a.RenderText("{\"success\":false}")
    }

    // Make sure the user has rights to create attachment
    filtered := acl.Filter(a.Args, []string{"admin","write"}, []models.Attachment{new_att}, true)
    if len(filtered) != 1 {
        revel.ERROR.Printf("Attempt to create attachment %+v/%+v without access rights: %+v, user: %+v", wiki, attachment, new_att.Filename, a.Session["username"])
        return a.RenderText("{\"success\":false}")
    }

    new_att.Create_user = a.Session["username"]
    new_att.Status = "ACTIVE"
    new_att.Save(true)

    return a.RenderText("{\"success\":true}")
}

// READ
func (a Attachments) Read(wiki string) revel.Result {
    revel.TRACE.Printf("Attachments Read(): %+v", wiki)

    // Make sure the id looks like one
    if (!IsUUID(wiki)) {
        revel.ERROR.Printf("Garbage wiki %+v received from %+v", wiki, a.Session["username"])
        return a.RenderText("{\"success\":false}")
    }

    attachments := models.GetAttachments(wiki)
    // Filter by acls...
    filtered := acl.Filter(a.Args, []string{"read", "write", "admin"}, attachments, true)

    revel.TRACE.Printf("Attachments Read() returning %+v", filtered)
    return a.RenderJson(filtered)
}

// Serve direct links!
func (a Attachments) Serve(wiki string, attachment string) revel.Result {
    revel.TRACE.Printf("Attachments Serve(): %+v, attachment: %+v", wiki, attachment)

    // Make sure the ids at least look like one
    if (!IsUUID(wiki)) || (!IsUUID(attachment)) {
        revel.ERROR.Printf("Garbage attachment %+v/%+v received from %+v", wiki, attachment, a.Session["username"])
        return a.RenderText("{\"success\":false}")
    }

    att := models.GetAttachment(wiki, attachment)

    // Make sure the user has rights to read the attachment
    filtered := acl.Filter(a.Args, []string{"read", "admin","write"}, []models.Attachment{att}, true)
    if len(filtered) != 1 {
        revel.ERROR.Printf("Attempt to read attachment %+v/%+v without access rights: %+v, user: %+v", wiki, attachment, att.Filename, a.Session["username"])
        return a.RenderText("{\"success\":false}")
    }

    data, err := base64.StdEncoding.DecodeString(att.Attachment)
    if err != nil {
        revel.ERROR.Printf("Serve(): Unable to base64 decode attachment! %+v", err)
    }

    return a.RenderBinary(bytes.NewReader(data), att.Filename, "inline", time.Now())
}

// DELETE
func (a Attachments) Delete(wiki string, attachment string) revel.Result {
    revel.TRACE.Printf("Attachments Delete(): %+v, attachment: %+v", wiki, attachment)

    // Make sure the ids at least look like one
    if (!IsUUID(wiki)) || (!IsUUID(attachment)) {
        revel.ERROR.Printf("Garbage attachment %+v/%+v received from %+v", wiki, attachment, a.Session["username"])
        return a.RenderText("{\"success\":false}")
    }

    att := models.GetAttachment(wiki, attachment)

    // Make sure the user has rights to delete the attachment
    filtered := acl.Filter(a.Args, []string{"admin","write"}, []models.Attachment{att}, true)
    if len(filtered) != 1 {
        revel.ERROR.Printf("Attempt to delete attachment %+v/%+v without access rights: %+v, user: %+v", wiki, attachment, att.Filename, a.Session["username"])
        return a.RenderText("{\"success\":false}")
    }

    att.Create_user = a.Session["username"]
    att.Status = "DELETED"
    att.Save(true)

    return a.RenderText("{\"success\":true}")
}