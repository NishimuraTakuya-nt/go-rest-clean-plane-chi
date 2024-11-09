package telemetry

//func InitMetrics() error {
//	meter := global.Meter("your-service-metrics")
//
//	// リクエスト数のカウンター
//	requestCounter, err = meter.Int64Counter(
//		"http.server.request_count",
//		metric.WithDescription("Total number of HTTP requests"),
//	)
//	if err != nil {
//		return err
//	}
//
//	// レスポンス時間の測定
//	responseLatency, err = meter.Float64Histogram(
//		"http.server.duration",
//		metric.WithDescription("HTTP request duration"),
//	)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
