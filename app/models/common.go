package models

import (
    "github.com/revel/revel"
    "fmt"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "github.com/mikkolehtisalo/revel/acl"
    "strings"
)

var (
    db_user string
    db_password string
    db_name string
)

// Initialize settings from app.conf
func init() {
    revel.OnAppStart(func() {
        var ok bool
        if db_user, ok = revel.Config.String("db.user"); !ok {
            panic(fmt.Errorf("Unable to read db.user"))
        }
        if db_password, ok = revel.Config.String("db.password"); !ok {
            panic(fmt.Errorf("Unable to read db.password"))
        }
        if db_name , ok = revel.Config.String("db.name"); !ok {
            panic(fmt.Errorf("Unable to read db.name"))
        }
    })
}

func checkError(err error, s string) {
    // Syslogger
    //logger, _ := syslog.New(syslog.LOG_ERR, "SyncServer")
    //defer logger.Close()

    if err != nil {
            //logger.Err(fmt.Sprintf("%s: %s", s, err))
            panic(fmt.Sprintf("%s: %s", s, err))
    }
}

// Open new database connection
func get_db() *sqlx.DB {
    connstring := fmt.Sprintf("user=%s password='%s' dbname=%s sslmode=disable", db_user, db_password, db_name)
    db, err := sqlx.Open("postgres", connstring)
    checkError(err, "sqlx.Open")
    return db
}

// Common helper for building acl entries specific to this application
func entry_helper(read string, write string, admin string, reference string, tgt acl.Filterable) acl.ACLEntry {
    entry := acl.ACLEntry{}

    read_acl := acl.BuildPermissionACLs("read", strings.Split(read, ","))
    write_acl := acl.BuildPermissionACLs("write", strings.Split(write, ","))
    admin_acl := acl.BuildPermissionACLs("admin", strings.Split(admin, ","))
    acls := append(read_acl, write_acl...)
    acls = append(acls, admin_acl...)
    entry.ACLs = acls

    entry.ObjReference = reference
    entry.Inheritation = tgt.BuildACLInheritation()
    entry.Parent = tgt.BuildACLParent()

    return entry
}
