## Générer une clé privée :

$openssl genpkey -algorithm RSA -out server.key

Cela générera une clé privée nommée server.key

## Créer une demande de certificat

$openssl req -new -key server.key -out server.csr

## Auto-signer le certificat avec la clé privée pour obtenir le certificat :

$openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
