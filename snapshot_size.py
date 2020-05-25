import boto3
import argparse

def snapshots_size(ebs, snap1, snap2):
    '''
    small function to separate the snapshot size calculation
    return raw size difference between 2 snapshots
    params:
    :param ebs: the boto3 ebs client object
    :type ebs: boto3 session.client('ebs')
    :param snap[12]: snapshot ids
    :type snap[12]: string
    :type return: int 
    '''
    blocks = ebs.list_changed_blocks(FirstSnapshotId = snap1, SecondSnapshotId = snap2)
    size = len(blocks['ChangedBlocks'])
    block_size = blocks["BlockSize"]
    return size * block_size


def format_bytes(size):
    '''
    quick convert an integer into a more human readable format
    returns a string human readable size
    :param size: the "raw" size to conver
    :type size: int
    :type return: string
    '''
    power = 2**10
    n = 0
    power_labels = {0 : '', 1: 'K', 2: 'M', 3: 'G', 4: 'T'}
    while size > power:
        size /= power
        n += 1
    return f'{size} {power_labels[n]}b'


def vol_snapshots_size(**kwargs):
    '''
    Main function - will print the size of each snapshots based on a volume id
    Arguments:
    :param session: the boto3 session to use for api calls
    :type session: boto3.Session
    :param region: the aws region to make the calls against
    :type region: string
    :param volume_id: the volume id to query snapshots size
    :type volume_id: string
    '''
    ec2 = session.client('ec2', region_name=kwargs['region'])
    ebs = session.client('ebs', region_name=kwargs['region'])
    snapshots = ec2.describe_snapshots(
                    Filters=[
                        {
                            'Name': 'volume-id',
                            'Values': [
                                kwargs['volume_id'],
                            ]
                        }
                    ],
                    OwnerIds=['self']
                )
    sorted_snapshots = sorted(snapshots['Snapshots'], key=lambda snap: snap['StartTime'], reverse=True)
    for i in range(len(sorted_snapshots)-1):
        snap = sorted_snapshots[i]
        diff_size = snapshots_size(ebs, sorted_snapshots[i]["SnapshotId"], sorted_snapshots[i+1]["SnapshotId"])
        print(f'Date: {snap["StartTime"]} - ID: {snap["SnapshotId"]} - Size: {format_bytes(diff_size)}')


def arg_parse():
    '''
    argument parser function
    returns argparse object
    '''
    parser = argparse.ArgumentParser(description='''
                    Get snapshots size based on an original volume ID''',
                                     epilog='''
                    Example:
                    python[3] snapshot_size.py [-r ca-central-1] [-v vol-foo42bar31baz]
                            ''')
    parser.add_argument('-r','--region',
                        metavar='Region Name',
                        type=str,
                        default='us-east-1',
                        help='AWS region name')
    parser.add_argument('-p','--profile',
                        metavar='Profile Name',
                        type=str,
                        default='production',
                        help='AWS credentials profile name to use')
    parser.add_argument('-v','--volume-id',
                        metavar='Volume ID',
                        type=str,
                        required=True,
                        help='Volume ID used to fetch snapshots that are based on it')
    return parser.parse_args()


if __name__ == '__main__':

    # get arguments
    args = arg_parse()

    # create session
    session = boto3.Session(profile_name=args.profile)
    vol_snapshots_size(session=session, region=args.region, volume_id=args.volume_id)
