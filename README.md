# aws_snapshot_size
Small script to get snapshot size based on a volume id

The script takes an aws volume id as argument at minima.
It will output the list of snapshot based on this volume
and the size of each of the snapshots using boto3 ebs.list_changed_blocks function.

## usage

```
$ python volume_snapshot_size.py -h
usage: volume_snapshot_size.py [-h] [-r Region Name] [-p Profile Name] -v
                               Volume ID [-o Output format]

Get snapshots size based on an original volume ID

optional arguments:
  -h, --help            show this help message and exit
  -r Region Name, --region Region Name
                        AWS region name
  -p Profile Name, --profile Profile Name
                        AWS credentials profile name to use
  -v Volume ID, --volume-id Volume ID
                        Volume ID used to fetch snapshots that are based on it
  -o Output format, --output Output format
                        The output format in: [text, json]

Example: python[3] snapshot_size.py [-r ca-central-1] [-v vol-foo42bar31baz]
```

Example
```
# with python
$ python volume_snapshot_size.py -r us-east-1 -p dev -v vol-42123456789
Date: 2020-05-24 09:31:17.040000+00:00 - ID: snap-0cbd26ced80126b16 - Size: 5.5 Mb
Date: 2020-05-24 07:53:01.767000+00:00 - ID: snap-0d27999a3afffc6db - Size: 65.0 Mb

---
# with go
./volume_snapshot_size -r us-east-1 -p dev -v vol-42123456789 -o json
{"snapshot-date":"2020-05-25T12:32:25.06Z","snapshot-id":"snap-0cbd26ced80126b16","snapshot-size":"3.0 MiB"}
{"snapshot-date":"2020-05-25T07:46:02.454Z","snapshot-id":"snap-0d27999a3afffc6db","snapshot-size":"41.5 MiB"}
```

## build the go version

This can be done via the small bash [script](build_binary.sh)
Remove the line with `GOOS=darwin` if not compiling for mac
Otherwise just
```
go get ./...
go build -o volume_snapshot_size volume_snapshot_size.go

# or with the dockerfile
bash build_binary.sh volume_snapshot_size
```
