# Hybroid grammar (LaTeX)

A _program_ is made up of _statements_. _Statements_ are: variable declarations, function definitions, if statements, directive declarations, use statements as well as environment statement.

_Expressions_ are made up of several sub-expression types. There are: binary, literal, unary, group and identifier expressions. Each have different use cases.

There are terms and factors. A _term_ can either be an addition, or subtraction operation. A _factor_ can either be a multiplication or division operation.

$$
\begin{align}
    [prog] &\to
        \begin{cases}
            \text{stmt} \\
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
            \text{use} \\
            \text{fn}
        \end{cases}\\
    [expr] &\to
        \begin{cases}
            [\text{bin}] \\
            [\text{lit}] \\
            [\text{unary}] \\
            [\text{group}] \\
            [\text{ident}] \\
            \text{anon fn}
        \end{cases} \\
    [unary] &\to \; <operand>[expr] \\
    [group]^* &\to
        \begin{cases}
            [\text{expr}] \\
            ...
        \end{cases} \\
    [term] &\to \; <+> or <-> \\
    [factor] &\to \; <*> or </> \\
    [bin] &\to
        \begin{cases}
            [\text{expr}] <term> [\text{expr}] \\
            [\text{expr}] <factor> [\text{expr}] \\
        \end{cases} \\
    \\
    &\small{\text{*group are used for higher precedence}}
\end{align}
$$
