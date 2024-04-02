package xls

import (
	"os"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// Exportable 可以被导出的数据结构
type Exportable interface {
	Len() int               // 返回数据长度
	GetHead() []any         // 获取表头
	GetRow(index int) []any // 根据索引获取行数据
}

// Generate 生成 Excel 文件
func Generate(data Exportable, filename ...string) (string, string, error) {
	// 如果数据长度为 0，则返回空字符串和 nil
	if data.Len() == 0 {
		return "", "", nil
	}

	// 创建新的 Excel 文件
	f := excelize.NewFile()
	defer f.Close() // 延迟关闭文件流

	// 创建工作表数据流
	streamWriter, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		return "", "", err
	}

	// 获取表头数据
	head := data.GetHead()

	// 设置表头数据到 Excel 文件中
	cell1, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		return "", "", err
	}
	if err := streamWriter.SetRow(cell1, head); err != nil {
		return "", "", err
	}

	// 循环遍历数据，将每行数据写入 Excel 文件
	for rowID, i := 2, 0; rowID <= data.Len(); rowID++ {
		cell, _ := excelize.CoordinatesToCellName(1, rowID)
		if err := streamWriter.SetRow(cell, data.GetRow(i)); err != nil {
			return "", "", err
		}
		i++
	}

	// 刷新工作表流中的数据
	if err := streamWriter.Flush(); err != nil {
		return "", "", err
	}

	fn := ""

	if len(filename) > 0 && filename[0] != "" {
		fn = filename[0]
	} else {
		// 生成随机文件名并保存 Excel 文件
		fn = uuid.New().String() + ".xlsx"
	}

	if err := f.SaveAs(fn); err != nil {
		return "", "", err
	}

	// 获取当前工作目录路径
	path, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	// 返回文件名和文件路径
	return fn, path + "/" + fn, nil
}
