# Exercise #1: Quiz Game

Run timed quizzes via the command line. [Exercise details](https://github.com/gophercises/quiz)

## How to run

```
$ go build . && ./quiz
```

```
Usage of ./quiz:
    -csv string
        a csv file in the format of 'question,answer' (default "problems.csv")
    -limit int
        the time limit for the quiz in seconds (default 30)
```