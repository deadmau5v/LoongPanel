package Files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testDir(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	os.Mkdir(filepath.Join(tempDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(tempDir, "subdir", "file2.txt"), []byte("content2"), 0644)

	tests := []struct {
		name     string
		path     string
		expected int
	}{
		{"RootDir", tempDir, 2},
		{"SubDir", filepath.Join(tempDir, "subdir"), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := Dir(tt.path)
			if err != nil {
				t.Fatalf("Dir() error = %v", err)
			}
			if len(files) != tt.expected {
				t.Errorf("Dir() = %v, expected %v", len(files), tt.expected)
			}
		})
	}
}

func testContent(t *testing.T) {
	// 创建临时文件
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file.txt")
	content := "Hello, World!"
	os.WriteFile(filePath, []byte(content), 0644)

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"ValidFile", filePath, content},
		{"InvalidFile", filepath.Join(tempDir, "nonexistent.txt"), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Content(tt.path)
			if tt.name == "InvalidFile" {
				if err == nil {
					t.Errorf("Content() error = %v, expected an error", err)
				}
			} else {
				if err != nil {
					t.Fatalf("Content() error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("Content() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}

func testDelete(t *testing.T) {
	// 创建临时文件
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "file.txt")
	os.WriteFile(filePath, []byte("Hello, World!"), 0644)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"ValidFile", filePath, false},
		{"InvalidFile", filepath.Join(tempDir, "nonexistent.txt"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Delete(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if _, err := os.Stat(tt.path); !os.IsNotExist(err) {
					t.Errorf("Delete() = file still exists, expected it to be deleted")
				}
			}
		})
	}
}

func testCopy(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	srcFilePath := filepath.Join(tempDir, "file.txt")
	os.WriteFile(srcFilePath, []byte("Hello, World!"), 0644)
	srcDirPath := filepath.Join(tempDir, "srcDir")
	os.Mkdir(srcDirPath, 0755)
	destFilePath := filepath.Join(tempDir, "destFile.txt")
	destDirPath := filepath.Join(tempDir, "destDir")

	tests := []struct {
		name    string
		src     string
		dest    string
		wantErr bool
	}{
		{"CopyFileToFile", srcFilePath, destFilePath, false},
		{"CopyFileToDir", srcFilePath, destDirPath, false},
		{"CopyDirToDir", srcDirPath, destDirPath, false},
		{"InvalidSrcFile", filepath.Join(tempDir, "nonexistent.txt"), destFilePath, true},
		{"InvalidDestDir", srcFilePath, filepath.Join(tempDir, "nonexistentDir", "file.txt"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Copy(tt.src, tt.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if isDir(tt.src) {
					if !isDir(tt.dest) {
						t.Errorf("Copy() = destination is not a directory, expected it to be a directory")
					}
				} else {
					if _, err := os.Stat(tt.dest); os.IsNotExist(err) {
						t.Errorf("Copy() = destination file does not exist, expected it to be copied")
					}
				}
			}
		})
	}
}

func testMove(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	srcFilePath := filepath.Join(tempDir, "file.txt")
	os.WriteFile(srcFilePath, []byte("Hello, World!"), 0644)
	srcDirPath := filepath.Join(tempDir, "srcDir")
	os.Mkdir(srcDirPath, 0755)
	destFilePath := filepath.Join(tempDir, "destFile.txt")
	destDirPath := filepath.Join(tempDir, "destDir")

	tests := []struct {
		name    string
		src     string
		dest    string
		wantErr bool
	}{
		{"MoveFileToFile", srcFilePath, destFilePath, false},
		{"MoveFileToDir", srcFilePath, destDirPath, false},
		{"MoveDirToDir", srcDirPath, destDirPath, false},
		{"InvalidSrcFile", filepath.Join(tempDir, "nonexistent.txt"), destFilePath, true},
		{"InvalidDestDir", srcFilePath, filepath.Join(tempDir, "nonexistentDir", "file.txt"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Move(tt.src, tt.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if isDir(tt.src) {
					if !isDir(tt.dest) {
						t.Errorf("Move() = destination is not a directory, expected it to be a directory")
					}
				} else {
					if _, err := os.Stat(tt.dest); os.IsNotExist(err) {
						t.Errorf("Move() = destination file does not exist, expected it to be moved")
					}
					if _, err := os.Stat(tt.src); !os.IsNotExist(err) {
						t.Errorf("Move() = source file still exists, expected it to be moved")
					}
				}
			}
		})
	}
}

