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
  - Infos as a list of successful messages
  - Warns as a list of warning messages
  - Skips as a list of skipped messages
  - Errors as a list of typed errors
- Lots of RunTasks are needed
- Executing RunTasks conditionally
- Streaming template functions
  - Map invocation against a list
  - Filter invocation against a list
  - Use of Predicate
- Template Functions without the use of method calls
  - {{- $myNewRTObjList := $myRTObjList | Map | Filter p | Any p -}}
  - {{- $myRTObjList | Map | Filter p | Any p | Store id123 -}}
  - {{- $aMap := .Stores.id123 | Items -}}
  - {{- $myRTObjList | Select ".spec.abc" ".spec.def" ".spec.xyz" -}}
  - {{- $myRTObj | Select ".spec.abc" ".spec.def" ".spec.xyz" -}}
- Values as feeds from user, runtime, engine, etc
- Stores as storage during execution

### Things that went good !!!
- Template Functions
- Pipes
