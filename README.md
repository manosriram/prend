# prend (vAlpha)
vendoring tool for protocol buffers. customizable vendoring for proto files, 

# prend config
All the proto file source must be defined in the **prend.yaml** file

**prend.yaml example**
```yaml
name: test config name
description: test config description
sources:
  - repo_url: https://github.com/manosriram-youtube/reddit_backend.git
    paths:
    - source_proto_path: gateway/protos
      destination_proto_path: vendors/gateway_protos
    - source_proto_path: post_service/protos
      destination_proto_path: vendors/post_service_protos
  - repo_url: https://github.com/manosriram-youtube/grpc-heartbeat.git
    paths:
    - source_proto_path: heartbeat_pb
      destination_proto_path: pb/heartbeat_protos
```
`repo_url`: the `.git` url of the git repository

`source_proto_path`: the proto file(s) path in the repository

`destination_proto_path`: path where the fetched proto files are to be placed

Proto file vendoring

## Installing
1. `go install github.com/manosriram/prend`
2. `prend fetch` to fetch all protos from config


TODO
- [ ] ssh support
- [ ] `prend cleanup` command
- [ ] commit hash for each prend unique fetch
- [ ] other version control support (tentative)
- [ ] cleanup vendor folder before fetch
