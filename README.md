# Paca

Paca is a converter tool that takes Selenium IDE source code and
converts it into http://agouti.org/ tests.

The goal of this tool is to allow non developers to write browser tests
and have them converted into something that can be run against different
webdrivers/browsers and be maintained by developers.

Because Go is a compiled language and because the API is simple, a QA
person who might not yet be a developer should be able to tweak test
suites.

## Setup

You will need Go and the following 2 Go packages:

```
$ go get github.com/sclevine/agouti
$ go get github.com/mattetti/paca
```

This tool will convert your test cases but you will need to edit the
code and compile it again to run the test suite.

## Dependencies

You will need a webdriver, it's recommended that you install the
following:

```
$ brew install phantomjs
$ brew install chromedriver
$ brew install selenium-server-standalone
```

## Convert

```
$ go run cmd/main.go -source=<path of your Selenium IDE test case file>
```

This will create a `seltest` directory containing a test helper file and
the converted test case. However the test case won't be running
automatically, you will need to edit the helper.

Edit `helper_test.go` and set `TargetHost` to the host you want to test
against (note: you might want to use env variable to change/set that value).

You will notice a scenario called `TestFirstScenario`, feel free to
rename this scenario and at the bottom of the function's body call into
the test case you converted.
Imagine you converted a test case file called `login`, the converter
will have created a file called `Login_test.go` and a new function
called `Login(t *testing.T, page *agouti.Page)`. At the bottom of`TestFirstScenario`,
call this new function: `Login(t, page)`, the scenario will run and will
execute the test case.
You can create multiple scenarios and chain multiple test cases within
each. You can also add helper functions to `helper_test.go` (or a new
test file you add yourself) to generate random data for instance.


## Execute / Run test suite

```
$ cd seltest
$ go test -v
```

## TODOs

* Convert missing Selenium actions
* Env variable to overwrite the targetHost to use
* Env variable to set the webdriver to use
* Env variable to set a list of browsers to test
* helper functions (random password, email, names, credit cards....)
