package main

import (
    "log/syslog"
    "time"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    "net/http"
    "fmt"
    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx"
    "iw/replication/src/common"
    "encoding/json"
    "reflect"
)

func checkError(err error, s string) {
    if err != nil {
            panic(fmt.Sprintf("%s: %s", s, err))
    }
}

func countRows (rows *sqlx.Rows) int {
    var count = 0
    for rows.Next() {
        count ++
    }
    return count
}

func getDB(settings common.ServiceConfiguration) *sqlx.DB {
    // Open database connection
    connstring := fmt.Sprintf("user=%s password='%s' dbname=%s sslmode=disable", settings.DBUser, settings.DBPassword, settings.DBName)
    db, err := sqlx.Open("postgres", connstring)
    checkError(err, "sql.Open")
    return db
}

func itemExists (item interface{}, settings common.ServiceConfiguration) bool {
    db := getDB(settings)
    defer db.Close()

    // Query for data
    var rows *sqlx.Rows
    // Error
    var err error

    // Prepare for reflect
    val := reflect.ValueOf(item)

    // Query the correct item type
    switch item.(type) {
        case common.Wiki:
            modified := val.FieldByName("Modified").Interface()
            rows, err = db.Queryx("select wiki_id, modified from wikis "+
                "where wiki_id=uuid_in($1) and modified=$2", val.FieldByName("Wiki_id").String(), modified)
        case common.Page:
            modified := val.FieldByName("Modified").Interface()
            rows, err = db.Queryx("select page_id, wiki_id, modified from pages "+
                "where page_id=uuid_in($1) and wiki_id=uuid_in($2) and modified=$3", val.FieldByName("Page_id").String(), val.FieldByName("Wiki_id").String(), modified)
        case common.ContentField:
            modified := val.FieldByName("Modified").Interface()
            rows, err = db.Queryx("select contentfield_id, wiki_id, modified from contentfields "+
                "where contentfield_id=uuid_in($1) and wiki_id=uuid_in($2) and modified=$3", val.FieldByName("Contentfield_id").String(), val.FieldByName("Wiki_id").String(), modified)
        case common.Attachment:
            modified := val.FieldByName("Modified").Interface()
            rows, err = db.Queryx("select attachment_id, wiki_id, modified from attachments "+
                "where attachment_id=uuid_in($1) and wiki_id=uuid_in($2) and modified=$3", val.FieldByName("Attachment_id").String(), val.FieldByName("Wiki_id").String(), modified)
        case common.FavoriteWiki:
            modified := val.FieldByName("Modified").Interface()
            rows, err = db.Queryx("select username, wiki_id, modified from favoritewikis "+
                "where username=$1 and wiki_id=uuid_in($2) and modified=$3", val.FieldByName("Username").String(), val.FieldByName("Wiki_id").String(), modified)
        case common.Activity:
            rows, err = db.Queryx("select activity_id from activities "+
                "where activity_id=uuid_in($1)", val.FieldByName("Activity_id").String())
        case common.Lock:
            rows, err = db.Queryx("select target_id, wiki_id from locks "+
                "where target_id=uuid_in($1) and wiki_id=uuid_in($2)", val.FieldByName("Target_id").String(), val.FieldByName("Wiki_id").String())
    }
    
    checkError(err, "Queryx")

    count := countRows(rows)
    if count == 1 {
        return true
    }

    return false
}

