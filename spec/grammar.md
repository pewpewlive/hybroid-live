# Hybroid grammar (LaTeX)

A _program_ is made up of _statements_. _Statements_ are: variable declarations, function definitions, if statements, directive and use statements.

_Expressions_ are made up of several sub-expression types. There are: binary, literal, unary, group and identifier expressions. Each have different use cases.

There are terms and factors. A _term_ can either be an addition, or subtraction operation. A _factor_ can either be a multiplication or division operation.

$$
\begin{align}
    [prog] &\to
        \begin{cases}
            [stmt] \\
            ...
        \end{cases} \\
    [stmt] &\to
        \begin{cases}
            \text{let}\\
            \text{pub}\\
            \text{fn} \\
            \text{if} \\
            \text{else} \\
            \text{dir} \\
            \text{use}
        \end{cases}\\
    [expr] &\to
        \begin{cases}
            [bin] \\
            [lit] \\
            [unary] \\
            [group] \\
            [ident]
        \end{cases} \\
    [unary] &\to \; <operand>[expr] \\
    [group]^* &\to
        \begin{cases}
            [expr] \\
            ...
        \end{cases} \\
    [term] &\to \; <+> or <-> \\
    [factor] &\to \; <*> or </> \\
    [bin] &\to
        \begin{cases}
            [expr] <term> [expr] \\
            [expr] <factor> [expr] \\
        \end{cases} \\
    \\
    &\small{\text{*group are used for higher precedence}}
\end{align}
$$
