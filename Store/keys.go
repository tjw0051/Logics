package store

import models "github.com/tjw0051/log-go/Models"

func CreateKeys(keys models.KeysModel) error {
	for i := 0; i < len(keys); i++ {
		err := db.Create(&keys[i]).Error
		return err
	}
	return nil
}

func DeleteKeys(keys models.KeysModel) error {
	for i := 0; i < len(keys); i++ {
		//err := db.Delete(&keys[i]).Error
		err := db.Where("key = ?", keys[i].Key).Delete(models.KeyModel{}).Error
		return err
	}
	return nil
}

func GetKeys() (models.KeysModel, error) {
	allKeys := models.KeysModel{}
	err := db.Find(&allKeys).Error
	return allKeys, err
}
