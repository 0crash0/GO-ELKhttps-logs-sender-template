# GO-ELKhttps-logs-sender-template
 Go template for sending zap logs over http



how to get ssl fingerprint:
`openssl s_client \
-connect demos.es.us-east1.gcp.elastic-cloud.com:9243 \
-servername demos.es.us-east1.gcp.elastic-cloud.com \
-showcerts < /dev/null 2>/dev/null | \
openssl x509 -in /dev/stdin -sha256 -noout -fingerprint | \
sed 's/://g'  `
