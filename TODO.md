<!-- skip -->
# Main Requirement 0
- [x] Make extension plugin start server up
- [x] Make go server return all red text upon first key.
- [x] Refactor into an engine
- [x] Basic data structure for
- [x] Use go structs for request and responses
- [x] Update scanner test case and implementation now that vocab is its own language.
    - [x] Make scanner recognize date expression.
    - [x] Remove all markdown reference
- [x] Continue writing parser, all syntactic errors are thrown here.
    - [x] Finish implementing first parser version
    - [x] Write test for all the tiny cases
# 04/10/2025
    - [x] Write small cases
    - [x] Fix utterance parsing
    - [x] Stop including parens in language identifier!
    - [x] Make TestFullSectionParsing then 
    - [x] Then compiler!
        - [x] The compiler should be incremental in that 
            - [x] it accepts ast and turn it into an IR tree -- a hashmap of words to the date section and location / file they appear in. 
            - [x] The IR trees can be compiled and produce diagnostics independently and then merge. Every time they merge, new diagnostics should be produced based on newly available information. This means multicore-power!
    - [x] Words need to be graded  (writing parser test)
        You're implementing harvest
        - [ ] Then go back to implementing the harvest function!


```markdown
# 04/09/2025
>> (it) `inoltre`
> (de) `schön`, der Berg
Inoltre, questo plugin sarà fantastico. Sono sicuro.
```

# Main Requirement 1

Basic vocab capture

```markdown
# 04/09/2025
>> (it) `inoltre`
> (de) `schön`, der Berg
Inoltre, questo plugin sarà fantastico. Sono sicuro.
```
- [ ] Match absolute text within > or >> section, if those matched text does not appear in the following section, underline the correct point
- [ ] Don't clear diagnostics on document close. We need to keep it for spaced repetition later.

# Main Requirement 1.5 

- [ ] Token colors

# Main Requirement 2

- [ ] Handle multiple opened files in parallel (great opportunity to try out go's parallel power)
We can build the ast trees -> word tree in parallel
aggregate to final single word tree
Then fork again and apply in parallel sm2 

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

**Spaced repetition (compilation)** 

This requires [`interFileDependencies`](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#diagnosticOptions) server capability.

# Main requirement 6

**Sync to phone app and once reviewed make a commit to update**

# Possible Requirement 1

Markiert es wenn ein neues Wort nicht in seiner Infinitivform ist:

```
# 18/09/2025
> entrambe // wurde rot markiert weil die Infinitivform "Entrambi" ist.
Devi fare entrambe le cose.
```
# Side quests
- [ ] Command click for word occurences.
- [ ] Add a comment case.
- [ ] Hover to show definition in English
- [ ] Parallelize parsing of multiple vocab files with goroutine (see ts-go).
- [ ] Make pull mode work
- [ ] Show how much time remaining for each individual word.
