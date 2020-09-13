package node

import (
	"simulator/pkg/consts"
	"simulator/pkg/job"
)

type Node struct {
	ID				int
	// resource definition
	CPU				int
	FreeCPU			int

	Mem 			int
	FreeMem			int

	GPU 			int
	FreeGPU			int

	SwitchCnt		int

	AllocatedGPU	[]bool
	NVLinkStatus	[][]bool
}

func NewNode(_id int, _CPU int, _Mem int, _GPU int, _switchCnt int) Node {
	n := Node{
		ID: _id,
		CPU: _CPU,
		FreeCPU: _CPU,
		Mem: _Mem,
		FreeMem: _Mem,
		GPU: _GPU,
		FreeGPU: _GPU,
		SwitchCnt: _switchCnt,
	}
	n.AllocatedGPU = make([]bool, _GPU)
	n.NVLinkStatus = make([][]bool, _GPU)
	for i := range n.NVLinkStatus {
		n.NVLinkStatus[i] = make([]bool, _GPU)
	}
	return n
}


func (node *Node) SetNVLink(x int, y int) {
	node.NVLinkStatus[x][y] = true
	node.NVLinkStatus[y][x] = true
}

func (node *Node) HasNVLink(x int, y int) bool {
	return node.NVLinkStatus[x][y]
}

func (node *Node) HasPCIeSwitch(x int, y int) bool {
	groupCnt := node.GPU / node.SwitchCnt
	return x / groupCnt == y / groupCnt
}

func (node *Node) GetGPUConnectivityFromIdx(x int, y int) int {
	// if x->y has NVLink
	if node.HasNVLink(x, y) {
		return NVLinkBenefit
	}
	// if x->y on the same PCIe switch
	if node.HasPCIeSwitch(x, y) {
		return PCIeSwitchBenefit
	}
	// if x-> not on the same switch
	return NoPCIeSwitchBenefit
}

func (node *Node) CanAllocate(job *job.Job) bool {
	return job.CPU <= node.FreeCPU && job.GPU <= node.FreeGPU && job.Mem <= node.FreeMem
}

func (node *Node) Allocate(job *job.Job) {
	node.FreeCPU -= job.CPU
	node.FreeGPU -= job.GPU
	node.FreeMem -= job.Mem

}

func (node *Node) GetGlobalGPUIdxFromLocal(idx int) int {
	return node.ID * consts.GPUperNode + idx
}

func GetLocalGPUIdxFromGlobal(idx int) (int, int) {
	// node id, GPU id
	return idx / consts.GPUperNode, idx % consts.GPUperNode
}

func (node *Node) AllocateGPUIdx(id int) {
	if node.AllocatedGPU[id] == true {
		panic("wrong GPU allocation")
	}
	node.AllocatedGPU[id] = true
	node.FreeGPU--
}

func (node *Node) FreeGPUIdx(id int) {
	if node.AllocatedGPU[id] == false {
		panic("wrong GPU free")
	}
	node.AllocatedGPU[id] = false
	node.FreeGPU++
}