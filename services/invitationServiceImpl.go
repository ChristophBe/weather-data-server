package services

import (
	"de.christophb.wetter/data/models"
	"de.christophb.wetter/data/repositories"
	"fmt"
)

type invitationServiceImpl struct {
	measuringNodeRepository repositories.MeasuringNodeRepository
}

func (i invitationServiceImpl) HandleInvitation(user models.User, invitationId int64) (err error) {

	nodes, err := i.measuringNodeRepository.FetchAllNodesByInvitationId(invitationId)

	if err != nil {
		err = fmt.Errorf("failed to fetch nodes by invitation invitationId:%d cause: %w", invitationId, err)
		return
	}
	for _, node := range nodes {
		err = i.measuringNodeRepository.CreateAuthorisationRelation(node, user)

		if err != nil {
			err = fmt.Errorf("failed to add auth relation nodeID: %d ,cause: %w}", node.Id, err)
			return
		}
	}
	return
}
