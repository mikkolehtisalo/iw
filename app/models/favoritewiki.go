package models

import (
    "time"
    "github.com/revel/revel"
)

type FavoriteWiki struct {
    Username string
    Wiki_id string
    Modified time.Time
    Status string
}

// Sets modified time
func (f FavoriteWiki) Save() {
    revel.TRACE.Printf("FavoriteWiki Save(): %+v", f)
    db := get_db()
    defer db.Close()

    _, err := db.Exec("insert into favoritewikis(username, wiki_id, modified, status) values ($1, $2, $3, $4)", f.Username, f.Wiki_id, time.Now(), f.Status)

    if err != nil {
        revel.ERROR.Printf("FavoriteWiki Save(): error %+v", err)
    }
}

// List all favorites of the user
func ListFavoriteWikis(user string) []FavoriteWiki {
    revel.TRACE.Printf("ListFavoriteWikis(): %+v", user)
    favorites := []FavoriteWiki{}
    db := get_db()
    defer db.Close()

    err := db.Select(&favorites, "select * from favoritewikis w1 where username=$1 and status='ACTIVE' and not exists (select * from favoritewikis w2 where w1.wiki_id=w2.wiki_id and w1.username=w2.username and w2.modified>w1.modified)", user)
    if err != nil {
        revel.ERROR.Printf("ListFavoriteWikis(): error %+v", err)
    }

    revel.TRACE.Printf("ListFavoriteWikis() returning: %+v", favorites)
    return favorites
}

func DeleteFavorites(wiki string) {
    revel.TRACE.Printf("DeleteFavorites(): %+v", wiki)
    favorites := []FavoriteWiki{}
    db := get_db()
    defer db.Close()

    err := db.Select(&favorites, "select * from favoritewikis w1 where wiki_id=uuid_in($1) and status='ACTIVE' and not exists (select * from favoritewikis w2 where w1.wiki_id=w2.wiki_id and w1.username=w2.username and w2.modified>w1.modified)", wiki)
    if err != nil {
        revel.ERROR.Printf("DeleteFavorites(): error %+v", err)
    }

    // Delete!
    for _, item := range favorites {
        item.Status = "DELETED"
        item.Save()
    }
}

// Is the wiki of user already favorited?
// No point in optimizing probably, won't get called very often
func IsFavoriteWiki(wiki string, user string) bool {
    fav := false
    favs := ListFavoriteWikis(user)
    for _, item := range favs {
        if item.Wiki_id==wiki {
            fav = true
        }
    }
    return fav
}
