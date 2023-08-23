# Kubernetes CronJob Example

This job has been configured to run daily at 1 AM.

You need to change:
- user-auth.env
  - update the values for `USERNAME`, `PASSWORD`, and `PROVIDER` to your values.
- volume.yml
  - configure `nfs.path` and `nfs.server` to your values.

### Optional:

You may wish to add another volume mount targeting `${HOME}/.tadpoles-backup` for
the account that the job runs under. This will allow caching of event data and
speed up the job (tadpoles only).

### Notes:

The download command will scan any files present in the target backup directory
to prevent duplication. You should not remove files from this location to speed
up subsequent download operations.
