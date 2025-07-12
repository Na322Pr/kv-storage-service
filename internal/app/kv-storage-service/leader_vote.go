package kv_storage_service

import (
	"context"
	desc "github.com/Na322Pr/kv-storage-service/pkg/api"
)

func (s *Implementation) LeaderVote(
	ctx context.Context,
	req *desc.LeaderVoteRequest,
) (*desc.LeaderVoteResponse, error) {
	return &desc.LeaderVoteResponse{
		VoteGranted: s.nodeService.HandleVoteRequest(
			req.GetCandidateAddress(),
			req.GetTerm(),
		),
		Term: s.nodeService.GetTerm(),
	}, nil
}
