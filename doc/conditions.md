# Conditions

A condition is a block of parameters that evaluates to a single number. This output is always a number, even if boolean operators are used.

## Table of operators

Operator | Description | Precedence
--- | --- | ---
== | Equality, checks if two numbers are equal, evalutes to 0 or 1 | 0
!= | Inequality, checks if two numbers are not equal, evaluates to 0 or 1 | 0
< | Less than | 0
> | Greater than | 0
<= | Less than or equal | 0
>= | Greater than or equal | 0
+ | Addition | 1
- | Subtraction | 1
/ | Division | 2
* | Multiplication | 2
% | Modulo, returns the remainder of a division | 2
** | To the power of, e.g. 4**2 = 4² | 3
- | Unary minus, negates a number | 4


Brackets are supported, expressions are evaluated using bracket and precedence rules as defined by the table above (higher means first).