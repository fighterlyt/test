package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestNewData(t *testing.T) {
	sectionCount := 10
	d := NewData(sectionCount)
	e := NewElement("test")
	s := NewSection(10)
	s.AddElement(e)
	for i := 0; i < sectionCount*1000; i++ {
		d.Add(s)
		if i < sectionCount {
			require.Equal(t, 0, d.lower)
			require.Equal(t, (i+1)%sectionCount, d.upper)
			require.Equal(t, (i+1)%sectionCount, d.index)
			require.Equal(t, 0, d.sectionIndexLower(), strconv.Itoa(i))
			require.Equal(t, i+1, d.sectionIndexUpper())
		} else {
			require.Equal(t, (i+1)%sectionCount, d.lower)
			require.Equal(t, (i+1-sectionCount)%sectionCount, d.upper)
			require.Equal(t, (i+1-sectionCount)%sectionCount, d.index)
			require.Equal(t, d.sectionIndex-sectionCount, d.sectionIndexLower(), strconv.Itoa(i))
			require.Equal(t, d.sectionIndex, d.sectionIndexUpper())
		}
		require.Equal(t, i+1, d.sectionIndex)
	}

}

func TestDataGet(t *testing.T) {
	sectionCount := 10
	d := NewData(sectionCount)

	for i := 0; i < sectionCount*1000; i++ {
		e := NewElement(strconv.Itoa(i))
		s := NewSection(10)
		s.AddElement(e)
		d.Add(s)
	}
	for i:=sectionCount*999;i<sectionCount*1000;i++{
		s:=d.Get(i)
		if !assert.NotNil(t,s,strconv.Itoa(i)){
			t.Log(d.sectionIndexLower(),d.sectionIndexUpper())
			t.FailNow()
		}
		require.Equal(t,strconv.Itoa(i),s.tops[0].data)
	}
	for i:=0;i<sectionCount*999;i++{
		require.Nil(t,d.Get(i))
	}
	for i:=sectionCount*1000;i<sectionCount*1001;i++{
		require.Nil(t,d.Get(i),strconv.Itoa(i))
	}
}
