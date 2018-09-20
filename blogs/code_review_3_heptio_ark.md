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
- It defines various BackUp lifecycle phases such as InProgress, Failed, New, Deleting, FailedValidation
- It has a well defined specification of BackupStatus e.g. among other things BackupStatus understands:
  - backup format version
  - backup expiration i.e. time to garbage collection the backup
  - audit fields
```

#### Backup Storage Location Specifications
```yaml
ref: https://github.com/heptio/ark/blob/master/pkg/apis/ark/v1/backup_storage_location.go

- Object store is currently the only supported type to store a backup
- The rest of the specification inherits the object store model e.g.: one finds properties like
  - bucket, provider, provider config, 
  - provider location availability, provider location access mode
```
