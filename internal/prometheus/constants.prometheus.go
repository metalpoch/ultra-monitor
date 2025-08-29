package prometheus

const QUERY_STATS_IN string = `sum(rate(hwGponOltEthernetStatisticReceivedBytes_count{instance="%s"}[15m]) * 8) by (ponPortIndex)
  * on(ponPortIndex) group_left(ifName)
  label_replace(ifName{instance="%s"}, "ponPortIndex", "$1", "ifIndex", "(.*)")`

const QUERY_STATS_OUT string = `sum(rate(hwGponOltEthernetStatisticSendBytes_count{instance="%s"}[15m]) * 8) by (ponPortIndex)
  * on(ponPortIndex) group_left(ifName)
  label_replace(ifName{instance="%s"}, "ponPortIndex", "$1", "ifIndex", "(.*)")`

const QUERY_STATS_IFSPEED string = `sum(ifSpeed{instance="%s"}) by (ifIndex)`
