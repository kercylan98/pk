package xlsxer

import (
	"errors"
	"github.com/tealeg/xlsx"
)

// xlsx文件处理模型
type Xlsxer struct {

}

func NewXlsxer() *Xlsxer {
	return &Xlsxer{
		
	}
}

// 获取所有数据
func (slf *Xlsxer) GetAll(filepath string, sheetNumber, startRow int) ([]map[int]*xlsx.Cell, error) {
	if x, err := xlsx.OpenFile(filepath); err != nil {
		return nil, err
	}else {
		if (sheetNumber > len(x.Sheets) - 1) || sheetNumber < 0 {
			return nil, errors.New("exceeds the maximum number of \"sheets\"")
		}
		sheet := x.Sheets[sheetNumber]
		maxCell := len(sheet.Rows[startRow-1].Cells)
		if (startRow > sheet.MaxRow - 1) || startRow < 0 {
			return nil, errors.New("exceeds the maximum number of \"rows\"")
		}
		var result []map[int]*xlsx.Cell
		for i, row := range sheet.Rows {
			if i >= startRow {
				data := map[int]*xlsx.Cell{}
				if len(row.Cells) < maxCell {
					for ad := 0; ad < maxCell - len(row.Cells); ad ++ {
						row.AddCell()
					}
				}
				nullNum := 0
				for _, cell := range row.Cells {
					if cell.String() == "" {
						nullNum++
					}
				}
				if nullNum == len(row.Cells) {
					continue
				}
				for ci, cell := range row.Cells {
					data[ci] = cell
				}
				result = append(result, data)
			}
		}
		return result, nil
	}
}