package models

import (
    "time"
    "github.com/revel/revel"
)

type Lock struct {
    Target_id string
    Wiki_id string
    Username string
    Realname string
    Modified time.Time
}

// Sets modified time
func (l Lock) Save() {
    revel.TRACE.Printf("Lock Save(): %+v", l)
    db := get_db()
    defer db.Close()

    _, err := db.Exec("insert into locks(wiki_id, target_id, username, realname, modified) values ($1, $2, $3, $4, $5)",
        l.Wiki_id, l.Target_id, l.Username, l.Realname, time.Now())

    if err != nil {
        revel.ERROR.Printf("Lock Save(): failed with %+v", err)
    }
}

// Deletes by target and wiki id
func (l Lock) Delete() {
    revel.TRACE.Printf("Lock Delete(): %+v", l)
    db := get_db()
    defer db.Close()

    _, err := db.Exec("delete from locks where wiki_id=uuid_in($1) and target_id=uuid_in($2)",
        l.Wiki_id, l.Target_id)

    if err != nil {
        revel.ERROR.Printf("Lock Delete(): failed with %+v", err)
    }
}

func GetLock(wiki string, target string) Lock {
    revel.TRACE.Printf("GetLock(): target: %+v wiki: %+v", target, wiki)
    locks := []Lock{}
    lock := Lock{}
    db := get_db()
    defer db.Close()

    db.Select(&locks, "select * from locks where wiki_id=uuid_in($1) and target_id=uuid_in($2)",
        wiki, target)
    if len(locks)>0 {
        lock = locks[0]
    }
    return lock
}
