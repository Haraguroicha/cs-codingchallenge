# cs-codingchallenge

[![](https://travis-ci.org/Haraguroicha/cs-codingchallenge.svg?branch=master)](https://travis-ci.org/Haraguroicha/cs-codingchallenge)

## Digg / Reddit clone with upvote and downvotes

demo site is available at https://carousellcodingchallenge.herokuapp.com/

### Features

1. Any users can submit a topic, which is simply a string that does not exceed 255 characters
2. Any users can upvote or downvote a topic
3. Only top 20 topics will display on home page, which sorted by upvotes, descending
4. In-memory data structure, shared by the same process and without using data persistence

### Dependencies

1. go
2. dep
3. heroku CLI
4. make

### How to local run

1. Clone this project by `git clone https://github.com/Haraguroicha/cs-codingchallenge`
2. Make symbolic link to `$GOPATH` by following commands

    ```bash
    [[ -z "$GOPATH" ]] && export GOPATH=$HOME/go
    mkdir -p $GOPATH/src/github.com/Haraguroicha
    ln -s $PWD/cs-codingchallenge $GOPATH/src/github.com/Haraguroicha/cs-codingchallenge
    ```

3. Start local web server for debugging by `make debug`
4. Browse at http://localhost:5000

### APIs

TODO
