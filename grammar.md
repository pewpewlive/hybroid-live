$$
\begin{align}

    [term] &\to <+> or <-> \\ 

    [factor] &\to <*> or </> \\
    \\

    [stmt] &\to 
        \begin{cases}
            \text{var declaration}\\
            \text{functions} \\
            \text{if} \\
            \text{else} \\
        \end{cases}\\

    [prog] &\to 
        \begin{cases}
            [stmt] \\
            ...
        \end{cases} \\
    
    \\

    [expr] &\to 
        \begin{cases}
            [binExpr] \\
            [lit] \\
            [unary] \\
            [group] \\
            [ident] 
        \end{cases} \\

    [unary] &\to <operand>[expr] \\
    [group]^* &\to 
        \begin{cases}
            [Expr] 
        \end{cases} \\

    [binExpr] &\to 
        \begin{cases}
            [expr] <term> [expr] \\
            [expr] <factor> [expr] \\
        \end{cases} \\

    \\

    &\small{\text{*group are used for higher precedence}}
\end{align}