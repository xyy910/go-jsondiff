## JsonDiff Tool 

The main purpose of the tool is calculating the differences between two jsons, 
and return a json format result.

code reference: [jsondiff](https://github.com/nsf/jsondiff)

The standard library nsf/jsondiff can return readable and colorful
json diff results for console and html, but what I need is a json
format result, so I reformed the library. 

TEST

```
$ git clone https://github.com/xyy910/go-jsondiff.git

$ cd go-jsondiff

$ go run main.go

```

USAGE

```
JSON1: 
{
    "list":[
        "a",
        "b",
        "c"
    ],
    "same":"same",
    "xfy":123,
    "xyf":{
        "disanceng":{
            "c1ab":"c111",
            "c1cd":"c211"
        },
        "hoby":"dance",
        "milk":"多多"
    },
    "xyff":"aichifan"
}

JSON2:
{
    "hahaha":"笑点低",
    "list":[
        "a",
        "n",
        "c"
    ],
    "same":"same",
    "xyf":{
        "disanceng":{
            "c2ab":"c111",
            "c2cd":"c211"
        },
        "hoby":"chifan",
        "milk":"不多"
    },
    "xyff":"aichifan"
}


diff result:

{
    "hahaha":"+笑点低",
    "list":[
        "a",
        "b => n", //changed
        "c"
    ],
    "same":"same",
    "xfy":"-123",
    "xyf":{
        "disanceng":{
            "c1ab":"-c111", //removed
            "c1cd":"-c211",
            "c2ab":"+c111", //added
            "c2cd":"+c211"
        },
        "hoby":"dance => chifan",
        "milk":"多多 => 不多"
    },
    "xyff":"aichifan"
}

```

For ease of use，the params are two map[string]interface{}, and the result is bytes[],
more details refer to main.go
