package database

import (
	"errors"
	"github.com/ChristophBe/weather-data-server/data/models"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"time"
)

type invitationRepositoryImpl struct{}

func (invitationRepositoryImpl) parseInvitationFromRecord(record neo4j.Record) (interface{}, error) {

	nodeData, ok := record.Get("i")
	if !ok {
		err := errors.New("can not parse measuring form record")
		return nil, err
	}
	node := nodeData.(neo4j.Node)
	props := node.Props()

	invitation := models.Invitation{
		Id:           node.Id(),
		Email:        parseStringProp(props["email"], ""),
		CreationTime: parseTimeProp(props["creationTime"], time.Unix(0, 0)),
	}
	return invitation, nil
}

func (i invitationRepositoryImpl) SaveInvitation(invitation models.Invitation) (savedInvitation models.Invitation, err error) {
	params := map[string]interface{}{
		"email":        invitation.Email,
		"creationTime": invitation.CreationTime.Unix(),
	}
	if invitation.Id != 0 {
		params["id"] = invitation.Id
	}

	insertUserStatement := "CREATE (i:Invitation {creationTime: $creationTime,email: $email}) RETURN i"
	updateUserStatement := "Match (i:Invitation) WHERE id(i) = $id SET i.email = $email RETURN i"

	result, err := saveNode(insertUserStatement, updateUserStatement, params, parseSingleItemFromResult(i.parseInvitationFromRecord))

	if err != nil {
		return
	}
	savedInvitation = result.(models.Invitation)
	return
}

func (i invitationRepositoryImpl) FetchInvitationByEmail(email string) (invitation models.Invitation, err error) {

	params := map[string]interface{}{"email": email}

	stmt := "Match (i:Invitation) WHERE i.email = $email return i limit 1"
	results, err := doReadTransaction(stmt, params, parseSingleItemFromResult(i.parseInvitationFromRecord))

	if err != nil {
		return
	}

	invitation = results.(models.Invitation)
	return
}

func (i invitationRepositoryImpl) FetchInvitationById(invitationId int64) (invitation models.Invitation, err error) {

	params := map[string]interface{}{"id": invitationId}

	stmt := "Match (i:Invitation) WHERE id(i) = $id Return i"
	results, err := doReadTransaction(stmt, params, parseSingleItemFromResult(i.parseInvitationFromRecord))

	if err != nil {
		return
	}

	invitation = results.(models.Invitation)
	return
}

func (i invitationRepositoryImpl) AddNodeToInvitation(invitation models.Invitation, measuringNode models.MeasuringNode) (err error) {

	params := map[string]interface{}{
		"invitationId": invitation.Id,
		"nodeId":       measuringNode.Id,
	}

	stmt := "MATCH (i:Invitation), (n:MeasuringNode) WHERE id(i) = $invitationId and id(n) = $nodeId CREATE (i)-[r:INVITATION_FOR]->(n) RETURN r"
	_, err = doWriteTransaction(stmt, params, func(result neo4j.Result) (res interface{}, err error) {
		return nil, result.Err()
	})
	return
}
