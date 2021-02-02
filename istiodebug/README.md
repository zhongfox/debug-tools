docker build -t zhongfox/istiodebug .

grpc-dump --port 8080 --interface 0.0.0.0 --destination istiod-1-8-1.istio-system:15010 --log_level debug -proto_roots ./api | jq -Rrc '. as $line | try (fromjson ) catch $line'

grpc-dump --port 8080 --interface 0.0.0.0 --destination istiod-1-8-1.istio-system:15010 --log_level debug -proto_roots ./api \
| jq -Rr '. as $line | try (fromjson | del(.node)) catch $line'