func saveItem(item interface{}, settings common.ServiceConfiguration) {
    db := getDB(settings)
    defer db.Close()
    var err error

    switch item.(type) {
        case common.Wiki:
            _, err = db.NamedExec("INSERT INTO wikis (wiki_id, title, description, create_user, readacl, writeacl, adminacl, status, modified) VALUES "+
                "(:wiki_id, :title, :description, :create_user, :readacl, :writeacl, :adminacl, :status, :modified)", item)
        case common.Page:
            _, err = db.NamedExec("insert into pages (page_id, wiki_id, path, title, create_user, readacl, writeacl, adminacl, stopinheritation, index, depth, status, modified) values "+
                "(:page_id, :wiki_id, :path, :title, :create_user, :readacl, :writeacl, :adminacl, :stopinheritation, :index, :depth, :status, :modified)", item)
        case common.ContentField:
            _, err = db.NamedExec("insert into contentfields (contentfield_id, wiki_id, content, modified, status, create_user) values "+
                "(:contentfield_id, :wiki_id, :content, :modified, :status, :create_user)", item)
        case common.Attachment:
            _, err = db.NamedExec("insert into attachments (attachment_id, wiki_id, attachment, mime, filename, modified, status, create_user) values "+
                "(:attachment_id, :wiki_id, decode(:attachment, 'base64'), :mime, :filename, :modified, :status, :create_user)", item)
        case common.FavoriteWiki:
           _, err = db.NamedExec("insert into favoritewikis (username, wiki_id, modified, status) values "+
                "(:username, :wiki_id, :modified, :status)", item)
        case common.Activity:
            _, err = db.NamedExec("insert into activities (activity_id, timestamp, user_id, user_name, activity_type, target_type, target_title, target_id, readacl, writeacl, adminacl) values "+
                "(:activity_id, :timestamp, :user_id, :user_name, :activity_type, :target_type, :target_title, :target_id, :readacl, :writeacl, :adminacl)", item)
        case common.Lock:
            _, err = db.NamedExec("insert into locks (target_id, wiki_id, username, realname, modified values "+
                "(:target_id, :wiki_id, :username, :realname, :modified", item)
    }
    checkError(err, "execInsert")
}

func mapToStruct(input map[string] interface{}, it interface{}) interface{} {
    itType := reflect.ValueOf(it).Type()
    v := reflect.New(itType).Elem()

    for key, value := range input {
        field := v.FieldByName(key)
        if !field.IsValid() {
            // or handle as error if you don't expect unknown values
            continue
        }
        if !field.CanSet() {
            // or return an error on private fields
            continue
        }

        switch key {
            case "Depth", "Index":
                fl := reflect.ValueOf(value).Float()
                fi := int(fl)
                field.Set(reflect.ValueOf(fi))
            case "Modified", "Timestamp":
                s := reflect.ValueOf(value).String()
                modified, terr := time.Parse(time.RFC3339Nano, s)
                checkError(terr, "Unable to parse time")
                field.Set(reflect.ValueOf(modified))
            default:
                field.Set(reflect.ValueOf(value))
        }
       
    }
    return v.Interface()
}

func syncItems(tr *http.Transport, settings common.ServiceConfiguration, it interface{}, queryListUrl string, queryUrl string) error {
    // Syslogger
    logger, _ := syslog.New(syslog.LOG_ERR, "SyncClient")
    defer logger.Close()

    // Get list of items
    client := &http.Client{Transport: tr}
    result, err := client.Get("https://" + settings.ServiceAddress + queryListUrl)
    checkError(err, "syncItems Get")
    defer result.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(result.Body)
    checkError(err, "syncItems ReadAll")

    // Marshal the response JSON to map
    var tmp []map[string] interface{}
    err = json.Unmarshal(body, &tmp)
    checkError(err, "syncItems Unmarshal")

    // Parse each item from map
    for _, item := range tmp {
        // Convert to struct
        target := mapToStruct(item, it)
        val := reflect.ValueOf(target)
        // If item doesn't exist already, sync full item
        if !itemExists(target, settings) {
            logger.Info(fmt.Sprintf("Item %+v : %+v did not exist, syncing...", reflect.ValueOf(target).Type(), target))

            var resp *http.Response
            // Get the whole item
            switch target.(type) {
                case common.Wiki:
                    iface := val.FieldByName("Modified").Interface()
                    modified := iface.(time.Time)
                    formated := modified.Format(time.RFC3339Nano)
                    resp, err = client.Get("https://" + settings.ServiceAddress + queryUrl + val.FieldByName("Wiki_id").String() + "/" + formated)  
                case common.Page:
                    iface := val.FieldByName("Modified").Interface()
                    modified := iface.(time.Time)
                    formated := modified.Format(time.RFC3339Nano)
                    resp, err = client.Get("https://" + settings.ServiceAddress + queryUrl + val.FieldByName("Page_id").String() + "/" + val.FieldByName("Wiki_id").String() + "/" + formated)  
                case common.ContentField:
                    iface := val.FieldByName("Modified").Interface()
                    modified := iface.(time.Time)
                    formated := modified.Format(time.RFC3339Nano)
                    resp, err = client.Get("https://" + settings.ServiceAddress + queryUrl + val.FieldByName("Contentfield_id").String() + "/" + val.FieldByName("Wiki_id").String() + "/" + formated)  
                case common.Attachment:
                    iface := val.FieldByName("Modified").Interface()
                    modified := iface.(time.Time)
                    formated := modified.Format(time.RFC3339Nano)
                    resp, err = client.Get("https://" + settings.ServiceAddress + queryUrl + val.FieldByName("Attachment_id").String() + "/" + val.FieldByName("Wiki_id").String() + "/" + formated)  
                case common.FavoriteWiki:
                    iface := val.FieldByName("Modified").Interface()
                    modified := iface.(time.Time)
                    formated := modified.Format(time.RFC3339Nano)
                    resp, err = client.Get("https://" + settings.ServiceAddress + queryUrl + val.FieldByName("Username").String() + "/" + val.FieldByName("Wiki_id").String() + "/" + formated)  
                case common.Activity:
                    resp, err = client.Get("https://" + settings.ServiceAddress + queryUrl + val.FieldByName("Activity_id").String())
                case common.Lock:
                    resp, err = client.Get("https://" + settings.ServiceAddress + queryUrl + val.FieldByName("Wiki_id").String() + "/" + val.FieldByName("Target_id").String())  
            }

            checkError(err, "syncItems GetFull")
            defer resp.Body.Close()
            // Read the body of whole item 
            fullBody, err := ioutil.ReadAll(resp.Body)
            checkError(err, "syncWikis ReadFullBody")

            // Marshal to map and convert to item
            var temp2 map[string] interface{}
            err = json.Unmarshal(fullBody, &temp2)
            checkError(err, "fullItem Unmarshal")
            fullItem := mapToStruct(temp2, it)

            logger.Debug(fmt.Sprintf("Saving item %+v", target))
            // Save to DB
            saveItem(fullItem, settings)

        }
    } 

    return nil
}

