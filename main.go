package main

import "PvP-Go/db/daos"

func main() {
	cpms := []float64{0.094, 0.135137432, 0.16639787, 0.192650919, 0.21573247, 0.236572661, 0.25572005, 0.273530381, 0.29024988, 0.306057378, 0.3210876, 0.335445036, 0.34921268, 0.362457751, 0.3752356, 0.387592416, 0.39956728, 0.411193551, 0.4225, 0.432926409, 0.44310755, 0.453059959, 0.4627984, 0.472336093, 0.48168495, 0.4908558, 0.49985844, 0.508701765, 0.51739395, 0.525942511, 0.5343543, 0.542635738, 0.5507927, 0.558830586, 0.5667545, 0.574569133, 0.5822789, 0.589887907, 0.5974, 0.604823665, 0.6121573, 0.619404122, 0.6265671, 0.633649143, 0.64065295, 0.647580967, 0.65443563, 0.661219252, 0.667934, 0.674581896, 0.6811649, 0.687684904, 0.69414365, 0.70054287, 0.7068842, 0.713169109, 0.7193991, 0.725575614, 0.7317, 0.734741009, 0.7377695, 0.740785594, 0.74378943, 0.746781211, 0.74976104, 0.752729087, 0.7556855, 0.758630368, 0.76156384, 0.764486065, 0.76739717, 0.770297266, 0.7731865, 0.776064962, 0.77893275, 0.781790055, 0.784637, 0.787473608, 0.7903}
	level := float64(1.0)
	for _, cpm := range cpms {
		err, _ := daos.CP_DAO.Upsert(level, cpm)
		daos.CheckError(err)
		level += 0.5
	}
}
