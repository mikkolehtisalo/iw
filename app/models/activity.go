package models

import (
    "time"
    "github.com/revel/revel"
    "github.com/twinj/uuid"
    "github.com/mikkolehtisalo/revel/ldapuserdetails"
    "github.com/mikkolehtisalo/revel/acl"
    "regexp"
    "strings"
)

type Activity struct {
    Activity_id string
    Timestamp time.Time
    User_id string
    User_name string
    Activity_type string
    Target_type string
    Target_title string
    Target_id string
    MatchedPermissions []string
    // For the "flattened" acls...
    Readacl string
    Writeacl string
    Adminacl string
}

// Oldest: 1 month
func GetActivities() []Activity {
    revel.TRACE.Printf("GetActivities()")
    activities := []Activity{}
    db := get_db()
    defer db.Close()

    err := db.Select(&activities, "select * from activities where timestamp > now() - interval '1 month' order by timestamp desc")

    if err != nil {
        revel.ERROR.Printf("GetActivities() err: %+v", err)
    }

    revel.TRACE.Printf("GetActivities() returning: %+v", activities)
    return activities
}

func GetActivity(activity_id string) Activity {
    revel.TRACE.Printf("GetActivity() id: %+v", activity_id)
    activities := []Activity{}
    activity := Activity{}
    db := get_db()
    defer db.Close()

    err := db.Select(&activities, "select * from activities where activity_id=uuid_in($1)", activity_id)

    if err != nil {
        revel.ERROR.Printf("Unable to get activity %+v: %+v", activity_id, err)
    }

    if len(activities) > 0 {
        activity = activities[0]
    }

    revel.TRACE.Printf("GetActivity() returning: %+v", activity)

    return activity
}

// Is there a similar activity within one minute range?
func ActivityExists(a Activity) bool {
    exists := false

    activities := []Activity{}
    db := get_db()
    defer db.Close()

    err := db.Select(&activities, "select * from activities where $1 - timestamp < interval '1 minute' and user_id=$2 and activity_type=$3 and target_type=$4 and target_id=$5",
        a.Timestamp, a.User_id, a.Activity_type, a.Target_type, a.Target_id)

    if err != nil {
        revel.ERROR.Printf("ActivityExists() err: %+v", err)
    }

    if len(activities) > 0 {
        exists = true
    }


    return exists
}

// Flatten the inherited ACL structure, for performance reasons
// If the ACLs of the parent items are edited, the ACLs of activity will not change..
func update_acl(a *Activity) {
    filterable := acl.Get_filterable([]Activity{*a})
    acls := acl.GetACLEntry("activity:" + a.Activity_id, filterable[0], true)

    // Remove the duplicates from previous (it will generate them) and flatten the result
    var read, write, admin []string
    for _, item := range acls.ACLs {
        if item.Permission == "read" && item.Principal != "" {
            read = append(read, item.Principal)
        }
        if item.Permission == "write" && item.Principal != "" {
            write = append(write, item.Principal)
        }
        if item.Permission == "admin" && item.Principal != "" {
            admin = append(admin, item.Principal)
        }
    }

    a.Readacl = strings.Join(read, ",")
    a.Writeacl = strings.Join(write, ",")
    a.Adminacl = strings.Join(admin, ",")
}

// Flattens the ACL before saving
func (a *Activity) Save() {
    revel.TRACE.Printf("Activity Save(): %+v", a)

    db := get_db()
    defer db.Close()

    update_acl(a)

    _, err := db.Exec("insert into activities(activity_id, timestamp, user_id, user_name, activity_type, target_type, target_title, target_id, readacl, writeacl, adminacl) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
        a.Activity_id, a.Timestamp, a.User_id, a.User_name, a.Activity_type, a.Target_type, a.Target_title, a.Target_id, a.Readacl, a.Writeacl, a.Adminacl)

    if err != nil {
        revel.ERROR.Printf("Activity Save(): error %+v", err)
    }

}

func SaveActivity(target interface{}) {
    revel.TRACE.Printf("target: %+v", target)
    a := Activity{}
    a.Activity_id = uuid.NewV4().String()

    switch target.(type) {
    case *Wiki:
        w := target.(*Wiki)
        a.Timestamp = w.Modified
        a.User_id = w.Create_user
        dets := ldapuserdetails.Get_user_details(w.Create_user)
        a.User_name = dets.Visiblename
        a.Activity_type = w.Status
        a.Target_type = "WIKI"
        a.Target_title = w.Title
        a.Target_id = w.Wiki_id
        if !ActivityExists(a) {
            a.Save()
        }
    case *Page:
        p := target.(*Page)
        a.Timestamp = p.Modified
        a.User_id = p.Create_user
        dets := ldapuserdetails.Get_user_details(p.Create_user)
        a.User_name = dets.Visiblename
        a.Activity_type = p.Status
        a.Target_type = "PAGE"
        a.Target_title = p.Title

        // We have to save the parent node for treeStore, or "" if it's on first level....
        parent := ""
        if p.Depth > 0 {
            split := strings.Split(p.Path, "/")
            parent = split[len(split)-2]
        }
        a.Target_id = p.Wiki_id + "/" +parent + "/" + p.Page_id

        if !ActivityExists(a) {
            a.Save()
        }
    case *Attachment:
        x := target.(*Attachment)
        a.Timestamp = x.Modified
        a.User_id = x.Create_user
        dets := ldapuserdetails.Get_user_details(x.Create_user)
        a.User_name = dets.Visiblename
        a.Activity_type = x.Status
        a.Target_type = "ATTACHMENT"
        a.Target_title = x.Filename
        a.Target_id = x.Wiki_id + "/" + x.Attachment_id
        if !ActivityExists(a) {
            a.Save()
        }
    }
}

// ACL stuff

// Never inherit the parent!
func (a Activity) BuildACLInheritation() bool {
    return true
}

// Activity + activity_id
func (a Activity) BuildACLReference() string {
    return "activity:" + a.Activity_id
}

// No parent..
func (a Activity) BuildACLParent() string {
    return strings.ToLower(a.Target_type) + ":" + a.Target_id
}

// Set the matched permissions to a variable
func (a Activity) SetMatched(permissions []string) interface{} {
    a.MatchedPermissions = permissions
    return a
}

func (a Activity) BuildACLEntry(reference string) acl.ACLEntry {
    revel.TRACE.Printf("BuildACLEntry() %+v", reference)
    var entry acl.ACLEntry

    if reference != ("activity:"+a.Activity_id) {

        if strings.Index(reference, "page") == 0 {
            // We are not working on this copy, get from database
            re := regexp.MustCompile("page:([^/]*)/([^/]*)/(.*)")
            m := re.FindStringSubmatch(reference)
            wref := m[1]
            pref := m[3]
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
        entry = entry_helper(a.Readacl, a.Writeacl, a.Adminacl, reference, a)
    }

    revel.TRACE.Printf("BuildACLEntry() returning %+v", entry)

    return entry
}
