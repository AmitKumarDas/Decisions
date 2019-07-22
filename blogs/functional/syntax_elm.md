```elm
-- learn syntax from https://elm-lang.org/docs/syntax
```

```elm
-- lists
-- all are same
[1,2,3,4]
1 :: [2,3,4]
1 :: 2 :: 3 :: 4 :: []
```

```elm
-- conditional behavior based on the structure of algebraic data types and literals
case maybeList of
  Just xs -> xs
  Nothing -> []

case xs of
  [] ->
    Nothing
  first :: rest ->
    Just (first, rest)

case n of
  0 -> 1
  1 -> 1
  _ -> fib (n-1) + fib (n-2)
```

```elm
-- record
-- create records
origin = { x = 0, y = 0 }
point = { x = 3, y = 4 }

-- access fields
origin.x == 0
point.x == 3

-- field access function
List.map .x [ origin, point ] == [ 0, 3 ]

-- update a field
{ point | x = 6 } == { x = 6, y = 4 }

-- update many fields
{ point | x = point.x + 1, y = point.y + 1 }
```

```elm
-- functions
square n =
  n^2

hypotenuse a b =
  sqrt (square a + square b)

distance (a,b) (x,y) =
  hypotenuse (a - x) (b - y)
```

```elm
-- anonymous functions
square =
  \n -> n^2

squares =
  List.map (\n -> n^2) (List.range 1 100)
```

```elm
viewNames1 names =
  String.join ", " (List.sort names)

viewNames2 names =
  names
    |> List.sort
    |> String.join ", "

-- (arg |> func) is the same as (func arg)
-- Just keep repeating that transformation!
-- if you have written functions within go template then its similar
-- good that here indentation is possible
```

```elm
-- let lets you do extra stuff
-- notice the destructuring assignment

let
  ( three, four ) =
    ( 3, 4 )

  hypotenuse a b =
    sqrt (a^2 + b^2)
in
hypotenuse three four
```

```elm
-- more let
let =
 twentyFour =
  3 * 8
 
 sixteen =
  4 ^ 2
 
in
twentyFour + sixteen
```

```elm
-- let with type info

let
  name : String
  name =
    "Hermann"

  increment : Int -> Int
  increment n =
    n + 1
in
increment 10
```

```elm
-- modules

module MyModule exposing (..)

-- qualified imports
import List                            -- List.map, List.foldl
import List as L                       -- L.map, L.foldl

-- open imports
import List exposing (..)              -- map, foldl, concat, ...
import List exposing ( map, foldl )    -- map, foldl

import Maybe exposing ( Maybe )        -- Maybe
import Maybe exposing ( Maybe(..) )    -- Maybe, Just, Nothing
```

```elm
-- type annotations
-- IMO this is IMP

answer : Int
answer =
  42

factorial : Int -> Int
factorial n =
  List.product (List.range 1 n)

distance : { x : Float, y : Float } -> Float
distance {x,y} =
  sqrt (x^2 + y^2)
```

```elm
-- type aliases
-- IMO this is IMP

type alias Name = String
type alias Age = Int

info : (Name,Age)
info =
  ("Steve", 28)

type alias Point = { x:Float, y:Float }

origin : Point
origin =
  { x = 0, y = 0 }
```

```elm
-- custom types

type User
  = Regular String Int
  | Visitor String
```
