# Terraform Team Token Cleanup

This script deletes Terraform team tokens based on the criteria provided via the command line flags.
By default, the criteria used to delete the team tokens is:

- The token has expired.

OR

- The token has never been used or has not been used in more than 30 days AND the token is older than 7 days.

These criteria can be modified via command line flags. Additionally, the tokens can be filtered to
a specific team using the team name with `--team`. The script requires two confirmations for deletion,
providing the `--delete` flag and then inputting `y` or `yes`. See [Command Line Flags](#command-line-flags)
for more details.

## Usage

### Configuration

The following configurations are required via environment variables:

```
$ export TFE_ORGANIZATION=<my-org-name>
$ export TFE_TOKEN=<my-token>
```

### Command Line Flags

```
Usage of ./team-token-cleanup:
  -created-at-days-ago int
      Duration of time in days for how long ago a resource should have been created before deleting. Used in conjunction with last-used-days-ago, where the token is deleted only if it is old enough AND has not been used recently. Set to 0 to ignore created time. (default 7)
  -delete
      Deletes the team tokens that fit the provided criteria for deletion. Defaults to false.
  -expired
      Marks expired tokens for deletion, regardless of created_at or last_used_at. (default true)
  -last-used-days-ago int
      Duration of time in days for how long ago a resource should have been last used before deleting. Used in conjunction with created-at-days-ago, where the token is deleted only if it has not been used recently AND is old enough. Set to 0 to ignore last used time. (default 30)
  -team string
      The team name to delete tokens for. If not provided, tokens from all teams will be considered for deletion.
```


### Running

Execute the cleanup with by building the script and running it with the optional arguments.

```
$ go build
$ ./team-token-cleanup --delete
$ ./team-token-cleanup --delete --team my-team --last-used-days-ago 14
```

Or run via go directly:

```
$ go run main.go
$ go run main.go --delete
```

Example Output:

```
./team-token-cleanup
Marking token for deletion because expired: 'old token' in team 'old-team' expired_at=2025-05-01 15:22:05.571 +0000 UTC
Marking token for deletion because last used and created too long ago: 'CI token' in team 'ci-team' last_used_at=0001-01-01 00:00:00 +0000 UTC
Marking token for deletion because last used and created too long ago: 'at-Ry2qKmUvTa3DnBoL' in team 'ci-team' last_used_at=0001-01-01 00:00:00 +0000 UTC
Marking token for deletion because last used and created too long ago: 'at-Xn7pGcQoZrB8LsMd' in team 'test-team' last_used_at=0001-01-01 00:00:00 +0000 UTC
Marking token for deletion because because expired: 'test token' in team 'test-team' expired_at=2025-03-14 13:01:09.460 +0000 UTC
Marking token for deletion because last used and created too long ago: 'prod token' in team 'prod-team' last_used_at=0001-01-01 00:00:00 +0000 UTC
Marking token for deletion because last used and created too long ago: 'test token 2' in team 'test-team' last_used_at=0001-01-01 00:00:00 +0000 UTC

7 tokens marked for deletion.
Use the --delete flag to delete the tokens that fit the specified criteria.
```

```
$ ./team-token-cleanup --delete
Marking token for deletion because expired: 'old token' in team 'old-team' expired_at=2025-05-01 15:22:05.571 +0000 UTC
Marking token for deletion because last used and created too long ago: 'CI token' in team 'ci-team' last_used_at=0001-01-01 00:00:00 +0000 UTC
Marking token for deletion because last used and created too long ago: 'at-Ry2qKmUvTa3DnBoL' in team 'ci-team' last_used_at=0001-01-01 00:00:00 +0000 UTC
Marking token for deletion because last used and created too long ago: 'at-Xn7pGcQoZrB8LsMd' in team 'test-team' last_used_at=0001-01-01 00:00:00 +0000 UTC
Marking token for deletion because because expired: 'test token' in team 'test-team' expired_at=2025-03-14 13:01:09.460 +0000 UTC
Marking token for deletion because last used and created too long ago: 'prod token' in team 'prod-team' last_used_at=0001-01-01 00:00:00 +0000 UTC
Marking token for deletion because last used and created too long ago: 'test token 2' in team 'test-team' last_used_at=0001-01-01 00:00:00 +0000 UTC

7 tokens marked for deletion.
Are you sure you want to delete these team tokens? (y/n):
y
Deleting token: at-Fk3dJvNpXzB7LmTw (expiration date)
Deleting token: at-Vc6gZnLpWqT9HbEk (CI token)
Deleting token: at-Ry2qKmUvTa3DnBoL
Deleting token: at-Xn7pGcQoZrB8LsMd
Deleting token: at-Lz4mHvEtWi6NkPyQ (test token)
Deleting token: at-Dk9sTfYmXa1VbReO (prod token)
Deleting token: at-Nj2hPxKoLd7MwBiT (test token 2)
Team tokens deleted.
```
