# Scanner

Does not emit any diagnostics as we are not yet sure at this level whether we're parsing markdown element or vocab elements.

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