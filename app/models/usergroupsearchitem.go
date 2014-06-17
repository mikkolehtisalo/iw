package models

import (
    "github.com/mikkolehtisalo/revel/ldapuserdetails"
    "fmt"
    "strings"
    "github.com/revel/revel"
    //"time"
       // _ "github.com/lib/pq"
       //"github.com/jmoiron/sqlx"
)

var (
    ldap_user_filter     string   = "&(objectClass=*)"
    ldap_user_uid_attr string = "*"
    ldap_user_cn_attr string = "*"
    ldap_user_base     string   = "dc=*,dc=*"
    ldap_group_filter string = "&(objectClass=*)"
    ldap_group_cn_attr string = "*"
    ldap_group_dn_attr string = "*"
    ldap_group_base string = "dc=*,dc=*"
)

func get_c_str(name string) string {
    if tmp, ok := revel.Config.String(name); !ok {
        panic(fmt.Errorf("%s invalid", name))
    } else {
        return tmp
    }
}

func init() {
    revel.OnAppStart(func() {
        ldap_user_filter = get_c_str("ldap.user_filter")
        ldap_user_base = get_c_str("ldap.user_base")
        ldap_user_uid_attr = get_c_str("ldap.user_uid_attr")
        ldap_user_cn_attr = get_c_str("ldap.user_cn_attr")
        ldap_group_filter = get_c_str("ldap.group_filter")
        ldap_group_base = get_c_str("ldap.group_base")
        ldap_group_cn_attr = get_c_str("ldap.group_cn_attr")
        ldap_group_dn_attr = get_c_str("ldap.group_dn_attr")
    })
}


type UserGroupSearchItem struct {
    Id string
    Name string
    Type string
}

func ListUserGroupSearchItems(query string) []UserGroupSearchItem {
    items := []UserGroupSearchItem{}
    l := ldapuserdetails.Get_connection()
    defer l.Close()

    sru := ldapuserdetails.QueryLdap(ldap_user_base, strings.Replace(ldap_user_filter, "*", 
        fmt.Sprintf("*%s*",query), -1), []string{ldap_user_uid_attr, ldap_user_cn_attr})
    srg := ldapuserdetails.QueryLdap(ldap_group_base, strings.Replace(ldap_group_filter, "*", 
        fmt.Sprintf("*%s*",query), -1), []string{ldap_group_cn_attr, ldap_group_dn_attr}) 
    
    for _, user := range sru.Entries {
        item := UserGroupSearchItem {
            Id: fmt.Sprintf("u:%s",user.GetAttributeValue(ldap_user_uid_attr)),
            Name: fmt.Sprintf("%s (u:%s)", user.GetAttributeValue(ldap_user_cn_attr), user.GetAttributeValue(ldap_user_uid_attr)),
            Type: "user"}
        items = append(items, item)
    }
    for _, group := range srg.Entries {
        item := UserGroupSearchItem {
            Id: fmt.Sprintf("g:%s", group.GetAttributeValue(ldap_group_cn_attr)),
            Name: fmt.Sprintf("%s (g:%s)", group.GetAttributeValue(ldap_group_cn_attr), group.GetAttributeValue(ldap_group_cn_attr)),
            Type: "group"}
        items = append(items, item)
    }

    return items
}

