Thumuht 后端服务

Run `godoc -http=:6060` to read document. ( if you installed godoc )

TODO:
  foreign key via hooks
  upload & download files
  
curl localhost:8899/query   -F operations='{ "query": "mutation($req: Upload!) { fileUpload(input: {parentId: 1, parentType: post, upload: $req}) }", "variables": { "req": null } }'   -F map='{ "0": ["variables.req"] }'   -F 0=@amdmanual.pdf