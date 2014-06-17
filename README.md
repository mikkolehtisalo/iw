Wiki. Kind of.
==============

Handles more like a client application, replication, authenticaton and acl systems. Target group is organization internal usage, writing notes and such.

Features:
* GSSAPI/Kerberos authentication
* LDAP backend for user & group information
* Three role (read, write, admin) ACL system for content, with inheritation support
* Simple macros (TOC, include children pages)
* ExtJS based user interface
* Non-destructive versioning data model
* Wysiwyg html editor, improved image uploading & gallery
* Input validation, including html whitelisting
* Rudimentary replication between instances

[Mandatory screenshot][1]

Installation
============

Prerequisities
* PostgreSQL
* Working Kerberos setup
* LDAP server 
* Golang
* Memcached recommended

Most important things
---------------------

This should pull the required packages

```
go get github.com/mikkolehtisalo/iw
```

Change app.secret from app.conf

```
app.secret=OWd2Zg7A5wWomHFWkMNIwGNvS7qCAbGY8NKrSADg50bAaSU4hXemJTslFVV3Ah3Q
```

Database
--------

Create the user, database, and import the schema

```
createuser wiki
createdb -U wiki --owner=wiki wiki
psql wiki wiki
\i db/schema.sql
```

Then review app.conf

```
# DB
db.name=wiki
db.user=wiki
db.password=password
```

LDAP
----

Create bind account, and add photo;binary for users' avatar images.

Review app.conf, and fill up the LDAP server information, filters and the attributes
```
# ldap
ldap.server=freeipa.localdomain
ldap.port=389
ldap.user_base=cn=users,cn=accounts,dc=localdomain
ldap.user_filter=(&(uid=*)(objectClass=inetUser))
ldap.user_uid_attr=uid
ldap.user_cn_attr=cn
ldap.user_photo_attr=photo;binary
ldap.user_group_attr=memberOf
ldap.group_filter=(&(cn=*)(objectClass=groupOfNames))
ldap.group_cn_attr=cn
ldap.group_dn_attr=dn
ldap.user=uid=admin,cn=users,cn=accounts,dc=localdomain
ldap.passwd=password
ldap.group_base=cn=groups,cn=accounts,dc=localdomain
ldap.group_regexp=cn=([^,]+)
```

Kerberos
--------

Make sure /etc/krb5.conf is sane, for example

```
[logging]
 default = FILE:/var/log/krb5libs.log
 kdc = FILE:/var/log/krb5kdc.log
 admin_server = FILE:/var/log/kadmind.log

[libdefaults]
 default_realm = LOCALDOMAIN
 dns_lookup_realm = false
 dns_lookup_kdc = true
 rdns = false
 ticket_lifetime = 24h
 forwardable = yes
 allow_weak_crypto = false

[realms]
 LOCALDOMAIN = {
  kdc = freeipa.localdomain:88
  master_kdc = freeipa.localdomain:88
  admin_server = freeipa.localdomain:749
  default_domain = localdomain
  pkinit_anchors = FILE:/etc/ipa/ca.crt
}

[domain_realm]
 .localdomain = LOCALDOMAIN
 localdomain = LOCALDOMAIN

[dbmodules]
  LOCALDOMAIN = {
    db_library = ipadb.so
  }

```

Create the keytab

```
# Create the principal on kdc
kadmin
addprinc -randkey HTTP/dev.localdomain@LOCALDOMAIN

# Add it to the server's /etc/krb5.keytab
kadmin
ktadd HTTP/dev.localdomain@LOCALDOMAIN
```

Replication
-----------

Replication directory contains client and server. They are very KISS, but working. Build them. 

Accompany with server/client.json file for settings. Create the mentioned CA certificate, and certificates for all servers. Configure valid peers, and other settings. Start the servers and clients. 

All logging goes to syslog.

Starting up
-----------

Kick the tires by

```
revel run github.com/mikkolehtisalo/iw dev 80
```

The good and the bad
====================

Most parts are pretty solid. Would need test(s/ing) to get rid of some corner case bugs, and proper documentation.

It's not a ready product, but wouldn't take much to polish it.

[1]:https://raw.githubusercontent.com/mikkolehtisalo/iw/master/docs/screenshot.png
