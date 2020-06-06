package main

//    1 块钱可以买 3 个桃子吃，吃完后 3 个桃核可以换 1 个桃子，请问 135142857 元可以最多吃到多少个桃子。
func main() {
	println(count(135142857))
	println(count(27))

}

func count(money int64) int64 {
	peaches := money * 3 // 桃
	count := peaches     // 总数量
	core := count        // 桃核

	for core >= 3 { // 3个桃核换一个桃
		peaches = core / 3 // 换桃
		core = core % 3    // 剩余的核
		core += peaches    // 每个桃对应1个桃核
		count += peaches   // 桃数量
	}

	return count
}
