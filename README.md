![Testkube Logo](https://raw.githubusercontent.com/kubeshop/testkube/main/assets/testkube-color-gray.png)

# Welcome to Testkube SoapUI Executor

Testkube SoapUI Executor is a test executor module for [Testkube](https://testkube.io).

# Details

## Importing the SoapUI Executor into your Testkube cluster

In order to be able to run tests with this executor, it first needs to be imported as a Custom Resource into your Kubernetes cluster.
This can be achieved by cloning the repository and running:

```bash
$ kubectl apply -f executor.yaml
executor.executor.testkube.io/soapui-executor configured
```

## Running a SoapUI test

In order to run a SoapUI test using Testkube, it is necessary to create a Testkube Test.

### Using files as input

Testkube and the SoapUI executor accepts a project file as input.

```bash
$ kubectl testkube create test --file REST-Project-1-soapui-project.xml --type soapui/rest --name example-test

████████ ███████ ███████ ████████ ██   ██ ██    ██ ██████  ███████ 
   ██    ██      ██         ██    ██  ██  ██    ██ ██   ██ ██      
   ██    █████   ███████    ██    █████   ██    ██ ██████  █████   
   ██    ██           ██    ██    ██  ██  ██    ██ ██   ██ ██      
   ██    ███████ ███████    ██    ██   ██  ██████  ██████  ███████ 
                                           /tɛst kjub/ by Kubeshop


Test created  / example-test 🥇

```

### Using strings as input

```bash
$ cat REST-Project-2-soapui-project.xml | kubectl testkube create test --type soapui/rest --name example-test-string

████████ ███████ ███████ ████████ ██   ██ ██    ██ ██████  ███████ 
   ██    ██      ██         ██    ██  ██  ██    ██ ██   ██ ██      
   ██    █████   ███████    ██    █████   ██    ██ ██████  █████   
   ██    ██           ██    ██    ██  ██  ██    ██ ██   ██ ██      
   ██    ███████ ███████    ██    ██   ██  ██████  ██████  ███████ 
                                           /tɛst kjub/ by Kubeshop


Test created  / example-test-string 🥇

```

### Using params and args in your tests
TODO

To run the test, use:

```bash
$ kubectl testkube run test example-test

████████ ███████ ███████ ████████ ██   ██ ██    ██ ██████  ███████ 
   ██    ██      ██         ██    ██  ██  ██    ██ ██   ██ ██      
   ██    █████   ███████    ██    █████   ██    ██ ██████  █████   
   ██    ██           ██    ██    ██  ██  ██    ██ ██   ██ ██      
   ██    ███████ ███████    ██    ██   ██  ██████  ██████  ███████ 
                                           /tɛst kjub/ by Kubeshop


Type          : soapui/rest
Name          : example-test
Execution ID  : 624eedd443ed8485ae9289e2
Execution name: illegally-credible-mouse



Test execution started

Watch test execution until complete:
$ kubectl testkube watch execution 624eedd443ed8485ae9289e2


Use following command to get test execution details:
$ kubectl testkube get execution 624eedd443ed8485ae9289e2

```

## Plugins and extensions

In order to be able to add plugins and extensions the way [SoapUI docs](https://www.soapui.org/docs/test-automation/running-in-docker/) describe it, is currently not supported by Testkube.
In case you need this feature, please create an issue in the Testkube repository as described below, and we will proceed with the implementation.

## Reports

Getting complex reports on your run requires running tests with Docker volumes, which is not part of Testkube at the moment.
In case you need this feature, please create an issue in the Testkube repository as described below, and we will proceed with the implementation

# API

Testkube Executor SoapUI implements [testkube OpenAPI for executors](https://kubeshop.github.io/testkube/openapi/#operations-tag-executor) (look at executor tag)

# Issues and enchancements

Please follow the main [Testkube repository](https://github.com/kubeshop/testkube) for reporting any [issues](https://github.com/kubeshop/testkube/issues) or [discussions](https://github.com/kubeshop/testkube/discussions)