func syncAll(tr *http.Transport, settings common.ServiceConfiguration) {

    logger, _ := syslog.New(syslog.LOG_ERR, "SyncClient")
    defer logger.Close()

    defer func() {
            if e := recover(); e != nil {
                    logger, _ := syslog.New(syslog.LOG_ERR, "SyncClient")
                    defer logger.Close()
                    switch x := e.(type) {
                    case error:
                            err := x
                            logger.Err(fmt.Sprintf("Error: %s\n", err))
                    default:
                            err := fmt.Errorf("%v", x)
                            logger.Err(fmt.Sprintf("Error: %s\n", err))
                    }
            }
    }()

    logger.Info("Syncing wikis...")
    syncItems(tr, settings, common.Wiki{}, "/repl/wikis", "/repl/wiki/")
    logger.Info("Syncing pages...")
    syncItems(tr, settings, common.Page{}, "/repl/pages", "/repl/page/")
    logger.Info("Syncing contentfields...")
    syncItems(tr, settings, common.ContentField{}, "/repl/contentfields", "/repl/contentfield/")
    logger.Info("Syncing attachments...")
    syncItems(tr, settings, common.Attachment{}, "/repl/attachments", "/repl/attachment/")
    logger.Info("Syncing activities...")
    syncItems(tr, settings, common.Activity{}, "/repl/activities", "/repl/activity/")
    logger.Info("Syncing favoritewikis...")
    syncItems(tr, settings, common.FavoriteWiki{}, "/repl/favoritewikis", "/repl/favoritewiki/")
    logger.Info("Syncing page locks...")
    syncItems(tr, settings, common.Lock{}, "/repl/locks", "/repl/lock/")
}

func main() {
    // Load the settings
    var settings common.ServiceConfiguration
    settings.ConfigFrom("client.json")

    // Load my SSL key and certificate
    cert, err := tls.LoadX509KeyPair(settings.MyCertificateFile, settings.MyKeyFile)
    checkError(err, "LoadX509KeyPair")

    // Load the CA certificate for server certificate validation
    capool := x509.NewCertPool()
    cacert, err := ioutil.ReadFile(settings.CAKeyFile)
    checkError(err, "loadCACert")
    capool.AppendCertsFromPEM(cacert)

    // Prepare config and transport
    config := tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: capool}
    tr := &http.Transport{
        TLSClientConfig: &config,
    }

    // Prepare timer
    ticker := time.NewTicker(time.Duration(settings.Interval) * time.Second)
    quit := make(chan struct{})

    go func() {
        for {
           select {
            case <- ticker.C:
                // Attempt to synchronize all...
                syncAll (tr, settings)

            case <- quit:
                ticker.Stop()
                return
            }
        }
     }()


    select {}
}

