package common

import (
        "io/ioutil"
        "log"
        "encoding/json"
        "errors"
        "net/http"
        "time"
)

// Checks whether string can be found from slice
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

// Defines the configuration file format
type ServiceConfiguration struct {
    CAKeyFile string
    MyCertificateFile string
    MyKeyFile string
    ValidPeers []string
    ServiceAddress string
    DBUser string
    DBPassword string
    DBName string
    Interval int
}

// Used to load the configuration from file
func (l *ServiceConfiguration) ConfigFrom(path string) (error) {
        b, err := ioutil.ReadFile(path)
        if err != nil {
            log.Fatalf("ReadFile: %s", err)
            return err
        }
        err = json.Unmarshal(b, &l)
        if err != nil {
            log.Fatalf("Bad json: %s", err)
            return err
        }
        return nil
}

// Peer CN check
func CheckPeer(r *http.Request, settings ServiceConfiguration) error {
    var valid bool = false

    for _,chain := range r.TLS.VerifiedChains {
       for _,certificate := range chain {
            cn := certificate.Subject.CommonName
            if stringInSlice(cn, settings.ValidPeers) {
                valid = true
            }
       } 
    }

    if !valid {
        return errors.New("No valid certificate received from peer!")
    }

    return nil
}

// Types for the replication

type Attachment struct {
    Attachment_id string
    Wiki_id string
    // pg driver works probably better with BASE64 encoding instead of handling bytea hex string
    //Attachment []byte `json:",omitempty"`
    Attachment string `json:",omitempty"`
    Mime string `json:",omitempty"`
    Filename string `json:",omitempty"`
    Create_user string `json:",omitempty"`
    Modified time.Time
    Status string `json:",omitempty"`
}

type FavoriteWiki struct {
    Username string
    Wiki_id string
    Modified time.Time
    Status string `json:",omitempty"`
}

type ContentField struct {
    Contentfield_id string
    Wiki_id string
    Content string `json:",omitempty"`
    Contentwithmacros string `json:",omitempty"`
    Modified time.Time
    Status string `json:",omitempty"`
    Create_user string `json:",omitempty"`
}

type Wiki struct {
    Wiki_id string
    Title string `json:",omitempty"`
    Description string `json:",omitempty"`
    Create_user string `json:",omitempty"`
    Readacl string `json:",omitempty"`
    Writeacl string `json:",omitempty"`
    Adminacl string `json:",omitempty"`
    Status string `json:",omitempty"`
    Modified time.Time
}

type Page struct {
    Page_id string
    Wiki_id string
    Path string `json:",omitempty"`
    Title string `json:",omitempty"`
    Create_user string `json:",omitempty"`
    Readacl string `json:",omitempty"`
    Writeacl string `json:",omitempty"`
    Adminacl string `json:",omitempty"`
    Stopinheritation bool `json:",omitempty"`
    Index int `json:",omitempty"`
    Depth int `json:",omitempty"`
    Status string `json:",omitempty"`
    Modified time.Time
}

type Activity struct {
    Activity_id string
    Timestamp time.Time `json:",omitempty"`
    User_id string `json:",omitempty"`
    User_name string `json:",omitempty"`
    Activity_type string `json:",omitempty"`
    Target_type string `json:",omitempty"`
    Target_title string `json:",omitempty"`
    Target_id string `json:",omitempty"`
    Readacl string `json:",omitempty"`
    Writeacl string `json:",omitempty"`
    Adminacl string `json:",omitempty"`
}

type Lock struct {
    Target_id string
    Wiki_id string
    Username string `json:",omitempty"`
    Realname string `json:",omitempty"`
    Modified time.Time `json:",omitempty"`
}
