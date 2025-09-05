package prometheus

const QUERY_STATES_BY_REGION_IN string = `sum(rate(hwGponOltEthernetStatisticReceivedBytes_count{region="%s"}[15m]) * 8) by (state)`
const QUERY_STATES_BY_REGION_OUT string = `sum(rate(hwGponOltEthernetStatisticSendBytes_count{region="%s"}[15m]) * 8) by (state)`
const QUERY_STATES_BY_REGION_IFSPEED string = `sum(ifSpeed{region="%s"}) by (state)`

const QUERY_OLT_BY_STATES_IN string = `sum(rate(hwGponOltEthernetStatisticReceivedBytes_count{state="%s"}[15m]) * 8) by (instance,job,state,region)
	* on (instance,job,state,region) group_left(sysName) sysName`
const QUERY_OLT_BY_STATES_OUT string = `sum(rate(hwGponOltEthernetStatisticSendBytes_count{state="%s"}[15m]) * 8) by (instance,job,state,region)
	* on (instance,job,state,region) group_left(sysName) sysName`
const QUERY_OLT_BY_STATES_IFSPEED string = `sum(ifSpeed{state="%s"}) BY (instance,job,state,region) * on (instance,job,state,region) group_left(sysName) sysName`

const QUERY_GPON_STATS_IN string = `sum(rate(hwGponOltEthernetStatisticReceivedBytes_count{instance="%s"}[15m]) * 8) by (ponPortIndex)
  * on(ponPortIndex) group_left(ifName)
  label_replace(ifName{instance="%s"}, "ponPortIndex", "$1", "ifIndex", "(.*)")`

const QUERY_GPON_STATS_OUT string = `sum(rate(hwGponOltEthernetStatisticSendBytes_count{instance="%s"}[15m]) * 8) by (ponPortIndex)
  * on(ponPortIndex) group_left(ifName)
  label_replace(ifName{instance="%s"}, "ponPortIndex", "$1", "ifIndex", "(.*)")`

const QUERY_GPON_STATS_IFSPEED string = `sum(ifSpeed{instance="%s"}) by (ifIndex)`
