local utils = import 'mixin-utils/utils.libsonnet';

(import 'dashboard-utils.libsonnet') {
  'alertmanager-resources.json':
    local filterNodeDiskByAlertmanager = |||
      ignoring(pod) group_right() (label_replace(count by(pod, instance, device) (container_fs_writes_bytes_total{%s,container="alertmanager",device!~".*sda.*"}), "device", "$1", "device", "/dev/(.*)") * 0)
    ||| % $.namespaceMatcher();
    ($.dashboard('Cortex / Alertmanager Resources') + { uid: '68b66aed90ccab448009089544a8d6c6' })
    .addClusterSelectorTemplates()
    .addRow(
      $.row('Gateway')
      .addPanel(
        $.containerCPUUsagePanel('CPU', 'cortex-gw'),
      )
      .addPanel(
        $.containerMemoryWorkingSetPanel('Memory (workingset)', 'cortex-gw'),
      )
      .addPanel(
        $.goHeapInUsePanel('Memory (go heap inuse)', 'cortex-gw'),
      )
    )
    .addRow(
      $.row('Alertmanager')
      .addPanel(
        $.containerCPUUsagePanel('CPU', 'alertmanager'),
      )
      .addPanel(
        $.containerMemoryWorkingSetPanel('Memory (workingset)', 'alertmanager'),
      )
      .addPanel(
        $.goHeapInUsePanel('Memory (go heap inuse)', 'alertmanager'),
      )
    )
    .addRow(
      $.row('Instance Mapper')
      .addPanel(
        $.containerCPUUsagePanel('CPU', 'alertmanager-im'),
      )
      .addPanel(
        $.containerMemoryWorkingSetPanel('Memory (workingset)', 'alertmanager-im'),
      )
      .addPanel(
        $.goHeapInUsePanel('Memory (go heap inuse)', 'alertmanager-im'),
      )
    )
    .addRow(
      $.row('Network')
      .addPanel(
        $.panel('Receive Bandwidth') +
        $.queryPanel('sum by(pod) (rate(container_network_receive_bytes_total{%s,pod=~"alertmanager.*"}[$__interval]))' % $.namespaceMatcher(), '{{pod}}') +
        $.stack +
        { yaxes: $.yaxes('Bps') },
      )
      .addPanel(
        $.panel('Transmit Bandwidth') +
        $.queryPanel('sum by(pod) (rate(container_network_transmit_bytes_total{%s,pod=~"alertmanager.*"}[$__interval]))' % $.namespaceMatcher(), '{{pod}}') +
        $.stack +
        { yaxes: $.yaxes('Bps') },
      )
    )
    .addRow(
      $.row('Disk')
      .addPanel(
        $.panel('Writes') +
        $.queryPanel('sum by(instance, device) (rate(node_disk_written_bytes_total[$__interval])) + %s' % filterNodeDiskByAlertmanager, '{{pod}} - {{device}}') +
        $.stack +
        { yaxes: $.yaxes('Bps') },
      )
      .addPanel(
        $.panel('Reads') +
        $.queryPanel('sum by(instance, device) (rate(node_disk_read_bytes_total[$__interval])) + %s' % filterNodeDiskByAlertmanager, '{{pod}} - {{device}}') +
        $.stack +
        { yaxes: $.yaxes('Bps') },
      )
    )
    .addRow(
      $.row('')
      .addPanel(
        $.panel('Disk Space Utilization') +
        $.queryPanel('max by(persistentvolumeclaim) (kubelet_volume_stats_used_bytes{%s} / kubelet_volume_stats_capacity_bytes{%s}) and count by(persistentvolumeclaim) (kube_persistentvolumeclaim_labels{%s,label_name="alertmanager"})' % [$.namespaceMatcher(), $.namespaceMatcher(), $.namespaceMatcher()], '{{persistentvolumeclaim}}') +
        { yaxes: $.yaxes('percentunit') },
      )
    ),
}
