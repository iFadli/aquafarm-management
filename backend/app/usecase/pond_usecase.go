package usecase

import (
	"aquafarm-management/app/model"
	"aquafarm-management/app/repository"
	"errors"
)

// PondUsecase menangani proses bisnis item
type PondUsecase struct {
	Repository     *repository.PondRepository
	FarmRepository *repository.FarmRepository
}

// NewPondUsecase membuat instance baru PondUsecase
func NewPondUsecase(r *repository.PondRepository, f *repository.FarmRepository) *PondUsecase {
	return &PondUsecase{
		Repository:     r,
		FarmRepository: f,
	}
}

// Fetch mengambil semua item dari repository
func (u *PondUsecase) Fetch() ([]model.Pond, error) {
	return u.Repository.Fetch()
}

// GetById mengambil item dari repository berdasarkan ID
func (u *PondUsecase) GetById(farmId string, pondId string) (*model.Pond, bool, error) {
	farmKey, _ := u.FarmRepository.GetFarmKeyById(farmId)
	if farmKey == nil {
		return nil, false, errors.New("farm id '" + farmId + "' tidak terdaftar atau telah dinonaktifkan")
	}
	return u.Repository.GetById(*farmKey, pondId)
}

// Store menyimpan item baru ke repository
func (u *PondUsecase) Store(pond *model.Pond) (*model.Pond, bool, error) {
	var err error
	var getPond *model.Pond
	var isNotFound bool

	// memanggil usecase untuk memvalidasi farm
	getFarmKey, _ := u.FarmRepository.GetFarmKeyById(pond.FarmID)
	if getFarmKey == nil {
		return nil, false, errors.New("farm id '" + pond.FarmID + "' tidak terdaftar atau telah dinonaktifkan")
	}
	pond.FarmKey = *getFarmKey

	_, isNotFound, err = u.Repository.GetById(*getFarmKey, pond.ID)
	if !isNotFound {
		// menangani duplicate entry
		return nil, true, nil
	}
	getPond, err = u.Repository.Store(pond)

	return getPond, false, err
}

// UpdateById memperbarui data Pond dari repository
//	@return	:    - bool : isUpdate=false | isCreate=true
//			    - error : apabila ada masalah
func (u *PondUsecase) UpdateById(pond *model.Pond) (bool, error) {
	// memanggil usecase untuk memvalidasi farm
	farmKey, _ := u.FarmRepository.GetFarmKeyById(pond.FarmID)
	if farmKey == nil {
		return false, errors.New("farm id '" + pond.FarmID + "' tidak terdaftar atau telah dinonaktifkan")
	}
	pond.FarmKey = *farmKey

	// memanggil usecase untuk memvalidasi farm
	_, isNotFound, err := u.Repository.GetById(*farmKey, pond.ID)
	if isNotFound {
		_, err = u.Repository.Store(pond)
		if err != nil {
			return true, err
		} else {
			return true, nil
		}
	}
	_, err = u.Repository.UpdateById(pond)
	return false, err
}

// SoftDeleteById menghapus Pond dari repository
//	@return	:    - bool : isDataNotFound=true | isDataFound=true
//			    - error : apabila ada masalah
func (u *PondUsecase) SoftDeleteById(farmId string, pondId string) (bool, error) {
	// memanggil usecase untuk memvalidasi farm
	farmKey, _ := u.FarmRepository.GetFarmKeyById(farmId)
	if farmKey == nil {
		return false, errors.New("farm id '" + farmId + "' tidak terdaftar atau telah dinonaktifkan")
	}

	// memanggil usecase untuk memvalidasi farm
	_, isNotFound, _ := u.Repository.GetById(*farmKey, pondId)
	if isNotFound {
		return true, nil
	}
	return false, u.Repository.SoftDeleteById(farmKey, pondId)
}
