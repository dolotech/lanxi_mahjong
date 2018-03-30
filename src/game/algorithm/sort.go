package algorithm


// 对牌值从小到大排序，采用快速排序算法
func Sort(arr []byte, start, end int) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for arr[i] < key {
				i++
			}
			for arr[j] > key {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}
		if start < j {
			Sort(arr, start, j)
		}
		if end > i {
			Sort(arr, i, end)
		}
	}
}
