# gauge-go
[![Gauge Badge](https://cdn.rawgit.com/renjithgr/gauge-js/72f332d11f54e16b74aedb875f702643708156f7/Gauge_Badge_1.svg)](http://getgauge.io)

Go language plugin for ThoughtWorks Gauge.

[![Build Status](https://travis-ci.org/manuviswam/gauge-go.svg?branch=master)](https://travis-ci.org/manuviswam/gauge-go)

## Build from Source
* Checkout the repository
* Run commands:
```sh
go run build/make.go
go run build/make.go --install
```
This installs the gauge-go plugin to its default location based on OS.

### Run Unit tests

`go test ./...`

### Generate gauge-go distributables

`go run build/make.go --distro`

This generates the gauge-go plugin as zip distributable. In order to generate it for all platforms, add `--all-platforms` flag in above command.

## Usage

If you are new to Gauge, please consult the [Gauge documentation](http://getgauge.io/documentation/user/current/) to know about how Gauge works.

**Initialize:** To initialize a project with gauge-go, in an empty directory run:

```sh
$ gauge --init go
```
Note: Create your project in `GOPATH`.

**Run specs:**

```sh
$ gauge specs/
```

## Methods

### Step implementation

**`var _ = gauge.Step(<step-text>, fn)`**

Use the `gauge.Step()` method to implement your steps. For example:

```go
import "github.com/manuviswam/gauge-go/gauge"

var _ = gauge.Step("Vowels in English language are <vowels>.", func (vowelsGiven string) {
  assert.Equal(vowelsGiven, "aeiou")
})
```

### Execution Hooks

gauge-go supports tagged [execution hooks](http://getgauge.io/documentation/user/current/execution/execution_hooks.html). These methods are available for each type of hook:

"Before" hooks:

- **`BeforeSuite(fn, tags, operator)`** - Executed before the test suite begins
- **`BeforeSpec(fn, tags, operator)`** - Executed before each specification
- **`BeforeScenario(fn, tags, operator)`** - Executed before each scenario
- **`BeforeStep(fn, tags, operator)`**- Execute before each step

"After" hooks:

- **`AfterSuite(fn, tags, operator)`** - Executed after the test suite begins
- **`AfterSpec(fn, tags, operator)`** - Executed after each specification
- **`AfterScenario(fn, tags, operator)`** - Executed after each scenario
- **`AfterStep(fn, tags, operator)`**- Execute after each step

Here's an example of a hook that is executed before each scenario:

```go
var _ = gauge.BeforeScenario(function () {
  assert.equal(vowels.join(""), "aeiou")
}, []string{"foo", "bar"}, testsuit.AND)
```

- *`tags`*
An array of strings for the tags for which to execute the current hook. These are only useful at specification or scenario level. If tags is empty, the provided hook is executed on each occurrence of the hook.
- *`operator`*: Valid values: `"AND"`, `"OR"`.
This controls whether the current hook is executed when all of the tags match (in case of `"AND"`), or if any of the tags match (in case of `OR`).

### Data Stores

Step implementations can share custom data across scenarios, specifications and suites using data stores.

There are 3 different types of data stores based on the lifecycle of when it gets cleared. These are present in the `github.com/manuviswam/gauge-go/gauge` package of gauge-go plugin.

#### Scenario store

This data store keeps values added to it in the lifecycle of the scenario execution. Values are cleared after every scenario executes.

**Store a value:**

```go
GetScenarioStore()[key] = value
```

**Retrieve a value:**

```go
value := GetScenarioStore()[key]
```

#### Specification store

This data store keeps values added to it in the lifecycle of the specification execution. Values are cleared after every specification executes.

**Store a value:**

```go
GetSpecStore()[key] = value
```

**Retrieve a value:**

```go
value := GetSpecStore()[key]
```

#### Suite store

This data store keeps values added to it in the lifecycle of the entire suite's execution. Values are cleared after entire suite executes.

**Store a value:**

```go
GetSuiteStore()[key] = value
```

**Retrieve a value:**

```go
value := GetSuiteStore()[key]
```

**Note:** Suite Store is not advised to be used when executing specs in parallel. The values are not retained between parallel streams of execution.


## License

![GNU Public License version 3.0](http://www.gnu.org/graphics/gplv3-127x51.png)
Gauge-GO is released under [GNU Public License version 3.0](http://www.gnu.org/licenses/gpl-3.0.txt)
