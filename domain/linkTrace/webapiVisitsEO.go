package linkTrace

type WebapiVisitsEO struct {
	VisitsNode    string  // 访问节点
	SpeedMinMs    int64   // 最小访问速度
	SpeedMaxMs    int64   // 最大访问速度
	SpeedAvgMs    float64 // 平均访问速度
	Speed95LineMs int64   // 95线访问速度
	Speed99LineMs int64   // 99线访问速度
	ErrorCount    int     // 错误数量
	TotalCount    int     /// 总的数量
	QPS           float64 // 平均并发
}
