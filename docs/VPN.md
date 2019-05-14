# VPN

OpenVPN setup from Docker container.

## Generate Configuration

```bash
docker run -v $OVPN_DATA:/etc/openvpn \
  --rm godinezj/openvpn-ldap ovpn_genconfig \
  -u tcp://solidly.io:443
```

## Initialize PKI

```bash
docker run -v $OVPN_DATA:/etc/openvpn \
    --rm -it godinezj/openvpn-ldap ovpn_initpki
```

## Configure LDAP

```bash
docker run -v $OVPN_DATA:/etc/openvpn \
    --rm -it godinezj/openvpn-ldap ovpn_config_ldap \
    ldap://192.168.1.15 "cn=admin,dc=solidly,dc=io" "8RlDnnSb1Kce"
```

## Run the Container

Note: container always exposes port 1194, regardless of the specified protocol:

```bash
docker run -v $OVPN_DATA:/etc/openvpn \
    -d -p 443:1194/tcp \
    --cap-add=NET_ADMIN godinezj/openvpn-ldap
```

## Generate Client Configuration

```bash
docker run -v $OVPN_DATA:/etc/openvpn \
  --log-driver=none --rm -it godinezj/openvpn-ldap easyrsa \
  build-client-full solidclient2 nopass
```

## Retrieve the Client Config

```bash
docker run -v $OVPN_DATA:/etc/openvpn \
  --log-driver=none --rm \
  godinezj/openvpn-ldap ovpn_getclient solidclient2 > solidclient2.ovpn
```
