### CASTemplate Engine Design - Chapter 2
Exposing CASTemplate along with RunTasks enables exposing logic in a declarative fashion. It helped the team in numerous ways.
One of the typical use of these templates was the ability to embed logic into RunTasks and solve the problem statement without
taking a dip into compilation and so on. There as other rosy aspects with the use of template. However, other side of the coin
has a lot of pain points. In this design which will form the basis of next CAST engine, I will list down my thoughts of what
should go in to make the engine better.

### Things to improve !!!
- Avoid as much as possible the need to write specific CAST engines
  - Volume, Clone, Pool, Upgrade and so on needs to be written to influence generic CAST engine
- Ability to Unit Test a RunTask
  - This helps in verifying the validity of the RunTask with various data options
- Avoid as much as possible the need to write conditional logic
- Avoid as much as possible the need to write meta info in RunTasks
- Avoid get from specific paths and setting at specific paths
  - User/Dev/Operator does not care about this
- Errors are hard to debug
  - Improvements:
    - Infos as a list of successful runtasks
    - Warns as a list of warning runtasks
    - Skips as a list of skipped runtasks
    - Errors as a list of errored runtasks
- Lots of RunTasks are needed
- Executing RunTasks conditionally
- Streaming based template functions
  - Map invocation against a list
  - Filter invocation against a list
  - Use of Predicate
- Template Functions without the use of method calls
- All validations should be triggered before invoking all/any runtasks
  - data related validations
  - template related validations
  - old vs. new templating style

### Things that went good !!!
- Template Functions
- Pipes

### New Design
- Values as feeds from user, runtime, engine, etc
- Stores as storage during execution of one or more runtasks
- Docs stores the yamls as unstruct instances indexed with name of unstruct
  - i.e. `[]*unstruct`
```yaml
runs:
  - yaml: |
      kind: cool
      apiVersion: v1
    if: 
    run: |
    - {{- .Docs.abc-123 | -}}
    - {{- $myNewRTObjList := $myRTObjList | Map | Filter p | Any p -}}
    - {{- $myRTObjList | Map | Filter p | Any p | Store id123 -}}
    - {{- $aMap := .Stores.id123 | Items -}}
    - {{- $myRTObjList | Select ".spec.abc" ".spec.def" ".spec.xyz" -}}
    - {{- $myRTObj | Select ".spec.abc" ".spec.def" ".spec.xyz" -}}
    onErr: 
onErr:
```
