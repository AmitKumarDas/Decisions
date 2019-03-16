### Meta Info
- Version - 1
- Last Updated On - 16-Mar-2019

### Motivation
What things in software are most common? What is the most simplistic version of any software product. IMO it is its CLI.
Somehow I tend to agree that given a CLI that autocompletes & provides me the help at the terminal itself is the coolest
thing and at the same time the most most simplest thing that has ever happened to software. Recently same can be told for
a declarative specification that can be saved, versioned and injected into workflows.

The motivation of this article is to find a balance between action oriented CLI to specification driven yamls.

### Behaviours == Actions == Verbs
I shall take the example of openebs storage as the software. What can be various actions related to openebs? We can think
of below:
- create
- delete
- get
- identify
- list
- resize
- backup
- snapshot
- etc...

### Nouns == Entities
Lets think of various nouns used in openebs:
- cstor volume
- cstor replica
- cstor target
- cstor pool
- jiva volume
- jiva target
- jiva pool
- etc...

### Mix them up == CLI
- create cstor volume
- delete cstor pool
- resize cstor volume
- backup cstor volume
- get cstor volume
- etc...

### Devil == CLI in YAML
This is considered as the devil by most fellow programmers. However, IMO if it helps users to get the requirement done
in a static compiled language and can be saved & versioned then why not. Being declarative is not a option but a need.
In other words, there would not have any `github actions` to begin with

### Prototype
