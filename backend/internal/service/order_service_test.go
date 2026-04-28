package service

import (
	"math"
	"testing"
)

// TestRound2 验证金额二位小数的四舍五入
func TestRound2(t *testing.T) {
	cases := []struct {
		in   float64
		want float64
	}{
		{0, 0},
		{1.234, 1.23},
		{1.235001, 1.24},
		{99.994, 99.99},
		{99.996, 100.00},
	}
	for _, c := range cases {
		got := round2(c.in)
		if math.Abs(got-c.want) > 1e-9 {
			t.Errorf("round2(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

// TestAmountFenConversion 验证元->分的转换策略（避免 0.1+0.2 != 0.3 浮点误差）
// PayCallback 使用 math.Round(order.PayAmount * 100) 比对回调 payer_total（单位分）
func TestAmountFenConversion(t *testing.T) {
	cases := []struct {
		pay  float64
		want int64
	}{
		{0.01, 1},
		{0.1, 10},
		{1.00, 100},
		{99.99, 9999},
		{0.3, 30}, // 0.1 + 0.2 浮点陷阱
		{123.45, 12345},
	}
	for _, c := range cases {
		got := int64(math.Round(c.pay * 100))
		if got != c.want {
			t.Errorf("fen(%v) = %d, want %d", c.pay, got, c.want)
		}
	}
}
