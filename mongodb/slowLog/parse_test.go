package slowLog

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseLine(t *testing.T){
	data:=[]byte(`297573ms 2019-04-04T10:55:50.236+0800 I WRITE    [conn280462] update trade.user_assets command: { q: { uid: 6, freeze_asset.usdt.m_zec_usdt_1554201224549521160: { $lt: 1E-22 } }, u: { $unset: { freeze_asset.usdt.m_zec_usdt_1554201224549521160: 1 } }, multi: false, upsert: false } planSummary: IXSCAN { uid: 1, asset.xmr: 1 } keysExamined:1 docsExamined:1 fromMultiPlanner:1 nMatched:1 nModified:1 writeConflicts:1219 numYields:1230 locks:{ Global: { acquireCount: { r: 1231, w: 1231 } }, Database: { acquireCount: { w: 1231 }, acquireWaitCount: { w: 2 }, timeAcquiringMicros: { w: 46578 } }, Collection: { acquireCount: { w: 1231 } } } 297573ms`)
	elem,err:=parseLine(data)
	require.NoError(t,err)
	t.Log(elem)
}