## Code Review
In today's episode, I shall do a code review of [heptio's ark project](https://github.com/heptio/ark). 

### Resources
This is in-line with kubernetes pattern. All of the resource schemas are defined at:

- `https://github.com/heptio/ark/blob/master/pkg/apis/`

#### Backup Specifications
```yaml
ref: https://github.com/heptio/ark/blob/master/pkg/apis/ark/v1/backup.go

- It represents a config to be used during backup
- It can also be treated as a:
  - `template` or 
  - `a set of strategies` or
  - `a set of rules` that defines the backup's implementation
- It provides hooks to different phases of taking a backup
- Hooks provide a way to inject custom behaviour
- It defines various BackUp lifecycle phases such as:
  - InProgress
  - Failed
  - New
  - Deleting
  - FailedValidation
- It has a well defined specification of BackupStatus e.g.:
  - backup format version
  - backup expiration i.e. time to garbage collection the backup
  - audit fields
```

#### Backup Storage Location Specifications
```yaml
ref: https://github.com/heptio/ark/blob/master/pkg/apis/ark/v1/backup_storage_location.go

- Object store is currently the only supported type to store a backup
- The rest of the specification inherits the object store model e.g.:
  - Bucket
  - Provider
  - Provider Config
  - Provider Location Availability
  - Provider Location Access Mode
```

#### Constants vs. Labels & Annotations
```yaml
ref: https://github.com/heptio/ark/blob/master/pkg/apis/ark/v1/constants.go
ref: https://github.com/heptio/ark/blob/master/pkg/apis/ark/v1/labels_annotations.go

- This has the project level constants e.g.
  - namespace where ark components will get installed
  - label key(s) e.g. 'ark.heptio.com/restore-name'
```

#### Backup Request Specifications
```yaml
ref: https://github.com/heptio/ark/blob/master/pkg/apis/ark/v1/delete_backup_request.go

- This has the spec related to backup request argument
- It has a List of BackUpRequest as well that can be used for bulk deletion of backups
- There is also a BackupRequestStatus which indicates the response of backup request
- This seems to be a kubernetes custom resource
  - If this is used as a custom resource then an audit of all requests can be done
  - kubectl can be used to get a history of requests as well as individual status
  - If reconcile logic or operator implements the custom resource then:
    - the request gets handled eventually i.e. asynchronous
```

#### Download Request Specifications
```yaml
ref: https://github.com/heptio/ark/blob/master/pkg/apis/ark/v1/download_request.go

- This is spec of download request
- It understand various type of download requests e.g:
  - BackupContents
  - BackupLogs
  - RestoreLogs
  - RestoreResults
```

#### Restore Specifications
```yaml
ref: https://github.com/heptio/ark/blob/master/pkg/apis/ark/v1/restore.go

- It understands restore from backup name
- It can restore from a particular backup's schedule name
- It has various options / rules to include exclude resources for restore
- It can restore PV from snapshot (via cloudprovider)
```

#### Others
```yaml
- Pod Volume Backup Specs
- Pod Volume Restore Specs
- Restic Repository Specs
```
