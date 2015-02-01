ImageStore  [![GoDoc](https://godoc.org/github.com/ian-kent/imagestore?status.svg)](https://godoc.org/github.com/ian-kent/imagestore) [![Build Status](https://travis-ci.org/ian-kent/imagestore.svg?branch=master)](https://travis-ci.org/ian-kent/imagestore)
==========

An incredibly simple image store.

**WARNING** Don't run this in the wild! Images can be
uploaded and deleted by unauthenticated users.

Current features:
- Add files using HTTP POST
- Get files using HTTP GET
- Remove files using HTTP DELETE
- Images are stored in an S3 bucket

Requirements:
- Have AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables set
- Have the correct permissions for the bucket you're using
  - You can disable the "upload" endpoint by removing S3 bucket write permissions

Install imagestore:
```bash
go get github.com/ian-kent/imagestore
```

Start the server:
```bash
# Default port 5253
imagestore -bucket="bucket-name" -prefix="optional/prefix"

# Custom interface/port
imagestore -bind=:5353 -bucket="bucket-name" -prefix="optional/prefix"
```

Upload something:
```bash
curl -v -X POST --data-binary "@Makefile" http://localhost:5253/images/Makefile
```

Get something:
```bash
curl -v http://localhost:5253/images/Makefile
```

Delete something:
```bash
curl -v -X DELETE http://localhost:5253/images/Makefile
```

### Using Marathon

If you're using Mesos and Marathon, you can easily start imagestore:

Command:

`AWS_ACCESS_KEY_ID=your-id AWS_SECRET_ACCESS_KEY=your-key ./imagestore -bucket="bucket-name" -bind=:$PORT`

URI:

`https://github.com/ian-kent/imagestore/releases/download/v1.0.0/imagestore-1.0.0_linux_amd64.zip`

### Licence

Copyright ©‎ 2015, Ian Kent (http://www.iankent.eu).

Released under MIT license, see [LICENSE](LICENSE.md) for details.
