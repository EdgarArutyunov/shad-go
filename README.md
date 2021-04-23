# shad-go
Тут я ботаю шадовские задачи по го


- group:    Low level
  tasks:
    - task: illegal
      score: 100
    - task: blowfish
      score: 100
    - task: tarstreamtest
      score: 100
      watch:
        - distbuild
    - task: filecachetest
      score: 100
      watch:
        - distbuild
    - task: artifacttest
      score: 100
      watch:
        - distbuild

- group:    Reflect
  tasks:
    - task: reversemap
      score: 100
    - task: jsonlist
      score: 100
    - task: jsonrpc
      score: 100
    - task: structtags
      score: 100

- group: April Fools' Day
  tasks:
    - task: foolsday1
      score: 0
    - task: foolsday2
      score: 0
    - task: foolsday3
      score: 0

- group:    HTTP
  tasks:
    - task: urlshortener
      score: 100
    - task: digitalclock
      score: 100
    - task: coverme
      score: 200
    - task: olympics
      score: 200
    - task: firewall
      score: 200

- group:    Dockertest
  tasks:
    - task: dockertest
      score: 0

- group:    Concurrency with shared memory
  tasks:
    - task: dupcall
      score: 200
    - task: keylock
      score: 200
    - task: batcher
      score: 200
    - task: pubsub
      score: 300

- group:    Testing
  tasks:
    - task: cond
      score: 100
    - task: testequal
      score: 100
    - task: fileleak
      score: 100
    - task: tabletest
      score: 100
    - task: tparallel
      score: 200

- group:    Goroutines
  tasks:
    - task: tour1
      score: 100
    - task: once
      score: 100
    - task: rwmutex
      score: 100
    - task: waitgroup
      score: 100
    - task: ratelimit
      score: 100

- group:    "[HW] Gitfame"
  tasks:
    - task: gitfame
      score: 0

- group:    Interfaces
  tasks:
    - task: otp
      score: 100
    - task: lrucache
      score: 100
    - task: externalsort
      score: 100
    - task: retryupdate
      score: 100

- group:    Basics
  tasks:
    - task: hotelbusiness
      score: 100
    - task: hogwarts
      score: 100
    - task: utf8
      score: 100
    - task: varfmt
      score: 100
    - task: speller
      score: 100
    - task: ciletters
      score: 100
    - task: forth
      score: 100

- group:    Hello World
  tasks:
    - task: sum
      score: 100
    - task: tour0
      score: 100
    - task: wordcount
      score: 100
    - task: urlfetch
      score: 100
    - task: fetchall
      score: 100
