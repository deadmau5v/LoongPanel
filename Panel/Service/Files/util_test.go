package Files

import (
	"io"
	"os"
	"testing"
)

type test []struct {
	name     string
	input    string
	expected interface{}
}

// region 测试工具函数

func testCheckFileName(t *testing.T) {
	cases := test{
		{"ValidName", "example.txt", true},
		{"InvalidNameWithSlash", "ex/ample.txt", false},
		{"InvalidNameWithBackslash", "ex\\ample.txt", false},
		{"InvalidNameWithColon", "ex:ample.txt", false},
		{"InvalidNameWithAsterisk", "ex*ample.txt", false},
		{"InvalidNameWithQuestionMark", "ex?ample.txt", false},
		{"InvalidNameWithQuote", "ex\"ample.txt", false},
		{"InvalidNameWithLessThan", "ex<ample.txt", false},
		{"InvalidNameWithGreaterThan", "ex>ample.txt", false},
		{"InvalidNameWithPipe", "ex|ample.txt", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := CheckFileName(c.input)
			if result != c.expected {
				t.Errorf("CheckFileName(%q) == %t, expected %t", c.input, result, c.expected)
			}
		})
	}
}

func testGetFilePath(t *testing.T) {
	cases := test{
		{"RootPath", "/example.txt", "/"},
		{"NestedPath", "/home/user/example.txt", "/home/user"},
		{"RelativePath", "user/example.txt", "user"},
		{"EmptyPath", "", "."},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := getFilePath(c.input)
			if result != c.expected {
				t.Errorf("getFilePath(%q) == %q, expected %q", c.input, result, c.expected)
			}
		})
	}
}

func testGetFileName(t *testing.T) {
	cases := test{
		{"RootPath", "/example.txt", "example.txt"},
		{"NestedPath", "/home/user/example.txt", "example.txt"},
		{"RelativePath", "user/example.txt", "example.txt"},
		{"EmptyPath", "", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := getFileName(c.input)
			if result != c.expected {
				t.Errorf("getFileName(%q) == %q, expected %q", c.input, result, c.expected)
			}
		})
	}
}

func testCopyFileConflictRename(t *testing.T) {
	cases := test{
		{"SimpleFile", "/home/user/example.txt", "/home/user/example_copy.txt"},
		{"NestedFile", "/home/user/docs/example.txt", "/home/user/docs/example_copy.txt"},
		{"FileWithMultipleDots", "/home/user/docs/example.test.txt", "/home/user/docs/example.test_copy.txt"},
		{"FileWithoutExtension", "/home/user/docs/example", "/home/user/docs/example_copy"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := copyFileConflictRename(c.input)
			if result != c.expected {
				t.Errorf("copyFileConflictRename(%q) == %q, expected %q", c.input, result, c.expected)
			}
		})
	}
}

func testCheckFilePath(t *testing.T) {
	cases := test{
		{"ValidPath", "/home/user/example.txt", true},
		{"PathWithBackslash", "C:\\user\\example.txt", false},
		{"PathWithForwardSlash", "/home/user/example.txt", true},
		{"PathWithColon", "/home/user:example.txt", false},
		{"PathWithAsterisk", "/home/user/example*.txt", false},
		{"PathWithQuestionMark", "/home/user/example?.txt", false},
		{"PathWithDoubleQuote", "/home/user/example\".txt", false},
		{"PathWithLessThan", "/home/user/example<.txt", false},
		{"PathWithGreaterThan", "/home/user/example>.txt", false},
		{"PathWithPipe", "/home/user/example|.txt", false},
		{"EmptyPath", "", false},
		{"WhitespacePath", "   ", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := checkFilePath(c.input)
			if result != c.expected {
				t.Errorf("checkFilePath(%q) == %v, expected %v", c.input, result, c.expected)
			}
		})
	}
}

func testIsDir(t *testing.T) {
	cases := test{
		{"ExistingDir", "testdir", true},
		{"NonExistingDir", "nonexistingdir", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.name == "ExistingDir" {
				err := os.Mkdir(c.input, 0755)
				if err != nil {
					t.Fatalf("无法创建目录: %v", err)
				}
				defer os.Remove(c.input) // 测试完成后删除
			}

			result := isDir(c.input)
			if result != c.expected {
				t.Errorf("isDir(%q) == %v, expected %v", c.input, result, c.expected)
			}
		})
	}
}

func testCopyFile(t *testing.T) {
	srcPath := "temp.txt"
	targetPath := "temp_copy.txt"

	// 创建临时文件
	srcFile, err := os.Create(srcPath)
	if err != nil {
		t.Fatalf("无法创建源文件: %v", err)
	}
	defer os.Remove(srcPath) // 测试完成后删除

	// 写入一些内容到临时文件
	_, err = srcFile.WriteString("这是一个测试文件")
	if err != nil {
		t.Fatalf("无法写入源文件: %v", err)
	}
	srcFile.Close()

	// 测试 copyFile 函数
	err = copyFile(srcPath, targetPath)
	if err != nil {
		t.Errorf("copyFile 失败: %v", err)
	}
	defer os.Remove(targetPath) // 测试完成后删除

	// 检查目标文件是否存在
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Errorf("目标文件未创建: %v", targetPath)
	}

	// 检查目标文件内容是否正确
	targetFile, err := os.Open(targetPath)
	if err != nil {
		t.Fatalf("无法打开目标文件: %v", err)
	}
	defer targetFile.Close()

	content, err := io.ReadAll(targetFile)
	if err != nil {
		t.Fatalf("无法读取目标文件: %v", err)
	}

	expectedContent := "这是一个测试文件"
	if string(content) != expectedContent {
		t.Errorf("目标文件内容不正确，期望: %q，实际: %q", expectedContent, string(content))
	}
}

func testIsTar(t *testing.T) {
	cases := test{
		{"TarFile", "example.tar", true},
		{"NonTarFile", "example.txt", false},
		{"EmptyPath", "", false},
		{"TarGzFile", "example.tar.gz", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := isTar(c.input)
			if result != c.expected {
				t.Errorf("isTar(%q) == %t, expected %t", c.input, result, c.expected)
			}
		})
	}
}

// endregion

func TestUtil(t *testing.T) {
	testCheckFileName(t)
	testGetFilePath(t)
	testGetFileName(t)
	testCopyFileConflictRename(t)
	testCheckFilePath(t)
	testIsDir(t)
	testCopyFile(t)
	testIsTar(t)
}
