package ussn

type cpReg struct {
	add  bool
	ussn *ussn
}

type cpWater struct{}

type pcWater struct {
	ct    int
	water []string
}

type pcSc struct {
	msg interface{}
}

type pcUpdateInfo struct{}