func testRename(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	srcFilePath := filepath.Join(tempDir, "file.txt")
	os.WriteFile(srcFilePath, []byte("Hello, World!"), 0644)
	newFileName := "renamed_file.txt"
	newFilePath := filepath.Join(tempDir, newFileName)

	tests := []struct {
		name    string
		src     string
		newName string
		wantErr bool
	}{
		{"RenameFile", srcFilePath, newFileName, false},
		{"InvalidSrcFile", filepath.Join(tempDir, "nonexistent.txt"), newFileName, true},
		{"InvalidNewName", srcFilePath, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Rename(tt.src, tt.newName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rename() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
					t.Errorf("Rename() = new file does not exist, expected it to be renamed")
				}
				if _, err := os.Stat(tt.src); !os.IsNotExist(err) {
					t.Errorf("Rename() = source file still exists, expected it to be renamed")
				}
			}
		})
	}
}

func testCompress(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	srcFilePath := filepath.Join(tempDir, "file.txt")
	os.WriteFile(srcFilePath, []byte("Hello, World!"), 0644)
	srcDirPath := filepath.Join(tempDir, "testDir")
	os.Mkdir(srcDirPath, 0755)
	os.WriteFile(filepath.Join(srcDirPath, "file1.txt"), []byte("File 1"), 0644)
	os.WriteFile(filepath.Join(srcDirPath, "file2.txt"), []byte("File 2"), 0644)

	tests := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{"CompressFile", srcFilePath, false},
		{"CompressDirectory", srcDirPath, false},
		{"InvalidSrcPath", filepath.Join(tempDir, "nonexistent"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Compress(tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				tarFilePath := tt.src + ".tar"
				if _, err := os.Stat(tarFilePath); os.IsNotExist(err) {
					t.Errorf("Compress() = tar file does not exist, expected it to be created")
				}
			}
		})
	}
}

func testDecompress(t *testing.T) {
	// 创建临时目录和文件
	tempDir := t.TempDir()
	srcFilePath := filepath.Join(tempDir, "file.txt")
	os.WriteFile(srcFilePath, []byte("Hello, World!"), 0644)
	srcDirPath := filepath.Join(tempDir, "testDir")
	os.Mkdir(srcDirPath, 0755)
	os.WriteFile(filepath.Join(srcDirPath, "file1.txt"), []byte("File 1"), 0644)
	os.WriteFile(filepath.Join(srcDirPath, "file2.txt"), []byte("File 2"), 0644)

	// 压缩文件和目录
	Compress(srcFilePath)
	Compress(srcDirPath)

	tests := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{"DecompressFile", srcFilePath + ".tar", false},
		{"DecompressDirectory", srcDirPath + ".tar", false},
		{"InvalidSrcPath", filepath.Join(tempDir, "nonexistent.tar"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Decompress(tt.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decompress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				decompressedPath := strings.TrimSuffix(tt.src, ".tar")
				if _, err := os.Stat(decompressedPath); os.IsNotExist(err) {
					t.Errorf("Decompress() = decompressed path does not exist, expected it to be created")
				}
			}
		})
	}
}

func TestAll(t *testing.T) {
	testDir(t)
	testContent(t)
	testDelete(t)
	testCopy(t)
	testMove(t)
	testRename(t)
	testCompress(t)
	testDecompress(t)
}
