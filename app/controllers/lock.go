package controllers

import (
    "github.com/revel/revel"
    "iw/app/models"
    //"encoding/json"
    //"strings"
    //"regexp"
    . "github.com/mikkolehtisalo/revel/common"
    "github.com/mikkolehtisalo/revel/ldapuserdetails"
    //"github.com/mikkolehtisalo/revel/acl"
)

type Locks struct {
    *revel.Controller
}

// CREATE
func (l Locks) Create(wiki string, target string) revel.Result {
    revel.TRACE.Printf("Locks Create() wiki: %+v, target: %+v", wiki, target)

    // Make sure the ids at least look like one
    if (!IsUUID(wiki)) || (!IsUUID(target)) {
        revel.ERROR.Printf("Garbage lock %+v/%+v received from %+v", wiki, target, l.Session["username"])
        return l.RenderText("{\"success\":false}")
    }

    // Make sure the lock doesn't pre-exist
    exists_test := models.GetLock(wiki, target)
    if exists_test.Wiki_id == wiki && exists_test.Target_id == target {
        revel.ERROR.Printf("Attempt to rewrite pre-existing lock %+v/%+v by user %+v", wiki, target, l.Session["username"])
        return l.RenderText("{\"success\":false}")
    }

    // Make sure the wiki exists!
    wiki_exists := models.GetWiki(wiki)
    if wiki_exists.Wiki_id != wiki {
        revel.ERROR.Printf("Attempt to add lock to non-existent wiki: %+v/%+v", wiki, target)
        return l.RenderText("{\"success\":false}")
    }

    // We could check for rights for page, but as obeying lock is non-forcing convenience function, meh

    lock := models.Lock{}
    lock.Target_id = target
    lock.Wiki_id = wiki
    lock.Username = l.Session["username"]
    dets := l.Args["user_details"].(ldapuserdetails.User_details)
    lock.Realname = dets.Visiblename
    lock.Save()

    return l.RenderText("{\"success\":true}")
}

// READ
func (l Locks) Read(wiki string, target string) revel.Result {
    revel.TRACE.Printf("Locks Read() wiki: %+v, target: %+v", wiki, target)

    // Make sure the ids at least look like one
    if (!IsUUID(wiki)) || (!IsUUID(target)) {
        revel.ERROR.Printf("Garbage lock %+v/%+v received from %+v", wiki, target, l.Session["username"])
        return l.RenderText("{\"success\":false}")
    }

    // Get the lock
    lock := models.GetLock(wiki, target)

    return l.RenderJson(lock)
}

// DELETE
func (l Locks) Delete(wiki string, target string) revel.Result {
    revel.TRACE.Printf("Locks Delete() wiki: %+v, target: %+v", wiki, target)

    // Make sure the ids at least look like one
    if (!IsUUID(wiki)) || (!IsUUID(target)) {
        revel.ERROR.Printf("Garbage lock %+v/%+v received from %+v", wiki, target, l.Session["username"])
        return l.RenderText("{\"success\":false}")
    }

    // Get the lock
    lock := models.GetLock(wiki, target)
    lock.Delete()

    return l.RenderText("{\"success\":true}")

}