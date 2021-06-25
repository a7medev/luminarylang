# Luminary ☄

A very minimalist programming language built just for fun!

## Docs

### 1. Data Types

Luminary has a few set of data types, which are numbers (which includes booleans), string, functions, lists and null

```
1.5                 # Number
"Luminary"          # String
fun() = "Hello"     # Function
["A", "B", "C"]     # List
null                # Null
```

### 2. Comments

Comments in Luminary are created using the # symbol followed by any text

```
# This is a comment
```

### 3. Vairables

You can declare/assign vairables using the assignment operator `=`

```
name = "Luminary"   # This is a variable
```

### 4. Functions

You can declare functions in Luminary using the `fun` keyword

```
# Lambda function
fun sum(x, b) = x + b

# Regular/Multi-statement function
fun sum(x, b) {
  println(x + b)    # Function call - NOTE: println() is a builtin function in Luminary
  return x + b      # Return a value from the function
}
```

Functions are values, which means you can pass it as an argument to another function, store it in a variable, etc

```
myFun = fun(x, b) = x + b      # Store a function into a variable
println(myFun(1, 5))

# Pass a function as an argument
newList = map(
  [1, 2, 3],
  fun(val, index) = index + ": " + val
)
```

### 5. Lists

Lists are just a list of data which can store any data types in it

```
list = [
  54,
  "A",
  "C",
  fun() = "Welcome"
]

println(list)
```

### 6. Null

Null is a value that means nothing or an empty value

```
null     # This is null
```

### 7. Booleans

Booleans are just the values of `true` or `false`, they are represented in Luminary as `1` for `true` and `0` for false, they can be used in control flows for example

```
true     # This is a boolean value of true
false    # This is a boolean value of false
```

### 8. Comparison operators

Comparison operators are just operators which are resolved to a boolean value based on thier truthy

```
a == b
a != b
a > b
a >= b
a <= b
not a == b
a and b
a or b
```

### If statements

If statements are used to execute some code if a condition is true

```
if condition {
  # do something ...
} elif otherCondition {
  # do another thing ...
} else {
  # do another thig ...
}
```

A condition can be a boolean value or a variable that stores a boolean value (which contains comparison operators as they are resolved to a boolean value)

### For loops

For loops are used to execute some code repeatedly

```
# Execute this code for a variable i set from 0 to 20
for i = 0 : 20 {
  # do something
}

# Execute this code for a variable i set from 0 to 20 increasing by 2 in each iteration
for i = 0 : 20 by 2 {
  # do something
}
```

### While loops

While loops are used to execute some code repeatedly while a condition is true

```
while condition {
  # do something
}
```

### Break/Continue statements

- Break statement is used inside a loop to the execution of it

- Continue statement is used inside a loop to stop the execution of the current iteration and continue to the next one

```
while condition {
  # do something
  if condition {
    # do something
    break
  }
  # do something
  if condition {
    # do something
    continue
  }
  # do something
}
```

### Builtin Functions

There are some builtin functions in Luminary, which are:

#### 1. print(...values)

Which takes any number of arguments and prints them to the terminal

#### 2. println(...values)

Which takes any number of arguments and prints them to Stdout and adds a new line after then

#### 3. scan(prompt)

Which takes an optional argument of type string and scans a string value the user enters in Stdin

#### 4. len(value)

Which takes one argument of type list or string and returns a number value of it's length

> There are other builtin functions that will be added soon to the documentation
