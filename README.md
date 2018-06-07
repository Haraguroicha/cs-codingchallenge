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

### Project structure

* `conf`

    there only a configure file named `config.yaml` inside, the config file structure may reference to `Configs/Configs.go`

* `Configs`

    the config structure, and config store

* `Error`

    the error message structure and its messages

* `static`

    frontend used files, e.g. css and javascript files

* `templates`

    the frontend page templates

* `Topic`

    the topic data structure and vote data structure are inside

* `Utilities`

    map, reduce and filter utility functions, in this project only used map function

* `vendor`

    the packages which we used inside

### APIs

* `GET /api/getTopics/*page`

    for getting the topics by paging, default is the first page

    Response will likes following

    ```json
    {
        "data": [
            {
                "topicId": 0,
                "topicTitle": "Sample",
                "votes": {
                    "upVotes": 0,
                    "downVotes": 0,
                    "sumVotes": 0
                }
            }
        ],
        "success": true
    }
    ```

* `POST /api/newTopic`

    for submit the new topic, the request body need following structure

    ```json
    {
        "topicTitle": "Title"
    }
    ```

    And Success will Response likes following

    ```json
    {
        "data": [
            {
                "topicId": 0,
                "topicTitle": "Title",
                "votes": {
                    "upVotes": 0,
                    "downVotes": 0,
                    "sumVotes": 0
                }
            }
        ],
        "success": true
    }
    ```

    And exceeded 255 characters will failed likes following

    ```json
    {
        "title": "ExceededTopicLengthError",
        "message": "The Topic length exceeded to 255 characters",
        "success": false
    }
    ```

* `POST /api/upVote/:topic/*page`

    submit a UpVote to specified topic id and response specified page to get topics

    Success and error response reference to `getTopics`

    addition error message by topic id was not found are following

    ```json
    {
        "title": "NoTopicError",
        "message": "The Topic Not Found for TopicID xxx",
        "success": false
    }
    ```

    and vote success but page invalid are following

    ```json
    {
        "title": "PageInvalidError",
        "message": "There is an invalid Page of xxx",
        "success": false
    }
    ```

* `POST /api/downVote/:topic/*page`

    submit a DownVote to specified topic id and response specified page to get topics

    Success and error response reference to `upVote`

### Implementation and design

1. Only sorting during add new topic or vote some topic
2. No need to sort the data when we get the topics data
3. Get the first 20 is just slice first 20 topics data is okay
4. Topic Title length and Topics per page is configurable by `conf/config.yaml`

### Knowing issues

1. No Lock implementation for add or voting the topics data during its sorting by previous request
2. Underlying the `sort` package to sort the topics data may bound the performance during large data transation
