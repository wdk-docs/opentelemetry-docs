---
title: 资源
weight: 70
---

[resource][]表示作为资源属性产生遥测的实体。
例如，在Kubernetes上的容器中运行的产生遥测的进程有一个Pod名称、一个名称空间，可能还有一个部署名称。
所有这三个属性都可以包含在资源中。

在可观察性后端，您可以使用资源信息来更好地调查有趣的行为。
例如，如果您的跟踪或度量数据表明系统中的延迟，则可以将其缩小到特定的容器、Pod或Kubernetes部署。

下面你会发现一些关于如何使用Node.js SDK设置资源检测的介绍。

## 设置

按照[入门-Node.js](./getting-started/nodejs.md)中的说明，这样你就有了 `package.json`， `app.js`和`tracing.js`文件。

## 进程和环境资源检测

Node.js SDK开箱即用，检测[进程和进程运行时资源][]并从环境变量 `OTEL_RESOURCE_ATTRIBUTES` 中获取属性。
你可以通过在 `trace.js` 中打开诊断日志来验证它检测到什么:

```javascript
// For troubleshooting, set the log level to DiagLogLevel.DEBUG
diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.DEBUG);
```

运行应用程序，将一些值设置为`OTEL_RESOURCE_ATTRIBUTES`，例如我们设置`host.name`来标识[host][]:

```console
$ env OTEL_RESOURCE_ATTRIBUTES="host.name=localhost" \
  node --require ./tracing.js app.js
@opentelemetry/api: Registered a global for diag v1.2.0.
...
Listening for requests on http://localhost:8080
EnvDetector found resource. Resource { attributes: { 'host.name': 'localhost' } }
ProcessDetector found resource. Resource {
  attributes: {
    'process.pid': 12345,
    'process.executable.name': 'node',
    'process.command': '/app.js',
    'process.command_line': '/bin/node /app.js',
    'process.runtime.version': '16.17.0',
    'process.runtime.name': 'nodejs',
    'process.runtime.description': 'Node.js'
  }
}
...
```

## 使用环境变量添加资源

在上面的示例中，SDK检测到进程，并通过环境变量自动添加了 `host.name=localhost` 属性集。

您将在下面找到为您自动检测资源的说明。
但是，您可能会遇到这样的情况，即您需要的资源不存在检测器。
在这种情况下，您可以使用环境 `OTEL_RESOURCE_ATTRIBUTES` 来注入您需要的任何内容。
例如，下面的脚本添加[Service][]， [Host][]和[OS][]资源属性:

```console
$ env OTEL_RESOURCE_ATTRIBUTES="service.name=app.js,service.namespace=tutorial,service.version=1.0,service.instance.id=`uuidgen`,host.name=${HOSTNAME:},host.type=`uname -m`,os.name=`uname -s`,os.version=`uname -r`" \
  node --require ./tracing.js app.js
...
EnvDetector found resource. Resource {
  attributes: {
    'service.name': 'app.js',
    'service.namespace': 'tutorial',
    'service.version': '1.0',
    'service.instance.id': '46D99F44-27AB-4006-9F57-3B7C9032827B',
    'host.name': 'myhost',
    'host.type': 'arm64',
    'os.name': 'linux',
    'os.version': '6.0'
  }
}
...
```

## 在代码中添加资源


自定义资源也可以在代码中配置。
`NodeSDK`提供了一个配置选项，您可以在其中设置它们。
例如，你可以像下面这样更新`tracing.js`，设置`service.*`属性:

```javascript
...
const { Resource } = require('@opentelemetry/resources');
const { SemanticResourceAttributes } = require('@opentelemetry/semantic-conventions');
...
const sdk = new opentelemetry.NodeSDK({
  ...
  resource: new Resource({
    [ SemanticResourceAttributes.SERVICE_NAME ]: "yourServiceName",
    [ SemanticResourceAttributes.SERVICE_NAMESPACE ]: "yourNameSpace",
    [ SemanticResourceAttributes.SERVICE_VERSION ]: "1.0",
    [ SemanticResourceAttributes.SERVICE_INSTANCE_ID ]: "my-instance-id-1",
  })
  ...
});
...
```

**Note**: 如果通过环境变量和代码设置资源属性，则通过环境变量设置的值优先。

## 容器资源检测

使用相同的设置(`package.json`, `app.js` and `tracing.js`打开调试)和`Dockerfile`在同一目录下具有以下内容:

```Dockerfile
FROM node:latest
WORKDIR /usr/src/app
COPY package.json ./
RUN npm install
COPY . .
EXPOSE 8080
CMD [ "node", "--require", "./tracing.js", "app.js" ]
```

要确保您可以使用<kbd>Ctrl + C</kbd>(`SIGINT`)将以下内容添加到`app.js`的底部:

