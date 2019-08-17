# imports-rename
Tool to change import paths

## Important case

* Important case is to append suffix to existing path, i.e. to grow `gitlab.example.com/common/utils` into `gitlab.example.com/common/utils/v2`

    use 
    
    ```shell script
    go-imports-rename 'gitlab.example.com/common/utils/ ++' 
    ```
    to grow `gitlab.example.com/common/utils/package` into `gitlab.example.com/common/utils/v2/package`

## Usage examples

* Simple usage. Please do not use it for import paths upgrade from version below 2. Use `++` operator instead.
    ```shell script
    go-imports-rename 'github.com/rsz/ => github.com/rs/' 
    ```
    will list all possible import path changes with given prefix `github.com/rsz/` in Go files in current directory and 
    all its subdirectories. No saves will be done.
     
* Use `--save` flag to commit these changes:
    ```shell scriptgitk
    go-imports-rename --save 'github.com/rsz/ => github.com/rs/' 
    ```
* Use `--root` flag to diagnose or make changes in specific directory:
    ```shell script
    go-imports-rename --root $GOPATH/src 'github.com/rsz/ => github.com/rs/' 
    ```
    Will diagnose changes in `$GOPATH/src`
* There is a shortcut for the fresh v1.x.y ⇒ v2.α.β migrations:
    ```shell script
    go-imports-rename 'github.com/user/project ++' 
    ```
    It is an equivalent for
    ```shell script
    go-imports-rename 'github.com/user/project/ => github.com/user/project/v2/'
    ```
    In this case the utility takes care of paths and will not replace `github.com/user/projectNext`. Also, 
    `github.com/user/project/v2` won't change.
    You can also move from `v2` to `v3` in the same way
    ```shell script
    go-imports-rename 'github.com/user/project/v2 ++'
    ```
* You can jump through several migrations:
    ```shell script
    go-imports-rename 'github.com/user/project += 5'
    ```
    It is an equivalent for
    ```shell script
    go-imports-rename 'github.com/user/project/ => github.com/user/project/v6'
    ```
* Use `//` operator to imports rename with regular expression
    ```shell script
    go-imports-rename --regexp '^gen/(.*)$ // gitlab.example.com/common/schema/$1' 
    ```
    will diagnose possible `gen/common` and similar imports changes into `gitlab.example.com/common/schema/common`, etc
     