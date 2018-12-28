# OpenLDAP
OpenLDAP setup from Docker container.

# Run it

```
docker run --env LDAP_ORGANISATION="Solidly" \
  --env LDAP_DOMAIN="solidly.io" \
  --env LDAP_ADMIN_PASSWORD="ADMIN_PASSWORD" \
  --detach --publish 389:389/tcp --publish 636:636/tcp \
  osixia/openldap:1.2.2
```

# Search it

`docker exec SILLY_CONTAINER_NAME ldapsearch -x -H ldap://localhost -b dc=solidly,dc=io -D "cn=admin,dc=solidly,dc=io" -w ADMIN_PASSWORD`

# Install OpenLDAP password policy overlay

## Step-1 : Enable ppolicy overlay
Note: you will have to be local to the ldap daemon (since the command uses local auth).

`ldapadd -Y EXTERNAL -H ldapi:/// -f /etc/ldap/schema/ppolicy.ldif`

## Step-2 Add oupolicy.ldif OU

`ldapadd -D cn=admin,dc=solidly,dc=io -w ADMIN_PASSWORD -f ldap/policy/oupolicy.ldif`

## Step-3 load the module ppmodule.ldif
Note: you will have to be local to the ldap daemon (since the command uses local auth).

`ldapadd -Y EXTERNAL -H ldapi:/// -f ldap/policy/ppmodule.ldif`

## Step-4 Add ppolicyoverlay.ldif using ldapadd command

`ldapadd -Y EXTERNAL -H ldapi:/// -f ldap/policy/ppolicyoverlay.ldif`

## Step-5 Add passwordpolicy.ldif in LDAP

`ldapadd -D cn=admin,dc=solidly,dc=io -w ADMIN_PASSWORD -f ldap/policy/passwordpolicy.ldif`

## Step-6 add the OU users_ou.ldif

`ldapadd -D cn=admin,dc=solidly,dc=io -w ADMIN_PASSWORD -f ldap/policy/userou.ldif`

## Step-7 Add a user in ldap and generate the passwoed

`ldapadd -D cn=admin,dc=solidly,dc=io -w ADMIN_PASSWORD -f ldap/policy/users.ldif`

`ldappasswd -D cn=admin,dc=solidly,dc=io -w ADMIN_PASSWORD 'uid=testuser,ou=users,dc=solidly,dc=io'`

## Step-8 Verify configuration

`ldapsearch -Y EXTERNAL -H ldapi:/// -b olcDatabase={1}mdb,cn=config`