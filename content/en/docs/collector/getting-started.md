---
title: 入门
spelling:
  - cSpell:ignore darwin dpkg GOARCH journalctl kubectl linux otelcorecol pprof
  - cSpell:ignore zpages
weight: 1
---

如果您不熟悉适用于收集器的部署模型、组件和存储库，请首先查看[数据收集]和[开发方法]页面。

## 演示

部署负载生成器、代理和网关以及 Jaeger、Zipkin 和 Prometheus 后端。更多信息可以在
演示[README.md][]中找到。

```sh
git clone git@github.com:open-telemetry/opentelemetry-collector-contrib.git --depth 1; \
  cd opentelemetry-collector-contrib/examples/demo; \
  docker compose up -d
```

## Docker

提取 docker 映像并在容器中运行收集器。将`{{% param collectorVersion %}}`替换为要
运行的收集器的版本。

=== "DockerHub"

    ```sh
    docker pull otel/opentelemetry-collector-contrib:{{% param collectorVersion %}}
    docker run otel/opentelemetry-collector-contrib:{{% param collectorVersion %}}
    ```

=== "ghcr.io"

    ```sh
    docker pull ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib:{{% param collectorVersion %}}
    docker run ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib:{{% param collectorVersion %}}
    ```

要从当前工作目录加载自定义配置 `config.yaml` ，将该文件挂载为卷:

=== "DockerHub"

    ```sh
    docker run -v $(pwd)/config.yaml:/etc/otelcol-contrib/config.yaml otel/opentelemetry-collector-contrib:{{% param collectorVersion %}}
    ```

=== "ghcr.io"

    ```sh
    docker run -v $(pwd)/config.yaml:/etc/otelcol-contrib/config.yaml ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib:{{% param collectorVersion %}}
    ```

## Docker Compose

你可以像下面这样添加 OpenTelemetry 收集器到你现有的`docker-compose.yaml`中:

```yaml
otel-collector:
  image: otel/opentelemetry-collector-contrib
  command: [--config=/etc/otel-collector-config.yaml]
  volumes:
    - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
  ports:
    - 1888:1888 # pprof extension
    - 8888:8888 # Prometheus metrics exposed by the collector
    - 8889:8889 # Prometheus exporter metrics
    - 13133:13133 # health_check extension
    - 4317:4317 # OTLP gRPC receiver
    - 4318:4318 # OTLP http receiver
    - 55679:55679 # zpages extension
```

## Kubernetes

将代理部署为守护进程和单个网关实例。

```sh
kubectl apply -f https://raw.githubusercontent.com/open-telemetry/opentelemetry-collector/main/examples/k8s/otel-config.yaml
```

上面的示例旨在作为一个起点，在实际生产使用之前进行扩展和定制。对于生产就绪的定制
和安装，请参见[OpenTelemetry Helm Charts][]。

[OpenTelemetry Operator][]也可用于提供和维护 OpenTelemetry Collector 实例，具有
自动升级处理、基于 OpenTelemetry 配置的“服务”配置、自动侧车注入部署等功能。

## Nomad

将 Collector 部署为代理、网关和完整演示的参考作业文件可以在[Getting Started with
OpenTelemetry on HashiCorp Nomad][]中找到。

## Linux 包

每个 Collector 版本都包含 Linux amd64/arm64/i386 系统的 APK、DEB 和 RPM 包。包中
包含一个默认配置，可以在“/etc/otelcol/config”中找到。yaml 的安装后。

> Note: 自动服务配置需要 `systemd`。

### APK 安装

要在 alpine 系统上开始运行，请运行以下命令，将 `v{{% param collectorVersion %}}`
替换为您希望运行的 Collector 版本。

=== "AMD64"

    ```sh
    apk update
    apk add wget shadow
    wget https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_amd64.apk
    apk add --allow-untrusted otelcol_{{% param collectorVersion %}}_linux_amd64.apk
    ```

=== "ARM64"

    ```sh
    apk update
    apk add wget shadow
    wget https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_arm64.apk
    apk add --allow-untrusted otelcol_{{% param collectorVersion %}}_linux_arm64.apk
    ```

=== "i386"

    ```sh
    apk update
    apk add wget shadow
    wget https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_386.apk
    apk add --allow-untrusted otelcol_{{% param collectorVersion %}}_linux_386.apk
    ```

### DEB 安装

To get started on Debian systems run the following replacing
`v{{% param collectorVersion %}}` with the version of the Collector you wish to
run and `amd64` with the appropriate architecture.

<!-- prettier-ignore-start -->

=== "AMD64"

    ```sh
    sudo apt-get update
    sudo apt-get -y install wget systemctl
    wget https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_amd64.deb
    sudo dpkg -i otelcol_{{% param collectorVersion %}}_linux_amd64.deb
    ```

=== "ARM64"

    ```sh
    sudo apt-get update
    sudo apt-get -y install wget systemctl
    wget https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_arm64.deb
    sudo dpkg -i otelcol_{{% param collectorVersion %}}_linux_arm64.deb
    ```

=== "i386"

    ```sh
    sudo apt-get update
    sudo apt-get -y install wget systemctl
    wget https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_386.deb
    sudo dpkg -i otelcol_{{% param collectorVersion %}}_linux_386.deb
    ```

### RPM 安装

To get started on Red Hat systems run the following replacing
`v{{% param collectorVersion %}}` with the version of the Collector you wish to
run and `x86_64` with the appropriate architecture.

