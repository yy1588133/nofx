package llm

import (
	"os"
	"testing"
)

// TestMain skips all llm live tests when Qwen credentials are not provided.
//
// 说明：
// - 该仓库的 Qwen 测试属于“真实网络调用 + 需要密钥”的集成测试；
// - 在 CI / 本地未配置密钥的情况下应跳过，避免因缺少环境变量导致失败。
func TestMain(m *testing.M) {
	if os.Getenv("QWEN_APP_ID") == "" || os.Getenv("QWEN_API_KEY") == "" {
		os.Exit(0)
	}
	os.Exit(m.Run())
}

