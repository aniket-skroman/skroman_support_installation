package helper

import (
	"github.com/aniket-skroman/skroman_support_installation/utils"
)

func SetPaginationData(page int, total int64) {
	if page == 0 {
		utils.PREVIOUS_IDX = 0
	} else {
		utils.PREVIOUS_IDX = page - 1
	}

	utils.CURRENT_IDX = page
	utils.TOTALCOUNT = total
}
