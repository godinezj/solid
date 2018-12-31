# LDAP

OpenLDAP setup from Docker container.

## Run it

```bash
docker run --env LDAP_ORGANISATION="Solidly" \
  --env LDAP_DOMAIN="solidly.io" \
  --env LDAP_ADMIN_PASSWORD="ADMIN_PASSWORD" \
  --detach --publish 389:389/tcp --publish 636:636/tcp \
  osixia/openldap:1.2.2
```

## Search it

```bash
docker exec SILLY_CONTAINER_NAME ldapsearch -x -H ldap://localhost \
    -b dc=solidly,dc=io -D "cn=admin,dc=solidly,dc=io" -w ADMIN_PASSWORD
```

## Install OpenLDAP Password Policy Overlay

The commands below should be run from the top level directory.

### Enable ppolicy overlay

Note: you will have to be local to the ldap daemon (since the command uses local auth).

```bash
ldapadd -Y EXTERNAL -H ldapi:/// -f /etc/ldap/schema/ppolicy.ldif
```

### Add oupolicy.ldif OU

```bash
ldapadd -D cn=admin,dc=solidly,dc=io -w ADMIN_PASSWORD -f ldap/policy/oupolicy.ldif
```

### Load the Module ppmodule.ldif

Note: you will have to be local to the ldap daemon (since the command uses local auth).

```bash
ldapadd -Y EXTERNAL -H ldapi:/// -f ldap/policy/ppmodule.ldif
```

### Add ppolicyoverlay.ldif using ldapadd command

```bash
ldapadd -Y EXTERNAL -H ldapi:/// -f ldap/policy/ppolicyoverlay.ldif
```

### Add passwordpolicy.ldif in LDAP

```bash
ldapadd -D cn=admin,dc=solidly,dc=io -w ADMIN_PASSWORD -f ldap/policy/passwordpolicy.ldif
```

### Add the OU users_ou.ldif

```bash
ldapadd -D cn=admin,dc=solidly,dc=io -w ADMIN_PASSWORD -f ldap/policy/userou.ldif
```

### Add a user in ldap and generate the passwoed

```bash
ldapadd -D cn=admin,dc=solidly,dc=io -w ADMIN_PASSWORD -f ldap/policy/users.ldif
```

```bash
ldappasswd -D cn=admin,dc=solidly,dc=io -w ADMIN_PASSWORD 'uid=testuser,ou=users,dc=solidly,dc=io'
```

### Verify Configuration

```bash
ldapsearch -Y EXTERNAL -H ldapi:/// -b olcDatabase={1}mdb,cn=config
```