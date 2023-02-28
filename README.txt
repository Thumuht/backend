Thumuht 后端服务

TODO:
  foreign key via hooks
  upload & download files
  
curl localhost:8899/query   -F operations='{ "query": "mutation($req: Upload!) { fileUpload(input: {postId: 1, upload: $req}) }", "variables": { "req": null } }'   -F map='{ "0": ["variables.req"] }'   -F 0=@amdmanual.pdf