package usecase

import (
	"aquafarm-management/app/model"
	"aquafarm-management/app/repository"
)

// FarmUsecase menangani proses bisnis item
type FarmUsecase struct {
	Repository *repository.FarmRepository
}

// NewFarmUsecase membuat instance baru FarmUsecase
func NewFarmUsecase(r *repository.FarmRepository) *FarmUsecase {
	return &FarmUsecase{
		Repository: r,
	}
}

// Fetch mengambil semua item dari repository
func (u *FarmUsecase) Fetch() ([]model.Farm, error) {
	return u.Repository.Fetch()
}

// GetById mengambil item dari repository berdasarkan ID
func (u *FarmUsecase) GetById(id string) (*model.Farm, bool, error) {
	return u.Repository.GetById(id)
}

// Store menyimpan item baru ke repository
func (u *FarmUsecase) Store(farm *model.Farm) (*model.Farm, bool, error) {
	var err error
	var getFarm *model.Farm
	var isNotFound bool

	// memanggil usecase untuk memvalidasi farm
	_, isNotFound, err = u.Repository.GetById(farm.ID)
	if !isNotFound {
		// menangani duplicate entry
		return nil, true, nil
	}
	getFarm, err = u.Repository.Store(farm)

	return getFarm, false, err
}

// UpdateById memperbarui data Farm dari repository
// @return :    - bool : isUpdate=false | isCreate=true
//			    - error : apabila ada masalah
func (u *FarmUsecase) UpdateById(farm *model.Farm) (bool, error) {
	// memanggil usecase untuk memvalidasi farm
	_, isNotFound, err := u.Repository.GetById(farm.ID)
	if isNotFound {
		_, err = u.Repository.Store(farm)
		if err != nil {
			return true, err
		} else {
			return true, nil
		}
	}
	_, err = u.Repository.UpdateById(farm)
	return false, err
}

// SoftDeleteById menghapus Farm dari repository
// @return :    - bool : isDataNotFound=true | isDataFound=true
//			    - error : apabila ada masalah
func (u *FarmUsecase) SoftDeleteById(id string) (bool, error) {
	// memanggil usecase untuk memvalidasi farm
	_, isNotFound, _ := u.Repository.GetById(id)
	if isNotFound {
		return true, nil
	}
	return false, u.Repository.SoftDeleteById(id)
}
