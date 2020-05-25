# aws_snapshot_size
Small script to get snapshot size based on a volume id

The script takes an aws volume id as argument at minima.
It will output the list of snapshot based on this volume
and the size of each of the snapshots using boto3 ebs.list_changed_blocks function.

## usage

```
$ python volume_snapshot_size.py -h
usage: volume_snapshot_size.py [-h] [-r Region Name] [-p Profile Name] -v
                               Volume ID

Get snapshots size based on an original volume ID

optional arguments:
  -h, --help            show this help message and exit
  -r Region Name, --region Region Name
                        AWS region name
  -p Profile Name, --profile Profile Name
                        AWS credentials profile name to use
  -v Volume ID, --volume-id Volume ID
                        Volume ID used to fetch snapshots that are based on it

Example: python[3] snapshot_size.py [-r ca-central-1] [-v vol-foo42bar31baz]
```

Example
```
$ python volume_snapshot_size.py -r us-east-1 -p dev -v vol-42123456789
Date: 2020-05-24 09:31:17.040000+00:00 - ID: snap-0cbd26ced80126b16 - Size: 5.5 Mb
Date: 2020-05-24 07:53:01.767000+00:00 - ID: snap-0d27999a3afffc6db - Size: 65.0 Mb
```
