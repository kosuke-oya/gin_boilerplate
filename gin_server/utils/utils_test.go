package utils

import "testing"

func TestUniqueID(t *testing.T) {
	// テストケース
	tests := []struct {
		name  string
		digit int
	}{
		{
			name:  "10桁のIDを生成する",
			digit: 10,
		},
		{
			name:  "20桁のIDを生成する",
			digit: 20,
		},
		{
			name:  "30桁のIDを生成する",
			digit: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト対象関数の実行
			got := UniqueID(tt.digit)
			// 生成されたIDの桁数が正しいか検証
			if len(got) != tt.digit {
				t.Errorf("UniqueID() = %v, want %v", len(got), tt.digit)
			}
		})
	}
}
