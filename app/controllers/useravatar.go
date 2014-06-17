package controllers

import (
    "github.com/revel/revel"
    "bytes"
    "github.com/mikkolehtisalo/revel/ldapuserdetails"
    "time"
    "regexp"
    "github.com/revel/revel/cache"
    "fmt"
)

type UserAvatars struct {
    *revel.Controller
}

func (u UserAvatars) Read(user string) revel.Result {
    revel.TRACE.Printf("UserAvatars Read() user:%+v", user)
    var avatar []byte 

    err := cache.Get(fmt.Sprintf("avatar:%s", user), &avatar)

    if err != nil {
        // Not in cache, generate
        re := regexp.MustCompile("^(\\w*)\\.(jpeg|jpg|png)")

        if !re.MatchString(user) {
            revel.TRACE.Printf("UserAvatars Read() invalid user: %+v", user)
            return u.RenderText("{\"success\":false}")
        }

        username := re.FindStringSubmatch(user)[1]
        dets := ldapuserdetails.Get_user_details(username)
        avatar = dets.Photo
        go cache.Set(fmt.Sprintf("avatar:%s", user), avatar, cache.DEFAULT)
    }

    return u.RenderBinary(bytes.NewReader(avatar), user, "inline", time.Now())

}
