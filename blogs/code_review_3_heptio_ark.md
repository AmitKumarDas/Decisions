## Code Review
In today's episode, I shall do a code review of [heptio's ark project](https://github.com/heptio/ark). 

### Resources
This is in-line with kubernetes pattern. All of the resource schemas are defined at:

- `https://github.com/heptio/ark/blob/master/pkg/apis/`

#### BackupSpec
```
ref - https://github.com/heptio/ark/blob/master/pkg/apis/ark/v1/backup.go

- It represents a config to be used during backup
- It can also be treated as a:
  - `template` or 
  - `a set of strategies` or
  - `a set of rules` that defines the backup's implementation
- It provides hooks to different phases of taking a backup
- Hooks provide a way to inject custom behaviour
```
