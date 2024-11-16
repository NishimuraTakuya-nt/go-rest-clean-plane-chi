package datadog

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/iancoleman/strcase"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var spanNameCache sync.Map

// StartOperation は呼び出し元の情報から自動的にspan名を生成してトレースを開始します
func StartOperation(ctx context.Context, opts ...tracer.StartSpanOption) (context.Context, tracer.Span) {
	// 呼び出し元の情報を取得
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	// span名を生成（キャッシュからの取得または新規生成）
	spanName := getOrCreateSpanName(funcName)

	// 基本オプションの設定
	baseOpts := []tracer.StartSpanOption{
		tracer.SpanType("custom"),
		tracer.ResourceName(spanName),
	}

	// 追加のオプションを結合
	options := append(baseOpts, opts...)

	// spanの開始
	span, ctx := tracer.StartSpanFromContext(ctx, spanName, options...)

	return ctx, span
}

// キャッシュからspan名を取得または新規生成
func getOrCreateSpanName(funcName string) string {
	if cached, ok := spanNameCache.Load(funcName); ok {
		return cached.(string)
	}

	spanName := generateSpanName(funcName)
	spanNameCache.Store(funcName, spanName)
	return spanName
}

// span名の生成
func generateSpanName(funcName string) string {
	parts := strings.Split(funcName, ".")
	if len(parts) < 2 {
		return funcName
	}

	// メソッド名の取得
	methodName := parts[len(parts)-1]

	// 構造体名の取得（存在する場合）
	var structName string
	if len(parts) > 2 {
		// (*HealthcheckHandler) のような形式から "HealthcheckHandler" を抽出
		structPart := parts[len(parts)-2]
		structName = cleanStructName(structPart)
	}

	// レイヤー名の推測
	layer := inferLayer(structName)
	operation := strcase.ToSnake(methodName)

	return fmt.Sprintf("%s.%s.%s", layer, strcase.ToSnake(structName), operation)
}

func cleanStructName(structPart string) string {
	structPart = strings.TrimPrefix(structPart, "(")
	structPart = strings.TrimSuffix(structPart, ")")
	structPart = strings.TrimPrefix(structPart, "*")
	return structPart
}

// レイヤーの推測
func inferLayer(structName string) string {
	switch {
	case strings.HasSuffix(structName, "Handler"):
		return "handler"
	case strings.HasSuffix(structName, "UseCase"):
		return "usecase"
	case strings.HasSuffix(structName, "Service"):
		return "service"
	case strings.HasSuffix(structName, "Repository"):
		return "repository"
	case strings.HasSuffix(structName, "Client"):
		return "client"
	default:
		return "other"
	}
}
