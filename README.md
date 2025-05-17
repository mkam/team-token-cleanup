# Terraform Team Token Cleanup

This script deletes Terraform team tokens based on the criteria provided via the command line flags.
By default, the criteria used to delete the team tokens is:

- The token has expired.
- The token has not been used in more than 30 days.

The script requires two confirmations for deletion, providing the `--delete` flag and then inputting `y` or `yes`.

## Usage

### Configuration

The following configurations are required via environment variables:

```
$ export TFE_ORGANIZATION=<my-org-name>
$ export TFE_TOKEN=<my-token>
```

### Running

Execute the cleanup with by building the script and running it with the optional arguments.

```
$ go build
$ ./tfc_cleanup --delete
$ ./tfc_cleanup --delete --team my-team --last-used-days-ago 14
```

```
Usage of ./team-token-cleanup:
  -created-at-days-ago int
      Duration of time in days for how long ago a resource should have been created before deleting.
  -delete
      Deletes the team tokens that fit the provided criteria for deletion. Defaults to false.
  -expired
      Marks expired tokens for deletion, regardless of created_at or last_used_at. (default true)
  -last-used-days-ago int
      Duration of time in days for how long ago a resource should have been last used before deleting. (default 30)
  -team string
      The team name to delete tokens for. If not provided, tokens from all teams will be considered for deletion.
```

Or run via go directly:

```
$ go run main.go
$ go run main.go --delete
```
