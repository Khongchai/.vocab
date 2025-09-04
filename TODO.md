- [x] Make extension plugin start server up
- [x] Make go server return all red text upon first key.
- [ ] Refactor into an engine

# Main Requirement 0

- [ ] Basic data structure for

```markdown
# 04/09/2025
>> `inoltre`, something
> `meglio`
Inoltre, questo plugin sarà fantastico. Sono sicuro.
```

- [ ] Token colors

# Main Requirement 1

Basic vocab capture

```markdown
# 04/09/2025
>> `inoltre`
> `meglio`
Inoltre, questo plugin sarà fantastico. Sono sicuro.
```
- [ ] Match absolute text within > or >> section, if those matched text does not appear in the following section, underline the correct point
- [ ] Don't clear diagnostics on document close. We need to keep it for spaced repetition later.

# Main Requirement 2

- [ ] Handle multiple opened files in parallel (great opportunity to try out go's parallel power)

# Main Requirement 3 

**Remove articles**

- [ ] Given `la` `le` `die` `das`, etc, they should be skipped

```markdown
# 04/09/2025
> Das hund
// ...
```

# Main requirement 4

**Lemmatization**

```markdown
# 04/09/2025
> sono
Sarò lì
```
- [ ] Plugin should know `sarò` is future form of `sono` and match against that.

# Main requirement 5

**Spaced repetition** 

This requires [`interFileDependencies`](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#diagnosticOptions) server capability.


# Side quests
- [ ] Hover to show definition in English
- [ ] Parallelize parsing of multiple vocab files with goroutine (see ts-go).