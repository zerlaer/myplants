package database

import (
	"myplants/internal/logger"
	"myplants/internal/model"
)

func SeedDefaultPots() error {
	var count int64
	if err := DB.Model(&model.Pot{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	pots := []model.Pot{
		{Type: "青山盆", Diameter: 16, Height: 18, Size: "⌀16×18cm", Name: "青山盆16cm", Status: "空闲"},
		{Type: "青山盆", Diameter: 18, Height: 20, Size: "⌀18×20cm", Name: "青山盆18cm", Status: "空闲"},
		{Type: "青山盆", Diameter: 21, Height: 22, Size: "⌀21×22cm", Name: "青山盆21cm", Status: "空闲"},
		{Type: "青山盆", Diameter: 23, Height: 24, Size: "⌀23×24cm", Name: "青山盆23cm", Status: "空闲"},

		{Type: "加仑盆", Gallon: 1, Size: "1加仑", Name: "加仑盆1加仑", Status: "空闲"},
		{Type: "加仑盆", Gallon: 1.5, Size: "1.5加仑", Name: "加仑盆1.5加仑", Status: "空闲"},
		{Type: "加仑盆", Gallon: 2, Size: "2加仑", Name: "加仑盆2加仑", Status: "空闲"},
		{Type: "加仑盆", Gallon: 2.5, Size: "2.5加仑", Name: "加仑盆2.5加仑", Status: "空闲"},

		{Type: "透气盆", Diameter: 10, Height: 10, Size: "⌀10×10cm", Name: "透气盆10cm", Status: "空闲"},
		{Type: "透气盆", Diameter: 12, Height: 12, Size: "⌀12×12cm", Name: "透气盆12cm", Status: "空闲"},
		{Type: "透气盆", Diameter: 14, Height: 14, Size: "⌀14×14cm", Name: "透气盆14cm", Status: "空闲"},
		{Type: "透气盆", Diameter: 16, Height: 16, Size: "⌀16×16cm", Name: "透气盆16cm", Status: "空闲"},
	}

	for i := range pots {
		if err := DB.Create(&pots[i]).Error; err != nil {
			return err
		}
	}

	logger.S().Infof("种子数据初始化完成 预置花盆数=%d", len(pots))
	return nil
}
