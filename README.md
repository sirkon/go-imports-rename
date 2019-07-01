# imports-rename
Tool to change import paths

#Usage examples

* 
    ```shell script
    go-imports-rename 'github.com/rsz/ => github.com/rs/' 
    ```
    will list all possible import path changes with given prefix `github.com/rsz/` in Go files in current directory and 
    all its subdirectories. No saves will be done.
     
* Use `--save` flag to commit these changes:
    ```shell script
    go-imports-rename --save 'github.com/rsz/ => github.com/rs/' 
    ```
* Use `--root` flag to diagnose or make changes in specific directory:
    ```shell script
    go-imports-rename --root $GOPATH/src 'github.com/rsz/ => github.com/rs/' 
    ```
    Will diagnose changes in `$GOPATH/src`
* Use `--regexp` flag to use regular expression for import path changes:
    ```shell script
    go-imports-rename --regexp '^gen/(.*)$ => gitlab.example.com/common/schema/$1' 
    ```
    will diagnose possible `gen/common` and similar imports changes into `gitlab.example.com/common/schema/common`, etc
     
