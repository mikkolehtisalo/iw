package main
import (
        "fmt"
        "crypto/tls"
        "crypto/rand"
        "crypto/x509"
        "time"
        "io/ioutil"
        "net/http"
        "net"
        "encoding/json"
        _ "github.com/lib/pq"
        "github.com/jmoiron/sqlx"
        "regexp"
        "reflect"
        "iw/replication/src/common"
        "log/syslog"
)

func checkError(err error, s string) {
    // Syslogger
    logger, _ := syslog.New(syslog.LOG_ERR, "SyncServer")
    defer logger.Close()

    if err != nil {
            logger.Err(fmt.Sprintf("%s: %s", s, err))
            panic(fmt.Sprintf("%s: %s", s, err))
    }
}

func myHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "ALIVE")
}

func buildListJSON(rows *sqlx.Rows, tgt interface{}) ([]byte, error) {
    tgtType := reflect.ValueOf(tgt).Type()
    targets := reflect.MakeSlice(reflect.SliceOf(tgtType), 0, 10)
    for rows.Next() {
        target := reflect.New(tgtType).Interface()
        err := rows.StructScan(target)
        checkError(err, "rows.Scan")
        targets = reflect.Append(targets, reflect.ValueOf(target).Elem())
    }
    checkError(rows.Err(), "rows.Err")
    jsontext, err := json.Marshal(targets.Interface())
    return jsontext, err
}

func buildJSON(rows *sqlx.Rows, tgt interface{}) ([]byte, error) {
    tgtType := reflect.ValueOf(tgt).Type()
    target := reflect.New(tgtType).Interface()
    for rows.Next() {
        err := rows.StructScan(target)
        checkError(err, "rows.Scan")
    }
    checkError(rows.Err(), "rows.Err")
    jsontext, err := json.Marshal(target)
    return jsontext, err
}

