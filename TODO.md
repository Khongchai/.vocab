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
        You're testing harvest
        - [x] Assert that: given known inputs/outputs map (from sm2), word tree produces the correct remaining time that matches the inputs/outputs map
        - [ ] Then make lsp work
            - [x] Get go debugger to work on windows with vscode
            - [x] Fix the error
            - [ ] starting a new section does not get rid of the error...
                ```
                20/05/2025
                > (it) la magia, bene, scorprire
                lskjfljalfkjsl jlsjf ljsdf 
                23/10/2025
                >> (it) la magia
                laskdjflkasjdf lksajdf
                ```
        - [x] Handle document deletion
            - [x] see main.go
        - [ ] When diagnostics disappear, should still return registered documents.
        - [x] Now with pull mode, diagnostics of even just one file is not updated correctly...?
        - [ ] Then also check diagnostics of two files.
        - [x] comment on vocab line not handled properly
        - [x] Error not re-added when changing in one document.
- [x] tree.replace method for only replacing the ast of certain files!
- [x] test word with special char for parser and compiler
- [x] Multi-threading
    - [x] Go routine dispatch parsing
- [x] While harvesting, do a global lock (nothing should change!)
- [ ] There are still duplicate diagnostics...somehow
- [x] Normalization change
    - [x] move normlaization to word tree
    - [x] make normalization turn to lower case only for italian
    - [x] normalization should strip definite and indefinite articles of both languages.
- [-] Support hover.
    - [-] Write test for hover
    - [x] Implement hover logic (Pick)
        - [-] continue from here "textDocument/hover": func(rm lsproto.RequestMessage) (any, err
    - [-] Connect hover to lsp

## Last two...then done
- [x] Trigger whole workspace root parse immediately upon opening any .vocab file.
- [ ] Harvest command for collecting all words needs review and create a new section. This would need a workspace command.
    - [x] Collecting from all files and this file
    - [x] Collect params and collect response
    - [ ] At this point, we'll probably need to refactor the stuff inside main.go into something more structured
    
## Cleanup
- [ ] Lsproto position seems wonky now? (could be related to the next problem)
    - This is not related to walkDir, but just how errors are somehow mixed when all files diagnostics are joined in the end.
    - Does not happen with one file. 
- [ ] File scheme for windows when collecting all diagnostics on startup incorrect:
    Expect file:///c%3A/Users/world/Desktop/vocab/test.vocab
    Got file://c:\\Users\\world\\Desktop\\vocab\\test.vocab
## Bonus
- [ ] Syntax highlighting
```
20/05/2025
> (it) thing

20/06/2025
> (it) thing
```
- [ ] Syntax highlighting
- [ ] Lemmatization
- [ ] Inline word highlighting
- [ ] Error offset wrong
    -	16/10/2025
		>> (de) gewöhnlich, ewig | right here

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
- [x] Don't clear diagnostics on document close. We need to keep it for spaced repetition later.

# Main Requirement 1.5 

- [ ] Token colors

# Main Requirement 2

- [x] Handle multiple opened files in parallel (great opportunity to try out go's parallel power)
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
- [ ] Audio Pronunciation Integration: Add a command or CodeLens link next to words that, when clicked, fetches and plays the pronunciation (using an online API like Forvo or browser speech synthesis).
- [ ] Statistics Dashboard: Create a custom webview panel within VS Code that shows learning statistics:
Number of words learned per language.
Number of words due today/this week.
A graph showing learning progress over time.
Words causing the most difficulty (lowest average grades).
Audio Pronunciation Integration: Add a command or CodeLens link next to words that, when clicked, fetches and plays the pronunciation (using an online API like Forvo or browser speech synthesis).
