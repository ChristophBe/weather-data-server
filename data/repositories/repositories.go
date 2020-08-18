package repositories

import (
	"de.christophb.wetter/data/models"
)

type MeasuringRepository interface {
	CreateMeasurement(stationId int64, measurement models.Measurement) (savedMeasurement models.Measurement, err error)
	FetchAllMeasuringsByNodeId(nodeId int64) (measurements []models.Measurement, err error)
	FetchLastMeasuringsByNodeId(nodeId int64, hours int64) (measurements []models.Measurement, err error)
}

type MeasuringNodeRepository interface {
	FetchAllPublicNodes() ([]models.MeasuringNode, error)
	FetchNodesOwnedByUserId(ownerId int64) ([]models.MeasuringNode, error)
	CreateMeasuringNode(node models.MeasuringNode, ownerId int64) (models.MeasuringNode, error)
	FetchAllVisibleNodesByUserId(userId int64) ([]models.MeasuringNode, error)
	FetchAllNodesByInvitationId(invitationId int64) ([]models.MeasuringNode, error)
	FetchMeasuringNodeById(nodeId int64) (models.MeasuringNode, error)
	FetchAllMeasuringNodeUserRelations(nodeId int64, userId int64) ([]string, error)
	CreateAuthorisationRelation(node models.MeasuringNode, user models.User) error
}

type NodeAuthTokenRepository interface {
	InsertNodeAuthToken(nodeId int64, token models.NodeAuthToken) (models.NodeAuthToken, error)
	FetchAuthTokenByNodeId(nodeId int64) (models.NodeAuthToken, error)
}

type UserRepository interface {
	SaveUser(user models.User) (models.User, error)
	FetchUserById(userId int64) (models.User, error)
	FetchOwnerByMeasuringNode(nodeId int64) (models.User, error)
	FetchUserByEmail(email string) (models.User, error)
	HasUserWithEmail(email string) bool
	HasUserWithUsername(username string) bool
}

type InvitationRepository interface {
	SaveInvitation(invitation models.Invitation) (models.Invitation, error)
	FetchInvitationById(invitationId int64) (models.Invitation, error)
	AddNodeToInvitation(invitation models.Invitation, measuringNode models.MeasuringNode) error
	FetchInvitationByEmail(email string) (models.Invitation, error)
}