```javascript
process.on('SIGINT', function () {
  process.exit();
});
```

要获得自动检测到的容器id，安装以下附加依赖项:

```sh
npm install @opentelemetry/resource-detector-docker
```

接下来，像下面这样更新你的 `tracing.js`:

```javascript
const opentelemetry = require('@opentelemetry/sdk-node');
const {
  getNodeAutoInstrumentations,
} = require('@opentelemetry/auto-instrumentations-node');
const { diag, DiagConsoleLogger, DiagLogLevel } = require('@opentelemetry/api');
const {
  dockerCGroupV1Detector,
} = require('@opentelemetry/resource-detector-docker');

// For troubleshooting, set the log level to DiagLogLevel.DEBUG
diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.DEBUG);

const sdk = new opentelemetry.NodeSDK({
  traceExporter: new opentelemetry.tracing.ConsoleSpanExporter(),
  instrumentations: [getNodeAutoInstrumentations()],
  resourceDetectors: [dockerCGroupV1Detector],
});

sdk.start();
```

构建docker镜像:

```sh
docker build . -t nodejs-otel-getting-started
```

运行你的docker容器:

```sh
$ docker run --rm -p 8080:8080 nodejs-otel-getting-started
@opentelemetry/api: Registered a global for diag v1.2.0.
...
Listening for requests on http://localhost:8080
DockerCGroupV1Detector found resource. Resource {
  attributes: {
    'container.id': 'fffbeaf682f32ef86916f306ff9a7f88cc58048ab78f7de464da3c320ldb5c54'
  }
}
```

The detector has extracted the `container.id` for you. However you might
recognize that in this example, the process attributes and the attributes set
via an environment variable are missing! To resolve this, when you set the
`resourceDetectors` list you also need to specify the `envDetector` and
`processDetector` detectors:
检测器已经为您提取了`container.id`。
然而，您可能会发现，在本例中，流程属性和通过环境变量设置的属性缺失了!
为了解决这个问题，当你设置`resourceDetectors`列表时，你还需要指定`envDetector`和`processDetector`探测器:
```javascript
const opentelemetry = require('@opentelemetry/sdk-node');
const {
  getNodeAutoInstrumentations,
} = require('@opentelemetry/auto-instrumentations-node');
const { diag, DiagConsoleLogger, DiagLogLevel } = require('@opentelemetry/api');
const {
  dockerCGroupV1Detector,
} = require('@opentelemetry/resource-detector-docker');
const { envDetector, processDetector } = require('@opentelemetry/resources');

// For troubleshooting, set the log level to DiagLogLevel.DEBUG
diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.DEBUG);

const sdk = new opentelemetry.NodeSDK({
  traceExporter: new opentelemetry.tracing.ConsoleSpanExporter(),
  instrumentations: [getNodeAutoInstrumentations()],
  // Make sure to add all detectors you need here!
  resourceDetectors: [envDetector, processDetector, dockerCGroupV1Detector],
});

sdk.start();
```

重新构建镜像并再次运行容器:

```shell
docker run --rm -p 8080:8080 nodejs-otel-getting-started
@opentelemetry/api: Registered a global for diag v1.2.0.
...
Listening for requests on http://localhost:8080
EnvDetector found resource. Resource { attributes: {} }
ProcessDetector found resource. Resource {
  attributes: {
    'process.pid': 1,
    'process.executable.name': 'node',
    'process.command': '/usr/src/app/app.js',
    'process.command_line': '/usr/local/bin/node /usr/src/app/app.js',
    'process.runtime.version': '18.9.0',
    'process.runtime.name': 'nodejs',
    'process.runtime.description': 'Node.js'
  }
}
DockerCGroupV1Detector found resource. Resource {
  attributes: {
    'container.id': '654d0670317b9a2d3fc70cbe021c80ea15339c4711fb8e8b3aa674143148d84e'
  }
}
...
```

## 下一步

您可以在配置中添加更多的资源检测器，例如获取有关您的[Cloud]环境或[Deployment][]的详细信息。
你可以在[这里](https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/detectors/node)找到一份清单.

[resource]: ../../specs/otel/resource/sdk/
[getting started - node.js]: ../../instrumentation/js/getting-started/nodejs/
[process and process runtime resources]: ../../specs/otel/resource/semantic_conventions/process/
[host]: ../../specs/otel/resource/semantic_conventions/host/
[otlp exporter]: ../../instrumentation/js/exporters/#otlp-endpoint
[cloud]: ../../specs/otel/resource/semantic_conventions/cloud/
[deployment]: ../../specs/otel/resource/semantic_conventions/deployment_environment/
[service]: ../../specs/otel/resource/semantic_conventions/#service
[os]: ../../specs/otel/resource/semantic_conventions/os/
