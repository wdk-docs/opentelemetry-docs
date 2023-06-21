# 监控

Collector 为其监视提供了许多指标。下面列出了一些用于警报和监视的关键建议。

## 关键的监控

### 数据丢失

Use rate of `otelcol_processor_dropped_spans > 0` and
`otelcol_processor_dropped_metric_points > 0` to detect data loss, depending on
the requirements set up a minimal time window before alerting, avoiding
notifications for small losses that are not considered outages or within the
desired reliability level.

### CPU 资源不足

This depends on the CPU metrics available on the deployment, eg.:
`kube_pod_container_resource_limits_cpu_cores` for Kubernetes. Let's call it
`available_cores` below. The idea here is to have an upper bound of the number
of available cores, and the maximum expected ingestion rate considered safe,
let's call it `safe_rate`, per core. This should trigger increase of resources/
instances (or raise an alert as appropriate) whenever
`(actual_rate/available_cores) < safe_rate`.

The `safe_rate` depends on the specific configuration being used. // TODO:
Provide reference `safe_rate` for a few selected configurations.

## 二级监控

### 队列长度

Most exporters offer a
[queue/retry mechanism](../exporter/exporterhelper/README.md) that is
recommended as the retry mechanism for the Collector and as such should be used
in any production deployment.

The `otelcol_exporter_queue_capacity` indicates the capacity of the retry queue
(in batches). The `otelcol_exporter_queue_size` indicates the current size of
retry queue. So you can use these two metrics to check if the queue capacity is
enough for your workload.

The `otelcol_exporter_enqueue_failed_spans`,
`otelcol_exporter_enqueue_failed_metric_points` and
`otelcol_exporter_enqueue_failed_log_records` indicate the number of span/metric
points/log records failed to be added to the sending queue. This may be cause by
a queue full of unsettled elements, so you may need to decrease your sending
rate or horizontally scale collectors.

The queue/retry mechanism also supports logging for monitoring. Check the logs
for messages like `"Dropping data because sending_queue is full"`.

### 接收失败

Sustained rates of `otelcol_receiver_refused_spans` and
`otelcol_receiver_refused_metric_points` indicate too many errors returned to
clients. Depending on the deployment and the client’s resilience this may
indicate data loss at the clients.

Sustained rates of `otelcol_exporter_send_failed_spans` and
`otelcol_exporter_send_failed_metric_points` indicate that the Collector is not
able to export data as expected. It doesn't imply data loss per se since there
could be retries but a high rate of failures could indicate issues with the
network or backend receiving the data.

## 数据流

### 数据导入

The `otelcol_receiver_accepted_spans` and
`otelcol_receiver_accepted_metric_points` metrics provide information about the
data ingested by the Collector.

### 出口数据

The `otecol_exporter_sent_spans` and
`otelcol_exporter_sent_metric_points`metrics provide information about the data
exported by the Collector.
