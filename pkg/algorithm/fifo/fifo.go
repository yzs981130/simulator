package fifo

import (
	"simulator/pkg/consts"
	"simulator/pkg/job"
	"simulator/pkg/node"
)

func reconcileJobStatus(jobs []job.Job) {
	for i, currentJob := range jobs {
		// last term: running
		// add running time
		if currentJob.Status == job.Running {
			jobs[i].RunningTime += consts.SessionInterval
			// if job complete, change to finishing
			if jobs[i].RunningTime > jobs[i].CompletionTime {
				jobs[i].Status = job.Finishing
			}
		}
	}
}

func releaseFinishingGPU(jobs []job.Job, nodes []node.Node) {
	for i, currentJob := range jobs {
		if currentJob.Status == job.Finishing {
			// release currentJob's GPU
			for _, GPUIdx := range currentJob.GPUIndex {
				nodeIdx, GPUid := node.GetLocalGPUIdxFromGlobal(GPUIdx)
				// free GPU
				nodes[nodeIdx].FreeGPUIdx(GPUid)
			}
			// change job status to completed
			jobs[i].Status = job.Completed
		}
	}
}

func getPendingJob(currentTime int, jobs []job.Job) []job.Job {
	var pendingJobs []job.Job
	for _, currentJob := range jobs {
		if currentJob.SubmitTime > currentTime {
			break
		}
		if currentJob.Status == job.Pending {
			pendingJobs = append(pendingJobs, currentJob)
		}
	}
	return pendingJobs
}

func allocatePendingJobs(currentTime int, jobs []job.Job, nodes []node.Node) {
	for _, currentJob := range jobs {
		if currentJob.SubmitTime > currentTime {
			break
		}
		if currentJob.Status == job.Pending {
			// try to allocate currentJob
			if ok, idx := tryAllocateJob(currentJob, nodes); ok {
				// can allocate
				allocateJobToNode(nodes, idx)
			} else {
				// fifo return when cannot allocate currentJob
				return
			}
		}
	}
}

func tryAllocateJob(currentJob job.Job, nodes []node.Node) (bool, []int) {
	restGPU := currentJob.GPU
	var gpuIdx []int
	for _, currentNode := range nodes {
		if currentNode.FreeGPU > 0 {
			if currentNode.FreeGPU >= restGPU {
				// try allocate restGPU
				for j := range currentNode.AllocatedGPU {
					if currentNode.AllocatedGPU[j] == false {
						gpuIdx = append(gpuIdx, currentNode.GetGlobalGPUIdxFromLocal(j))
						restGPU--
					}
					if restGPU == 0 {
						return true, gpuIdx
					}
				}

			} else {
				// try allocate freeGPU
				for j := range currentNode.AllocatedGPU {
					if currentNode.AllocatedGPU[j] == false {
						gpuIdx = append(gpuIdx, currentNode.GetGlobalGPUIdxFromLocal(j))
						restGPU--
					}
				}
			}
		}
	}
	return false, nil
}

func allocateJobToNode(nodes []node.Node, idx []int) {
	for _, currentIdx := range idx {
		nodeId, GPUId := node.GetLocalGPUIdxFromGlobal(currentIdx)
		nodes[nodeId].AllocateGPUIdx(GPUId)
	}
}

func SessionRun(currentTime int, jobs []job.Job, nodes []node.Node) {
	reconcileJobStatus(jobs)
	releaseFinishingGPU(jobs, nodes)
	allocatePendingJobs(currentTime, jobs, nodes)
}