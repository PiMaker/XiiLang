# Statements

Statements are the core principle of XiiLang. Every line holds one statement.
A statement always has the following format: ``` <statement> [parameter]* ```
Below is a list of all statements available.

## number

Format: ``` number <varname> ```
The ```number``` statement creates a new 64-bit float number and initializes it with 0. This command is used in conjunction with the set command and conditions to evaluate dynamic expressions using the concept of variables.
The command takes one parameter, which is the variable name that will be given to the new variable.

## string

Format: ``` string <varname> ```
Creates a new literal variable. See ```number``` above for more info about variables.

## in

Format: ``` in <varname> ```
The ```in``` statement reads input from the user using stdin. It takes the name of a previously created variable as its only parameter.

## out

Format: ``` out [parameter]* ```
The out statement is used for user output. You can use literals surrounded by double quotes (```"```), numbers, variables or any combination of these as parameters. Output will be formated according to the passed type.
Example: ``` out "The number " 6 " is " x ``` where x is an initialized variable

## if

Format: ``` if <condition> ```
The if statement is used for branching. It supports the standard condition syntax defined in "conditions.md". The following block (ended by the ```end``` statement) is only executed if the condition passed as parameters evaluates to something other than 0. If this is not the case, code execution will continue immediately after the next ```end```.

## while

Format: ``` while <condition> ```
The while command is used for loops. It works the same as the ```if``` statement above, with the exception that after reaching the next ```end``` block, the condition is checked again and if it still holds true the block will execute from the beginning again.

## end

Format ``` end ```
The end statement does not take any parameters. It is only used in conjunction with ```if``` and ```while```. For good readability it is recommended that a block between ```if```/```while``` and ```end``` is indented.