func uniHandler(w http.ResponseWriter, r *http.Request, t interface{}, l bool) {
    // Handle panics - for logging mostly
    defer func() {
            if e := recover(); e != nil {
                    logger, _ := syslog.New(syslog.LOG_ERR, "SyncServer")
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

    // Load settings
    var settings common.ServiceConfiguration
    settings.ConfigFrom("server.json")

    // Check the peer
    err := common.CheckPeer(r, settings)
    checkError(err, "checkPeer")

    // Open database connection
    connstring := fmt.Sprintf("user=%s password='%s' dbname=%s sslmode=disable", settings.DBUser, settings.DBPassword, settings.DBName)
    db, err := sqlx.Open("postgres", connstring)
    checkError(err, "sql.Open")
    defer db.Close()

    // Query for data
    var rows *sqlx.Rows

    if l {
        // Lists
        switch t.(type) {
        case common.Wiki:
            rows, err = db.Queryx("select wiki_id, modified from wikis")
        case common.Page:
            rows, err = db.Queryx("select page_id, wiki_id, modified from pages")
        case common.ContentField:
            rows, err = db.Queryx("select contentfield_id, wiki_id, modified from contentfields")
        case common.FavoriteWiki:
            rows, err = db.Queryx("select username, wiki_id, modified from favoritewikis")
        case common.Attachment:
            rows, err = db.Queryx("select attachment_id, wiki_id, modified from attachments")
        case common.Activity:
            rows, err = db.Queryx("select activity_id from activities")
        case common.Lock:
            rows, err = db.Queryx("select target_id, wiki_id from locks")
        }
    } else {
        // Single objects
        switch t.(type) {
        case common.Wiki:
            re := regexp.MustCompile("/repl/wiki/([0-9a-f-]+)/([0-9-+:.T ]+)")
            vars := re.FindStringSubmatch(r.URL.Path)
            modified, terr := time.Parse(time.RFC3339Nano, vars[2])
            checkError(terr, "Unable to parse wiki time")
            rows, err = db.Queryx("select wiki_id, title, description, create_user, readacl, writeacl, adminacl, status, modified " + 
                "from wikis " +
                "where wiki_id=uuid_in($1) and modified=$2", vars[1], modified)
        case common.Page:
            re := regexp.MustCompile("/repl/page/([0-9a-f-]+)/([0-9a-f-]+)/([0-9-+:.T ]+)")
            vars := re.FindStringSubmatch(r.URL.Path)
            modified, terr := time.Parse(time.RFC3339Nano, vars[3])
            checkError(terr, "Unable to parse page time")
            rows, err = db.Queryx("select page_id, wiki_id, path, title, create_user, readacl, writeacl, adminacl, stopinheritation, index, depth, status, modified " +
                "from pages " +
                "where page_id=uuid_in($1) and wiki_id=uuid_in($2) and modified=$3", vars[1], vars[2], modified)
        case common.ContentField:
            re := regexp.MustCompile("/repl/contentfield/([0-9a-f-]+)/([0-9a-f-]+)/([0-9-+:.T ]+)")
            vars := re.FindStringSubmatch(r.URL.Path)
            modified, terr := time.Parse(time.RFC3339Nano, vars[3])
            checkError(terr, "Unable to parse contentfield time")
            rows, err = db.Queryx("select contentfield_id, wiki_id, content, modified, status, create_user "+
                "from contentfields "+
                "where contentfield_id=uuid_in($1) and wiki_id=uuid_in($2) and modified=$3", vars[1], vars[2], modified)
        case common.FavoriteWiki:
            re := regexp.MustCompile("/repl/favoritewiki/([^/]+)/([0-9a-f-]+)/([0-9-+:.T ]+)")
            vars := re.FindStringSubmatch(r.URL.Path)
            modified, terr := time.Parse(time.RFC3339Nano, vars[3])
            checkError(terr, "Unable to parse favoritewiki time")
            rows, err = db.Queryx("select username, wiki_id, modified, status "+
                "from favoritewikis "+
                "where username=$1 and wiki_id=uuid_in($2) and modified=$3", vars[1], vars[2], modified)
        case common.Attachment:
            re := regexp.MustCompile("/repl/attachment/([0-9a-f-]+)/([0-9a-f-]+)/([0-9-+:.T ]+)")
            vars := re.FindStringSubmatch(r.URL.Path)
            modified, terr := time.Parse(time.RFC3339Nano, vars[3])
            checkError(terr, "Unable to parse attachment time")
            rows, err = db.Queryx("select attachment_id, wiki_id, encode(attachment, 'base64') as attachment, mime, filename, create_user, modified, status "+
                "from attachments "+
                "where attachment_id=uuid_in($1) and wiki_id=uuid_in($2) and modified=$3", vars[1], vars[2], modified)
        case common.Activity:
            re := regexp.MustCompile("/repl/activity/([0-9a-f-]+)")
            vars := re.FindStringSubmatch(r.URL.Path)
            rows, err = db.Queryx("select activity_id, timestamp, user_id, user_name, activity_type, target_type, target_title, target_id, readacl, writeacl, adminacl " + 
                "from activities " +
                "where activity_id=uuid_in($1)", vars[1])
        case common.Lock:
            re := regexp.MustCompile("/repl/lock/([0-9a-f-]+)/([0-9a-f-]+)")
            vars := re.FindStringSubmatch(r.URL.Path)
            rows, err = db.Queryx("select target_id, wiki_id, username, realname, modified " + 
                "from locks " +
                "where wiki_id=uuid_in($1) and target_id=($2)", vars[1], vars[2])
        }
    }

    checkError(err, "db.Query")

    // Build the JSON
    var jsontext []byte

    if l {
        jsontext, err = buildListJSON(rows, t)
    } else {
        jsontext, err = buildJSON(rows, t)
    }

    checkError(err, "json.buildJSON")

    // Print the JSON
    fmt.Fprintf(w, "%s", jsontext)

}

func main() {
        var settings common.ServiceConfiguration
        settings.ConfigFrom("server.json")
        
        // Load my SSL key and certificate
        cert, err := tls.LoadX509KeyPair(settings.MyCertificateFile, settings.MyKeyFile)
        checkError(err, "LoadX509KeyPair")

        // Load the CA certificate for client certificate validation
        capool := x509.NewCertPool()
        cacert, err := ioutil.ReadFile(settings.CAKeyFile)
        checkError(err, "loadCACert")
        capool.AppendCertsFromPEM(cacert)

        // Prepare server configuration
        config := tls.Config{Certificates: []tls.Certificate{cert}, ClientCAs: capool, ClientAuth: tls.RequireAndVerifyClientCert}
        config.NextProtos = []string{"http/1.1"}
        config.Rand = rand.Reader

        // Web server and the handler methods on the paths
        myTLSWebServer := &http.Server{Addr: settings.ServiceAddress, TLSConfig: &config, Handler: nil}
        http.HandleFunc("/", myHandler)
        http.HandleFunc("/repl/wikis", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.Wiki{},true) })
        http.HandleFunc("/repl/wiki/", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.Wiki{},false) })
        http.HandleFunc("/repl/pages", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.Page{},true) })
        http.HandleFunc("/repl/page/", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.Page{},false) })
        http.HandleFunc("/repl/contentfields", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.ContentField{},true) })
        http.HandleFunc("/repl/contentfield/", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.ContentField{},false) })
        http.HandleFunc("/repl/favoritewikis", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.FavoriteWiki{},true) })
        http.HandleFunc("/repl/favoritewiki/", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.FavoriteWiki{},false) })
        http.HandleFunc("/repl/attachment/", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.Attachment{},false) })
        http.HandleFunc("/repl/attachments", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.Attachment{},true) })
        http.HandleFunc("/repl/activity/", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.Activity{},false) })
        http.HandleFunc("/repl/activities", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.Activity{},true) })
        http.HandleFunc("/repl/lock/", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.Lock{},false) })
        http.HandleFunc("/repl/locks", func(w http.ResponseWriter, r *http.Request) { uniHandler(w,r,common.Lock{},true) })

        // Bind to port
        conn, err := net.Listen("tcp", settings.ServiceAddress)
        checkError(err, "Listen")

        // Start the web server and serve until the doomsday
        tlsListener := tls.NewListener(conn, &config)
        err = myTLSWebServer.Serve(tlsListener)
        checkError(err, "Serve")
}