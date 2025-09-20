# Terminals

- DateExpression: `nn/nn/nnnn` where `n` can be any number

- LanguageExpression: `\`(it)\`` with two values: de, and it for German and Italian respectively.

- `,`: separates reviewed or new vocabulary.

- New line: mark the end of a section (date, new vocab, reviewed vocab, and sentences).
 
# Scanner

Does not emit any diagnostics.

# Parser

Can emit some syntactic error, as we are at this stage quite sure whether we are really in a vocab section or not.

This will result in 100% confident level that it is a vocab section.

```
02/04/2025
> 
```
But this is just a todo section

```
02/04/2025

- [ ] Something

```


# Compiler

The compiler finds out diagnostics error from spaced-repetition.