=== "AMD64"

    ```sh
    sudo yum update
    sudo yum -y install wget systemctl
    wget https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_amd64.rpm
    sudo rpm -ivh otelcol_{{% param collectorVersion %}}_linux_amd64.rpm
    ```

=== "ARM64"

    ```sh
    sudo yum update
    sudo yum -y install wget systemctl
    wget https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_arm64.rpm
    sudo rpm -ivh otelcol_{{% param collectorVersion %}}_linux_arm64.rpm
    ```

=== "i386"

    ```sh
    sudo yum update
    sudo yum -y install wget systemctl
    wget https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_386.rpm
    sudo rpm -ivh otelcol_{{% param collectorVersion %}}_linux_386.rpm
    ```

### 手动安装

Linux [releases][] are available for various architectures. It's possible to
download the archive containing the binary and install it on your machine
manually:

<!-- prettier-ignore-start -->

{{< tabpane lang=shell persistLang=false >}} {{< tab AMD64 >}}

=== "AMD64"

    ```sh
    curl --proto '=https' --tlsv1.2 -fOL https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_amd64.tar.gz
    tar -xvf otelcol_{{% param collectorVersion %}}_linux_amd64.tar.gz
    ```

=== "ARM64"

    ```sh
    curl --proto '=https' --tlsv1.2 -fOL https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_arm64.tar.gz
    tar -xvf otelcol_{{% param collectorVersion %}}_linux_arm64.tar.gz
    ```

=== "i386"

    ```sh
    curl --proto '=https' --tlsv1.2 -fOL https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_386.tar.gz
    tar -xvf otelcol_{{% param collectorVersion %}}_linux_386.tar.gz
    ```

=== "ppc64le"

    ```sh
    curl --proto '=https' --tlsv1.2 -fOL https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_linux_ppc64le.tar.gz
    tar -xvf otelcol_{{% param collectorVersion %}}_linux_ppc64le.tar.gz
    ```

### 业务自动配置

By default, the `otelcol` systemd service will be started with the
`--config=/etc/otelcol/config.yaml` option after installation. To customize
these options, modify the `OTELCOL_OPTIONS` variable in the
`/etc/otelcol/otelcol.conf` systemd environment file with the appropriate
command-line options (run `/usr/bin/otelcol --help` to see all available
options). Additional environment variables can also be passed to the `otelcol`
service by adding them to this file.

If either the Collector configuration file or `/etc/otelcol/otelcol.conf` are
modified, restart the `otelcol` service to apply the changes by running:

```sh
sudo systemctl restart otelcol
```

To check the output from the `otelcol` service, run:

```sh
sudo journalctl -u otelcol
```

## MacOS 包

MacOS [releases][] are available for Intel- & ARM-based systems. They are
packaged as gzipped tarballs (`.tar.gz`) and will need to be unpacked with a
tool that supports this compression format:

=== "Intel"

    ```sh
    curl --proto '=https' --tlsv1.2 -fOL https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_darwin_amd64.tar.gz
    tar -xvf otelcol_{{% param collectorVersion %}}_darwin_amd64.tar.gz
    ```

=== "ARM"

    ```sh
    curl --proto '=https' --tlsv1.2 -fOL https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/v{{% param collectorVersion %}}/otelcol_{{% param collectorVersion %}}_darwin_arm64.tar.gz
    tar -xvf otelcol_{{% param collectorVersion %}}_darwin_arm64.tar.gz
    ```

Every Collector release includes an `otelcol` executable that you can run after
unpacking.

## Windows 包

Windows [releases][] are packaged as gzipped tarballs (`.tar.gz`) and will need
to be unpacked with a tool that supports this compression format.

Every Collector release includes an `otelcol.exe` executable that you can run
after unpacking.

## 本地

Builds the latest version of the collector based on the local operating system,
runs the binary with all receivers enabled and exports all the data it receives
locally to a file. Data is sent to the container and the container scrapes its
own Prometheus metrics. The following example uses two terminal windows to
better illustrate the collector. In the first terminal window run the following:

```sh
git clone https://github.com/open-telemetry/opentelemetry-collector.git
cd opentelemetry-collector
make install-tools
make otelcorecol
./bin/otelcorecol_* --config ./examples/local/otel-config.yaml
```

In a second terminal window, you can test the newly built collector by doing the
following:

```sh
git clone https://github.com/open-telemetry/opentelemetry-collector-contrib.git
cd opentelemetry-collector-contrib/examples/demo/server
go build -o main main.go; ./main & pid1="$!"
cd ../client
go build -o main main.go; ./main
```

To stop the client, use type <kbd>Ctrl-C</kbd>. To stop the server, use the
`kill $pid1` command. To stop the collector, type <kbd>Ctrl-C</kbd> in its
terminal window as well.

!!! Note

    上面显示的命令在bash shell中演示了这个过程。
    对于其他shell，这些命令可能略有不同。

[数据收集]: ../concepts/data-collection.md
[开发方法]: ./deployment/index.md
[readme.md]:
  https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/examples/demo
[opentelemetry helm charts]:
  https://github.com/open-telemetry/opentelemetry-helm-charts
[opentelemetry operator]:
  https://github.com/open-telemetry/opentelemetry-operator
[getting started with opentelemetry on hashicorp nomad]:
  https://github.com/hashicorp/nomad-open-telemetry-getting-started
[releases]:
  https://github.com/open-telemetry/opentelemetry-collector-releases/releases
