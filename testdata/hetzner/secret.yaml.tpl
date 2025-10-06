apiVersion: v1
kind: Secret
metadata:
  name: hetzner
data:
  token: $HETZNER_TOKEN_BASE64
